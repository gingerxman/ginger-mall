package order

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/gingerxman/eel"
)

type EncodeInvoiceService struct {
	eel.ServiceBase
}

func NewEncodeInvoiceService(ctx context.Context) *EncodeInvoiceService {
	service := new(EncodeInvoiceService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeInvoiceService) Encode(invoice *Invoice) *RInvoice {
	if invoice == nil {
		return nil
	}
	
	totalProductPrice := 0.0
	
	rProducts := make([]*ROrderProduct, 0)
	for _, orderProduct := range invoice.Products {
		rProduct := &ROrderProduct{
			Id: orderProduct.ProductId,
			SupplierId: orderProduct.SupplierId,
			Name: orderProduct.Name,
			Price: eel.Decimal(orderProduct.Price),
			Weight: eel.Decimal(orderProduct.Weight),
			Thumbnail: orderProduct.Thumbnail,
			Sku: orderProduct.Sku,
			SkuDisplayName: orderProduct.SkuDisplayName,
			Count: orderProduct.PurchaseCount,
		}
		
		rProducts = append(rProducts, rProduct)
		totalProductPrice += rProduct.Price * float64(orderProduct.PurchaseCount)
	}
	
	//编码RShipInfo
	area := invoice.ShipInfo.GetArea()
	areaName := fmt.Sprintf("%s %s %s", area.Province.Name, area.City.Name, area.District.Name)
	rShipInfo := &RShipInfo{
		Name: invoice.ShipInfo.Name,
		Phone: invoice.ShipInfo.Phone,
		Address: invoice.ShipInfo.Address,
		AreaCode: invoice.ShipInfo.AreaCode,
		AreaName: areaName,
		Area: &RArea{
			Province: &RAreaItem{
				Id: area.Province.Id,
				Name: area.Province.Name,
			},
			City: &RAreaItem{
				Id: area.City.Id,
				Name: area.City.Name,
			},
			District: &RAreaItem{
				Id: area.District.Id,
				Name: area.District.Name,
			},
		},
	}
	
	//编码RInvoiceLogistics
	rLogistics := &RInvoiceLogistics{}
	if invoice.Logistics != nil {
		rLogistics.EnableLogistics = invoice.Logistics.EnableLogistics
		rLogistics.ExpressCompanyName = invoice.Logistics.ExpressCompanyName
		rLogistics.ExpressNumber = invoice.Logistics.ExpressNumber
		rLogistics.Shipper = invoice.Logistics.Shipper
	}
	
	var resources []map[string]interface{}
	err := json.Unmarshal([]byte(invoice.Resources), &resources)
	if err != nil {
		eel.Logger.Error(err)
	}

	return &RInvoice{
		Id: invoice.Id,
		Bid: invoice.Bid,
		Status: invoice.GetStatusText(),
		Postage: eel.Decimal(invoice.Money.Postage),
		FinalMoney: eel.Decimal(invoice.Money.FinalMoney),
		ProductPrice: eel.Decimal(totalProductPrice),
		IsCleared: invoice.IsCleared,
		Products: rProducts,
		LogisticsInfo: rLogistics,
		ShipInfo: rShipInfo,
		Resources: resources,
		CreatedAt: invoice.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeInvoiceService) EncodeMany(products []*Invoice) []*RInvoice {
	rDatas := make([]*RInvoice, 0)
	for _, product := range products {
		rDatas = append(rDatas, this.Encode(product))
	}
	
	return rDatas
}

func init() {
}
