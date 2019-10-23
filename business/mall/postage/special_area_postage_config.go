package postage

import (
	"context"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"time"
	
	"github.com/gingerxman/eel"
)

type SpecialAreaPostageConfig struct {
	eel.EntityBase
	Id int
	CorpId int
	FirstWeight float64
	FirstWeightPrice float64
	AddedWeight float64
	AddedWeightPrice float64
	Destination string
	CreatedAt time.Time

	//foreign key
}

//根据model构建对象
func NewSpecialAreaPostageConfigFromModel(ctx context.Context, model *m_mall.SpecialAreaPostageConfig) *SpecialAreaPostageConfig {
	instance := new(SpecialAreaPostageConfig)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.FirstWeight = model.FirstWeight
	instance.FirstWeightPrice = model.FirstWeightPrice
	instance.AddedWeight = model.AddedWeight
	instance.AddedWeightPrice = model.AddedWeightPrice
	instance.Destination = model.Destination
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
