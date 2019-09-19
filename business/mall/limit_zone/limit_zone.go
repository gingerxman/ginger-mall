package limit_zone

import (
	"context"
	"encoding/json"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
	"github.com/gingerxman/gorm"
	"time"
)

type limitCity struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type LimitProvince struct {
	Id int `json:"id"`
	Name string `json:"name"`
	IsSelectAllCity bool `json:"is_select_all_city"`
	Zone string `json:"zone"`
	Cities []*limitCity `json:"cities"`
}

type LimitZone struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	Provinces []*LimitProvince
	CreatedAt time.Time
}

func (this *LimitZone) UpdateName(name string) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.LimitZone{}).Where("id", this.Id).Update(gorm.Params{
		"name": name,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *LimitZone) UpdateAreas(strAreas string) error {
	provinces := make([]*LimitProvince, 0)
	err := json.Unmarshal([]byte(strAreas), &provinces)
	if err != nil {
		eel.Logger.Error(err)
		return err
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_mall.LimitZone{}).Where("id", this.Id).Update(gorm.Params{
		"areas": strAreas,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

func (this *LimitZone) Delete() error {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Where("id", this.Id).Delete(&m_mall.LimitZone{})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

func CreateLimitZone(ctx context.Context, corp business.ICorp, name string) *LimitZone {
	o := eel.GetOrmFromContext(ctx)
	
	model := &m_mall.LimitZone{
		CorpId: corp.GetId(),
		Name: name,
	}
	db := o.Create(model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	return NewLimitZoneFromModel(ctx, model)
}

//根据model构建对象
func NewLimitZoneFromModel(ctx context.Context, model *m_mall.LimitZone) *LimitZone {
	instance := new(LimitZone)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.CreatedAt = model.CreatedAt
	
	if model.Areas != "" {
		provinces := make([]*LimitProvince, 0)
		err := json.Unmarshal([]byte(model.Areas), &provinces)
		if err != nil {
			eel.Logger.Error(err)
			return nil
		}
		instance.Provinces = provinces
	}
	
	return instance
}

func init() {
}
