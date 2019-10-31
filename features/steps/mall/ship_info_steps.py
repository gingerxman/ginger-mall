# -*- coding: utf-8 -*-
import json

from behave import *

from features.bdd import util as bdd_util
from features.bdd.client import RestClient

def get_ship_info_id_by_address(address):
	objs = bdd_util.exec_sql("select * from mall_ship_info where address = %s", [address])
	return objs[0]['id']

def get_area_code_by_name(client, name):
	resp = client.get('area.area_code', {
		'name': name
	})
	bdd_util.assert_api_call_success(resp)

	return resp.data

@Then(u"{user}能看到收货地址列表")
def step_impl(context, user):
	expected = json.loads(context.text)
	resp = context.client.get("mall.ship_infos")
	actual = resp.data["ship_infos"]

	bdd_util.assert_api_call_success(resp)
	bdd_util.assert_list(expected, actual)

@When(u"{user}创建收货地址")
def step_impl(context, user):
	input_datas = json.loads(context.text)
	for input_data in input_datas:
		area_name = input_data.get("area")
		if area_name:
			area_code = get_area_code_by_name(context.client, area_name)
		else:
			area_code = '320104' #江苏省 南京市 秦淮区

		data = {
			"name": input_data.get('name', user),
			"phone": input_data.get('phone', '13811223344'),
			"area_code": area_code,
			"address": input_data.get('address', '国创园')
		}
		resp = context.client.put("mall.ship_info", data)
		bdd_util.assert_api_call_success(resp)

@When(u"{user}删除收货地址'{address}'")
def step_impl(context, user, address):
	id = get_ship_info_id_by_address(address)
	resp = context.client.delete("mall.ship_info", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}修改收货地址'{address}'的信息")
def step_impl(context, user, address):
	params = json.loads(context.text)
	id = get_ship_info_id_by_address(address)
	params['id'] = id

	#处理area
	area_name = params["area"]
	area_code = get_area_code_by_name(context.client, area_name)
	params['area_code'] = area_code

	resp = context.client.post("mall.ship_info", params)
	bdd_util.assert_api_call_success(resp)

@When(u"{user}修改收货地址'{address}'的排序")
def step_impl(context, user, address):
	params = json.loads(context.text)
	id = get_ship_info_id_by_address(address)
	params['id'] = id
	params['action'] = json.loads(context.text)['action']

	resp = context.client.post("mall.ship_info_display_index", params)
	bdd_util.assert_api_call_success(resp)

@When(u"{user}启用收货地址'{address}'")
def step_impl(context, user, address):
	id = get_ship_info_id_by_address(address)
	resp = context.client.delete("mall.disabled_ship_info", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}禁用收货地址'{address}'")
def step_impl(context, user, address):
	id = get_ship_info_id_by_address(address)
	resp = context.client.put("mall.disabled_ship_info", {"id": id})
	bdd_util.assert_api_call_success(resp)

@Then(u"{user}能获得默认收货地址")
def step_impl(context, user):
	resp = context.client.get("mall.default_ship_info")
	bdd_util.assert_api_call_success(resp)
	actual = resp.data

	expected = json.loads(context.text)
	bdd_util.assert_dict(expected, actual)

@When(u"{user}设置收货地址'{address}'为默认收货地址")
def step_impl(context, user, address):
	id = get_ship_info_id_by_address(address)
	resp = context.client.put("mall.default_ship_info", {"id": id})
	bdd_util.assert_api_call_success(resp)
