package limit_zone

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"
)

type LimitZoneRepository struct {
	eel.RepositoryBase
}

func NewLimitZoneRepository(ctx context.Context) *LimitZoneRepository {
	repository := new(LimitZoneRepository)
	repository.Ctx = ctx
	return repository
}

func (this *LimitZoneRepository) GetLimitZones(filters eel.Map, orderExprs ...string) []*LimitZone {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.LimitZone{})
	
	var models []*m_mall.LimitZone
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	db = db.Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	postageConfigs := make([]*LimitZone, 0)
	for _, model := range models {
		postageConfigs = append(postageConfigs, NewLimitZoneFromModel(this.Ctx, model))
	}
	return postageConfigs
}

func (this *LimitZoneRepository) GetPagedLimitZones(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*LimitZone, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.LimitZone{})
	
	var models []*m_mall.LimitZone
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil, paginateResult
	}
	
	postageConfigs := make([]*LimitZone, 0)
	for _, model := range models {
		postageConfigs = append(postageConfigs, NewLimitZoneFromModel(this.Ctx, model))
	}
	return postageConfigs, paginateResult
}

//GetAllLimitZonesForCorp 获得所有LimitZone对象集合
func (this *LimitZoneRepository) GetLimitZonesForCorp(corp business.ICorp, filters eel.Map) []*LimitZone {
	filters["corp_id"] = corp.GetId()
	
	return this.GetLimitZones(filters, "id desc")
	
}

//GetLimitZoneInCorp 根据id和corp获得LimitZone对象
func (this *LimitZoneRepository) GetLimitZoneInCorp(corp business.ICorp, id int) *LimitZone {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	postageConfigs := this.GetLimitZones(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

//GetLimitZone 根据id和corp获得LimitZone对象
func (this *LimitZoneRepository) GetLimitZone(id int) *LimitZone {
	filters := eel.Map{
		"id": id,
	}
	
	postageConfigs := this.GetLimitZones(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

func init() {
}
