# -*- coding: utf-8 -*-
import json
import time

from behave import *

from features.bdd import util as bdd_util
from features.bdd.client import RestClient
from features.steps.product import product_steps

def get_ship_info_id_by_address(address):
	objs = bdd_util.exec_sql("select * from mall_ship_info where address = %s", [address])
	return objs[0]['id']

def get_pool_product_id_by_name(name, corp_id=None):
	product_model = bdd_util.exec_sql("select * from product_product where name = %s", [name])[0]
	if corp_id:
		pool_product_model = bdd_util.exec_sql("select * from product_pool_product where product_id = %s and corp_id=%s", [product_model['id'], corp_id])[0]
	else:
		pool_product_model = bdd_util.exec_sql("select * from product_pool_product where product_id = %s order by id asc", [product_model['id']])[0]
	return pool_product_model['id']

def get_shopping_cart_item_id_by_product_name(product_name):
	pool_product_id = get_pool_product_id_by_name(product_name)
	shopping_cart_item = bdd_util.exec_sql("select * from mall_shopping_cart where product_id = %s", [pool_product_id])[0]
	return shopping_cart_item['id']

@when(u"{user}加入{corp_name}的商品到购物车")
def step_impl(context, user, corp_name):
	input_datas = json.loads(context.text)
	for product_data in input_datas:
		pool_product_id = get_pool_product_id_by_name(product_data['name'], context.corp['id'])
		sku_name = product_steps.get_sku_name_from_display_name(product_data.get('sku', 'standard'))

		data = {
			'pool_product_id': pool_product_id,
			'count': product_data['count'],
			'sku_name': sku_name
		}
		response = context.client.put('mall.shopping_cart_item', data)
		bdd_util.assert_api_call_success(response)
		time.sleep(0.1)

@then(u"{user}能获得购物车")
def step_impl(context, user):
	response = context.client.get("mall.shopping_cart", {})
	bdd_util.assert_api_call_success(response)

	actual = response.data
	for product_group in actual['product_groups']:
		#product_group['supplier'] = product_group['supplier']['name']
		for product in product_group['products']:
			product['count'] = product['purchase_count']
			product['sku'] = product['sku_display_name']
			product['price'] = bdd_util.format_price(product['price'])

	for product in actual['invalid_products']:
		product['count'] = product['purchase_count']
		product['sku'] = product['sku_display_name']

	expected = json.loads(context.text)
	bdd_util.assert_dict(expected, actual)

@when(u"{user}从购物车中删除商品")
def step_impl(context, user):
	input_datas = json.loads(context.text)
	for product_name in input_datas:
		shopping_cart_item_id = get_shopping_cart_item_id_by_product_name(product_name)

		data = {
			'id': shopping_cart_item_id
		}
		response = context.client.delete('mall.shopping_cart_item', data)
		bdd_util.assert_api_call_success(response)

@then(u"{user}能获得购物车中商品数量为'{count}'")
def step_impl(context, user, count):
	response = context.client.get("mall.shopping_cart_product_count", {})
	bdd_util.assert_api_call_success(response)

	actual = response.data['count']

	expected = int(count)
	assert expected == actual, 'expected(%d) != actual(%d)' % (expected, actual)
