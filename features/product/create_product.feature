Feature: 创建商品
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

	@ginger-mall @product @wip
	Scenario: 1. 管理员创建标准规格商品
		# 初始验证
		Given jobs登录系统
		When jobs添加商品
		#东坡肘子(有分类，上架，无限库存，多轮播图), 包含其他所有信息
		#叫花鸡(无分类，下架，有限库存，单轮播图)
		"""
		[{
			"name": "东坡肘子",
			"category": "分类1",
			"code": "zhouzi_1",
			"promotion_title": "促销的东坡肘子",
			"detail": "东坡肘子的详情",
			"images": [{
				"url": "/static/test_resource_img/hangzhou1.jpg"
			}, {
				"url": "/static/test_resource_img/hangzhou2.jpg"
			}, {
				"url": "/static/test_resource_img/hangzhou3.jpg"
			}],
			"skus": {
				"standard": {
					"price": 11.12,
					"cost_price": 1.11,
					"stocks": 99999
				}
			}
		}, {
			"name": "叫花鸡",
			"detail": "叫花鸡的详情",
			"images": [{
				"url": "/static/test_resource_img/hangzhou2.jpg"
			}],
			"skus": {
				"standard": {
					"price": 12.00,
					"cost_price": 3.14,
					"stocks": 3
				}
			}
		}]
		"""

		Then jobs能获取商品'东坡肘子'
		"""
		{
			"name": "东坡肘子",
			"shelve_type": "上架",
			"promotion_title": "促销的东坡肘子",
			"detail": "东坡肘子的详情",
			"medias": [{
				"type": "image",
				"url": "/static/test_resource_img/hangzhou1.jpg"
			}, {
				"type": "image",
				"url": "/static/test_resource_img/hangzhou2.jpg"
			}, {
				"type": "image",
				"url": "/static/test_resource_img/hangzhou3.jpg"
			}],
			"skus": {
				"standard": {
					"price": 11.12,
					"cost_price": 1.11,
					"stocks": 99999
				}
			}
		}
		"""
		Then jobs能获取商品'叫花鸡'
		"""
		{
			"name": "叫花鸡",
			"shelve_type": "上架",
			"medias": [{
				"type": "image",
				"url": "/static/test_resource_img/hangzhou2.jpg"
			}],
			"skus": {
				"standard": {
					"price": 12.00,
					"cost_price": 3.14,
					"stocks": 3
				}
			}
		}
		"""
		#待售列表按添加时间倒序排列
		Then jobs能获得'在售'商品列表
		"""
		[{
			"name": "叫花鸡",
			"skus": {
				"standard": {
					"stocks": 3,
					"cost_price": 3.14,
					"price": 12.00
				}
			},
			"thumbnail": "/static/test_resource_img/hangzhou2.jpg"
		}, {
			"name": "东坡肘子",
			"skus": {
				"standard": {
					"stocks": 99999,
					"cost_price": 1.11,
					"price": 11.12

				}
			},
			"thumbnail": "/static/test_resource_img/hangzhou1.jpg"
		}]
		"""

		# bill验证
		Given bill登录系统
		Then bill能获得'在售'商品列表
		"""
		[]
		"""

	@gpeanut @product
	Scenario:2 添加定制规格商品
		Given jobs登录系统
		When jobs添加商品
		#东坡肘子：多个定制规格，包含有限和无限库存
		#叫花鸡：单个定制规格
		"""
		[{
			"name": "东坡肘子",
			"skus": {
				"黑色 M": {
					"price": 11.12,
					"purchase_price": 1.1,
					"weight": 5.0,
					"stock_type": "无限"
				},
				"白色 S": {
					"price": 21.12,
					"purchase_price": 2.2,
					"weight": 25.0,
					"stock_type": "有限",
					"stocks": 99
				}
			}
		}, {
			"name": "叫花鸡",
			"skus": {
				"黑色 S": {
					"price": 3.14,
					"purchase_price": 1.3,
					"weight": 3.14,
					"stock_type": "无限"
				}
			}
		}]
		"""
		Then jobs能获取商品'东坡肘子'
		"""
		{
			"name": "东坡肘子",
			"skus": {
				"黑色 M": {
					"price": 11.12,
					"purchase_price": 1.1,
					"weight": 5.0,
					"stock_type": "无限"
				},
				"白色 S": {
					"price": 21.12,
					"purchase_price": 2.2,
					"weight": 25.0,
					"stock_type": "有限",
					"stocks": 99
				}
			}
		}
		"""
		Then jobs能获取商品'叫花鸡'
		"""
		{
			"name": "叫花鸡",
			"skus": {
				"黑色 S": {
					"price": 3.14,
					"purchase_price": 1.3,
					"weight": 3.14,
					"stock_type": "无限"
				}
			}
		}
		"""
		#待售列表按添加时间倒序排列
		Then jobs能获得'在售'商品列表
		"""
		[{
			"name": "叫花鸡",
			"skus": {
				"黑色 S": {
					"price": 3.14,
					"purchase_price": 1.3,
					"weight": 3.14,
					"stock_type": "无限"
				}
			}
		}, {
			"name": "东坡肘子",
			"skus": {
				"黑色 M": {
					"price": 11.12,
					"purchase_price": 1.1,
					"weight": 5.0,
					"stock_type": "无限"
				},
				"白色 S": {
					"price": 21.12,
					"purchase_price": 2.2,
					"weight": 25.0,
					"stock_type": "有限",
					"stocks": 99
				}
			}
		}]
		"""
