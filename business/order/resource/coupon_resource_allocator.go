package resource

import (
	"context"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
)

type CouponResourceAllocator struct {
	eel.ServiceBase
}

func NewCouponResourceAllocator(ctx context.Context) business.IResourceAllocator {
	service := new(CouponResourceAllocator)
	service.Ctx = ctx
	return service
}

//Allocate 申请优惠券资源，使用优惠券
func (this *CouponResourceAllocator) Allocate(resource business.IResource, newOrder business.IOrder) error {
	couponResource := resource.(*CouponResource)
	
	coupon := couponResource.GetCoupon()
	err := coupon.UseByOrder(newOrder)
	if err != nil {
		eel.Logger.Error(err)
		return err
	}

	return nil
}

//Release 释放冻结的资产
func (this *CouponResourceAllocator) Release(resource business.IResource) {
	imoneyResource := resource.(*IMoneyResource)
	reqResource := "imoney.unfrozen_record"
	for _, recordId := range imoneyResource.FrozenRecordIds{
		resp, err := eel.NewResource(this.Ctx).Put("gplutus", reqResource, eel.Map{
			"id": recordId,
		})
		if err != nil || !resp.IsSuccess(){
			 // TODO: 资源释放失败，发送钉钉消息
			 // dingMsg := fmt.Sprintf("> 释放冻结的虚拟资产失败 \n\n resouce:%s, recordid:%d \n\n ", reqResource, recordId)
			// common.UrgentMessage.Put(dingMsg)
		}
	}
}


func init() {
}
