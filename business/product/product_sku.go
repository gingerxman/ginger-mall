package product

import (
	"context"
	"github.com/gingerxman/eel"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"strconv"
	"strings"
	"time"
)

const _PRICE_LEVEL_SUPPLIER = "supplier"
const _PRICE_LEVEL_PLATFORM = "platform"
const _PRICE_LEVEL_CHANNEL = "channel"

type ProductSku struct {
	eel.EntityBase
	Id int
	CorpId int
	ProductId int
	Name string
	Code string
	IsStandard bool
	Price float64
	CostPrice float64
	Stocks int
	IsDeleted bool
	CreatedAt time.Time

	//foreign key
	PropertyValues []*ProductPropertyValue
	Property2Value map[string]string
}

func (this *ProductSku) IsStandardSku() bool {
	return this.Name == "standard"
}

func (this *ProductSku) FillProperties(id2property map[int]*ProductProperty, id2value map[int]*ProductPropertyValue) {
	/*
	获取sku关联的property信息
	sku.property_values = [{
		'property_id': 1,
		'property_name': '颜色',
		'id': 1,
		'value': '红'
	}, {
		'property_id': 2,
		'property_name': '尺寸',
		'id': 3,
		'value': 'S'
	}]

	sku.property2value = {
		'颜色': '红',
		'尺寸': 'S'
	}
	*/
	
	if id2property == nil {
		return
	}
	
	if this.IsStandardSku() {
		return
	}
	
	//商品规格名的格式为${property1_id}:${value1_id}_${property2_id}:${value2_id}
	ids := strings.Split(this.Name, "_")
	for _, id := range ids {
		// id的格式为${property_id}:${value_id}
		if strings.Index(id, ":") == -1 {
			//处理异常数据
			continue
		}
		
		items := strings.Split(id, ":")
		propertyId, _ := strconv.Atoi(items[0])
		valueId, _ := strconv.Atoi(items[1])
		
		if property, ok := id2property[propertyId]; ok {
			if value, ok2 := id2value[valueId]; ok2 {
				value.PropertyName = property.Name
				this.Property2Value[property.Name] = value.Text
				this.PropertyValues = append(this.PropertyValues, value)
			}
		}
	}
}

//CheckStock 检查库存是否满足
//TODO 处理并发购买
func (this *ProductSku) CanAffordStock(stock int) bool {
	if this.Stocks >= stock {
		return true
	}
	
	return false
}

func (this *ProductSku) HasStocks() bool {
	return this.CanAffordStock(1)
}

//GetDisplayName 获得sku name(3:1_5:8)对应的display name(黑色 M)
func (this *ProductSku) GetDisplayName() string {
	if this.Name == "standard" {
		return this.Name
	} else {
		names := make([]string, 0)
		propertyRepository := NewProductPropertyRepository(this.Ctx)
		items := strings.Split(this.Name, "_")
		for _, item := range items {
			valueId, _ := strconv.Atoi(strings.Split(item, ":")[1])
			propertyValue := propertyRepository.GetProductPropertyValue(valueId)
			names = append(names, propertyValue.Text)
		}
		return strings.Join(names, " ")
	}
}

//根据model构建对象
func NewProductSkuFromModel(ctx context.Context, model *m_product.ProductSku) *ProductSku {
	instance := new(ProductSku)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.ProductId = model.ProductId
	instance.Name = model.Name
	instance.Code = model.Code
	instance.IsStandard = (model.Name == "standard")
	instance.Price = model.Price
	instance.CostPrice = model.CostPrice
	instance.Stocks = model.Stocks
	instance.IsDeleted = model.IsDeleted
	instance.CreatedAt = model.CreatedAt
	
	instance.PropertyValues = make([]*ProductPropertyValue, 0)
	instance.Property2Value = make(map[string]string)

	return instance
}

func init() {
}
