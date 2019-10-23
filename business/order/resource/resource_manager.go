package resource

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/go-redsync/redsync"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
)

type sLock struct {
	key string
	mutex *redsync.Mutex
}

func (this *sLock) Unlock() {
	this.mutex.Unlock()
}

type ResourceManager struct {
	eel.ServiceBase
	
	Resources []business.IResource
	locks []*sLock
}

func NewResourceManager(ctx context.Context) *ResourceManager {
	service := new(ResourceManager)
	service.Ctx = ctx
	service.Resources = make([]business.IResource, 0)
	return service
}

func (this *ResourceManager) AddResource(resource business.IResource) {
	this.Resources = append(this.Resources, resource)
}

func (this *ResourceManager) AddResources(resources []business.IResource) {
	for _, resource := range resources {
		this.Resources = append(this.Resources, resource)
	}
}

func (this *ResourceManager) GroupResourceBySupplier() []*ResourceGroup {
	supplier2group := NewGroupResourceService(this.Ctx).Group(this.Resources)
	
	groups := make([]*ResourceGroup, 0)
	for _, group := range supplier2group {
		groups = append(groups, group)
	}
	
	return groups
}

func (this *ResourceManager) Validate() error {
	for _, resource := range this.Resources {
		error := resource.IsValid()
		if error != nil {
			return error
		}
	}
	
	return nil
}

func (this *ResourceManager) Lock() error {
	user := account.GetUserFromContext(this.Ctx)
	for _, resource := range this.Resources {
		if resource.IsNeedLockWhenConsume() {
			key := fmt.Sprintf("%s-%s-%d", resource.GetType(), resource.GetLockName(), user.Id)
			mutex, err := eel.Lock.Lock(key)
			if err != nil {
				eel.Logger.Error(err)
				return err
			} else {
				if mutex != nil {
					eel.Logger.Debug("[resource_manager]: lock resource: " + key)
					this.locks = append(this.locks, &sLock{
						key: key,
						mutex: mutex,
					})
				}
			}
		}
	}
	
	return nil
}

func (this *ResourceManager) Unlock() error {
	for _, lock := range this.locks {
		lock.Unlock()
		eel.Logger.Debug("[resource_manager]: unlock " + lock.key)
	}
	
	return nil
}

//Allocate 申请订单中涉及的资源
func (this *ResourceManager) AllocateForOrder(newOrder business.IOrder) error {
	return NewAllocateResourceService(this.Ctx).Allocate(this.Resources, newOrder)
}

//ReleaseAllocatedResources 释放订单中涉及的资源
func (this *ResourceManager) ReleaseAllocatedResources() {
	NewAllocateResourceService(this.Ctx).Release(this.Resources)
}

func (this *ResourceManager) GetProductResources() []business.IResource {
	productResources := make([]business.IResource, 0)
	for _, resource := range this.Resources {
		if resource.GetType() == RESOURCE_TYPE_PRODUCT {
			productResources = append(productResources, resource.(*ProductResource))
		}
	}
	
	return productResources
}

func (this *ResourceManager) GetNonProductResources() []business.IResource {
	resources := make([]business.IResource, 0)
	for _, resource := range this.Resources {
		if resource.GetType() == RESOURCE_TYPE_IMONEY {
			resources = append(resources, resource.(*IMoneyResource))
		}
	}
	
	return resources
}

func init() {
}
