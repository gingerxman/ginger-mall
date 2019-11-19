package resource

import (
	"context"
	"fmt"
	"errors"
	"github.com/gingerxman/eel"
	
)

type IMoneyResource struct {
	eel.EntityBase
	Resource
	
	Code string
	Count int
	SourceUserId int
	DestUserId int
	FrozenRecordIds []int
}

func (this *IMoneyResource) GetType() string {
	return RESOURCE_TYPE_IMONEY
}

func (this *IMoneyResource) CanSplit() bool {
	return false
}

func (this *IMoneyResource) GetDeductionMoney(deductableMoney int) int {
	return this.Count
}

func (this *IMoneyResource) GetPrice() int {
	return 0
}

func (this *IMoneyResource) GetPostage() int {
	return 0
}

func (this *IMoneyResource) GetRawResourceObject() interface{} {
	return this
}

func (this *IMoneyResource) IsNeedLockWhenConsume() bool {
	return true
}

func (this *IMoneyResource) GetLockName() string {
	return this.Code
}

func (this *IMoneyResource) AddFrozenRecord(frozenRecordId int) {
	this.FrozenRecordIds = append(this.FrozenRecordIds, frozenRecordId)
}

func (this *IMoneyResource) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type": this.GetType(),
		"code": this.Code,
		"count": this.Count,
		"price": this.GetPrice(),
		"deduction_money": this.GetDeductionMoney(0),
		"frozen_record_ids": this.FrozenRecordIds,
	}
}

func (this *IMoneyResource) IsValid() error {
	if this.Count <= 0 {
		eel.Logger.Error(fmt.Sprintf("使用数量错误(%d)", this.Count))
		return errors.New("invalid_purchase_count")
	}
	
	return nil
}

func NewIMoneyResource(ctx context.Context) *IMoneyResource {
	instance := &IMoneyResource{}
	instance.Ctx = ctx
	
	return instance
}

func init() {
}
