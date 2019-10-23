package order

import (
	"context"
	"fmt"
	
	"github.com/gingerxman/ginger-mall/business"
	common "github.com/gingerxman/ginger-mall/business/common/echart"
	"time"
	
	"github.com/gingerxman/eel"
	m_order "github.com/gingerxman/ginger-mall/models/order"
)

const STATS_RANGE_WEEK = 7
const STATS_RANGE_MONTH = 30

type _dateRange struct {
	StartDate string
	EndDate string
}

type OrderStatsService struct {
	eel.ServiceBase
}

func NewOrderStatsService(ctx context.Context) *OrderStatsService {
	service := new(OrderStatsService)
	service.Ctx = ctx
	return service
}

func (this *OrderStatsService) getDateRange(rangeType int) *_dateRange {
	today := time.Now()
	endDate := today.Format("2006-01-02 15:04:05")
	delta := eel.Timedelta{Days:30}
	a_month_ago := today.Add(0-delta.Duration())
	startDate := fmt.Sprintf("%s 00:00:00", a_month_ago.Format("2006-01-02"))
	
	return &_dateRange{
		StartDate: startDate,
		EndDate: endDate,
	}
}

func (this *OrderStatsService) getPlatformOrders(ctx context.Context, rangeType int) []*Order {
	orderRepository := NewOrderRepository(ctx)
	
	dateRange := this.getDateRange(rangeType)
	
	filters := eel.Map{
		"is_deleted": false,
		"created_at__gte": dateRange.StartDate,
		"created_at__lte": dateRange.EndDate,
		"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
		"status__in": []int{m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP, m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED},
	}
	
	//收集<date, money>
	orders := orderRepository.GetOrders(filters, "-id")
	
	return orders
}

func (this *OrderStatsService) getCorpOrders(ctx context.Context, corp business.ICorp, rangeType int) []*Order {
	orderRepository := NewOrderRepository(ctx)
	
	dateRange := this.getDateRange(rangeType)
	
	//获取作为渠道的订单
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"is_deleted": false,
		"created_at__gte": dateRange.StartDate,
		"created_at__lte": dateRange.EndDate,
		"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
		"status__in": []int{m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP, m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED},
	}
	orders := orderRepository.GetOrders(filters, "-id")
	
	//获取作为供货商的订单
	filters = eel.Map{
		"supplier_id": corp.GetId(),
		"is_deleted": false,
		"created_at__gte": dateRange.StartDate,
		"created_at__lte": dateRange.EndDate,
		"type": m_order.ORDER_TYPE_PRODUCT_INVOICE,
		"status__in": []int{m_order.ORDER_STATUS_WAIT_SUPPLIER_CONFIRM, m_order.ORDER_STATUS_PAYED_NOT_SHIP, m_order.ORDER_STATUS_PAYED_SHIPED, m_order.ORDER_STATUS_SUCCESSED},
	}
	supplierOrders := orderRepository.GetOrders(filters, "-id")
	
	//合并两类订单
	id2exist := make(map[int]bool)
	for _, order := range orders {
		id2exist[order.Id] = true
	}
	for _, supplierOrder := range supplierOrders {
		if _, ok := id2exist[supplierOrder.Id]; !ok {
			orders = append(orders, supplierOrder)
		}
	}
	
	return orders
}

func (this *OrderStatsService) getDates(rangeType int) []string {
	now := time.Now()
	oneDay := time.Duration(24 * time.Hour)
	
	var startDate time.Time
	if rangeType == STATS_RANGE_WEEK {
		startDate = now.Add(-STATS_RANGE_WEEK * oneDay)
	} else if rangeType == STATS_RANGE_MONTH {
		startDate = now.Add(-STATS_RANGE_MONTH * oneDay)
	}
	endDate := now.Add(oneDay)
	strEndDate := endDate.Format("01-02")
	
	dates := make([]string, 0)
	curDate := startDate
	for {
		strDate := curDate.Format("01-02")
		if strDate == strEndDate {
			dates = append(dates, strDate)
			break
		}
		
		dates = append(dates, strDate)
		curDate = curDate.Add(oneDay)
	}
	
	return dates
}

//GetIncrementMoneyTrend 订单金额增量趋势
func (this *OrderStatsService) GetIncrementMoneyTrend(corp business.ICorp, rangeType int) eel.Map {
	//收集<date, money>
	var orders []*Order
	if corp.IsPlatform() {
		orders = this.getPlatformOrders(this.Ctx, rangeType)
	} else {
		orders = this.getCorpOrders(this.Ctx, corp, rangeType)
	}
	date2money := make(map[string]float64, 0)
	for _, order := range orders {
		date := order.CreatedAt.Format("01-02")
		if money, ok := date2money[date]; ok {
			date2money[date] = money + order.Money.FinalMoney
		} else {
			date2money[date] = order.Money.FinalMoney
		}
	}
	
	//构建Points
	dates := this.getDates(rangeType)
	points := make([]*common.ChartPoint, 0)
	for _, date := range dates {
		if money, ok := date2money[date]; ok {
			points = append(points, &common.ChartPoint{
				X: date,
				Y: eel.Decimal(money),
			})
		} else {
			points = append(points, &common.ChartPoint{
				X: date,
				Y: 0,
			})
		}
	}
	
	//创建line chart
	chartInfo := &common.LineChartInfo{
		Title: "订单金额增量趋势",
		DataName: "金额",
		Points: points,
	}
	lineChart := common.CreateLineChart(chartInfo)
	
	return lineChart
}

//GetIncrementCountTrend 订单数量增量趋势
func (this *OrderStatsService) GetIncrementCountTrend(corp business.ICorp, rangeType int) eel.Map {
	//收集<date, count>
	var orders []*Order
	if corp.IsPlatform() {
		orders = this.getPlatformOrders(this.Ctx, rangeType)
	} else {
		orders = this.getCorpOrders(this.Ctx, corp, rangeType)
	}

	date2count := make(map[string]float64, 0)
	for _, order := range orders {
		date := order.CreatedAt.Format("01-02")
		eel.Logger.Error(date)
		if count, ok := date2count[date]; ok {
			date2count[date] = count + 1
		} else {
			date2count[date] = 1
		}
	}
	
	//构建Points
	dates := this.getDates(rangeType)
	points := make([]*common.ChartPoint, 0)
	for _, date := range dates {
		if count, ok := date2count[date]; ok {
			points = append(points, &common.ChartPoint{
				X: date,
				Y: eel.Decimal(count),
			})
		} else {
			points = append(points, &common.ChartPoint{
				X: date,
				Y: 0,
			})
		}
	}
	
	//创建line chart
	chartInfo := &common.LineChartInfo{
		Title: "订单增量趋势",
		DataName: "订单数",
		Points: points,
	}
	lineChart := common.CreateLineChart(chartInfo)
	
	return lineChart
}

func init() {
}
