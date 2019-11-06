Feature: 将商品添加到购物车
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
			"name": "尺寸",
			"values": [{
				"name": "M"
			}, {
				"name": "S"
			}]
		}, {
			"name": "颜色",
			"values": [{
				"name": "黑色"
			}, {
				"name": "蓝色"
			}]
		}]
		"""
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"price": 3.00
		}, {
			"name": "商品2",
			"price": 5.00
		}, {
			"name": "商品3",
			"skus": {
				"M": {
					"price": 7.00,
					"stocks": 2
				},
				"S": {
					"price": 8.00,
					"stocks": 999
				}
			}
		}, {
			"name": "商品4",
			"skus": {
				"M": {
					"price": 9.00,
					"stocks": 999
				}
			}
		}, {
			"name": "商品5",
			"skus": {
				"S": {
					"price": 10.00,
					"stocks": 999
				}
			}
		}, {
			"name": "商品6",
			"skus": {
				"S 黑色": {
					"price": 10.00,
					"stocks": 999
				},
				"M 蓝色": {
					"price": 11.00,
					"stocks": 999
				}
			}
		}]
		"""

	@ginger-mall @mall @shopping_cart @wip
	Scenario: 1. 手机用户放入单个商品到购物车
		# 初始验证
		Given lucy注册为App用户
		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 1
		}]
		"""
		# zhouuxn验证
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1",
					"price": 3.00,
					"count": 1
				}]
			}],
			"invalid_products": []
		}
		"""

		#再次加入，验证
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 2
		}]
		"""
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 3
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1",
					"price": 3.00,
					"count": 6
				}]
			}],
			"invalid_products": []
		}
		"""

		#yangmi验证
		Given yangmi访问'jobs'的商城
		Then yangmi能获得购物车
		"""
		{
			"product_groups": [],
			"invalid_products": []
		}
		"""

	@ginger-mall @mall @shopping_cart
	Scenario:2 放入多个商品到购物车
		jobs添加商品后
		1. lucy能在webapp中将jobs添加的商品放入购物车
		2. 多次放入不同商品会增加购物车中商品的条数

		Given lucy访问'jobs'的商城
		#初始验证购物车商品数量
		Then lucy能获得购物车中商品数量为'0'

		#添加商品到购物车
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 1
		}]
		"""
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品2",
			"count": 2
		}]
		"""
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品6",
			"sku": "M 蓝色",
			"count": 2
		}]
		"""

		#验证购物车商品列表
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1",
					"price": 3.00,
					"count": 1
				}, {
					"name": "商品2",
					"price": 5.00,
					"count": 2
				}, {
					"name": "商品6",
					"price": 11.00,
					"sku": "M 蓝色",
					"count": 2
				}]
			}],
			"invalid_products": []
		}
		"""
		#验证购物车商品数量
		Then lucy能获得购物车中商品数量为'3'

	@ginger-mall @mall @shopping_cart
	Scenario:4 商品添加到购物车后，后台对商品进行上下架管理
		lucy在webapp中将jobs的商品加入到购物车后，jobs对此商品进行删除操作
		1.lucy查看jobs的webapp购物车，此商品已无效
		2.不影响购物车的其他商品

		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 1
		}, {
			"name": "商品2",
			"count": 1
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1"
				}, {
					"name": "商品2"
				}]
			}],
			"invalid_products": []
		}
		"""
		Then lucy能获得购物车中商品数量为'2'

		#jobs部分下架商品
		Given jobs登录系统
		When jobs将商品移动到'待售'货架
		"""
		["商品1"]
		"""
		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品2"
				}]
			}],
			"invalid_products": [{
				"name": "商品1",
				"count": 1
			}]
		}
		"""
		Then lucy能获得购物车中商品数量为'1'

		#jobs全部下架商品
		Given jobs登录系统
		When jobs将商品移动到'待售'货架
		"""
		["商品2"]
		"""
		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [],
			"invalid_products": [{
				"name": "商品1",
				"count": 1
			}, {
				"name": "商品2",
				"count": 1
			}]
		}
		"""
		Then lucy能获得购物车中商品数量为'0'

	@ginger-mall @mall @shopping_cart
	Scenario:5.1 商品添加到购物车后，改变商品规格：库存减少为0
		lucy在webapp中将jobs的商品加入到购物车后，jobs将此商品的商品规格进行修改
		1.lucy查看jobs的webapp购物车，此商品已无效
		2.lucy可以清空无效商品

		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品3",
			"sku": "M",
			"count": 1
		}, {
			"name": "商品3",
			"sku": "S",
			"count": 2
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"price": 7.00,
					"count": 1,
					"sku": "M"
				}, {
					"name": "商品3",
					"price": 8.00,
					"count": 2,
					"sku": "S"
				}]
			}],
			"invalid_products": []
		}
		"""

		#更改规格M的库存为0
		Given jobs登录系统
		When jobs更新商品'商品3'
		"""
		{
			"name": "商品3",
			"skus":{
				"M": {
					"price": 7.00,
					"stocks": 0
				},
				"S": {
					"price": 8.00,
					"stocks": 999
				}
			}
		}
		"""
		Then jobs能获取商品'商品3'
		"""
		{
			"name": "商品3",
			"skus":{
				"M": {
					"price": 7.00,
					"stocks": 0
				},
				"S": {
					"price": 8.00,
					"stocks": 999
				}
			}
		}
		"""

		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"price": 8.00,
					"count": 2,
					"sku": "S"
				}]
			}],
			"invalid_products": [{
				"name": "商品3",
				"count": 1,
				"sku": "M"
			}]
		}
		"""

	@ginger-mall @mall @shopping_cart
	Scenario:5.2 商品添加到购物车后，改变商品规格：从定制规格变为标准规格
		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品3",
			"sku": "M",
			"count": 1
		}, {
			"name": "商品4",
			"sku": "M",
			"count": 1
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"sku": "M"
				}, {
					"name": "商品4",
					"sku": "M"
				}]
			}],
			"invalid_products": []
		}
		"""

		#更改商品规格
		Given jobs登录系统
		When jobs更新商品'商品4'
		"""
		{
			"name": "商品4",
			"skus":{
				"standard": {
					"price": 9.50
				}
			}
		}
		"""
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"sku": "M"
				}]
			}],
			"invalid_products": [{
				"name": "商品4",
				"sku": "M"
			}]
		}
		"""

	@ginger-mall @mall @shopping_cart
	Scenario:5.3 商品添加到购物车后，改变商品规格：从标准规格变为定制规格
		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 1
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1"
				}]
			}],
			"invalid_products": []
		}
		"""

		#更改商品规格
		Given jobs登录系统
		When jobs更新商品'商品1'
		"""
		{
			"name": "商品1",
			"skus":{
				"M": {
					"price": 9.50,
					"stocks": 2
				}
			}
		}
		"""
		#验证
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [],
			"invalid_products": [{
				"name": "商品1"
			}]
		}
		"""

	@ginger-mall @mall @shopping_cart
	Scenario:6 商品添加到购物车后，进行删除
		lucy加入jobs的商品到购物车后
		1.可以对购物车的商品进行删除

		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品1",
			"count": 1
		}, {
			"name": "商品2",
			"count": 1
		}, {
			"name": "商品4",
			"sku": "M",
			"count": 3
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品1"
				}, {
					"name": "商品2"
				}, {
					"name": "商品4"
				}]
			}],
			"invalid_products": []
		}
		"""

		#删除购物车商品
		When lucy从购物车中删除商品
		"""
		["商品1", "商品2"]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品4"
				}]
			}],
			"invalid_products": []
		}
		"""

	@ginger-mall @mall @shopping_cart
	Scenario:7 商品添加到购物车后，后台对商品的价格，库存进行修改（库存数量不为0）
		lucy在webapp中将jobs的商品加入到购物车后，jobs将此商品的商品规格进行修改
		1.lucy查看jobs的webapp购物车，此商品有效，价格与库存为更改后的值

		Given lucy访问'jobs'的商城
		When lucy加入jobs的商品到购物车
		"""
		[{
			"name": "商品3",
			"sku": "M",
			"count": 1
		}, {
			"name": "商品3",
			"sku": "S",
			"count": 1
		}]
		"""
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"price": 7.00,
					"count": 1,
					"sku": "M"
				}, {
					"name": "商品3",
					"price": 8.00,
					"count": 1,
					"sku": "S"
				}]
			}],
			"invalid_products": []
		}
		"""

		#更改规格M的库存为3 S的价格为10
		Given jobs登录系统
		When jobs更新商品'商品3'
		"""
		{
			"name": "商品3",
			"skus":{
				"M": {
					"price": 7.00,
					"stocks": 3
				},
				"S": {
					"price": 10.00,
					"stocks": 990
				}
			}
		}
		"""
		Then jobs能获取商品'商品3'
		"""
		{
			"name": "商品3",
			"skus":{
				"M": {
					"price": 7.00,
					"stocks": 3
				},
				"S": {
					"price": 10.00,
					"stocks": 990
				}
			}
		}
		"""

		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能获得购物车
		"""
		{
			"product_groups": [{
				"products": [{
					"name": "商品3",
					"price": 7.00,
					"count": 1,
					"sku": "M",
					"stocks": 3
				}, {
					"name": "商品3",
					"price": 10.00,
					"count": 1,
					"sku": "S"
				}]
			}],
			"invalid_products": []
		}
		"""