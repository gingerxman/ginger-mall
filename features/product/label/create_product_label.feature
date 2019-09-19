Feature: 创建标签
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
Scenario: 管理员创建标签
	# 初始验证
	Given jobs登录系统
	Then jobs能看到标签列表
	"""
	[]
	"""

    #  创建标签
	When jobs创建标签
		"""
		[{
			"name": "name_1"
		},{
			"name": "name_2"
		},{
			"name": "name_3"
		}]
		"""
	Then jobs能看到标签列表
		"""
		[{
			"name": "name_3"
		},{
			"name": "name_2"
		},{
			"name": "name_1"
		}]
		"""
	Then jobs能看到可用的标签列表
		"""
		[{
			"name": "name_3"
		},{
			"name": "name_2"
		},{
			"name": "name_1"
		}]
		"""

	# bill验证
	Given bill登录系统
	Then bill能看到标签列表
	"""
	[]
	"""
