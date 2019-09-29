Feature: 更新商品
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
			"type": "图片",
			"values": [{
				"name": "黑色",
				"image": "black.png"
			}, {
				"name": "白色",
				"image": "white.png"
			}]
		}, {
			"name": "尺寸",
			"type": "文字",
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
			"name": "东坡肘子",
			"promotion_title": "促销的东坡肘子",
			"detail": "东坡肘子的详情",
			"images": [{
				"url": "/static/test_resource_img/hangzhou1.jpg"
			}],
			"skus": {
				"standard": {
					"price": 11.12,
					"cost_price": 10.12,
					"stocks": 10
				}
			}
		}, {
			"name": "叫花鸡",
			"images": [{
				"url": "jiaohuaji1.jpg"
			}, {
				"url": "jiaohuaji2.jpg"
			}, {
				"url": "jiaohuaji3.jpg"
			}],
			"skus": {
				"standard": {
					"price": 12.00,
					"stocks": 3
				}
			}
		}, {
			"name": "莲藕排骨汤",
			"skus": {
				"黑色 M": {
					"price": 1.1
				},
				"白色 S": {
					"price": 2.2
				}
			}
		}]
		"""

	@ginger-mall @product
	Scenario: 1. 修改商品基本信息
		# 初始验证
		Given jobs登录系统
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"name": "东坡肘子",
			"promotion_title": "促销的东坡肘子",
			"detail": "东坡肘子的详情"
		}
		"""

		# 更新商品
		When jobs更新商品'东坡肘子'
		"""
		{
			"name": "东坡肘子*",
			"promotion_title": "促销的东坡肘子*",
			"detail": "东坡肘子的详情*"
		}
		"""
		Then jobs能获取商品'东坡肘子*'
		"""
		{
			"name": "东坡肘子*",
			"promotion_title": "促销的东坡肘子*",
			"detail": "东坡肘子的详情*"
		}
		"""

	@ginger-mall @product
	Scenario: 2. 修改商品的图片
		# 初始验证
		Given jobs登录系统
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"thumbnail": "/static/test_resource_img/hangzhou1.jpg",
			"medias": [{
				"url": "/static/test_resource_img/hangzhou1.jpg"
			}]
		}
		"""

		# 更新商品图片
		When jobs更新商品'东坡肘子'
		"""
		{
			"images": [{
				"url": "new_zhouzi1.jpg"
			}, {
				"url": "new_zhouzi2.jpg"
			}]
		}
		"""
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"thumbnail": "new_zhouzi1.jpg",
			"medias": [{
				"url": "new_zhouzi1.jpg"
			}, {
				"url": "new_zhouzi2.jpg"
			}]
		}
		"""

#	@ginger-mall @product
#	Scenario: 3. 修改商品的分类
#		# 初始验证
#		Given jobs登录系统
#		Then jobs能获取商品'东坡肘子'
#		"""
#		{
#			"categories": ["分类1"]
#		}
#		"""
#
#		# 更新分类
#		When jobs更新商品'东坡肘子'
#		"""
#		{
#			"categories": ["分类2", "分类3"]
#		}
#		"""
#		Then jobs能获取商品'东坡肘子'
#		"""
#		{
#			"categories": ["分类2", "分类3"]
#		}
#		"""
#
#		# 清空分类
#		When jobs更新商品'东坡肘子'
#		"""
#		{
#			"categories": []
#		}
#		"""
#		Then jobs能获取商品'东坡肘子'
#		"""
#		{
#			"categories": []
#		}
#		"""

	@ginger-mall @product
	Scenario: 4. 修改商品的标准规格
		# 初始验证
		Given jobs登录系统
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 11.12,
					"cost_price": 10.12,
					"stocks": 10
				}
			}
		}
		"""

		# 更新规格
		When jobs更新商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 22.00,
					"cost_price": 20.0,
					"stocks": 11
				}
			}
		}
		"""
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 22.00,
					"cost_price": 20.0,
					"stocks": 11
				}
			}
		}
		"""

	@ginger-mall @product @wip
	Scenario: 5. 商品规格在标准规格和定制规格之间切换
		# 初始验证
		Given jobs登录系统
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 11.12,
					"cost_price": 10.12,
					"stocks": 10
				}
			}
		}
		"""

		# 标准规格变为定制规格
		When jobs更新商品'东坡肘子'
		"""
		{
			"skus": {
				"黑色 M": {
					"price": 22.00,
					"cost_price": 20.0,
					"stocks": 1
				},
				"颜色:blue 尺寸:XS": {
					"price": 32.00,
					"cost_price": 30.0,
					"stocks": 2
				}
			}
		}
		"""
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"skus": {
				"黑色 M": {
					"price": 22.00,
					"cost_price": 20.0,
					"stocks": 1
				},
				"blue XS": {
					"price": 32.00,
					"cost_price": 30.0,
					"stocks": 2
				}
			}
		}
		"""

		# 定制规格变为标准规格
		When jobs更新商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 2.2,
					"cost_price": 1.1,
					"stocks": 1
				}
			}
		}
		"""
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"skus": {
				"standard": {
					"price": 2.2,
					"cost_price": 1.1,
					"stocks": 1
				}
			}
		}
		"""

	@ginger-mall @product
	Scenario: 6. 修改商品的定制规格
		Given jobs登录系统
		#删除规格"白色 S"，增加规格"白色 M", 修改规格"黑色 M"
		When jobs更新商品'莲藕排骨汤'
		"""
		{
			"skus": {
				"黑色 M": {
					"price": 22.00,
					"cost_price": 20,
					"stocks": 2
				},
				"白色 M": {
					"price": 32.00,
					"cost_price": 30,
					"stocks": 3
				}
			}
		}
		"""
		Then jobs能获取商品'莲藕排骨汤'
		"""
		{
			"skus": {
				"黑色 M": {
					"price": 22.00,
					"cost_price": 20,
					"stocks": 2
				},
				"白色 M": {
					"price": 32.00,
					"cost_price": 30,
					"stocks": 3
				}
			}
		}
		"""

