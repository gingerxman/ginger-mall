Feature: 更新收货地址
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
		Given lucy注册为App用户
		Given lily注册为App用户

@ginger-mall @mall
Scenario: 手机用户更新收货地址
	Given lucy访问'jobs'的商城
    #  创建收货地址
	When lucy创建收货地址
		"""
		[{
			"name": "Lucy1",
			"phone": "13811223355",
			"area": "江苏省 南京市 秦淮区",
			"address": "水平方"
		},{
			"name": "Lucy2",
			"phone": "13811223356",
			"area": "江苏省 无锡市 滨湖区",
			"address": "蠡园"
		}]
		"""
	Then lucy能看到收货地址列表
		"""
		[{
			"name": "Lucy1",
			"phone": "13811223355",
			"area": "江苏省 南京市 秦淮区",
			"address": "水平方"
		},{
			"name": "Lucy2",
			"phone": "13811223356",
			"area": "江苏省 无锡市 滨湖区",
			"address": "蠡园"
		}]
		"""

	#更新收货地址
	When lucy修改收货地址'水平方'的信息
		"""
		{
			"name": "Lucy1*",
			"phone": "13900000000",
			"area": "北京市 北京市 海淀区",
			"address": "嘉华大厦"
		}
		"""
	Then lucy能看到收货地址列表
		"""
		[{
			"name": "Lucy1*",
			"phone": "13900000000",
			"area": "北京市 北京市 海淀区",
			"address": "嘉华大厦"
		},{
			"name": "Lucy2",
			"phone": "13811223356",
			"area": "江苏省 无锡市 滨湖区",
			"address": "蠡园"
		}]
		"""
