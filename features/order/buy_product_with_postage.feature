Feature: 购买有运费的商品
	Background:
		Given ginger登录系统
		When ginger创建公司
		"""
		[{
			"name": "MIX",
			"username": "jobs"
		}, {
			"name": "BabyFace",
			"username": "bill"
		}, {
			"name": "Mocha",
			"username": "tom"
		}]
		"""
		Given jobs登录系统
		When jobs创建商品分类
		"""
		[{
			"分类1": []
		},{
			"分类2": []
		},{
			"分类3": []
		}]
		"""
		When jobs创建商品属性
		"""
		[{
			"name": "颜色",
			"values": [{
				"name": "black"
			}, {
				"name": "red"
			}]
		}, {
			"name": "尺寸",
			"values": [{
				"name": "M"
			}, {
				"name": "S"
			}]
		}]
		"""
		When jobs添加商品
		"""
		[{
			"name": "价格100重量1运费模板商品",
			"price": 100.0,
			"weight": 1,
			"postage": "系统"
		}, {
			"name": "价格20重量0.6运费模板商品",
			"price": 20.00,
			"weight": 0.6,
			"postage": "系统"
		}, {
			"name": "价格100重量1运费模板商品3",
			"price": 100.0,
			"weight": 1,
			"postage": "系统"
		}, {
			"name": "价格10重量1运费0元商品",
			"price": 10.0,
			"weight": 1,
			"postage": 0.0
		}, {
			"name": "价格10重量1运费15元商品",
			"price": 10.0,
			"weight": 1,
			"postage": 15.0
		}, {
			"name": "价格10重量1运费10元商品",
			"price": 10.0,
			"weight": 1,
			"postage": 10.0
		}, {
			"name": "价格50重量1运费模板的多规格商品",
			"postage": "系统",
			"skus":{
				"red M": {
					"price": 50.00,
					"weight": 1,
					"stocks": 99999
				},
				"black S": {
					"price": 50.00,
					"weight": 1,
					"stocks": 99999
				}
			}
		}, {
			"name": "价格50重量0.6运费10的多规格商品",
			"postage": 10.00,
			"skus":{
				"M": {
					"price": 50.00,
					"weight": 0.6,
					"stocks": 99999
				},
				"S": {
					"price": 50.00,
					"weight": 0.6,
					"stocks": 99999
				}
			}
		}]
		"""

	@ginger-mall @order @wip
	Scenario: 1. 购买单个商品，使用免运费商品
		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10重量1运费0元商品",
				"count": 2
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 20.00,
			"delivery_items": [{
				"final_money": 20.00,
				"product_price": 20.00,
				"postage": 0.00
			}]
		}
		"""

	@ginger-mall @order
	Scenario: 2. 购买单个商品，使用统一运费商品
		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10重量1运费15元商品",
				"count": 2
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 35.00,
			"delivery_items": [{
				"final_money": 35.00,
				"product_price": 20.00,
				"postage": 15.00
			}]
		}
		"""

	@ginger-mall @order
	Scenario:3. 购买多种商品，使用统一运费
		Given lucy访问'jobs'的商城
		#有免运费商品
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10重量1运费0元商品",
				"count": 1
			}, {
				"name": "价格10重量1运费15元商品",
				"count": 1
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 35.00,
			"delivery_items": [{
				"final_money": 35.00,
				"product_price": 20.00,
				"postage": 15.0
			}]
		}
		"""

		#没有免运费商品
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10重量1运费10元商品",
				"count": 1
			}, {
				"name": "价格10重量1运费15元商品",
				"count": 1
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 45.00,
			"delivery_items": [{
				"final_money": 45.00,
				"product_price": 20.00,
				"postage": 25.0
			}]
		}
		"""

