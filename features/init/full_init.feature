Feature: 初始化系统数据

@full_init
Scenario: 初始化系统数据
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
			"食品": []
		},{
			"女装": []
		},{
			"美妆": []
		},{
			"男装": []
		},{
			"亲子": []
		},{
			"运动户外": []
		}]
		"""
