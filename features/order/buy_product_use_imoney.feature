Feature: 使用imoney购买商品
	Background:
		Given 系统配置虚拟资产
		"""
		[{
			"code": "eth",
			"display_name": "以太坊",
			"exchange_rate": 1,
			"enable_fraction": false,
			"is_payable": true,
			"is_debtable": false
		}, {
			"code": "bitcoin",
			"display_name": "比特币",
			"exchange_rate": 1,
			"enable_fraction": false,
			"is_payable": true,
			"is_debtable": true
		}]
		"""

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
		When jobs添加商品
		"""
		[{
			"name": "价格10.5运费1商品",
			"price": 10.5,
			"weight": 1
		}, {
			"name": "价格20运费0商品",
			"price": 20.0
		}, {
			"name": "价格50运费0商品",
			"price": 50.00
		}]
		"""

		Given lucy注册为App用户

	@ginger-mall @order @order.imoney
	Scenario: 1 使用不足商品价格的虚拟资产进行购买
		#初始化bitcoin
		Given lucy访问'jobs'的商城
		When lucy充值'9'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 9
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10.5运费1商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 9
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 1.5,
			"delivery_items": [{
				"final_money": 1.50,
				"product_price": 10.50,
				"imoneys": [{
					"code": "bitcoin",
					"count": 9,
					"deduction_money": 9.0
				}]
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 9,
				"deduction_money": 9.0
			}]
		}
		"""

		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 0
		}
		"""

	@ginger-mall @order @order.imoney
	Scenario: 2 使用等于商品价格的虚拟资产进行购买
		Given lucy访问'jobs'的商城
		When lucy充值'50'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 50
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格50运费0商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 50
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"final_money": 0.0,
			"delivery_items": [{
				"final_money": 0.0,
				"product_price": 50.0,
				"imoneys": [{
					"code": "bitcoin",
					"count": 50,
					"deduction_money": 50.0
				}]
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 50,
				"deduction_money": 50.0
			}]
		}
		"""
		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 0
		}
		"""

	@ginger-mall @order @order.imoney
	Scenario: 3 使用等于商品价格的虚拟资产购买可自动完成的商品，订单自动完成
		Given lucy访问'jobs'的商城
		When lucy充值'50'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 50
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格50运费0商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 50
			}],
			"order_type": "order:auto_finish"
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "已完成",
			"final_money": 0.0,
			"delivery_items": [{
				"final_money": 0.0,
				"product_price": 50.0,
				"imoneys": [{
					"code": "bitcoin",
					"count": 50,
					"deduction_money": 50
				}]
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 50,
				"deduction_money": 50
			}]
		}
		"""
		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 0
		}
		"""

	@ginger-mall @order @order.imoney
	Scenario: 4 使用多于商品价格的虚拟资产进行购买
		Given lucy访问'jobs'的商城
		When lucy充值'12'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 12
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10.5运费1商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 11
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"final_money": 0.0,
			"delivery_items": [{
				"final_money": 0.0,
				"product_price": 10.5,
				"imoneys": [{
					"code": "bitcoin",
					"count": 11,
					"deduction_money": 11
				}]
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 11,
				"deduction_money": 11
			}]
		}
		"""
		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 1
		}
		"""


	@ginger-mall @order @order.imoney @wip
	Scenario: 5 使用超过余额的虚拟资产进行购买
		Given lucy访问'jobs'的商城
		When lucy充值'5'个'eth'
		Then lucy能获得虚拟资产'eth'
		"""
		{
			"balance": 5
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10.5运费1商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "eth",
				"count": 9
			}],
			"error": "not_enough_balance"
		}
		"""


		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'eth'
		"""
		{
			"balance": 5
		}
		"""

	@ginger-mall @order @order.imoney @wip
	Scenario: 6 使用多个虚拟资产进行购买
		Given lucy访问'jobs'的商城
		When lucy充值'5'个'eth'
		When lucy充值'5'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 5
		}
		"""
		Then lucy能获得虚拟资产'eth'
		"""
		{
			"balance": 5
		}
		"""

		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "价格10.5运费1商品",
				"count": 1
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 5
			}, {
				"code": "eth",
				"count": 5
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"final_money": 0.5,
			"delivery_items": [{
				"final_money": 0.50,
				"product_price": 10.50,
				"imoneys": [{
					"code": "bitcoin",
					"count": 5,
					"deduction_money": 5
				}, {
					"code": "eth",
					"count": 5,
					"deduction_money": 5
				}]
			}],
			"imoneys": [{
				"code": "bitcoin",
				"count": 5,
				"deduction_money": 5
			}, {
				"code": "eth",
				"count": 5,
				"deduction_money": 5
			}]
		}
		"""

		#验证虚拟资产使用情况
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 0
		}
		"""
		Then lucy能获得虚拟资产'eth'
		"""
		{
			"balance": 0
		}
		"""
