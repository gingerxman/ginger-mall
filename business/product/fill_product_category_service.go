package product

import (
	"context"
	"github.com/gingerxman/eel"
)

type FillProductCategoryService struct {
	eel.ServiceBase
}

func NewFillProductCategoryService(ctx context.Context) *FillProductCategoryService {
	service := new(FillProductCategoryService)
	service.Ctx = ctx
	return service
}

func (this *FillProductCategoryService) Fill(productCategories []*ProductCategory, option eel.FillOption) {
	if len(productCategories) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, productCategory := range productCategories {
		ids = append(ids, productCategory.Id)
	}
	return
}


func init() {
}
