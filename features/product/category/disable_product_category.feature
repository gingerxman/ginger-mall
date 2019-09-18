Feature: 禁用商品分类
	Background:
		Given 重置服务
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
	Scenario: 酒吧管理员禁用商品分类

		# 创建商品分类
		Given jobs登录系统
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

		#禁用分类
		When jobs禁用商品分类'分类23'
		Then jobs能看到'分类11'分类下的商品分类列表
		"""
		[{
			"name": "分类21"
		},{
			"name": "分类22"
		}]
		"""

		#lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能在手机上看到'分类11'分类下的商品分类列表
		"""
		[{
			"name": "分类21"
		},{
			"name": "分类22"
		}]
		"""

		#启用分类
		Given jobs登录系统
		When jobs启用商品分类'分类23'
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

		Given lucy访问'jobs'的商城
		Then lucy能在手机上看到'分类11'分类下的商品分类列表
		"""
		[{
			"name": "分类21"
		},{
			"name": "分类22"
		},{
			"name": "分类23"
		}]
		"""

