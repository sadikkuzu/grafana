// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by client-gen. DO NOT EDIT.

package v0alpha1

import (
	"context"

	v0alpha1 "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	alertingnotificationsv0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/alerting_notifications/v0alpha1"
	scheme "github.com/grafana/grafana/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ReceiversGetter has a method to return a ReceiverInterface.
// A group's client should implement this interface.
type ReceiversGetter interface {
	Receivers(namespace string) ReceiverInterface
}

// ReceiverInterface has methods to work with Receiver resources.
type ReceiverInterface interface {
	Create(ctx context.Context, receiver *v0alpha1.Receiver, opts v1.CreateOptions) (*v0alpha1.Receiver, error)
	Update(ctx context.Context, receiver *v0alpha1.Receiver, opts v1.UpdateOptions) (*v0alpha1.Receiver, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v0alpha1.Receiver, error)
	List(ctx context.Context, opts v1.ListOptions) (*v0alpha1.ReceiverList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.Receiver, err error)
	Apply(ctx context.Context, receiver *alertingnotificationsv0alpha1.ReceiverApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.Receiver, err error)
	ReceiverExpansion
}

// receivers implements ReceiverInterface
type receivers struct {
	*gentype.ClientWithListAndApply[*v0alpha1.Receiver, *v0alpha1.ReceiverList, *alertingnotificationsv0alpha1.ReceiverApplyConfiguration]
}

// newReceivers returns a Receivers
func newReceivers(c *NotificationsV0alpha1Client, namespace string) *receivers {
	return &receivers{
		gentype.NewClientWithListAndApply[*v0alpha1.Receiver, *v0alpha1.ReceiverList, *alertingnotificationsv0alpha1.ReceiverApplyConfiguration](
			"receivers",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v0alpha1.Receiver { return &v0alpha1.Receiver{} },
			func() *v0alpha1.ReceiverList { return &v0alpha1.ReceiverList{} }),
	}
}
