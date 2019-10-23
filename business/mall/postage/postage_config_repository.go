package postage

import (
	"context"
	"github.com/gingerxman/ginger-mall/business"
	m_mall "github.com/gingerxman/ginger-mall/models/mall"

	
	"github.com/gingerxman/eel"
)

type PostageConfigRepository struct {
	eel.RepositoryBase
}

func NewPostageConfigRepository(ctx context.Context) *PostageConfigRepository {
	repository := new(PostageConfigRepository)
	repository.Ctx = ctx
	return repository
}

func (this *PostageConfigRepository) GetPostageConfigs(filters eel.Map, orderExprs ...string) []*PostageConfig {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.PostageConfig{})
	
	var models []*m_mall.PostageConfig
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
	
	postageConfigs := make([]*PostageConfig, 0)
	for _, model := range models {
		postageConfigs = append(postageConfigs, NewPostageConfigFromModel(this.Ctx, model))
	}
	return postageConfigs
}

func (this *PostageConfigRepository) GetPagedPostageConfigs(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*PostageConfig, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_mall.PostageConfig{})
	
	var models []*m_mall.PostageConfig
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
	
	postageConfigs := make([]*PostageConfig, 0)
	for _, model := range models {
		postageConfigs = append(postageConfigs, NewPostageConfigFromModel(this.Ctx, model))
	}
	return postageConfigs, paginateResult
}

//GetEnabledPostageConfigsForCorp 获得启用的PostageConfig对象集合
func (this *PostageConfigRepository) GetEnabledPostageConfigsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PostageConfig, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_enabled"] = true
	filters["is_deleted"] = false
	
	return this.GetPagedPostageConfigs(filters, page, "display_index")
	
}

//GetAllPostageConfigsForCorp 获得所有PostageConfig对象集合
func (this *PostageConfigRepository) GetAllPostageConfigsForCorp(corp business.ICorp, filters eel.Map) []*PostageConfig {
	filters["corp_id"] = corp.GetId()
	filters["is_deleted"] = false
	
	return this.GetPostageConfigs(filters, "display_index")
	
}

//GetPostageConfigInCorp 根据id和corp获得PostageConfig对象
func (this *PostageConfigRepository) GetPostageConfigInCorp(corp business.ICorp, id int) *PostageConfig {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	postageConfigs := this.GetPostageConfigs(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

// GetActivePostageConfigInCorp 获得当前被选中的运费配置
func (this *PostageConfigRepository) GetActivePostageConfigInCorp(corp business.ICorp) *PostageConfig {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"is_used": true,
	}
	
	postageConfigs := this.GetPostageConfigs(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

//GetPostageConfig 根据id和corp获得PostageConfig对象
func (this *PostageConfigRepository) GetPostageConfig(id int) *PostageConfig {
	filters := eel.Map{
		"id": id,
	}
	
	postageConfigs := this.GetPostageConfigs(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

//GetPostageConfig 根据id和corp获得PostageConfig对象
func (this *PostageConfigRepository) GetUsedPostageConfigForCorp(corp business.ICorp) *PostageConfig {
	filters := eel.Map{
		"is_used": true,
		"corp_id": corp.GetId(),
	}
	
	postageConfigs := this.GetPostageConfigs(filters)
	
	if len(postageConfigs) == 0 {
		return nil
	} else {
		return postageConfigs[0]
	}
}

func init() {
}
