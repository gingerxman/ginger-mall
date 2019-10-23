package resource

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type ResourceGroup struct {
	eel.EntityBase
	SupplierId int
	Resources []business.IResource
	productResources []business.IResource
	imoneyResources []business.IResource
}

func (this *ResourceGroup) AddResource(resource business.IResource) {
	this.Resources = append(this.Resources, resource)
	
	if resource.GetType() == RESOURCE_TYPE_PRODUCT {
		this.productResources = append(this.productResources, resource)
	} else {
		this.imoneyResources = append(this.imoneyResources, resource)
	}
}

func (this *ResourceGroup) GetProductResources() []business.IResource {
	return this.productResources
}

func (this *ResourceGroup) GetImoneyResources() []business.IResource {
	return this.imoneyResources
}
