package postage

import (
	"context"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"time"
	
	"github.com/gingerxman/eel"
)

type FreePostageConfig struct {
	eel.EntityBase
	Id int
	CorpId int
	Destination string
	Condition string
	ConditionValue string
	CreatedAt time.Time

	//foreign key
}


//根据model构建对象
func NewFreePostageConfigFromModel(ctx context.Context, model *m_mall.FreePostageConfig) *FreePostageConfig {
	instance := new(FreePostageConfig)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Destination = model.Destination
	instance.Condition = model.Condition
	instance.ConditionValue = model.ConditionValue
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
