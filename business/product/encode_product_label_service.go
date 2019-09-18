package product

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeProductLabelService struct {
	eel.ServiceBase
}

func NewEncodeProductLabelService(ctx context.Context) *EncodeProductLabelService {
	service := new(EncodeProductLabelService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeProductLabelService) Encode(productLabel *ProductLabel) *RProductLabel {
	if productLabel == nil {
		return nil
	}

	return &RProductLabel{
		Id: productLabel.Id,
		Name: productLabel.Name,
		IsEnabled: productLabel.IsEnabled,
		CreatedAt: productLabel.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeProductLabelService) EncodeMany(productLabels []*ProductLabel) []*RProductLabel {
	rDatas := make([]*RProductLabel, 0)
	for _, productLabel := range productLabels {
		rDatas = append(rDatas, this.Encode(productLabel))
	}
	
	return rDatas
}

func init() {
}
