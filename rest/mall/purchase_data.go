package mall

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business/mall/ship_info"
	"github.com/gingerxman/ginger-mall/business/product"
	"strconv"
	"strings"
)

type PurchaseData struct {
	eel.RestResource
}

func (this *PurchaseData) Resource() string {
	return "mall.purchase_data"
}

func (this *PurchaseData) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"product_infos", "?ship_info_id:int"},
	}
}

func (this *PurchaseData) getShipInfoData(ctx *eel.Context) eel.Map {
	req := ctx.Request
	var shipInfoData eel.Map = nil
	shipInfoId, _ := req.GetInt("ship_info_id")
	if shipInfoId != 0 {
		shipInfo := ship_info.NewShipInfoRepository(ctx.GetBusinessContext()).GetShipInfo(shipInfoId)
		if shipInfo != nil {
			shipInfoData = eel.Map{
				"id":      shipInfo.Id,
				"name":    shipInfo.Name,
				"phone":   shipInfo.Phone,
				"area":    shipInfo.Area,
				"address": shipInfo.Address,
			}
		}
	}
	
	return shipInfoData
}


type _ProductInfo struct {
	Id int
	Count int
	SkuName string
}

//parseProductInfo 解析商品信息, 商品信息格式为"${product_id}.${count}.${model_name},..."
func (this *PurchaseData) parseProductInfo(productInfosStr string) []*_ProductInfo {
	productInfos := make([]*_ProductInfo, 0)
	
	productInfoItems := strings.Split(productInfosStr, ";")
	for _, productInfoItem := range productInfoItems {
		items := strings.Split(productInfoItem, ".")
		productId, err := strconv.Atoi(items[0])
		if err != nil {
			eel.Logger.Error(err)
			panic(err.Error())
		}
		
		productCount, err := strconv.Atoi(items[1])
		if err != nil {
			eel.Logger.Error(err)
			panic(err.Error())
		}
		
		skuName := items[2]
		productInfos = append(productInfos, &_ProductInfo{
			Id: productId,
			Count: productCount,
			SkuName: skuName,
		})
	}
	
	return productInfos
}

func (this *PurchaseData) getProductDatas(ctx *eel.Context) ([]*product.PoolProduct, []eel.Map) {
	bCtx := ctx.GetBusinessContext()
	req := ctx.Request
	
	productDatas := make([]eel.Map, 0)
	
	productInfosStr := req.GetString("product_infos", "")
	poolProducts := make([]*product.PoolProduct, 0)
	if productInfosStr != "" {
		productInfos := this.parseProductInfo(productInfosStr)
		id2info := make(map[int]*_ProductInfo)
		poolProductIds := make([]int, 0)
		for _, productInfo := range productInfos {
			poolProductIds = append(poolProductIds, productInfo.Id)
			id2info[productInfo.Id] = productInfo
		}
		
		//获取pool products
		poolProducts = product.GetGlobalProductPool(bCtx).GetPoolProductsByIds(poolProductIds)
		product.NewFillPoolProductService(bCtx).Fill(poolProducts, eel.FillOption{
			"with_sku": true,
			"with_logistics": true,
		})
		id2product := make(map[int]*product.PoolProduct, 0)
		for _, poolProduct := range poolProducts {
			id2product[poolProduct.Id] = poolProduct
		}
		
		//构建product datas
		encodedProducts := product.NewEncodePoolProductService(bCtx).EncodeMany(poolProducts)
		for _, encodedProduct := range encodedProducts {
			productInfo := id2info[encodedProduct.Id]
			poolProduct := id2product[encodedProduct.Id]
			productSku := poolProduct.GetSku(productInfo.SkuName)
			productDatas = append(productDatas, eel.Map{
				"id": encodedProduct.Id,
				"count": productInfo.Count,
				"name": encodedProduct.BaseInfo.Name,
				"price": productSku.Price,
				"thumbnail": encodedProduct.BaseInfo.Thumbnail,
				"sku": productInfo.SkuName,
				"sku_display_name": productSku.GetDisplayName(),
				//"weight": productSku.Weight,
				"logistics_info": encodedProduct.LogisticsInfo,
				"payable_imoneys": make([]int, 0),
			})
		}
	}
	
	return poolProducts, productDatas
}

func (this *PurchaseData) Get(ctx *eel.Context) {
	shipInfoData := this.getShipInfoData(ctx)
	_, productDatas := this.getProductDatas(ctx)
	//coupons := this.getUsableCoupons(bCtx, poolProducts)
	//imoneyDatas := this.getUsableIMoneyDatas(bCtx)
	
	//coupon.NewFillCouponService(bCtx).Fill(coupons, vanilla.FillOption{
	//	"with_rule": true,
	//})
	//couponDatas := coupon.NewEncodeCouponService(bCtx).EncodeMany(coupons)

	ctx.Response.JSON(eel.Map{
		"ship_info": shipInfoData,
		// "payable_imoneys": imoneyDatas,
		"products": productDatas,
		// "usable_coupons": couponDatas,
	})
}
