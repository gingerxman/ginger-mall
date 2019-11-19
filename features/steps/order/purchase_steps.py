# -*- coding: utf-8 -*-
import json

from behave import *

from features.bdd import util as bdd_util
from features.bdd.client import RestClient
from features.steps.product import product_steps
from features.steps import plutus_steps
from features.steps.mall import ship_info_steps

def get_product_category_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_category where name = %s", [name])
	return objs[0]['id']

def get_latest_order_bid():
	objs = bdd_util.exec_sql("select * from order_order where type = 1 order by id desc limit 1", [])
	if len(objs) == 0:
		objs = bdd_util.exec_sql("select * from order_order where type = 3 order by id desc limit 1", [])
	return objs[0]['bid']

def get_latest_invoice_bid():
	objs = bdd_util.exec_sql("select * from order_order where type = 2 order by id desc limit 1", [])
	return objs[0]['bid']

def get_sku_code_from_display_name(sku_name):
	target_pattern = '_%s' % sku_name
	objs = bdd_util.exec_sql("select * from product_sku", [])
	for obj in objs:
		if target_pattern in obj['code']:
			return obj['name']

def get_salesman_id_by_username(client, username):
	user_id = plutus_steps.get_user_id_by_user_name(username)
	objs = bdd_util.exec_sql("select * from mall_salesman where user_id = %s", [user_id])
	return objs[0]['id']


# def get_corp_id_by_corpuser_name(corpuser_name):
# 	client = RestClient()
# 	data = {
# 		"username": corpuser_name,
# 		"password": '55e421ee9bdc9d9f6b6c6518E590b0ee'
# 	}
# 	resp = client.put('ginger-account:login.logined_corp_user', data)
# 	return resp.data['cid']

STATUS2STR = {
	'wait_pay': u'待支付',
	'wait_ship': u'待发货',
	'finished': u'已完成',
	'shipped': u'已发货',
	'wait_confirm': u'待确认',
	'cancelled': u'已取消',
	'nonsense': 'nonsense'
}

@when(u"{app_user}通过分销员'{salesman_name}'购买'{corpuser_name}'的商品")
def step_impl(context, app_user, salesman_name, corpuser_name):
	input_data = json.loads(context.text)
	input_data['salesman_id'] = get_salesman_id_by_username(context.client, salesman_name)
	step = u"When {}购买'{}'的商品\n\"\"\"\n{}\n\"\"\"\n".format(app_user, corpuser_name, json.dumps(input_data))
	context.execute_steps(step)

@when(u"{app_user}购买'{corpuser_name}'的商品")
def step_impl(context, app_user, corpuser_name):
	input_data = json.loads(context.text)

	product_datas = []
	for product_data in input_data['products']:
		resp = context.client.get("product.corp_products", {
			"corp_id": context.corp['id'],
			"__f-name-contain": product_data['name']
		})
		bdd_util.assert_api_call_success(resp)

		product = resp.data['products'][0]
		sku_name = product_data.get('sku', 'standard')
		if sku_name != 'standard':
			sku_name = get_sku_code_from_display_name(sku_name)
			for sku in product['skus']:
				if sku['name'] == sku_name:
					sku_price = sku['price']
		else:
			sku_price = product['skus'][0]['price']

		product_datas.append({
			'id': product['id'],
			'count': product_data['count'],
			'sku': sku_name,
			'price': product_data.get('price', sku_price)
		})

	# #收货地址
	ship_info = {
		'phone': input_data.get('ship_tel', '13811223344'),
		'address': input_data.get('ship_address', u'103房'),
		'name': input_data.get('ship_name', u'默认姓名'),
		'area_code': ship_info_steps.get_area_code_by_name(context.client, input_data.get('ship_area', '江苏省 南京市 秦淮区'))
	}

	#imoney
	imoney_usages = []
	if 'imoneys' in input_data:
		imoney_usages = input_data['imoneys']
	for imoney_usage in imoney_usages:
		imoney_usage['count'] = int(round(imoney_usage['count']*100, 0))

	#message
	message = input_data.get('message', '')

	# if input_data.get('extra_data'):
	# 	if input_data['extra_data'].get('relevant_user'):
	# 		input_data['extra_data']["relevant_user_id"] = get_member_id_by_username(input_data['extra_data']['relevant_user'])
	# 		del input_data['extra_data']['relevant_user']

	data = {
		'products': json.dumps(product_datas),
		'ship_info': json.dumps(ship_info),
		'message': message,
		'imoney_usages': json.dumps(imoney_usages),
		'biz_code': 'bdd'
	}

	#coupon
	if 'coupon' in input_data:
		data['coupon_usage'] = json.dumps({
			"code": input_data['coupon'],
			"money": 0
		})

	if 'extra_data' in input_data:
		data['extra_data'] = input_data['extra_data']

	if 'order_type' in input_data:
		data['custom_order_type'] = input_data['order_type']

	if 'salesman_id' in input_data:
		data['salesman_id'] = input_data['salesman_id']

	url = 'order.order'
	response = context.client.put(url, data)
	bdd_util.assert_api_call(response, context)
	context.response = response


@Then(u"{webapp_user_name}成功创建订单")
def step_impl(context, webapp_user_name):
	latest_order_bid = get_latest_order_bid()

	response = context.client.get("order.order", {
		"bid": latest_order_bid
	})
	bdd_util.assert_api_call_success(response)

	order_data = response.data
	actual = {
		'bid': order_data['bid'],
		'status': STATUS2STR[order_data['status']],
		'message': order_data['message'],
		'final_money': bdd_util.format_price(order_data['final_money']),
		#'postage': order_data['postage'],
		'delivery_items': [],
		'imoneys': []
	}

	for resource in order_data['resources']:
		if resource['type'] != 'imoney':
			continue

		actual['imoneys'].append({
			'code': resource['code'],
			'count': bdd_util.format_price(resource['count']),
			'deduction_money': bdd_util.format_price(resource['deduction_money'])
		})

	for delivery_item in order_data['invoices']:
		ship_info = delivery_item['ship_info']
		area = ship_info['area']
		delivery_item_data = {
			'status': STATUS2STR[delivery_item['status']],
			'ship_name': ship_info['name'],
			'ship_tel': ship_info['phone'],
			'ship_address': ship_info['address'],
			'ship_area': '%s %s %s' % (area['province']['name'], area['city']['name'], area['district']['name']),
			'final_money': bdd_util.format_price(delivery_item['final_money']),
			'product_price': bdd_util.format_price(delivery_item['product_price']),
			'postage': bdd_util.format_price(delivery_item['postage'])
		}

		# build product data
		products = []
		for product_data in delivery_item['products']:
			products.append({
				'name': product_data['name'],
				'price': bdd_util.format_price(product_data['price']),
				'count': product_data['count'],
				'sku': product_data['sku_display_name']
			})
		delivery_item_data['products'] = products

		# build imoney data
		imoneys = []
		for resource in order_data['resources']:
			if resource['type'] != 'imoney':
				continue

			imoneys.append({
				'code': resource['code'],
				'count': bdd_util.format_price(resource['count']),
				'deduction_money': bdd_util.format_price(resource['deduction_money'])
			})
		delivery_item_data['imoneys'] = imoneys

		actual['delivery_items'].append(delivery_item_data)

	expected = json.loads(context.text)
	bdd_util.assert_dict(expected, actual)

@Then(u"{webapp_user_name}能获得最新订单的订单状态为'{status}'")
def step_impl(context, webapp_user_name, status):
	latest_order_bid = get_latest_order_bid()

	response = context.client.get("order.order_status", {
		"bid": latest_order_bid
	})
	bdd_util.assert_api_call_success(response)

	actual = STATUS2STR[response.data['status']]
	expected = status
	assert actual == expected, u"actual(%s) != expected(%s)" % (actual, expected)


@Then(u"{webapp_user_name}能获得最新出货单的状态为'{status}'")
def step_impl(context, webapp_user_name, status):
	latest_invoice_bid = get_latest_invoice_bid()

	response = context.client.get("order.order_status", {
		"bid": latest_invoice_bid
	})
	bdd_util.assert_api_call_success(response)

	actual = STATUS2STR[response.data['status']]
	expected = status
	assert actual == expected, u"actual(%s) != expected(%s)" % (actual, expected)


@When(u"{webapp_user_name}支付最新订单")
def step_impl(context, webapp_user_name):
	latest_order_bid = get_latest_order_bid()

	response = context.client.put("order.payed_order", {
		"bid": latest_order_bid
	})
	bdd_util.assert_api_call_success(response)


@When(u"{webapp_user_name}在手机上完成最新出货单")
def step_impl(context, webapp_user_name):
	latest_invoice_bid = get_latest_invoice_bid()

	response = context.client.put("order.finished_invoice", {
		"bid": latest_invoice_bid
	})
	bdd_util.assert_api_call_success(response)

	# 发起结算
	response = context.client.put("ginger-finance:clearance.order_clearance", {
		"bid": latest_invoice_bid,
		"sync": True
	})
	bdd_util.assert_api_call_success(response)


@When(u"{corp_user}确认最新发货单")
def step_impl(context, corp_user):
	latest_invoice_bid = get_latest_invoice_bid()

	response = context.client.put("order.confirmed_invoice", {
		"bid": latest_invoice_bid
	})
	bdd_util.assert_api_call_success(response)


@When(u"{corp_user}对最新发货单进行发货")
def step_impl(context, corp_user):
	latest_invoice_bid = get_latest_invoice_bid()
	data = {
		'bid': latest_invoice_bid
	}

	input_data = json.loads(context.text)
	data['enable_logistics'] = input_data.get('enable_logistics', False)
	data['express_company_name'] = input_data.get('express_company_name', '')
	data['express_number'] = input_data.get('express_number', '')
	data['shipper'] = input_data.get('shipper', u'默认发货人')


	response = context.client.put("order.shipped_invoice", {
		'ship_infos': json.dumps([data])
	})
	bdd_util.assert_api_call_success(response)
