package order

type ROrder struct {
	Id int `json:"id"`
	Bid string `json:"bid"`
	CorpId int `json:"corp_id"`
	UserId int `json:"user_id"`
	Status string `json:"status"`
	FinalMoney float64 `json:"final_money"`
	Invoices []*RInvoice `json:"delivery_items"`
	IsDeleted bool `json:"is_deleted"`
	Resources []map[string]interface{} `json:"resources"`
	CreatedAt string `json:"created_at"`
	OperationLogs []*ROperationLog `json:"operation_logs"`
	StatusLogs []*RStatusLog `json:"status_logs"`
	Remark string `json:"remark"`
	Message string `json:"message"`
	ExtraData map[string]interface{} `json:"extra_data"`
	
	ProductPrice float64 `json:"product_price"`
	Postage float64 `json:"postage"`
}

type RAreaItem struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type RArea struct {
	Province *RAreaItem `json:"province"`
	City *RAreaItem `json:"city"`
	District *RAreaItem `json:"district"`
}

type RShipInfo struct {
	Name string `json:"name"`
	Phone string `json:"phone"`
	Address string `json:"address"`
	Area *RArea `json:"area"`
}

type RInvoiceLogistics struct {
	EnableLogistics bool `json:"enable_logistics"`
	ExpressCompanyName string `json:"express_company_name"`
	ExpressNumber string `json:"express_number"`
	Shipper string `json:"leader_name"`
}

type RInvoice struct {
	Id int `json:"id"`
	Bid string `json:"bid"`
	Status string `json:"status"`
	PaymentType string `json:"payment_type"`
	PaymentTime string `json:"payment_time"`
	Postage float64 `json:"postage"`
	FinalMoney float64 `json:"final_money"`
	ProductPrice float64 `json:"product_price"`
	IsCleared bool `json:"is_cleared"`
	CreatedAt string `json:"created_at"`
	ShipInfo *RShipInfo `json:"ship_info"`
	Products []*ROrderProduct `json:"products"`
	LogisticsInfo *RInvoiceLogistics `json:"logistics_info"`
	Resources []map[string]interface{} `json:"resources"`
	Remark string `json:"remark"`
}

type ROrderProduct struct {
	Id int `json:"id"`
	SupplierId int `json:"supplier_id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Weight float64 `json:"weight"`
	Thumbnail string `json:"thumbnail""`
	Sku string `json:"sku_name"`
	SkuDisplayName string `json:"sku_display_name"`
	Count int `json:"count"`
}

type ROperationLog struct {
	Id int `json:"id"`
	OrderBid string `json:"order_bid"`
	Type string `json:"type"`
	Remark string `json:"remark"`
	Action string `json:"action"`
	Operator string `json:"operator"`
	CreatedAt string `json:"created_at"`
}

type RStatusLog struct {
	Id int `json:"id"`
	OrderBid string `json:"order_bid"`
	FromStatus string `json:"from_status"`
	ToStatus string `json:"to_status"`
	Remark string `json:"remark"`
	Operator string `json:"operator"`
	CreatedAt string `json:"created_at"`
}

type ROrderLogistics struct {
	Id int `json:"id"`
	OrderBid string `json:"order_bid"`
	EnableLogistics bool `json:"enable_logistics"`
	ExpressCompanyName string `json:"express_company_name"`
	ExpressNumber string `json:"express_number"`
	Shipper string `json:"shipper"`
}

type ROrderOutline struct {
	TotalMoney float64 `json:"total_money"`
	IncrementMoney float64 `json:"increment_money"`
	TotalOrderCount int `json:"total_order_count"`
	IncrementOrderCount int `json:"increment_order_count"`
	TotalUserCount int `json:"total_user_count"`
	IncrementUserCount int `json:"increment_user_count"`
}

func init() {
}
