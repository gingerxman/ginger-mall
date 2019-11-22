package consumption

import (
	"context"
	"github.com/gingerxman/ginger-mall/business"
	"github.com/gingerxman/ginger-mall/business/account"
	
	"github.com/gingerxman/eel"
	
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"time"
)

type ConsumptionRecord struct {
	eel.EntityBase
	
	Id int
	Money int
	ConsumeCount int
	
	UserId int
	User *account.User

	CorpId int
	
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewConsumptionRecord(ctx context.Context, user business.IUser) *ConsumptionRecord {
	model := &m_order.UserConsumptionRecord{
		UserId: user.GetId(),
		Money: 0,
		ConsumeCount: 0,
	}
	
	o := eel.GetOrmFromContext(ctx)
	db := o.Model(&m_order.UserConsumptionRecord{}).Create(model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(db.Error)
	}
	
	return NewConsumptionRecordFromModel(ctx, model)
}

func NewConsumptionRecordFromModel(ctx context.Context, model *m_order.UserConsumptionRecord) *ConsumptionRecord {
	instance := &ConsumptionRecord{}
	
	instance.Id = model.Id
	instance.UserId = model.UserId
	instance.CorpId = model.CorpId
	instance.Money = model.Money
	instance.ConsumeCount = model.ConsumeCount
	instance.CreatedAt = model.CreatedAt
	instance.UpdatedAt = model.UpdatedAt
	
	return instance
}

func init() {
}