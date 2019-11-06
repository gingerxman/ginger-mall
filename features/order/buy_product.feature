Feature: 购买商品
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
				"name": "黑色",
				"image": "black.png"
			}, {
				"name": "白色",
				"image": "white.png"
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

		Given lucy注册为App用户

	@ginger-mall @order
	Scenario: 1. App用户能购买单个商品
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"price": 9.90
		}]
		"""

		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "创新大厦",
			"pay_type":"微信支付",
			"products": [{
				"name": "商品1",
				"count": 2
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 19.8,
			"delivery_items": [{
				"status": "待支付",
				"final_money": 19.8,
				"products": [{
					"name": "商品1",
					"price": 9.90,
					"count": 2
				}],
				"ship_name": "lucy",
				"ship_tel": "13811223344",
				"ship_area": "江苏省 南京市 秦淮区",
				"ship_address": "创新大厦"
			}]
		}
		"""
		Then lucy能获得最新订单的订单状态为'待支付'

	@ginger-mall @order
	Scenario: 2. App用户能购买有规格的商品
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"skus": {
				"黑色 M": {
					"price": 10.00
				}
			}
		}]
		"""

		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"sku": "黑色 M",
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
				"products": [{
					"name": "商品1",
					"sku": "黑色 M",
					"price": 10.00,
					"count": 2
				}]
			}]
		}
		"""

	@ginger-mall @order
	Scenario: 3. App用户不能购买下架的商品
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"price": 9.90
		}]
		"""
		When jobs将商品移动到'待售'货架
		"""
		["商品1"]
		"""

		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 1
			}],
			"error": "create_order_fail:product_off_shelve"
		}
		"""

	@ginger-mall @order @wip
	Scenario: 4. 购买商品影响库存：购买数量等于库存数量
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"skus": {
				"standard": {
					"price": 5.00,
					"stocks": 2
				}
			}
		}]
		"""

		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 2
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付"
		}
		"""

		#jobs验证库存变化
		Given jobs登录系统
		Then jobs能获取商品'商品1'
		"""
		{
			"name": "商品1",
			"skus": {
				"standard": {
					"stocks": 0
				}
			}
		}
		"""

		#lucy再次购买，库存为0，获得错误消息提示
		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 1
			}],
			"error": "create_order_fail:not_enough_stocks"
		}
		"""

	@ginger-mall @order
	Scenario: 5. 购买商品影响库存：购买数量大于库存数量
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"skus": {
				"standard": {
					"price": 5.00,
					"stock_type": "有限",
					"stocks": 2
				}
			}
		}]
		"""

		#lucy购买数量超过库存，库存为0，获得错误消息提示
		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 3
			}],
			"error": "create_order_fail:not_enough_stocks"
		}
		"""

		#jobs验证库存无变化
		Given jobs登录系统
		Then jobs能获取商品'商品1'
		"""
		{
			"name": "商品1",
			"skus": {
				"standard": {
					"stocks": 2
				}
			}
		}
		"""

	@ginger-mall @order
	Scenario: 6. 购买价格发生变化的商品
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"skus": {
				"standard": {
					"price": 5.00,
					"stocks": 2
				}
			}
		}]
		"""

		#lucy购买，提交的价格与后台价格不一致，获得错误消息提示
		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 1,
				"price": 1
			}],
			"error": "create_order_fail:price_change"
		}
		"""

		#jobs验证库存无变化
		Given jobs登录系统
		Then jobs能获取商品'商品1'
		"""
		{
			"name": "商品1",
			"skus": {
				"standard": {
					"stocks": 2
				}
			}
		}
		"""

	@ginger-mall @order
	Scenario: 7. 购买价格为0的商品
		购买价格为0的自动完成订单，订单状态自动为"已完成"

		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"skus": {
				"standard": {
					"price": 0.0
				}
			}
		}]
		"""

		Given lucy访问'jobs'的商城
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 1
			}],
			"order_type": "order:auto_finish"
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "已完成"
		}
		"""
		Then lucy能获得最新订单的订单状态为'已完成'

