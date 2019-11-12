Feature: 购买商品后完成订单
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

	@ginger-mall @order
	Scenario: 1. 手机购买单个商品，完成订单，进行清算
		#xiaocheng验证初始平台收入
		Given jobs登录系统
		Then jobs能获得公司的虚拟资产'cash'
		"""
		{
			"balance": 0
		}
		"""

		#jobs添加商品，验证初始酒吧收入
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"price": 9.90
		}]
		"""
		Then jobs能获得公司的虚拟资产'cash'
		"""
		{
			"balance": 0
		}
		"""

		Given lucy访问'jobs'的商城
		Then lucy能获得虚拟资产'rmb'
		"""
		{
			"balance": 0
		}
		"""
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
			"status": "待支付",
			"final_money": 19.8,
			"delivery_items": [{
				"status": "待支付",
				"final_money": 19.8
			}]
		}
		"""
		Then lucy能获得最新订单的订单状态为'待支付'

		#支付订单
		When lucy支付最新订单
		Then lucy能获得最新订单的订单状态为'待确认'

		#jobs确认订单
		Given jobs登录系统
		When jobs确认最新发货单
		Then jobs能获得最新出货单的状态为'待发货'

		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得最新出货单的状态为'待发货'

		#jobs发货
		Given jobs登录系统
		When jobs对最新发货单进行发货
		"""
		{
			"enable_logistics": false
		}
		"""
		Then jobs能获得最新出货单的状态为'已发货'

		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得最新出货单的状态为'已发货'
		When lucy在手机上完成最新出货单
		Then lucy能获得最新出货单的状态为'已完成'
		Then lucy能获得虚拟资产'rmb'
		"""
		{
			"balance": -19.8
		}
		"""

		#jobs验证收入
		Given jobs登录系统
		Then jobs能获得公司的虚拟资产'cash'
		"""
		{
			"balance": 19.8
		}
		"""

