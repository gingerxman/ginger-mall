package shopping_cart

import (
	"errors"
	"fmt"
	"github.com/gingerxman/ginger-mall/business/account"
	"github.com/gingerxman/ginger-mall/business/product"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"time"
	
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
)

type ShoppingCart struct {
	eel.EntityBase
	Id int
	UserId int
	CorpId int
	CreatedAt time.Time

	//foreign key
	validProducts []*ShoppingCartProduct
	invalidProducts []*ShoppingCartProduct
	isProductsSplitted bool
}

func (this *ShoppingCart) fillProducts(shoppingCartProducts []*ShoppingCartProduct) {
	if len(shoppingCartProducts) == 0 {
		return
	}
	
	poolProductIds := make([]int, 0)
	for _, shoppingCartProduct := range shoppingCartProducts {
		poolProductIds = append(poolProductIds, shoppingCartProduct.PoolProductId)
	}
	
	//fill product
	corp := account.NewCorpFromOnlyId(this.Ctx, this.CorpId)
	poolProducts := product.GetProductPoolForCorp(this.Ctx, corp).GetPoolProductsByIds(poolProductIds)
	if len(poolProducts) == 0 {
		return
	}
	product.NewFillPoolProductService(this.Ctx).Fill(poolProducts, eel.FillOption{
		"with_sku": true,
	})
	
	//构建<ppid, poolProduct>，因为不同的shopping cart product可能对应相同的pool product，所以通过构建<ppid, poolProduct>完成填充
	ppid2pp := make(map[int]*product.PoolProduct, 0)
	for _, poolProduct := range poolProducts {
		ppid2pp[poolProduct.Id] = poolProduct
	}
	
	//填充ShoppingCartProduct.Product
	for _, shoppingCartProduct := range shoppingCartProducts {
		if poolProduct, ok := ppid2pp[shoppingCartProduct.PoolProductId]; ok {
			shoppingCartProduct.PoolProduct = poolProduct
		}
	}
}

func (this *ShoppingCart) splitProducts() error {
	this.isProductsSplitted = true
	o := eel.GetOrmFromContext(this.Ctx)
	
	models := make([]*m_mall.ShoppingCartItem, 0)
	db := o.Model(&m_mall.ShoppingCartItem{}).Where(eel.Map{
		"user_id": this.UserId,
		"corp_id": this.CorpId,
	}).Find(&models)
	err := db.Error
	
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	shoppingCartProducts := make([]*ShoppingCartProduct, 0)
	for _, model := range models {
		shoppingCartProducts = append(shoppingCartProducts, NewShoppingCartProductFromModel(this.Ctx, model))
	}
	
	this.fillProducts(shoppingCartProducts)
	
	//根据shoppingCartProduct.IsValid切分products
	for _, shoppingCartProduct := range shoppingCartProducts {
		if shoppingCartProduct.IsValid() {
			this.validProducts = append(this.validProducts, shoppingCartProduct)
		} else {
			this.invalidProducts = append(this.invalidProducts, shoppingCartProduct)
		}
	}
	
	return nil
}

// GetValidProducts 获得购物车中有效商品集合
func (this *ShoppingCart) GetValidProducts() []*ShoppingCartProduct {
	if !this.isProductsSplitted {
		err := this.splitProducts()
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	return this.validProducts
}

// GetInvalidProducts 获得购物车中无效商品集合
func (this *ShoppingCart) GetInvalidProducts() []*ShoppingCartProduct {
	if !this.isProductsSplitted {
		err := this.splitProducts()
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	return this.invalidProducts
}

// GetProductCount 获得购物车中商品的数量
func (this *ShoppingCart) GetProductCount() int {
	return len(this.GetValidProducts())
}

// AddProduct 向购物车中添加商品
func (this *ShoppingCart) AddProduct(poolProduct *product.PoolProduct, skuName string, count int) error {
	isValidProductSku := product.NewProductRepository(this.Ctx).IsValidProductSku(poolProduct.ProductId, skuName)
	if !isValidProductSku {
		return errors.New(fmt.Sprintf("invalid product sku(%d-%s)", poolProduct.ProductId, skuName))
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	isShoppingCartItemExists := o.Model(&m_mall.ShoppingCartItem{}).Where(eel.Map{
		"product_id": poolProduct.Id,
		"product_sku_name": skuName,
		"user_id": this.UserId,
	}).Exist()
	
	if isShoppingCartItemExists {
		//shopping cart item已存在，更新count
		db := o.Model(&m_mall.ShoppingCartItem{}).Where(eel.Map{
			"user_id": this.UserId,
			"corp_id": this.CorpId,
			"product_id": poolProduct.Id,
			"product_sku_name": skuName,
		}).Update("count", gorm.Expr("count + ?", count))
		
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			return db.Error
		}
	} else {
		//shopping cart item不存在，创建之
		product.NewFillPoolProductService(this.Ctx).FillOne(poolProduct, eel.FillOption{
			"with_sku": true,
		})
		skuDisplayName := poolProduct.GetSku(skuName).GetDisplayName()
		model := m_mall.ShoppingCartItem{
			UserId: this.UserId,
			CorpId: this.CorpId,
			ProductId: poolProduct.Id,
			ProductSkuName: skuName,
			ProductSkuDisplayName: skuDisplayName,
			Count: count,
		}
		
		db := o.Create(&model)
		err := db.Error
		if err != nil {
			eel.Logger.Error(err)
			return err
		}
	}
	
	return nil
}

// DeleteItems 从购物车中删除商品项
func (this *ShoppingCart) DeleteItems(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Where(eel.Map{
		"user_id": this.UserId,
		"corp_id": this.CorpId,
		"id__in": ids,
	}).Delete(&m_mall.ShoppingCartItem{})
	err := db.Error
	
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	return nil
}

// DeleteItem 从购物车中删除商品项
func (this *ShoppingCart) DeleteItem(id int) error {
	return this.DeleteItems([]int{id})
}

func init() {
}
