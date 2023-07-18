/*
Copyright 2023 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ClusterGroupController interface for managing ClusterGroup resources.
type ClusterGroupController interface {
	generic.ControllerInterface[*v1alpha1.ClusterGroup, *v1alpha1.ClusterGroupList]
}

// ClusterGroupClient interface for managing ClusterGroup resources in Kubernetes.
type ClusterGroupClient interface {
	generic.ClientInterface[*v1alpha1.ClusterGroup, *v1alpha1.ClusterGroupList]
}

// ClusterGroupCache interface for retrieving ClusterGroup resources in memory.
type ClusterGroupCache interface {
	generic.CacheInterface[*v1alpha1.ClusterGroup]
}

type ClusterGroupStatusHandler func(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) (v1alpha1.ClusterGroupStatus, error)

type ClusterGroupGeneratingHandler func(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) ([]runtime.Object, v1alpha1.ClusterGroupStatus, error)

func RegisterClusterGroupStatusHandler(ctx context.Context, controller ClusterGroupController, condition condition.Cond, name string, handler ClusterGroupStatusHandler) {
	statusHandler := &clusterGroupStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

func RegisterClusterGroupGeneratingHandler(ctx context.Context, controller ClusterGroupController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterGroupGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterGroupGeneratingHandler{
		ClusterGroupGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterGroupStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterGroupStatusHandler struct {
	client    ClusterGroupClient
	condition condition.Cond
	handler   ClusterGroupStatusHandler
}

func (a *clusterGroupStatusHandler) sync(key string, obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type clusterGroupGeneratingHandler struct {
	ClusterGroupGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterGroupGeneratingHandler) Remove(key string, obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.ClusterGroup{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterGroupGeneratingHandler) Handle(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) (v1alpha1.ClusterGroupStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterGroupGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}