package order

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-mall/business"
	m_order "github.com/gingerxman/ginger-mall/models/order"
	m_product "github.com/gingerxman/ginger-mall/models/product"
	"strings"
	"time"
)

type OrderRepository struct {
	eel.ServiceBase
}

func NewOrderRepository(ctx context.Context) *OrderRepository {
	service := new(OrderRepository)
	service.Ctx = ctx
	return service
}

func (this *OrderRepository) GetOrders(filters eel.Map, orderExprs ...string) []*Order {
	orders := make([]*Order, 0)
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{})
	
	var models []*m_order.Order
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
		return orders
	}
	
	for _, model := range models {
		orders = append(orders, NewOrderFromModel(this.Ctx, model))
	}
	return orders
}

func (this *OrderRepository) checkFieldNames(fields []string, fieldName string) bool {
	ok := false
	for _, field := range fields {
		if fieldName == field {
			ok = true
		}
	}
	return ok
}

func (this *OrderRepository) parseFilters(filters map[string]interface{}) map[string]interface{}{
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

func (this *OrderRepository) GetPagedOrders(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Order, eel.INextPageInfo) {
	orders := make([]*Order, 0)
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{})

	type2filters := this.parseFilters(filters)
	if productFilters, ok := type2filters["productFilters"]; ok && productFilters != nil {
		productIds := []int{0}
		var productModels []*m_product.Product
		db := o.Model(&m_product.Product{}).Where(productFilters).Limit(-1).Find(&productModels)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
		}
		for _, productModel := range productModels {
			productIds = append(productIds, productModel.Id)
		}
		orderIds := []int{0}
		var orderProductModels []*m_order.OrderHasProduct
		db = o.Model(&m_order.OrderHasProduct{}).Where(eel.Map{
			"product_id__in": productIds,
		}).Limit(-1).Find(&orderProductModels)
		err = db.Error
		if err != nil{
			eel.Logger.Error(err)
		}
		for _, orderProductModel := range orderProductModels {
			orderIds = append(orderIds, orderProductModel.OrderId)
		}
		db = db.Where("id__in", orderIds)
	}

	if orderFilters, ok := type2filters["orderFilters"]; ok && orderFilters != nil {
		db = db.Where(orderFilters)
	}

	var models []*m_order.Order
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	
	spew.Dump(page)
	paginateResult, db := eel.Paginate(db, page, &models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return orders, paginateResult
	}
	
	for _, model := range models {
		orders = append(orders, NewOrderFromModel(this.Ctx, model))
	}
	return orders, paginateResult
}

func (this *OrderRepository) GetPagedOrdersForCorp(corp business.ICorp, filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Order, eel.INextPageInfo) {
	orders := make([]*Order, 0)
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_order.Order{})
	
	db = db.Where("corp_id = ? or supplier_id = ?", corp.GetId(), corp.GetId())
	
	type2filters := this.parseFilters(filters)
	if productFilters, ok := type2filters["productFilters"]; ok && productFilters != nil {
		productIds := []int{0}
		var productModels []*m_product.Product
		db := o.Model(&m_product.Product{}).Where(productFilters).Limit(-1).Find(&productModels)
		err := db.Error
		if err != nil{
			eel.Logger.Error(err)
		}
		for _, productModel := range productModels {
			productIds = append(productIds, productModel.Id)
		}
		orderIds := []int{0}
		var orderProductModels []*m_order.OrderHasProduct
		db = o.Model(&m_order.OrderHasProduct{}).Where(eel.Map{
			"product_id__in": productIds,
		}).Limit(-1).Find(&orderProductModels)
		err = db.Error
		if err != nil{
			eel.Logger.Error(err)
		}
		for _, orderProductModel := range orderProductModels {
			orderIds = append(orderIds, orderProductModel.OrderId)
		}
		db = db.Where("id in (?)", orderIds)
	}
	
	if orderFilters, ok := type2filters["orderFilters"]; ok && orderFilters != nil {
		conditions := orderFilters.(map[string]interface{})
		for key, value := range conditions {
			db = db.Where(key, value)
		}
	}
	
	var models []*m_order.Order
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	
	paginateResult, db := eel.Paginate(db, page, &models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return orders, paginateResult
	}
	
	for _, model := range models {
		orders = append(orders, NewOrderFromModel(this.Ctx, model))
	}
	return orders, paginateResult
}

func (this *OrderRepository) GetPagedOrdersForUserInCorp(user business.IUser, corp business.ICorp, filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Order, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	filters["corp_id"] = corp.GetId()
	filters["type"] = m_order.ORDER_TYPE_PRODUCT_ORDER
	
	return this.GetPagedOrders(filters, page, orderExprs...)
}

func (this *OrderRepository) GetOrderByBid(bid string) *Order {
	filters := eel.Map{
		"bid": bid,
	}
	
	orders := this.GetOrders(filters)
	
	if len(orders) == 0 {
		return nil
	} else {
		return orders[0]
	}
}

func (this *OrderRepository) GetOrdersByBids(bids []string) []*Order{
	filters := eel.Map{
		"bid__in": bids,
	}
	return this.GetOrders(filters)
}

func (this *OrderRepository) GetOrderByBidForUser(user business.IUser, bid string) *Order {
	filters := eel.Map{
		"bid": bid,
		"user_id": user.GetId(),
	}
	
	orders := this.GetOrders(filters)
	
	if len(orders) == 0 {
		return nil
	} else {
		return orders[0]
	}
}

//GetOrderById 根据id获得Order对象
func (this *OrderRepository) GetOrderById(id int) *Order {
	filters := eel.Map{
		"id": id,
	}
	
	orders := this.GetOrders(filters)
	
	if len(orders) == 0 {
		return nil
	} else {
		return orders[0]
	}
}

//GetInvoicesByOrderIds 根据订单id集合，获取出货单对象集合
func (this *OrderRepository) GetInvoicesByOrderIds(orderIds []int) []*Invoice {
	if len(orderIds) == 0 {
		return make([]*Invoice, 0)
	}
	
	filters := eel.Map{
		"original_order_id__in": orderIds,
		"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
	}
	
	orders := this.GetOrders(filters, "id")
	
	invoices := make([]*Invoice, 0)
	for _, order := range orders {
		invoices = append(invoices, NewInvoiceFromOrder(this.Ctx, order))
	}
	
	return invoices
}

func (this *OrderRepository) GetInvoicesForOrder(orderId int) []*Invoice {
	return this.GetInvoicesByOrderIds([]int{orderId})
}

//GetInvoicesByIds 根据出货单id，获得出货单对象集合
func (this *OrderRepository) GetInvoicesByIds(ids []int) []*Invoice {
	if len(ids) == 0 {
		return make([]*Invoice, 0)
	}
	
	filters := eel.Map{
		"id__in": ids,
		//"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
	}
	
	orders := this.GetOrders(filters)
	
	invoices := make([]*Invoice, 0)
	for _, order := range orders {
		invoice := NewInvoiceFromOrder(this.Ctx, order)
		invoice.OriginalOrderId = invoice.Id
		invoices = append(invoices, invoice)
	}
	
	return invoices
}

// GetOrderCountForCorp 获得corp中status指定状态的订单的数量
func (this *OrderRepository) GetOrderCountForCorp(corp business.ICorp, status int) int64 {
	orderType := m_order.ORDER_TYPE_PRODUCT_ORDER
	if status != m_order.ORDER_STATUS_WAIT_PAY {
		orderType = m_order.ORDER_TYPE_PRODUCT_INVOICE
	}
	
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"status": status,
		"type": orderType,
		"is_deleted": false,
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	count, err := o.Model(&m_order.Order{}).Where(filters).Count()
	if err != nil {
		eel.Logger.Error(err)
		return 0
	}
	
	return count
}

// GetOrderCountForUser 获得属于user的status指定状态的订单的数量
func (this *OrderRepository) GetOrderCountForUser(user business.IUser, status int) int64 {
	orderType := m_order.ORDER_TYPE_PRODUCT_ORDER
	if status != m_order.ORDER_STATUS_WAIT_PAY {
		orderType = m_order.ORDER_TYPE_PRODUCT_INVOICE
	}
	
	filters := eel.Map{
		"user_id": user.GetId(),
		"status": status,
		"type": orderType,
		"is_deleted": false,
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	count, err := o.Model(&m_order.Order{}).Where(filters).Count()
	if err != nil {
		eel.Logger.Error(err)
		return 0
	}
	
	return count
}

// 根据出货单的bid获取出货单
func (this *OrderRepository) GetInvoiceByBidForCorp(corp business.ICorp, bid string) *Invoice {
	filters := eel.Map{
		"bid": bid,
	}
	if !corp.IsPlatform() {
		filters["supplier_id"] = corp.GetId()
	}
	orders := this.GetOrders(filters)

	if len(orders) == 0 {
		return nil
	} else {
		return NewInvoiceFromOrder(this.Ctx, orders[0])
	}
}

// 根据出货单的bids获取出货单集合
func (this *OrderRepository) GetInvoicesByBidsForCorp(corp business.ICorp, bids []string) []*Invoice {
	invoices := make([]*Invoice, 0)
	orders := this.GetOrders(eel.Map{
		"supplier_id": corp.GetId(),
		"bid__in": bids,
	})
	if orders != nil{
		for _, order := range orders {
			invoices = append(invoices, NewInvoiceFromOrder(this.Ctx, order))
		}
	}
	return invoices
}

// 根据出货单的bid获取出货单
func (this *OrderRepository) GetInvoiceByBidForUser(user business.IUser, bid string) *Invoice {
	filters := eel.Map{
		"bid": bid,
		"user_id": user.GetId(),
	}
	orders := this.GetOrders(filters)

	if len(orders) == 0 {
		return nil
	} else {
		return NewInvoiceFromOrder(this.Ctx, orders[0])
	}
}

// 根据出货单的bid获取订单
func (this *OrderRepository) GetOrderByBidForCorp(corp business.ICorp, bid string) *Order {
	filters := eel.Map{
		"bid": bid,
	}
	if !corp.IsPlatform(){
		filters["corp_id"] = corp.GetId()
	}
	orders := this.GetOrders(filters)

	if len(orders) == 0 {
		return nil
	} else {
		return orders[0]
	}
}


type OrderOutline struct {
	totalMoney float64
	incrementMoney float64
	totalOrderCount int
	incrementOrderCount int
	totalUserCount int
	incrementUserCount int
}

func (this *OrderRepository) GetOrderOutlineForCorp(corp business.ICorp) *OrderOutline {
	o := eel.GetOrmFromContext(this.Ctx)
	
	now := time.Now()
	yesterday := now.Add(time.Duration(-24 * time.Hour))
	strToday := fmt.Sprintf("%s 00:00:00", now.Format("2006-01-02"))
	strYesterday := fmt.Sprintf("%s 00:00:00", yesterday.Format("2006-01-02"))
	
	//total money
	totalMoney := 0.0
	{
		db := o.Raw("select sum(final_money) as total_money from order_order where type = ? and status in (?, ?, ?, ?)", m_order.ORDER_TYPE_PRODUCT_ORDER, m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP, m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED)
		err := db.Error
		if err != nil {
			eel.Logger.Error(err)
		}
		
		err = db.Row().Scan(&totalMoney)
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	
	//increment money
	incrementMoney := 0.0
	{
		db := o.Raw("select sum(final_money) as increment_money from order_order where type = ? and status in (?, ?, ?, ?) and created_at between ? and ?", m_order.ORDER_TYPE_PRODUCT_ORDER, m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP, m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED, strYesterday, strToday)
		err := db.Error
		if err != nil {
			eel.Logger.Error(err)
		}
		
		err = db.Row().Scan(&incrementMoney)
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	
	//total order count
	totalOrderCount := 0
	{
		db := o.Raw("select count(*) as count from order_order where type = ? and status = ?", m_order.ORDER_TYPE_PRODUCT_ORDER, m_order.ORDER_STATUS_SUCCESSED)
		err := db.Error
		if err != nil {
			eel.Logger.Error(err)
		}
		
		err = db.Row().Scan(&totalOrderCount)
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	
	//increment order count
	incrementOrderCount := 0
	{
		db := o.Raw("select count(*) as count from order_order where type = ? and status = ? and created_at between ? and ?", m_order.ORDER_TYPE_PRODUCT_ORDER, m_order.ORDER_STATUS_SUCCESSED, strYesterday, strToday)
		err := db.Error
		if err != nil {
			eel.Logger.Error(err)
		}
		
		err = db.Row().Scan(&incrementMoney)
		if err != nil {
			eel.Logger.Error(err)
		}
	}
	
	return &OrderOutline{
		totalMoney:          totalMoney,
		incrementMoney:      incrementMoney,
		totalOrderCount:     totalOrderCount,
		incrementOrderCount: incrementOrderCount,
		totalUserCount:      0,
		incrementUserCount:  0,
	}
}


func init() {
}
