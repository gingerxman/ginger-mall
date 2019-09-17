Feature: 创建商品分类
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

@ginger-mall @product @wip
Scenario: 酒吧管理员创建商品分类
	# 初始验证
	Given jobs登录系统
	Then jobs能看到'一级'分类下的商品分类列表
	"""
	[]
	"""

    #  创建商品分类
	When jobs创建商品分类
		"""
		[{
			"分类11": [{
				"分类21": []
			},{
				"分类22": []
			},{
				"分类23": [{
					"分类31": []
				}]
			}]
		},{
			"分类12": [{
				"分类24": []
			}]
		},{
			"分类13": []
		}]
		"""
	Then jobs能看到'一级'分类下的商品分类列表
		"""
		[{
			"name": "分类11"
		},{
			"name": "分类12"
		},{
			"name": "分类13"
		}]
		"""
	Then jobs能看到'分类11'分类下的商品分类列表
		"""
		[{
			"name": "分类21"
		},{
			"name": "分类22"
		},{
			"name": "分类23"
		}]
		"""

	# bill验证
	Given bill登录系统
	Then bill能看到'一级'分类下的商品分类列表
	"""
	[]
	"""
