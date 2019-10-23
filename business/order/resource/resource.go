package resource

import "github.com/gingerxman/ginger-mall/business"

type Resource struct {
	resourceType string
	isAllocated bool
}

func (this *Resource) GetType() string {
	return this.resourceType
}

func (this *Resource) IsAllocated() bool {
	return this.isAllocated
}

func (this *Resource) CanSplit() bool {
	panic("not implement")
}

func (this *Resource) ResetAllocation() {
	this.isAllocated = false
}

func (this *Resource) SetAllocated() {
	this.isAllocated = true
}

func (this *Resource) GetDeductionMoney() float64 {
	panic("not implement")
}

func (this *Resource) GetPrice() float64 {
	panic("not implement")
}

func (this *Resource) GetPostage() float64 {
	panic("not implement")
}

func (this *Resource) IsNeedLockWhenConsume() bool {
	return false
}

func (this *Resource) GetLockName() string {
	panic("not implement")
}

func (this *Resource) IsValid() error {
	panic("not implement")
}

func (this *Resource) ToMap() map[string]interface{} {
	panic("not implement")
}

func (this *Resource) SaveForOrder(order business.IOrder) error {
	return nil
}

func init() {
}
