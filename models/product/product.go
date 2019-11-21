package product

import (
	"github.com/gingerxman/eel"
	"time"
)

//Product Model
type Product struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Type string
	Name string
	PromotionTitle string
	Code string `gorm:"size:52;index"`
	DisplayIndex int
	CategoryId int //分类id
	Thumbnail string
	IsDeleted bool `gorm:"default:false"`
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
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Name string `gorm:"size:52"`
	NodeType int //分类节点类型（中间节点，叶节点）
	FatherId int //父节点id（0表示没有父节点)
	ProductCount int
	DisplayIndex int
	IsEnabled bool `gorm:"default:true"`
}
func (self *ProductCategory) TableName() string {
	return "product_category"
}


//ProductDescription Model
type ProductDescription struct {
	Id int `gorm:"primary_key"`
	ProductId int `gorm:"index"`
	Introduction string
	Detail string `gorm:"type:text"`
	Remark string `gorm:"type:text"`
}
func (self *ProductDescription) TableName() string {
	return "product_description"
}


//ProductMedia Model
const PRODUCT_MEDIA_TYPE_IMAGE = 1
const PRODUCT_MEDIA_TYPE_VIDEO = 2
var PRODUCTMEDIA2STR = map[int]string{
	PRODUCT_MEDIA_TYPE_IMAGE: "image",
	PRODUCT_MEDIA_TYPE_VIDEO: "video",
}
var STR2PRODUCTMEDIA = map[string]int {
	"image": PRODUCT_MEDIA_TYPE_IMAGE,
	"video": PRODUCT_MEDIA_TYPE_VIDEO,
}
type ProductMedia struct {
	Id int `gorm:"primary_key"`
	ProductId int `gorm:"index"`
	Type int
	Url string
}
func (self *ProductMedia) TableName() string {
	return "product_media"
}


//ProductUsableImoney Model
type ProductUsableImoney struct {
	Id int `gorm:"primary_key"`
	ProductId int `gorm:"index"`
	ImoneyCode string
	IsEnabled bool `gorm:"default:true"`
}
func (self *ProductUsableImoney) TableName() string {
	return "product_usable_imoney"
}


//ProductLogisticsInfo Model
const LIMITZONE_TYPE_NO_RESTRICT = 0
const LIMITZONE_TYPE_SALE = 1
const LIMITZONE_TYPE_FORBIDDEN = 2
type ProductLogisticsInfo struct {
	Id int `gorm:"primary_key"`
	ProductId int `gorm:"index"`
	PostageType string //运费类型
	UnifiedPostageMoney int //统一运费金额
	LimitZoneType int
	LimitZoneId int
	CreatedAt time.Time `gorm:"auto_now_add;type(datetime)"`
}
func (self *ProductLogisticsInfo) TableName() string {
	return "product_logistics"
}


//PoolProduct model
const PP_STATUS_DELETE = -1 //商品删除 不在当前供应商显示
const PP_STATUS_OFF = 0 //商品下架(待售)
const PP_STATUS_ON = 1 //商品上架
const PP_STATUS_ON_POOL = 2 //商品在商品池中显示
const PP_STATUS_WAIT_APPLICATION = 3 //等待审核
var PPSTATUS2STR = map[int]string{
	PP_STATUS_OFF: "off_shelf",
	PP_STATUS_ON: "on_shelf",
	PP_STATUS_DELETE: "deleted",
	PP_STATUS_ON_POOL: "on_pool",
	PP_STATUS_WAIT_APPLICATION: "wait_application",
}
var STR2PPSTATUS = map[string]int {
	"off_shelf": PP_STATUS_OFF,
	"on_shelf": PP_STATUS_ON,
	"deleted": PP_STATUS_DELETE,
	"on_pool": PP_STATUS_ON_POOL,
	"wait_application": PP_STATUS_WAIT_APPLICATION,
}
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
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	ProductId int
	ProductType string
	SupplierId int
	Status int
	Type int
	SoldCount int // 销量
	SourcePoolProductId int `gorm:"index"` //sync product的source pool product
	DisplayIndex int
	SyncAt time.Time `gorm:"type(datetime)"`
}
func (self *PoolProduct) TableName() string {
	return "product_pool_product"
}

//ProductLabel Model
type ProductLabel struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Name string
	IsEnabled bool `gorm:"default(true)"`
}
func (self *ProductLabel) TableName() string {
	return "product_label"
}


//ProductHasLabel Model
type ProductHasLabel struct {
	eel.Model
	ProductId int //foreign key product
	LabelId   int //foreign key product_label
}
func (self *ProductHasLabel) TableName() string {
	return "product_has_label"
}


//ProductProperty Model
type ProductProperty struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Name string
	IsDeleted bool `gorm:"default:false"`
}
func (self *ProductProperty) TableName() string {
	return "product_sku_property"
}


//ProductPropertyValue Model
type ProductPropertyValue struct {
	eel.Model
	PropertyId int `gorm:"index"` //foreign key for property
	Text string
	Image string
	IsDeleted bool `gorm:"default:false"`
}
func (self *ProductPropertyValue) TableName() string {
	return "product_sku_property_value"
}


//ProductSku Model
type ProductSku struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	ProductId int
	Name string //规格名
	Code string //规格编码
	Price int //商品价格
	CostPrice int //商品成本价
	Stocks int //库存
	IsDeleted bool `gorm:"default:false"`
}
func (self *ProductSku) TableName() string {
	return "product_sku"
}


type ProductSkuHasPropertyValue struct {
	Id int
	SkuId int //foreign key product_sku
	PropertyId int //foreign key product_property
	PropertyValueId int //foreign key product_property_value
	CreatedAt time.Time `gorm:"type:datetime"`
}
func (self *ProductSkuHasPropertyValue) TableName() string {
	return "product_sku_has_property"
}





func init() {
	eel.RegisterModel(new(Product))
	eel.RegisterModel(new(ProductCategory))
	eel.RegisterModel(new(ProductDescription))
	eel.RegisterModel(new(ProductMedia))
	eel.RegisterModel(new(ProductUsableImoney))
	eel.RegisterModel(new(ProductLogisticsInfo))
	eel.RegisterModel(new(PoolProduct))
	eel.RegisterModel(new(ProductLabel))
	eel.RegisterModel(new(ProductHasLabel))
	eel.RegisterModel(new(ProductProperty))
	eel.RegisterModel(new(ProductPropertyValue))
	eel.RegisterModel(new(ProductSku))
	eel.RegisterModel(new(ProductSkuHasPropertyValue))
}
