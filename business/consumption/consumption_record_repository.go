package consumption

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	"strings"
)

type ConsumptionRecordRepository struct {
	eel.ServiceBase
}

func NewConsumptionRecordRepository(ctx context.Context) *ConsumptionRecordRepository {
	service := new(ConsumptionRecordRepository)
	service.Ctx = ctx
	return service
}

func (this *ConsumptionRecordRepository) GetRecords(filters eel.Map, orderExprs ...string) []*ConsumptionRecord {
	records := make([]*ConsumptionRecord, 0)
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.UserConsumptionRecord{})
	
	var models []*m_order.UserConsumptionRecord
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	db = db.Find(&models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return records
	}
	
	for _, model := range models {
		records = append(records, NewConsumptionRecordFromModel(this.Ctx, model))
	}
	return records
}

func (this *ConsumptionRecordRepository) checkFieldNames(fields []string, fieldName string) bool {
	ok := false
	for _, field := range fields {
		if fieldName == field {
			ok = true
		}
	}
	return ok
}

func (this *ConsumptionRecordRepository) parseFilters(filters map[string]interface{}) map[string]interface{}{
	orderFilters := eel.Map{}

	productFilters := eel.Map{}
	productFields := []string{"product_name"}

	type2filters := eel.Map{}

	for key, value := range filters {
		keyString := strings.Split(key, "__")
		fieldName := keyString[0]
		match := ""
		if len(keyString) > 1{
			match = keyString[1]
		}
		if ok := this.checkFieldNames(productFields, fieldName); ok{
			tempMatch := ""
			if match != ""{
				tempMatch = fmt.Sprintf("__%s", match)
			}
			if fieldName == "product_name" {
				productFilters[fmt.Sprintf("name%s", tempMatch)] = value
			} else{
				productFilters[key] = value
			}
			type2filters["productFilters"] = productFilters
		} else {
			if fieldName == "status" {
				if match == "in" {
					status := make([]int, 0)
					for _, st := range value.([]interface{}) {
						status = append(status, m_order.STR2STATUS[st.(string)])
					}
					orderFilters[key] = status
				} else {
					strValue := value.(string)
					if strValue != "all" {
						orderFilters[key] = m_order.STR2STATUS[value.(string)]
					}
				}
			} else {
				orderFilters[key] = value
			}
			type2filters["orderFilters"] = orderFilters
		}
	}
	return type2filters
}

func (this *ConsumptionRecordRepository) GetPagedOrders(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*ConsumptionRecord, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.UserConsumptionRecord{})
	
	var models []*m_order.UserConsumptionRecord
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
	
	records := make([]*ConsumptionRecord, 0)
	for _, model := range models {
		records = append(records, NewConsumptionRecordFromModel(this.Ctx, model))
	}
	return records, paginateResult
}

func (this *ConsumptionRecordRepository) GetRecordsForUsers(userIds []int) []*ConsumptionRecord {
	filters := eel.Map{
		"user_id__in": userIds,
	}
	
	return this.GetRecords(filters)
}

func init() {
}
