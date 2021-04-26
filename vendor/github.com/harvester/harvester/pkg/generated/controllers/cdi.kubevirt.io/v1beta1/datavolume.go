/*
Copyright 2021 Rancher Labs, Inc.

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

package v1beta1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	v1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type DataVolumeHandler func(string, *v1beta1.DataVolume) (*v1beta1.DataVolume, error)

type DataVolumeController interface {
	generic.ControllerMeta
	DataVolumeClient

	OnChange(ctx context.Context, name string, sync DataVolumeHandler)
	OnRemove(ctx context.Context, name string, sync DataVolumeHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() DataVolumeCache
}

type DataVolumeClient interface {
	Create(*v1beta1.DataVolume) (*v1beta1.DataVolume, error)
	Update(*v1beta1.DataVolume) (*v1beta1.DataVolume, error)
	UpdateStatus(*v1beta1.DataVolume) (*v1beta1.DataVolume, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.DataVolume, error)
	List(namespace string, opts metav1.ListOptions) (*v1beta1.DataVolumeList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.DataVolume, err error)
}

type DataVolumeCache interface {
	Get(namespace, name string) (*v1beta1.DataVolume, error)
	List(namespace string, selector labels.Selector) ([]*v1beta1.DataVolume, error)

	AddIndexer(indexName string, indexer DataVolumeIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.DataVolume, error)
}

type DataVolumeIndexer func(obj *v1beta1.DataVolume) ([]string, error)

type dataVolumeController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewDataVolumeController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) DataVolumeController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &dataVolumeController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromDataVolumeHandlerToHandler(sync DataVolumeHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1beta1.DataVolume
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1beta1.DataVolume))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *dataVolumeController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1beta1.DataVolume))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateDataVolumeDeepCopyOnChange(client DataVolumeClient, obj *v1beta1.DataVolume, handler func(obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error)) (*v1beta1.DataVolume, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *dataVolumeController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *dataVolumeController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *dataVolumeController) OnChange(ctx context.Context, name string, sync DataVolumeHandler) {
	c.AddGenericHandler(ctx, name, FromDataVolumeHandlerToHandler(sync))
}

func (c *dataVolumeController) OnRemove(ctx context.Context, name string, sync DataVolumeHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromDataVolumeHandlerToHandler(sync)))
}

func (c *dataVolumeController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *dataVolumeController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *dataVolumeController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *dataVolumeController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *dataVolumeController) Cache() DataVolumeCache {
	return &dataVolumeCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *dataVolumeController) Create(obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error) {
	result := &v1beta1.DataVolume{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *dataVolumeController) Update(obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error) {
	result := &v1beta1.DataVolume{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *dataVolumeController) UpdateStatus(obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error) {
	result := &v1beta1.DataVolume{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *dataVolumeController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *dataVolumeController) Get(namespace, name string, options metav1.GetOptions) (*v1beta1.DataVolume, error) {
	result := &v1beta1.DataVolume{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *dataVolumeController) List(namespace string, opts metav1.ListOptions) (*v1beta1.DataVolumeList, error) {
	result := &v1beta1.DataVolumeList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *dataVolumeController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *dataVolumeController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1beta1.DataVolume, error) {
	result := &v1beta1.DataVolume{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type dataVolumeCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *dataVolumeCache) Get(namespace, name string) (*v1beta1.DataVolume, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1beta1.DataVolume), nil
}

func (c *dataVolumeCache) List(namespace string, selector labels.Selector) (ret []*v1beta1.DataVolume, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.DataVolume))
	})

	return ret, err
}

func (c *dataVolumeCache) AddIndexer(indexName string, indexer DataVolumeIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1beta1.DataVolume))
		},
	}))
}

func (c *dataVolumeCache) GetByIndex(indexName, key string) (result []*v1beta1.DataVolume, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1beta1.DataVolume, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1beta1.DataVolume))
	}
	return result, nil
}

type DataVolumeStatusHandler func(obj *v1beta1.DataVolume, status v1beta1.DataVolumeStatus) (v1beta1.DataVolumeStatus, error)

type DataVolumeGeneratingHandler func(obj *v1beta1.DataVolume, status v1beta1.DataVolumeStatus) ([]runtime.Object, v1beta1.DataVolumeStatus, error)

func RegisterDataVolumeStatusHandler(ctx context.Context, controller DataVolumeController, condition condition.Cond, name string, handler DataVolumeStatusHandler) {
	statusHandler := &dataVolumeStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromDataVolumeHandlerToHandler(statusHandler.sync))
}

func RegisterDataVolumeGeneratingHandler(ctx context.Context, controller DataVolumeController, apply apply.Apply,
	condition condition.Cond, name string, handler DataVolumeGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &dataVolumeGeneratingHandler{
		DataVolumeGeneratingHandler: handler,
		apply:                       apply,
		name:                        name,
		gvk:                         controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterDataVolumeStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type dataVolumeStatusHandler struct {
	client    DataVolumeClient
	condition condition.Cond
	handler   DataVolumeStatusHandler
}

func (a *dataVolumeStatusHandler) sync(key string, obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error) {
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

type dataVolumeGeneratingHandler struct {
	DataVolumeGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *dataVolumeGeneratingHandler) Remove(key string, obj *v1beta1.DataVolume) (*v1beta1.DataVolume, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1beta1.DataVolume{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *dataVolumeGeneratingHandler) Handle(obj *v1beta1.DataVolume, status v1beta1.DataVolumeStatus) (v1beta1.DataVolumeStatus, error) {
	objs, newStatus, err := a.DataVolumeGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}