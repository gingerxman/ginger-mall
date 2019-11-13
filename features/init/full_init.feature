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
			"美味": []
		},{
			"小吃": []
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
		},{
			"潮品": []
		}]
		"""
	When jobs创建商品属性
		"""
		[{
			"name": "辣度",
			"values": [{
				"name": "微辣"
			}, {
				"name": "中辣"
			}, {
				"name": "变态辣"
			}]
		}, {
			"name": "分量",
			"values": [{
				"name": "大碗"
			}, {
				"name": "小碗"
			}]
		}]
		"""
	When jobs添加商品
		"""
		[{
			"name": "东坡肘子",
			"category": "美味",
			"code": "zhouzi_1",
			"promotion_title": "促销的东坡肘子",
			"detail": "东坡肘子的详情",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/hangzhou1.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/hangzhou2.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/hangzhou3.jpg"
			}],
			"skus": {
				"standard": {
					"price": 9.9,
					"cost_price": 9.0,
					"weight": 5.0,
					"stocks": 9999999
				}
			}
		}, {
			"name": "松鼠桂鱼",
			"category": "美味",
			"detail": "松鼠桂鱼的详情<div style='color:red'>真好吃</div>",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/yu1.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/yu2.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/yu3.jpg"
			}],
			"skus": {
				"小碗 变态辣": {
					"price": 11.12,
					"cost_price": 1.1,
					"weight": 5.0,
					"stocks": 9999999
				},
				"大碗 中辣": {
					"price": 21.12,
					"cost_price": 2.2,
					"weight": 25.0,
					"stocks": 3
				}
			}
		}, {
			"name": "热干面",
			"category": "",
			"detail": "热干面的详情<div style='color:red'>真好吃</div>",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/mian2.jpg"
			}],
			"skus": {
				"standard": {
					"price": 10.00,
					"cost_price": 8.00,
					"weight": 0.0,
					"stocks": 10
				}
			}
		}, {
			"name": "十三香小龙虾",
			"category": "小吃",
			"detail": "小龙虾的详情<div style='color:red'>真好吃</div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/xia_2.jpg' /></div>",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/xia_1.jpg"
			}],
			"skus": {
				"standard": {
					"price": 30.00,
					"cost_price": 10.00,
					"weight": 0.0,
					"stocks": 10
				}
			}
		}]
		"""

	#bill创建商品，提交平台审核
	Given bill登录系统
	When bill创建商品属性
		"""
		[{
			"name": "口味",
			"values": [{
				"name": "地中海海盐味"
			}, {
				"name": "日式青芥味"
			}, {
				"name": "美式番茄味"
			}]
		}]
		"""
	When bill添加商品
		"""
		[{
			"name": "MT原味烟草 MOTI国烟烟油弹 魔笛电子烟套装无焦油小烟 辅助戒烟 替烟解瘾不漏油 封闭式换弹【4枚/盒】",
			"category": "潮品",
			"code": "modi_1",
			"promotion_title": "【京东211配送.时效有保障】无焦油二手烟.欧盟CE认证.1盒相当于一条烟.MT与MOTI烟弹新老包装随机配送【烟弹补充装.MOTI魔笛通用】",
			"detail": "<div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/detail1.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/detail2.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/detail3.png' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/detail4.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/detail5.jpg' /></div>",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/moti1.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/moti2.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/moti/moti3.jpg"
			}],
			"skus": {
				"standard": {
					"price": 168,
					"cost_price": 168,
					"stocks": 9999999
				}
			}
		}, {
			"name": "单身狗粮（SINGLE DOG） 地中海盐味马铃薯片71g 网红膨化小吃袋装",
			"category": "小吃",
			"detail": "<div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dogdetail1.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dogdetail2.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dogdetail3.jpg' /></div><div><img src='http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dogdetail4.jpg' />",
			"images": [{
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dog1.jpg"
			}, {
				"url": "http://vxiaocheng-jh.oss-cn-beijing.aliyuncs.com/peanut/singledog/dog2.jpg"
			}],
			"skus": {
				"地中海海盐味": {
					"price": 10.50,
					"cost_price": 10.50,
					"stocks": 9999999
				},
				"日式青芥味": {
					"price": 11.50,
					"cost_price": 10.50,
					"stocks": 9999999
				},
				"美式番茄味": {
					"price": 12.50,
					"cost_price": 10.50,
					"stocks": 9999999
				}
			}
		}]
		"""


	#lucy购买
	Given lucy注册为App用户
	Given lucy访问'jobs'的商城
	When lucy创建收货地址
		"""
		[{
			"name": "Lucy1",
			"phone": "13811223355",
			"area": "江苏省 南京市 秦淮区",
			"address": "夫子庙"
		},{
			"name": "Lucy2",
			"phone": "13811223356",
			"area": "江苏省 无锡市 滨湖区",
			"address": "海岸城"
		}]
		"""
	When lucy设置收货地址'夫子庙'为默认收货地址
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "松鼠桂鱼",
				"sku": "小碗 变态辣",
				"count": 3
			}, {
				"name": "热干面",
				"count": 2
			}]
		}
		"""
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"message": "上午上班，下午送",
			"products": [{
				"name": "东坡肘子",
				"count": 1
			}]
		}
		"""
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "东坡肘子",
				"count": 2
			}]
		}
		"""
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "东坡肘子",
				"count": 3
			}]
		}
		"""
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "东坡肘子",
				"count": 4
			}]
		}
		"""
	When lucy购买'jobs'的商品
		"""
		{
			"ship_name": "Lucy",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "东坡肘子",
				"count": 5
			}]
		}
		"""

	#lily购买
	Given lily注册为App用户
	Given lily访问'jobs'的商城
	When lily创建收货地址
		"""
		[{
			"name": "lily",
			"phone": "13811221122",
			"area": "北京市 北京市 海淀区",
			"address": "嘉华大厦"
		}]
		"""
	When lily设置收货地址'嘉华大厦'为默认收货地址
	When lily购买'jobs'的商品
		"""
		{
			"ship_name": "Lily",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "夫子庙",
			"products": [{
				"name": "东坡肘子",
				"count": 10
			}]
		}
		"""
	Given lily访问'bill'的商城
	When lily购买'bill'的商品
		"""
		{
			"ship_name": "Lily",
			"ship_tel": "13811223344",
			"ship_area": "江苏省 南京市 秦淮区",
			"ship_address": "桃花渡",
			"products": [{
				"sku": "地中海海盐味",
				"name": "单身狗粮（SINGLE DOG） 地中海盐味马铃薯片71g 网红膨化小吃袋装",
				"count": 1
			}]
		}
		"""

	#gail 购买
	Given gal注册为App用户
	Given gal访问'jobs'的商城
