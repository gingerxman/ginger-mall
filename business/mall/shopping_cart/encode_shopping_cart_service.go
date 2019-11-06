package shopping_cart

import (
	"context"
	
	"github.com/gingerxman/eel"
)

type EncodeShoppingCartService struct {
	eel.ServiceBase
}

func NewEncodeShoppingCartService(ctx context.Context) *EncodeShoppingCartService {
	service := new(EncodeShoppingCartService)
	service.Ctx = ctx
	return service
}

func (this *EncodeShoppingCartService) encodeProducts(shoppingCartProducts []*ShoppingCartProduct) []*RShoppingCartProduct {
	rValidProducts := make([]*RShoppingCartProduct, 0)
	for _, shoppingCartProduct := range shoppingCartProducts {
		price := 0
		stocks := 0
		if shoppingCartProduct.IsValid() {
			productSku := shoppingCartProduct.PoolProduct.GetSku(shoppingCartProduct.ProductSkuName)
			price = productSku.Price
			stocks = productSku.Stocks
		}
		
		rValidProduct := RShoppingCartProduct{
			Id: shoppingCartProduct.PoolProduct.Id,
			ShoppingCartItemId: shoppingCartProduct.Id,
			Name: shoppingCartProduct.PoolProduct.Product.Name,
			SkuName: shoppingCartProduct.ProductSkuName,
			SkuDisplayName: shoppingCartProduct.ProductSkuDisplayName,
			Thumbnail: shoppingCartProduct.PoolProduct.Product.Thumbnail,
			Price: price,
			Stocks: stocks,
			PurchaseCount: shoppingCartProduct.Count,
		}
		rValidProducts = append(rValidProducts, &rValidProduct)
	}
	
	return rValidProducts
}

//Encode 对单个实体对象进行编码
func (this *EncodeShoppingCartService) Encode(shoppingCart *ShoppingCart) *RShoppingCart {
	if shoppingCart == nil {
		return nil
	}

	//编码product_groups
	rProductGroups := make([]*RShoppingCartProductGroup, 0)
	validProducts := shoppingCart.GetValidProducts()
	if len(validProducts) > 0 {
		rProductGroups = append(rProductGroups, &RShoppingCartProductGroup{
			Products: this.encodeProducts(validProducts),
		})
	}
	
	//编码invalid products
	return &RShoppingCart{
		ProductGroups: rProductGroups,
		InvalidProducts: this.encodeProducts(shoppingCart.GetInvalidProducts()),
	}
}

func init() {
}
