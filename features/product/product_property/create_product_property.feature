Feature: 创建商品规格
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

@ginger-mall @product
Scenario: 管理员创建商品属性
	# 初始验证
	Given jobs登录系统
	Then jobs能看到商品属性列表
	"""
	[]
	"""

    # 创建商品属性
	When jobs创建商品属性
	"""
	[{
		"name": "颜色",
		"values": [{
			"name": "黑色",
			"image": "/static/test_resource_img/hangzhou1.jpg"
		}, {
			"name": "白色",
			"image": "/static/test_resource_img/hangzhou2.jpg"
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
	Then jobs能看到商品属性列表
	"""
	[{
		"name": "颜色",
		"values": [{
			"name": "黑色",
			"image": "/static/test_resource_img/hangzhou1.jpg"
		}, {
			"name": "白色",
			"image": "/static/test_resource_img/hangzhou2.jpg"
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

	# bill验证
	Given bill登录系统
	Then bill能看到商品属性列表
	"""
	[]
	"""

