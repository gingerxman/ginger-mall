package order
import (
	"context"
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
	b_order_params "github.com/gingerxman/ginger-mall/business/order/params"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"time"
)

type OrderLogistics struct {
	eel.EntityBase

	Id int
	OrderBid string
	EnableLogistics bool
	ExpressCompanyName string
	ExpressNumber string
	Shipper string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (this *OrderLogistics) Update(shipInfo *b_order_params.LogisticsParams){
	o := eel.GetOrmFromContext(this.Ctx)
	updateParams := gorm.Params{
		"EnableLogistics": shipInfo.EnableLogistics,
		"ExpressCompanyName": shipInfo.ExpressCompanyName,
		"ExpressNumber": shipInfo.ExpressNumber,
		"Shipper": shipInfo.Shipper,
		"UpdatedAt": time.Now(),
	}
	db := o.Model(&m_order.OrderLogistics{}).Where("id", this.Id).Update(updateParams)
	err := db.Error
	if err != nil{
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("order_logistics:update_failed", "更新物流信息失败"))
	}
}

func NewOrderLogisticsFromModel(ctx context.Context, dbModel *m_order.OrderLogistics) *OrderLogistics{
	instance := new(OrderLogistics)
	instance.Ctx = ctx
	instance.Id = dbModel.Id
	instance.OrderBid = dbModel.OrderBid
	instance.EnableLogistics = dbModel.EnableLogistics
	instance.ExpressCompanyName = dbModel.ExpressCompanyName
	instance.ExpressNumber = dbModel.ExpressNumber
	instance.Shipper = dbModel.Shipper
	instance.UpdatedAt = dbModel.UpdatedAt
	instance.CreatedAt = dbModel.CreatedAt
	return instance
}