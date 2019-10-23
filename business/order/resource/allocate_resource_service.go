package resource

import (
	"context"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type AllocateResourceService struct {
	eel.ServiceBase
	
	resource2allocator map[string]business.IResourceAllocator
}

func NewAllocateResourceService(ctx context.Context) *AllocateResourceService {
	service := new(AllocateResourceService)
	service.Ctx = ctx
	
	service.initAllocator()
	
	return service
}

func (this *AllocateResourceService) initAllocator() {
	this.resource2allocator = map[string]business.IResourceAllocator{
		RESOURCE_TYPE_PRODUCT: NewProductResourceAllocator(this.Ctx),
		RESOURCE_TYPE_IMONEY: NewIMoneyResourceAllocator(this.Ctx),
		RESOURCE_TYPE_COUPON: NewCouponResourceAllocator(this.Ctx),
	}
}

func (this *AllocateResourceService) Allocate(resources []business.IResource, newOrder business.IOrder) error {
	isSuccess := false
	defer func() {
		if !isSuccess {
			this.Release(resources)
		}
	}()
	
	for _, resource := range resources {
		if allocator, ok := this.resource2allocator[resource.GetType()]; ok {
			err := allocator.Allocate(resource, newOrder)
			if err != nil {
				eel.Logger.Error(err)
				return err
			}
			
			resource.SetAllocated()
			eel.Logger.Debug("[resource] allocate resource " + resource.GetType())
		}
	}
	isSuccess = true
	
	return nil
}

func (this *AllocateResourceService) Release(resources []business.IResource) {
	for _, resource := range resources {
		if resource.IsAllocated() {
			if allocator, ok := this.resource2allocator[resource.GetType()]; ok {
				allocator.Release(resource)
				resource.ResetAllocation()
				eel.Logger.Debug("[resource] release resource " + resource.GetType())
			}
		}
	}
}


func init() {
}
