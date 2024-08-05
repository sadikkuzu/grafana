package provisioning

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"unsafe"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/notifier/legacy_storage"
	"github.com/grafana/grafana/pkg/services/ngalert/provisioning/validation"
)

type TemplateService struct {
	configStore     alertmanagerConfigStore
	provenanceStore ProvisioningStore
	xact            TransactionManager
	log             log.Logger
	validator       validation.ProvenanceStatusTransitionValidator
}

func NewTemplateService(config alertmanagerConfigStore, prov ProvisioningStore, xact TransactionManager, log log.Logger) *TemplateService {
	return &TemplateService{
		configStore:     config,
		provenanceStore: prov,
		xact:            xact,
		validator:       validation.ValidateProvenanceRelaxed,
		log:             log,
	}
}

func (t *TemplateService) GetTemplates(ctx context.Context, orgID int64) ([]definitions.NotificationTemplate, error) {
	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return nil, err
	}

	templates := make([]definitions.NotificationTemplate, 0, len(revision.Config.TemplateFiles))
	for name, tmpl := range revision.Config.TemplateFiles {
		tmpl := definitions.NotificationTemplate{
			UID:             legacy_storage.NameToUid(name),
			Name:            name,
			Template:        tmpl,
			ResourceVersion: calculateTemplateFingerprint(tmpl),
		}

		provenance, err := t.provenanceStore.GetProvenance(ctx, &tmpl, orgID)
		if err != nil {
			return nil, err
		}
		tmpl.Provenance = definitions.Provenance(provenance)

		templates = append(templates, tmpl)
	}

	return templates, nil
}

func (t *TemplateService) GetTemplate(ctx context.Context, orgID int64, nameOrUid string) (definitions.NotificationTemplate, error) {
	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}

	name := nameOrUid
	content, ok := revision.Config.TemplateFiles[nameOrUid]
	if !ok {
		name, content, ok = getTemplateByUid(revision.Config.TemplateFiles, nameOrUid)
	}
	if !ok {
		return definitions.NotificationTemplate{}, ErrTemplateNotFound.Errorf("")
	}

	template := definitions.NotificationTemplate{
		UID:             legacy_storage.NameToUid(name),
		Name:            name,
		Template:        content,
		ResourceVersion: calculateTemplateFingerprint(content),
	}

	provenance, err := t.provenanceStore.GetProvenance(ctx, &template, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}
	template.Provenance = definitions.Provenance(provenance)

	return template, nil
}

func (t *TemplateService) CreateTemplate(ctx context.Context, orgID int64, tmpl definitions.NotificationTemplate) (definitions.NotificationTemplate, error) {
	err := tmpl.Validate()
	if err != nil {
		return definitions.NotificationTemplate{}, MakeErrTemplateInvalid(err)
	}

	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}
	return t.createTemplate(ctx, revision, orgID, tmpl)
}

func (t *TemplateService) createTemplate(ctx context.Context, revision *legacy_storage.ConfigRevision, orgID int64, tmpl definitions.NotificationTemplate) (definitions.NotificationTemplate, error) {
	if revision.Config.TemplateFiles == nil {
		revision.Config.TemplateFiles = map[string]string{}
	}

	_, found := revision.Config.TemplateFiles[tmpl.Name]
	if found {
		return definitions.NotificationTemplate{}, ErrTemplateExists.Errorf("")
	}

	revision.Config.TemplateFiles[tmpl.Name] = tmpl.Template

	err := t.xact.InTransaction(ctx, func(ctx context.Context) error {
		if err := t.configStore.Save(ctx, revision, orgID); err != nil {
			return err
		}
		return t.provenanceStore.SetProvenance(ctx, &tmpl, orgID, models.Provenance(tmpl.Provenance))
	})
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}

	return definitions.NotificationTemplate{
		UID:             legacy_storage.NameToUid(tmpl.Name),
		Name:            tmpl.Name,
		Template:        tmpl.Template,
		Provenance:      tmpl.Provenance,
		ResourceVersion: calculateTemplateFingerprint(tmpl.Template),
	}, nil
}

func (t *TemplateService) UpdateTemplate(ctx context.Context, orgID int64, tmpl definitions.NotificationTemplate) (definitions.NotificationTemplate, error) {
	err := tmpl.Validate()
	if err != nil {
		return definitions.NotificationTemplate{}, MakeErrTemplateInvalid(err)
	}

	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}
	return t.updateTemplate(ctx, revision, orgID, tmpl)
}

func (t *TemplateService) updateTemplate(ctx context.Context, revision *legacy_storage.ConfigRevision, orgID int64, tmpl definitions.NotificationTemplate) (definitions.NotificationTemplate, error) {
	if revision.Config.TemplateFiles == nil {
		revision.Config.TemplateFiles = map[string]string{}
	}

	var found bool
	var existingName, exisitingContent string
	if tmpl.UID != "" {
		existingName, exisitingContent, found = getTemplateByUid(revision.Config.TemplateFiles, tmpl.UID)
	} else {
		existingName = tmpl.Name
		exisitingContent, found = revision.Config.TemplateFiles[existingName]
	}
	if !found {
		return definitions.NotificationTemplate{}, ErrTemplateNotFound.Errorf("")
	}

	// check that provenance is not changed in an invalid way
	storedProvenance, err := t.provenanceStore.GetProvenance(ctx, &tmpl, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}
	if err := t.validator(storedProvenance, models.Provenance(tmpl.Provenance)); err != nil {
		return definitions.NotificationTemplate{}, err
	}

	err = t.checkOptimisticConcurrency(tmpl.Name, exisitingContent, models.Provenance(tmpl.Provenance), tmpl.ResourceVersion, "update")
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}

	revision.Config.TemplateFiles[tmpl.Name] = tmpl.Template

	err = t.xact.InTransaction(ctx, func(ctx context.Context) error {
		if existingName != tmpl.Name {
			delete(revision.Config.TemplateFiles, existingName)
			err := t.provenanceStore.DeleteProvenance(ctx, &definitions.NotificationTemplate{Name: existingName}, orgID)
			if err != nil {
				return err
			}
		}

		if err := t.configStore.Save(ctx, revision, orgID); err != nil {
			return err
		}
		return t.provenanceStore.SetProvenance(ctx, &tmpl, orgID, models.Provenance(tmpl.Provenance))
	})
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}

	return definitions.NotificationTemplate{
		UID:             legacy_storage.NameToUid(tmpl.Name),
		Name:            tmpl.Name,
		Template:        tmpl.Template,
		Provenance:      tmpl.Provenance,
		ResourceVersion: calculateTemplateFingerprint(tmpl.Template),
	}, nil
}

func (t *TemplateService) UpsertTemplate(ctx context.Context, orgID int64, tmpl definitions.NotificationTemplate) (definitions.NotificationTemplate, error) {
	err := tmpl.Validate()
	if err != nil {
		return definitions.NotificationTemplate{}, MakeErrTemplateInvalid(err)
	}

	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return definitions.NotificationTemplate{}, err
	}

	d, err := t.updateTemplate(ctx, revision, orgID, tmpl)
	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return t.createTemplate(ctx, revision, orgID, tmpl)
		}
	}
	return d, err
}

func (t *TemplateService) DeleteTemplate(ctx context.Context, orgID int64, nameOrUid string, provenance definitions.Provenance, version string) error {
	revision, err := t.configStore.Get(ctx, orgID)
	if err != nil {
		return err
	}

	if revision.Config.TemplateFiles == nil {
		return nil
	}

	existingName := nameOrUid
	existing, ok := revision.Config.TemplateFiles[nameOrUid]
	if !ok {
		existingName, existing, ok = getTemplateByUid(revision.Config.TemplateFiles, nameOrUid)
	}
	if !ok {
		return nil
	}

	err = t.checkOptimisticConcurrency(existingName, existing, models.Provenance(provenance), version, "delete")
	if err != nil {
		return err
	}

	// check that provenance is not changed in an invalid way
	storedProvenance, err := t.provenanceStore.GetProvenance(ctx, &definitions.NotificationTemplate{Name: existingName}, orgID)
	if err != nil {
		return err
	}
	if err = t.validator(storedProvenance, models.Provenance(provenance)); err != nil {
		return err
	}

	delete(revision.Config.TemplateFiles, existingName)

	return t.xact.InTransaction(ctx, func(ctx context.Context) error {
		if err := t.configStore.Save(ctx, revision, orgID); err != nil {
			return err
		}
		tgt := definitions.NotificationTemplate{
			Name: existingName,
		}
		return t.provenanceStore.DeleteProvenance(ctx, &tgt, orgID)
	})
}

func (t *TemplateService) checkOptimisticConcurrency(name, currentContent string, provenance models.Provenance, desiredVersion string, action string) error {
	if desiredVersion == "" {
		if provenance != models.ProvenanceFile {
			// if version is not specified and it's not a file provisioning, emit a log message to reflect that optimistic concurrency is disabled for this request
			t.log.Debug("ignoring optimistic concurrency check because version was not provided", "template", name, "operation", action)
		}
		return nil
	}
	currentVersion := calculateTemplateFingerprint(currentContent)
	if currentVersion != desiredVersion {
		return ErrVersionConflict.Errorf("provided version %s of template %s does not match current version %s", desiredVersion, name, currentVersion)
	}
	return nil
}

func calculateTemplateFingerprint(t string) string {
	sum := fnv.New64()
	_, _ = sum.Write(unsafe.Slice(unsafe.StringData(t), len(t))) //nolint:gosec
	return fmt.Sprintf("%016x", sum.Sum64())
}

func getTemplateByUid(templates map[string]string, uid string) (string, string, bool) {
	for n, tmpl := range templates {
		if legacy_storage.NameToUid(n) == uid {
			return n, tmpl, true
		}
	}
	return "", "", false
}
