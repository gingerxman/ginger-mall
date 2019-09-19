# -*- coding: utf-8 -*-
import json

from behave import *

from features.bdd import util as bdd_util

def get_product_property_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_sku_property where name = %s", [name])
	return objs[0]['id']

def get_product_property_value_ids():
	objs = bdd_util.exec_sql("select * from product_sku_property_value", [])
	return [obj['id'] for obj in objs]

def get_product_property_value_by_value(value):
	objs = bdd_util.exec_sql("select * from product_sku_property_value where Text = %s", [value])
	return objs[0]['id'], objs[0]['property_id']


def __add_product_property(context, product_property):
	# 创建product model property
	data = {
		"name": product_property['name']
	}
	response = context.client.put('product.property', data)
	bdd_util.assert_api_call_success(response)

	product_property_id = response.data['id']

	#处理value
	for value in product_property['values']:
		value['text'] = value['name']
		if 'image' in value:
			value['image'] = value['image']
		else:
			value['image'] = ''
		value['property_id'] = product_property_id
		response = context.client.put('product.property_value', value)
		bdd_util.assert_api_call_success(response)

@Then(u"{user}能看到商品属性列表")
def step_impl(context, user):
	expected = json.loads(context.text)
	resp = context.client.get("product.corp_product_properties")
	bdd_util.assert_api_call_success(resp)

	actual = resp.data["product_properties"]
	for property_item in actual:
		for value_item in property_item['values']:
			value_item['image'] = value_item['image']
			value_item['name'] = value_item['text']

	bdd_util.assert_list(expected, actual)

@When(u"{user}创建商品属性")
def step_impl(context, user):
	product_properties = json.loads(context.text)
	for property in product_properties:
		__add_product_property(context, property)

@When(u"{user}删除商品规格属性'{name}'")
def step_impl(context, user, name):
	expected = json.loads(context.text)
	id = get_product_property_id_by_name(name)
	resp = context.client.delete("product.property", {"id": id})
	if not expected.get("error_code"):
		bdd_util.assert_api_call_success(resp)
	else:
		bdd_util.assert_api_call_failed(resp, expected.get("error_code"))

@When(u"{user}删除商品规格属性值'{value}'")
def step_impl(context, user, value):
	expected = json.loads(context.text)
	id, property_id = get_product_property_value_by_value(value)
	resp = context.client.delete("product.property_value", {"id": id, "property_id": property_id})
	if not expected.get("error_code"):
		bdd_util.assert_api_call_success(resp)
	else:
		bdd_util.assert_api_call_failed(resp, expected.get("error_code"))

@When(u"{user}修改商品属性'{name}'的信息")
def step_impl(context, user, name):
	params = json.loads(context.text)
	id = get_product_property_id_by_name(name)
	params['id'] = id

	values = params.get('values')
	if values:
		del params['values']

	resp = context.client.post("product.property", params)
	bdd_util.assert_api_call_success(resp)

	#更新property value
	if values:
		existed_value_ids = get_product_property_value_ids()
		for existed_value_id in existed_value_ids:
			response = context.client.delete('product.property_value', {
				"property_id": id,
				"id": existed_value_id
			})
			bdd_util.assert_api_call_success(response)

		for value in values:
			value['text'] = value['name']
			if 'image' in value:
				value['image'] = value['image']
			else:
				value['image'] = ''
			value['property_id'] = id
			response = context.client.put('product.property_value', value)
			bdd_util.assert_api_call_success(response)


@When(u"{user}启用商品属性'{name}'")
def step_impl(context, user, name):
	id = get_product_property_id_by_name(name)
	resp = context.client.delete("product.disabled_property", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}禁用商品属性'{name}'")
def step_impl(context, user, name):
	id = get_product_property_id_by_name(name)
	resp = context.client.put("product.disabled_property", {"id": id})
	bdd_util.assert_api_call_success(resp)

