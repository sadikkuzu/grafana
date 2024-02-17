package v0alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	openapi "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

// Generic query request with shared time across all values
// Copied from: https://github.com/grafana/grafana/blob/main/pkg/api/dtos/models.go#L62
type GenericQueryRequest struct {
	metav1.TypeMeta `json:",inline"`

	// From Start time in epoch timestamps in milliseconds or relative using Grafana time units.
	// example: now-1h
	From string `json:"from,omitempty"`

	// To End time in epoch timestamps in milliseconds or relative using Grafana time units.
	// example: now
	To string `json:"to,omitempty"`

	// queries.refId – Specifies an identifier of the query. Is optional and default to “A”.
	// queries.datasourceId – Specifies the data source to be queried. Each query in the request must have an unique datasourceId.
	// queries.maxDataPoints - Species maximum amount of data points that dashboard panel can render. Is optional and default to 100.
	// queries.intervalMs - Specifies the time interval in milliseconds of time series. Is optional and defaults to 1000.
	// required: true
	// example: [ { "refId": "A", "intervalMs": 86400000, "maxDataPoints": 1092, "datasource":{ "uid":"PD8C576611E62080A" }, "rawSql": "SELECT 1 as valueOne, 2 as valueTwo", "format": "table" } ]
	Queries []GenericDataQuery `json:"queries"`

	// required: false
	Debug bool `json:"debug,omitempty"`
}

// GenericDataQuery is a replacement for `dtos.MetricRequest` that provides more explicit types
type GenericDataQuery struct {
	schema.CommonQueryProperties `json:",inline"`

	// Additional Properties (that live at the root)
	Additional map[string]any `json:",inline"`
}

// Produce an API definition that represents map[string]any
func (g GenericDataQuery) OpenAPIDefinition() openapi.OpenAPIDefinition {
	s := spec.Schema{}
	_ = json.Unmarshal(schema.GetCommonJSONSchema(), &s)
	s.SchemaProps.Type = []string{"object"}
	s.SchemaProps.AdditionalProperties = &spec.SchemaOrBool{Allows: true}
	s.VendorExtensible = spec.VendorExtensible{
		Extensions: map[string]interface{}{
			"x-kubernetes-preserve-unknown-fields": true,
		},
	}
	return openapi.OpenAPIDefinition{Schema: s}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericDataQuery.
func (g *GenericDataQuery) DeepCopy() *GenericDataQuery {
	if g == nil {
		return nil
	}
	out := new(GenericDataQuery)
	jj, err := json.Marshal(g)
	if err != nil {
		_ = json.Unmarshal(jj, out)
	}
	return out
}

func (g *GenericDataQuery) DeepCopyInto(out *GenericDataQuery) {
	clone := g.DeepCopy()
	*out = *clone
}

func NewGenericDataQuery(vals map[string]any) GenericDataQuery {
	q := GenericDataQuery{}
	_ = q.unmarshal(vals)
	return q
}

// MarshalJSON ensures that the unstructured object produces proper
// JSON when passed to Go's standard JSON library.
func (g GenericDataQuery) MarshalJSON() ([]byte, error) {
	vals := map[string]any{}
	if g.Additional != nil {
		for k, v := range g.Additional {
			vals[k] = v
		}
	}

	vals["refId"] = g.RefID
	if g.Datasource != nil && (g.Datasource.Type != "" || g.Datasource.UID != "") {
		vals["datasource"] = g.Datasource
	}
	if g.DatasourceId > 0 {
		vals["datasourceId"] = g.DatasourceId
	}
	if g.IntervalMS > 0 {
		vals["intervalMs"] = g.IntervalMS
	}
	if g.MaxDataPoints > 0 {
		vals["maxDataPoints"] = g.MaxDataPoints
	}
	return json.Marshal(vals)
}

// UnmarshalJSON ensures that the unstructured object properly decodes
// JSON when passed to Go's standard JSON library.
func (g *GenericDataQuery) UnmarshalJSON(b []byte) error {
	vals := map[string]any{}
	err := json.Unmarshal(b, &vals)
	if err != nil {
		return err
	}
	return g.unmarshal(vals)
}

func (g *GenericDataQuery) unmarshal(vals map[string]any) error {
	if vals == nil {
		g.Additional = nil
		return nil
	}

	key := "refId"
	v, ok := vals[key]
	if ok {
		g.RefID, ok = v.(string)
		if !ok {
			return fmt.Errorf("expected string refid (got: %t)", v)
		}
		delete(vals, key)
	}

	key = "datasource"
	v, ok = vals[key]
	if ok {
		wrap, ok := v.(map[string]any)
		if ok {
			g.Datasource = &schema.DataSourceRef{}
			g.Datasource.Type, _ = wrap["type"].(string)
			g.Datasource.UID, _ = wrap["uid"].(string)
			delete(vals, key)
		} else {
			// Old old queries may arrive with just the name
			name, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected datasource as object (got: %t)", v)
			}
			g.Datasource = &schema.DataSourceRef{}
			g.Datasource.UID = name // Not great, but the lookup function will try its best to resolve
			delete(vals, key)
		}
	}

	key = "intervalMs"
	v, ok = vals[key]
	if ok {
		g.IntervalMS, ok = v.(float64)
		if !ok {
			return fmt.Errorf("expected intervalMs as float (got: %t)", v)
		}
		delete(vals, key)
	}

	key = "maxDataPoints"
	v, ok = vals[key]
	if ok {
		count, ok := v.(float64)
		if !ok {
			return fmt.Errorf("expected maxDataPoints as number (got: %t)", v)
		}
		g.MaxDataPoints = int64(count)
		delete(vals, key)
	}

	key = "datasourceId"
	v, ok = vals[key]
	if ok {
		count, ok := v.(float64)
		if !ok {
			return fmt.Errorf("expected datasourceId as number (got: %t)", v)
		}
		g.DatasourceId = int64(count)
		delete(vals, key)
	}

	g.Additional = vals
	return nil
}
