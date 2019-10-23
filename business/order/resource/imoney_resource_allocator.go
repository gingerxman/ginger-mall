package resource

import (
	"context"
	"errors"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
)

type IMoneyResourceAllocator struct {
	eel.ServiceBase
}

func NewIMoneyResourceAllocator(ctx context.Context) business.IResourceAllocator {
	service := new(IMoneyResourceAllocator)
	service.Ctx = ctx
	return service
}

//Allocate 申请商品资源，减少库存
func (this *IMoneyResourceAllocator) Allocate(resource business.IResource, newOrder business.IOrder) error {
	imoneyResource := resource.(*IMoneyResource)
	
	users := account.NewUserRepository(this.Ctx).GetUsers([]int{imoneyResource.SourceUserId})
	if len(users) == 0 {
		return errors.New(fmt.Sprintf("invalid_imoney_source_user_id(%d)", imoneyResource.SourceUserId))
	}
	
	user := users[0]
	// TODO: 检查LoginAs
	resp, err := eel.NewResource(this.Ctx).LoginAs(user.Unionid).Put("gplutus", "imoney.frozen_record", eel.Map{
		"imoney_code": imoneyResource.Code,
		"amount": imoneyResource.Count,
		"type": 1,
		"remark": "使用虚拟资产购买商品",
	})
		
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	if resp.IsSuccess(){
		recordId, err := resp.RespData.Get("frozen_record_id").Int()
		if err == nil{
			imoneyResource.AddFrozenRecord(recordId)
		}
	}

	return nil
}

//Release 释放冻结的资产
func (this *IMoneyResourceAllocator) Release(resource business.IResource) {
	imoneyResource := resource.(*IMoneyResource)
	reqResource := "imoney.unfrozen_record"
	for _, recordId := range imoneyResource.FrozenRecordIds{
		resp, err := eel.NewResource(this.Ctx).Put("gplutus", reqResource, eel.Map{
			"id": recordId,
		})
		if err != nil || !resp.IsSuccess(){
			 // TODO: 资源释放失败，发送钉钉消息
			// dingMsg := fmt.Sprintf("> 释放冻结的虚拟资产失败 \n\n resouce:%s, recordid:%d \n\n ", reqResource, recordId)
			//common.UrgentMessage.Put(dingMsg)
		}
	}
}


func init() {
}
