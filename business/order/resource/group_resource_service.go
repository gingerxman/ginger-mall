package resource

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type GroupResourceService struct {
	eel.ServiceBase
}

func NewGroupResourceService(ctx context.Context) *GroupResourceService {
	service := new(GroupResourceService)
	service.Ctx = ctx
	return service
}

func (this *GroupResourceService) Group(resources []business.IResource) map[int]*ResourceGroup {
	productResources := make([]business.IResource, 0)
	imoneyResources := make([]business.IResource, 0)
	
	for _, resource := range resources {
		if resource.GetType() == RESOURCE_TYPE_PRODUCT {
			productResources = append(productResources, resource)
		} else {
			imoneyResources = append(imoneyResources, resource)
		}
	}
	
	supplier2group := this.groupProductResources(productResources)
	
	//将imoney resource全部加入到第一个group中
	//TODO 后续改进
	var firstGroup *ResourceGroup
	for _, group := range supplier2group {
		firstGroup = group
		if firstGroup != nil {
			break
		}
	}
	for _, imoneyResource := range imoneyResources {
		firstGroup.AddResource(imoneyResource)
	}
	
	return supplier2group
}

func (this *GroupResourceService) groupProductResources(productResources []business.IResource) map[int]*ResourceGroup {
	supplier2group := make(map[int]*ResourceGroup)
	
	for _, productResource := range productResources {
		rawProductResource := productResource.GetRawResourceObject().(*ProductResource)
		poolProduct := rawProductResource.GetPoolProduct()
		supplierId := poolProduct.SupplierId
		
		var group *ResourceGroup
		var ok bool
		if group, ok = supplier2group[supplierId]; !ok {
			group = &ResourceGroup{
				SupplierId: supplierId,
			}
			supplier2group[supplierId] = group
		}
		
		group.AddResource(productResource)
	}
	
	return supplier2group
}


func init() {
}
