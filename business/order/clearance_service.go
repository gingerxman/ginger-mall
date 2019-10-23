package order

import (
	"context"
	
	"github.com/gingerxman/eel"
)

const _PLATFORM_ACCOUNT_USER_ID = 0

type ClearanceService struct {
	eel.ServiceBase
}

func NewClearanceService(ctx context.Context) *ClearanceService {
	service := new(ClearanceService)
	service.Ctx = ctx
	return service
}

func (this *ClearanceService) ClearInvoice(invoice *Invoice) error {
	return nil
	////填充invoice中的OrderProducts
	//NewFillOrderService(this.Ctx).FillInvoiceProducts(invoice)
	//transfers := make([]*common.PlutusTransfer, 0)
	//
	////transfer: buyer -> platform
	//buyerUserId := invoice.UserId
	////platformId := account.GetPlatformId()
	////platformUser := account.NewUserRepository(this.Ctx).GetUserByCorpUserId(platformId)
	//transfers = append(transfers, &common.PlutusTransfer{
	//	Bid: invoice.Bid,
	//	SourceUserId: buyerUserId,
	//	DestUserId: _PLATFORM_ACCOUNT_USER_ID,
	//	SourceIMoneyCode: "rmb",
	//	DestIMoneyCode: "cash",
	//	Amount: invoice.Money.FinalMoney,
	//	Remark: fmt.Sprintf("invoice_clearance(%s):buyer-platform", invoice.Bid),
	//})
	//
	//for _, orderProduct := range invoice.Products {
	//	//获取pool product
	//	corp := account.NewCorpFromOnlyId(this.Ctx, invoice.CorpId)
	//	productPool := product.GetProductPoolForCorp(this.Ctx, corp)
	//	poolProduct := productPool.GetPoolProductByProductId(orderProduct.ProductId)
	//
	//	//填充pool product的sku信息
	//	product.NewFillPoolProductService(this.Ctx).FillOne(poolProduct, eel.FillOption{
	//		"with_sku": true,
	//	})
	//
	//	//获取商品的distribution config
	//	distributionConfig := product.NewDistributionConfigRepository(this.Ctx).GetDistributionConfigForProduct(poolProduct, orderProduct.Sku)
	//	distributionConfig.ProductCount = orderProduct.PurchaseCount
	//
	//	channelUser := account.NewUserRepository(this.Ctx).GetUserByCorpUserId(invoice.CorpId)
	//
	//	if distributionConfig.HasChannelProfit() {
	//		//transfer: platform -> channel
	//		transfers = append(transfers, &common.PlutusTransfer{
	//			Bid: invoice.Bid,
	//			SourceUserId: _PLATFORM_ACCOUNT_USER_ID,
	//			DestUserId: channelUser.Id,
	//			SourceIMoneyCode: "cash",
	//			DestIMoneyCode: "cash",
	//			Amount: distributionConfig.GetFloatChannelProfit(),
	//			Remark: fmt.Sprintf("product_clearance(%d):platform-channel", poolProduct.Id),
	//		})
	//	}
	//
	//	if distributionConfig.HasSalesmanProfit() {
	//		//transfer: platform -> salesman
	//		transfers = append(transfers, &common.PlutusTransfer{
	//			Bid: invoice.Bid,
	//			SourceUserId: _PLATFORM_ACCOUNT_USER_ID,
	//			DestUserId: channelUser.Id,
	//			SourceIMoneyCode: "cash",
	//			DestIMoneyCode: "cash",
	//			Amount: distributionConfig.GetFloatSalesmanProfit(),
	//			Remark: fmt.Sprintf("product_clearance(%d):platform-salesman", poolProduct.Id),
	//		})
	//	}
	//
	//	// transfer: platform -> supplier
	//	supplierUser := account.NewUserRepository(this.Ctx).GetUserByCorpUserId(invoice.SupplierId)
	//	transfers = append(transfers, &common.PlutusTransfer{
	//		Bid: invoice.Bid,
	//		SourceUserId: _PLATFORM_ACCOUNT_USER_ID,
	//		DestUserId: supplierUser.Id,
	//		SourceIMoneyCode: "cash",
	//		DestIMoneyCode: "cash",
	//		Amount: distributionConfig.GetFloatSupplierProfitForPrice(orderProduct.Price),
	//		Remark: fmt.Sprintf("product_clearance(%d):platform-supplier", poolProduct.Id),
	//	})
	//}
	//
	//common.NewPlutusService(this.Ctx).DoBatchTransfers(transfers)
	//return nil
}

func init() {
}
