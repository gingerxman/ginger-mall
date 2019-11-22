package order

import (
	"github.com/gingerxman/eel"
	"time"
)

const ORDER_STATUS_CANCEL = -1  // 已取消：取消订单(回退销量)
const ORDER_STATUS_WAIT_PAY = 0  // 待支付：已下单，未付款
const ORDER_STATUS_PAYING = 1 //开始支付
const ORDER_STATUS_WAIT_SUPPLIER_CONFIRM = 2  // 已支付，等待供货商确认订单
const ORDER_STATUS_PAYED_NOT_SHIP = 3  // 待发货：供货商确认订单后，发货过程中
const ORDER_STATUS_PAYED_SHIPED = 4  // 已发货：供货商已发货，等待消费者确认接收
const ORDER_STATUS_SUCCESSED = 5  // 已完成：消费者确认订单完成，或，自下单10日后自动置为已完成状态
const ORDER_STATUS_REFUNDING = 6  // 退款中
const ORDER_STATUS_PLATFORM_REFUNDING = 7  // 平台退款中
const ORDER_STATUS_REFUNDED = 8  // 退款完成(回退销量)
const ORDER_STATUS_NONSENSE = 999 //状态无意义
var STATUS2STR = map[int]string{
	ORDER_STATUS_WAIT_PAY: "wait_pay",
	ORDER_STATUS_PAYING: "paying",
	ORDER_STATUS_CANCEL: "canceled",
	ORDER_STATUS_PAYED_NOT_SHIP: "wait_ship",
	ORDER_STATUS_PAYED_SHIPED: "shipped",
	ORDER_STATUS_SUCCESSED: "finished",
	ORDER_STATUS_REFUNDING: "refunding",
	ORDER_STATUS_REFUNDED: "refunded",
	ORDER_STATUS_PLATFORM_REFUNDING: "platform_refunding",
	ORDER_STATUS_WAIT_SUPPLIER_CONFIRM: "wait_confirm",
	ORDER_STATUS_NONSENSE: "nonsense",
}
var STR2STATUS = map[string]int {
	"wait_pay": ORDER_STATUS_WAIT_PAY,
	"paying": ORDER_STATUS_PAYING,
	"canceled": ORDER_STATUS_CANCEL,
	"wait_ship": ORDER_STATUS_PAYED_NOT_SHIP,
	"shipped": ORDER_STATUS_PAYED_SHIPED,
	"finished": ORDER_STATUS_SUCCESSED,
	"refunding": ORDER_STATUS_REFUNDING,
	"refunded" : ORDER_STATUS_REFUNDED,
	"platform_refunding": ORDER_STATUS_PLATFORM_REFUNDING,
	"wait_confirm": ORDER_STATUS_WAIT_SUPPLIER_CONFIRM,
	"nonsense": ORDER_STATUS_NONSENSE,
}

const ORDER_BILL_TYPE_NONE = 0  // 无发票
const ORDER_BILL_TYPE_PERSONAL = 1  // 个人发票
const ORDER_BILL_TYPE_COMPANY = 2  // 公司发票

const ORDER_TYPE_PRODUCT_ORDER = 1
const ORDER_TYPE_PRODUCT_INVOICE = 2
const ORDER_TYPE_CUSTOM = 3
var ORDERTYPE2STR = map[int]string{
	ORDER_TYPE_PRODUCT_ORDER: "order",
	ORDER_TYPE_PRODUCT_INVOICE: "invoice",
	ORDER_TYPE_CUSTOM: "custom",
}
var STR2ORDERTYPE = map[string]int {
	"order": ORDER_TYPE_PRODUCT_ORDER,
	"invoice": ORDER_TYPE_PRODUCT_INVOICE,
	"custom": ORDER_TYPE_CUSTOM,
}

//Order Model
type Order struct {
	eel.Model
	Type int `gorm:"index"`
	CustomType string `gorm:"size:50;index"`
	Bid string `gorm:"unique;size:160"`
	BizCode string `gorm:"index;size:125"` //订单所属业务，用于搜索
	PlatformCorpId int
	CorpId int
	UserId int
	SupplierId int //供货商的corp id，用于拆单
	OriginalOrderId int `gorm:"index"` //出货单所属的订单id
	Status int //订单状态
	Remark string `gorm:"type:text"`//订单备注
	CancelReason string `gorm:size:256` // 取消订单的理由
	CustomerMessage string //客户留言
	PaymentType string
	Postage int
	PostageStrategy int `gorm:"default:0"`//运费策略
	FinalMoney int
	EnableLogistics bool //是否需要物流
	ShipName string //收货人姓名
	ShipPhone string //收货人电话
	ShipAddress string //收货人地址
	ShipAreaCode string //收货人地区编码
	IsDeleted bool `gorm:"default:false"` //是否删除
	IsCleared bool `gorm:"default:false"` //是否完成清算
	PaymentTime time.Time `gorm:"auto_now_add;type:datetime"`
	Resources string `gorm:"type:text"`//resource集合的json字符串
	ExtraData string `gorm:"type:text"`//订单的额外信息
}
func (self *Order) TableName() string {
	return "order_order"
}


type OrderHasProduct struct {
	eel.Model
	OrderId int `gorm:index`
	PoolProductId int `gorm:index`
	ProductId int `gorm:index`
	ProductName string
	ProductSkuName string
	ProductSkuDisplayName string
	Price int
	Weight float64
	Count int
	Thumbnail string //商品图片
	Code string
}
func (self *OrderHasProduct) TableName() string {
	return "order_has_product"
}


const ORDER_OPERATION_TYPE_MEMBER = 1
const ORDER_OPERATION_TYPE_MALL_OPERATOR = 2
const ORDER_OPERATION_TYPE_SYSTEM = 3
var OPERATONTYPE2CODE = map[int]string {
	ORDER_OPERATION_TYPE_MEMBER: "member_action",
	ORDER_OPERATION_TYPE_MALL_OPERATOR: "operator_action",
	ORDER_OPERATION_TYPE_SYSTEM: "system_action",
}

type OrderOperationLog struct {
	eel.Model
	OrderBid string `gorm:index;size:40`
	Type int
	Remark string `gorm:"type:text"`
	Action string
	Operator string
}
func (self *OrderOperationLog) TableName() string {
	return "order_operation_log"
}

type OrderStatusLog struct{
	eel.Model
	OrderBid string `gorm:index;size:160`
	FromStatus int
	ToStatus int
	Remark string `gorm:"type:text"`
	Operator string
}
func (self *OrderStatusLog) TableName() string {
	return "order_status_log"
}

type OrderLogistics struct {
	Id int
	OrderBid string
	EnableLogistics bool //是否需要物流
	ExpressCompanyName string //物流公司名
	ExpressNumber string // 快递单号
	Shipper string // 收货人
	UpdatedAt time.Time `gorm:"auto_now_add;type:datetime"`
	CreatedAt time.Time `gorm:"auto_now_add;type:datetime"`
}
func (self *OrderLogistics) TableName() string {
	return "order_logistics"
}

type UserConsumptionRecord struct {
	eel.Model
	UserId int `orm:index`
	CorpId int `orm:index`
	Money int // 消费金额
	ConsumeCount int //消费次数
}
func (self *UserConsumptionRecord) TableName() string {
	return "order_user_consumption_record"
}

func init() {
	eel.RegisterModel(new(Order))
	eel.RegisterModel(new(OrderHasProduct))
	eel.RegisterModel(new(OrderOperationLog))
	eel.RegisterModel(new(OrderStatusLog))
	eel.RegisterModel(new(OrderLogistics))
	eel.RegisterModel(new(UserConsumptionRecord))
}
