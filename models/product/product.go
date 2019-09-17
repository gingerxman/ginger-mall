package product

import (
	"github.com/gingerxman/eel"
	"time"
)

//Product Model
type Product struct {
	Id int
	UserId int
	CorpId int `orm:"index"` //foreign key for corp
	Type string
	Name string
	PromotionTitle string
	BarCode string `orm:"size(52);index"`
	MinLimit int
	DisplayIndex int
	CategoryId int //分类id
	PhysicalUnit string
	Thumbnail string
	IsDeleted bool `orm:"default(false)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now_add;type(datetime)"`
}
func (self *Product) TableName() string {
	return "product_product"
}

//ProductCategory Model
const PRODUCT_CATEGORY_NODE_TYPE_NODE = 1
const PRODUCT_CATEGORY_NODE_TYPE_LEAF = 2
var PRODUCTCATEGORYNODE2STR = map[int]string{
	PRODUCT_CATEGORY_NODE_TYPE_NODE: "node",
	PRODUCT_CATEGORY_NODE_TYPE_LEAF: "leaf",
}
var STR2PRODUCTCATEGORYNODE = map[string]int{
	"node": PRODUCT_CATEGORY_NODE_TYPE_NODE,
	"leaf": PRODUCT_CATEGORY_NODE_TYPE_LEAF,
}
type ProductCategory struct {
	Id int
	CorpId int `gorm:"index"` //foreign key for corp
	Name string `gorm:"size:52"`
	NodeType int //分类节点类型（中间节点，叶节点）
	FatherId int //父节点id（0表示没有父节点)
	ProductCount int
	DisplayIndex int
	IsEnabled bool `gorm:"default:true"`
	CreatedAt time.Time `gorm:"type:datetime"`
}
func (self *ProductCategory) TableName() string {
	return "product_category"
}

//PoolProduct model
const PP_TYPE_SYNC = 1 //从其他商品池同步而来的商品
const PP_TYPE_CREATE = 2 //商户自身创建的商品
var PPTYPE2STR = map[int]string{
	PP_TYPE_SYNC: "sync",
	PP_TYPE_CREATE: "create",
}
var STR2PPTYPE = map[string]int {
	"sync": PP_TYPE_SYNC,
	"create": PP_TYPE_CREATE,
}
type PoolProduct struct {
	Id int
	CorpId int `orm:"index"` //foreign key for corp
	UserId int
	ProductId int
	ProductType string
	SupplierId int
	Status int
	PlatformProductStatus int `orm:"index"`
	Type int
	SourcePoolProductId int `orm:"index"` //sync product的source pool product
	DisplayIndex int
	IsDistributionByPlatform bool `orm:"default(false)"`
	IsDistributionByCorp bool `orm:"default(false)"`
	SyncAt time.Time `orm:"type(datetime)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}
func (self *PoolProduct) TableName() string {
	return "product_pool_product"
}

func init() {
	eel.RegisterModel(new(Product))
	eel.RegisterModel(new(ProductCategory))
	eel.RegisterModel(new(PoolProduct))
}