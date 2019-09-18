package product

import (
	"context"
	"github.com/gingerxman/eel"
	
)

type FillProductLabelService struct {
	eel.ServiceBase
}

func NewFillProductLabelService(ctx context.Context) *FillProductLabelService {
	service := new(FillProductLabelService)
	service.Ctx = ctx
	return service
}

func (this *FillProductLabelService) FillOne(productLabel *ProductLabel, option eel.FillOption) {
	this.Fill([]*ProductLabel{productLabel}, option)
}

func (this *FillProductLabelService) Fill(productLabels []*ProductLabel, option eel.FillOption) {
	if len(productLabels) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, productLabel := range productLabels {
		ids = append(ids, productLabel.Id)
	}
	return
}


func init() {
}
