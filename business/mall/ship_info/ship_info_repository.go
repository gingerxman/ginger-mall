package ship_info

import (
	"context"
	"github.com/gingerxman/ginger-mall/business"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"

	
	"github.com/gingerxman/eel"
)

type ShipInfoRepository struct {
	eel.RepositoryBase
}

func NewShipInfoRepository(ctx context.Context) *ShipInfoRepository {
	repository := new(ShipInfoRepository)
	repository.Ctx = ctx
	return repository
}

func (this *ShipInfoRepository) GetShipInfos(filters eel.Map, orderExprs ...string) []*ShipInfo {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.ShipInfo{})
	
	var models []*m_mall.ShipInfo
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
	
	shipInfos := make([]*ShipInfo, 0)
	for _, model := range models {
		shipInfos = append(shipInfos, NewShipInfoFromModel(this.Ctx, model))
	}
	return shipInfos
}

func (this *ShipInfoRepository) GetPagedShipInfos(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*ShipInfo, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.ShipInfo{})
	
	var models []*m_mall.ShipInfo
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, err := eel.Paginate(db, page, &models)
	
	if err != nil {
		eel.Logger.Error(err)
		return nil, paginateResult
	}
	
	shipInfos := make([]*ShipInfo, 0)
	for _, model := range models {
		shipInfos = append(shipInfos, NewShipInfoFromModel(this.Ctx, model))
	}
	return shipInfos, paginateResult
}

//GetEnabledShipInfosForUser 获得启用的ShipInfo对象集合
func (this *ShipInfoRepository) GetEnabledShipInfosForUser(user business.IUser, page *eel.PageInfo, filters eel.Map) ([]*ShipInfo, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	filters["is_deleted"] = false
	
	return this.GetPagedShipInfos(filters, page, "id")
	
}

//GetAllShipInfosForUser 获得所有ShipInfo对象集合
func (this *ShipInfoRepository) GetAllShipInfosForUser(user business.IUser, page *eel.PageInfo, filters eel.Map) ([]*ShipInfo, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	filters["is_deleted"] = false
	
	return this.GetPagedShipInfos(filters, page, "id")
	
}

//GetShipInfoInCorp 根据id和user获得ShipInfo对象
func (this *ShipInfoRepository) GetShipInfoForUser(user business.IUser, id int) *ShipInfo {
	filters := eel.Map{
		"user_id": user.GetId(),
		"id": id,
	}
	
	shipInfos := this.GetShipInfos(filters)
	
	if len(shipInfos) == 0 {
		return nil
	} else {
		return shipInfos[0]
	}
}

//GetDefaultShipInfoInCorp 获得user的默认收货地址
func (this *ShipInfoRepository) GetDefaultShipInfoForUser(user business.IUser) *ShipInfo {
	filters := eel.Map{
		"user_id": user.GetId(),
		"is_default": true,
	}
	
	shipInfos := this.GetShipInfos(filters)
	
	if len(shipInfos) == 0 {
		return nil
	} else {
		return shipInfos[0]
	}
}

//GetShipInfo 根据id和corp获得ShipInfo对象
func (this *ShipInfoRepository) GetShipInfo(id int) *ShipInfo {
	filters := eel.Map{
		"id": id,
	}
	
	shipInfos := this.GetShipInfos(filters)
	
	if len(shipInfos) == 0 {
		return nil
	} else {
		return shipInfos[0]
	}
}

func init() {
}
