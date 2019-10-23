package resource

import (
	"context"
	"errors"
	"fmt"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	b_coupon "github.com/gingerxman/ginger-mall/business/coupon"
)

type CouponUsage struct {
	Code string
	Money float64
}

type CouponResource struct {
	eel.EntityBase
	Resource
	
	Code string
	poolProductIds []int
	coupon *b_coupon.Coupon
	deductionMoney float64
}

func (this *CouponResource) GetType() string {
	return RESOURCE_TYPE_COUPON
}

func (this *CouponResource) CanSplit() bool {
	return false
}

func (this *CouponResource) GetDeductionMoney(deductableMoney float64) float64 {
	deductionMoney, err := this.GetCoupon().GetDeductionMoney(deductableMoney)
	if err != nil {
		eel.Logger.Error(err)
		return 0.0
	}
	
	this.deductionMoney = deductionMoney
	return this.deductionMoney
}

func (this *CouponResource) GetPrice() float64 {
	return 0.0
}

func (this *CouponResource) GetPostage() float64 {
	return 0.0
}

func (this *CouponResource) GetRawResourceObject() interface{} {
	return this
}

func (this *CouponResource) IsNeedLockWhenConsume() bool {
	return true
}

func (this *CouponResource) GetLockName() string {
	return this.Code
}

func (this *CouponResource) ToMap() map[string]interface{} {
	resourceInfo := make(map[string]interface{})
	
	resourceInfo["type"] = this.GetType()
	resourceInfo["code"] = this.Code
	resourceInfo["deduction_money"] = this.deductionMoney
	
	return resourceInfo
}

func (this *CouponResource) SaveForOrder(order business.IOrder) error {
	
	return nil
}

func (this *CouponResource) GetCoupon() *b_coupon.Coupon {
	if (this.coupon == nil) {
		this.coupon = b_coupon.NewCouponRepository(this.Ctx).GetCouponByCode(this.Code)
	}
	
	return this.coupon
}

func (this *CouponResource) IsValid() error {
	user := account.GetUserFromContext(this.Ctx)
	coupon := this.GetCoupon()
	if coupon == nil {
		eel.Logger.Error("无效的coupon(%s)", this.Code)
		return errors.New(fmt.Sprintf("无效的coupon(%s)", this.Code))
	}
	
	isValid, err := coupon.CheckValidity(user, this.poolProductIds)
	
	if isValid {
		return nil
	} else {
		eel.Logger.Error(err)
		return err
	}
}

func (this *CouponResource) IsAllocated() bool{
	return this.Resource.IsAllocated()
}

func (this *CouponResource) SetAllocated(){
	this.Resource.SetAllocated()
}

func (this *CouponResource) ResetAllocation(){
	this.Resource.ResetAllocation()
}

func NewCouponResource(ctx context.Context) *CouponResource {
	instance := &CouponResource{}
	instance.Ctx = ctx
	
	return instance
}

func init() {
}
