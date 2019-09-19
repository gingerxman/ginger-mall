package product

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type FillPoolProductService struct {
	eel.ServiceBase
	Corp business.ICorp
}

func NewFillPoolProductService(ctx context.Context) *FillPoolProductService {
	service := new(FillPoolProductService)
	service.Ctx = ctx
	return service
}

func NewFillPoolProductServiceForCorp(ctx context.Context, corp business.ICorp) *FillPoolProductService {
	service := new(FillPoolProductService)
	service.Ctx = ctx
	service.Corp = corp
	return service
}


func (this *FillPoolProductService) FillOne(poolProduct *PoolProduct, option eel.FillOption) {
	this.Fill([]*PoolProduct{poolProduct}, option)
}

func (this *FillPoolProductService) Fill(poolProducts []*PoolProduct, option eel.FillOption) {
	if len(poolProducts) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, poolProduct := range poolProducts {
		ids = append(ids, poolProduct.Id)
	}
	
	this.fillProduct(poolProducts, ids)
	products := make([]*Product, 0)
	for _, poolProduct := range poolProducts {
		if poolProduct.Product != nil {
			products = append(products, poolProduct.Product)
		}
	}
	
	NewFillProductService(this.Ctx).Fill(products, option)
}

func (this *FillPoolProductService) fillProduct(poolProducts []*PoolProduct, ids []int) {
	firstPoolProduct := poolProducts[0]
	if firstPoolProduct.Product != nil {
		//已经fill过了，直接返回
		return
	}
	
	productIds := make([]int, 0)
	for _, poolProduct := range poolProducts {
		productIds = append(productIds, poolProduct.ProductId)
	}
	
	//构建<id, product>
	products := NewProductRepository(this.Ctx).GetProductsByIds(productIds)
	id2product := make(map[int]*Product)
	for _, product := range products {
		id2product[product.Id] = product
	}
	
	//填充pool product
	for _, poolProduct := range poolProducts {
		if product, ok := id2product[poolProduct.ProductId]; ok {
			poolProduct.Product = product
		}
	}
}

func init() {
}
