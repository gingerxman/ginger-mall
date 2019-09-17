package product

type RProduct struct {
	Id int `json:"id"`
	BaseInfo *RProductBaseInfo `json:"base_info"`
	IsDeleted bool `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
}

type RPoolProduct struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	Type string `json:"type"`
	BaseInfo *RProductBaseInfo `json:"base_info"`
	LogisticsInfo *RProductLogisticsInfo `json:"logistics_info"`
	Medias []*RProductMedia `json:"medias"`
	Skus []*RProductSku `json:"skus"`
	Labels []*RProductLabel `json:"labels"`
	Categories []*RLintProductCategory `json:"categories"`
	IsDeleted bool `json:"is_deleted"`
	Status string `json:"status"`
	PlatformProductStatus string `json:"platform_product_status"`
	IsDistributionByPlatform bool `json:"is_distribution_by_platform"`
	IsDistributionByCorp bool `json:"is_distribution_by_corp"`
	Commissions []*RProductCommission `json:"commissions"`
	CreatedAt string `json:"created_at"`
	OffshelfPlan *RProductOffshelfPlan `json:"offshelf_plan"`
}

type RProductOffshelfPlan struct {
	OffshelfAt string `json:"offshelf_at"`
}

type RProductLogisticsInfo struct {
	PostageType string `json:"postage_type"`
	UnifiedPostageMoney float64 `json:"unified_postage_money"`
	//PostageConfig *postage.RPostageConfig `json:"postage_config"`
	LimitZoneType int `json:"limit_zone_type"`
	LimitZoneTypeCode string `json:"limit_zone_type_code"`
	LimitZoneId int `json:"limit_zone_id"`
	LimitZoneAreas interface{} `json:"limit_zone_areas"`
}

type RProductBaseInfo struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	CreateType string `json:"create_type"`
	PromotionTitle string `json:"promotion_title"`
	BarCode string `json:"bar_code"`
	Detail string `json:"detail"`
	ShelveType string `json:"shelve_type"`
	MinLimit int `json:"min_limit"`
	DisplayIndex int `json:"display_index"`
	PhysicalUnit string `json:"physical_unit"`
	Thumbnail string `json:"thumbnail"`
}

type RProductMedia struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Url string `json:"url"`
}

type RProductUsableImoney struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	ProductId int `json:"product_id"`
	ImoneyCode string `json:"imoney_code"`
	IsEnabled bool `json:"is_enabled"`
	CreatedAt string `json:"created_at"`
}

type RProductLabel struct {
	Id int `json:"id"`
	Name string `json:"name"`
	IsEnabled bool `json:"is_enabled"`
	CreatedAt string `json:"created_at"`
}

type RProductCategory struct {
	Id int `json:"id"`
	Name string `json:"name"`
	NodeType string `json:"type"`
	ProductCount int `json:"product_count"`
	DisplayIndex int `json:"display_index"`
	IsEnabled bool `json:"is_enabled"`
	CreatedAt string `json:"created_at"`
}

type RProductCategoryTreeNode struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FatherId int `json:"father_id"`
	Level int `json:"level"`
	SubCategories []*RProductCategoryTreeNode `json:"sub_categories"`
}

type RLintProductCategory struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type RProductProperty struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	UserId int `json:"user_id"`
	Name string `json:"name"`
	Type string `json:"type"`
	IsDeleted bool `json:"is_deleted"`
	Values []*RProductPropertyValue `json:"values"`
	CreatedAt string `json:"created_at"`
}

type RProductPropertyValue struct {
	Id int `json:"id"`
	PropertyId int `json:"property_id"`
	PropertyType string `json:"property_type"`
	PropertyName string `json:"property_name"`
	Text string `json:"text"`
	Image string `json:"image"`
}

type RProductSku struct {
	Id int `json:"id"`
	Name string `json:"name"`
	DisplayName string `json:"display_name"`
	SkuCode string `json:"sku_code"`
	Price float64 `json:"price"`
	MarketPrice float64 `json:"market_price"`
	PurchasePrice float64 `json:"purchase_price"`
	Weight float64 `json:"weight"`
	StockType string `json:"stock_type"`
	Stocks int `json:"stocks"`
	PropertyValues []*RProductPropertyValue `json:"property_values"`
	Level2Price map[string]float64 `json:"level2price"`
}

type RPlatformProductApplication struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	PoolProductId int `json:"pool_product_id"`
	Product *RProductBaseInfo `json:"product"`
	Status string `json:"status"`
	AuditRecords []*RPlatformProductAuditRecord `json:"audit_records"`
	CreatedAt string `json:"created_at"`
}

type RPlatformProductAuditRecord struct {
	Id int `json:"id"`
	Result string `json:"result"`
	Reason string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

type RProductApplication struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	PoolProductId int `json:"pool_product_id"`
	Product *RProductBaseInfo `json:"product"`
	Status string `json:"status"`
	AuditRecords []*RProductAuditRecord `json:"audit_records"`
	CreatedAt string `json:"created_at"`
}

type RProductAuditRecord struct {
	Id int `json:"id"`
	Result string `json:"result"`
	Reason string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

type RProductCommission struct {
	Rate float64 `json:"rate"`
	Level string `json:"level"`
}

type RDistributionResult struct {
	PlatformProfit float64 `json:"platform_profit"`
	ChannelProfit float64 `json:"channel_profit"`
	SalesmanProfit float64 `json:"salesman_profit"`
	PoolProduct *RPoolProduct `json:"pool_product"`
}

func init() {
}