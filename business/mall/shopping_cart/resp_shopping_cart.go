package shopping_cart

type RShoppingCartProduct struct {
	Id int `json:"id"`
	ShoppingCartItemId int `json:"shopping_cart_item_id"`
	Name string `json:"name"`
	SkuName string `json:"sku_name"`
	SkuDisplayName string `json:"sku_display_name"`
	Thumbnail string `json:"thumbnail"`
	Price int `json:"price"`
	Stocks int `json:"stocks"`
	PurchaseCount int `json:"purchase_count"`
}

type RShoppingCartProductGroup struct {
	Products []*RShoppingCartProduct `json:"products"`
}

type RShoppingCart struct {
	ProductGroups []*RShoppingCartProductGroup `json:"product_groups"`
	InvalidProducts []*RShoppingCartProduct `json:"invalid_products"`
}


func init() {
}
