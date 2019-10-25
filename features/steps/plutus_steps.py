# -*- coding:utf-8 -*-
#pylint: disable=E0602,E0102

import json

from behave import *
from features.bdd import util as bdd_util
from features.bdd import client as bdd_client
from features.bdd.client import RestClient

def get_user_id_by_user_name(name):
	remote_client = RestClient()

	resp = remote_client.put('gskep:login.logined_bdd_user', {
		'name': name
	})
	bdd_util.assert_api_call_success(resp)

	return resp.data['id']

def get_user_id_by_corpuser_id(client, corpuser_id):
	resp = client.get('gskep:account.user', {
		'corp_user_id': corpuser_id
	})
	bdd_util.assert_api_call_success(resp)

	return resp.data['id']

def get_platform_user_id(client):
	resp = client.get('gskep:corp.platform_corps', {
	})
	bdd_util.assert_api_call_success(resp)

	platform_corpuser_id = resp.data['corps'][0]['corp_user']['id']
	return get_user_id_by_corpuser_id(client, platform_corpuser_id)

@given(u"系统配置虚拟资产")
def step_impl(context):
	#from features.bdd.client import RestClient
	#rest_client = RestClient()
	rest_client = bdd_client.login('backend', "xiaocheng", password=None, context=context)

	input_datas = json.loads(context.text)

	for input_data in input_datas:
		data = {
			'code': input_data['code'],
			'display_name': input_data.get('display_name', input_data['code']),
			'exchange_rate': input_data.get('exchange_rate', 1.0),
			'enable_fraction': input_data.get('enable_fraction', False),
			'is_payable': input_data.get('is_payable', True),
			'is_debtable': input_data.get('is_debtable', False),
		}

		resp = rest_client.put('gplutus:imoney.imoney', data)
		bdd_util.assert_api_call_success(resp)


@given(u"系统为'{user_name}'转账'{amount}'个'{imoney_code}'")
def step_impl(context, user_name, amount, imoney_code):
	client = bdd_client.login('app', 'platform_user_1', password=None, context=context)

	user_id = get_user_id_by_user_name(user_name)
	data = {
		"source_user_id": 0,#client.cur_user_id,
		"dest_user_id": user_id,
		"imoney_code": imoney_code,
		"amount": amount,
		"bid": "bdd"
	}
	resp = client.put('gplutus:imoney.transfer', data)
	bdd_util.assert_api_call_success(resp)


@Then(u"{user_name}能获得虚拟资产'{imoney_code}'")
def step_impl(context, user_name, imoney_code):
	data = {
		"imoney_code": imoney_code,
		"_v": 2
	}
	resp = context.client.get('gplutus:imoney.balance', data)
	bdd_util.assert_api_call_success(resp)

	actual = resp.data['valid_balance']

	expected = json.loads(context.text)['balance']
	assert expected == actual, "e(%d) != a(%d)" % (expected, actual)

@Then(u"{corpuser_name}能获得酒吧的虚拟资产'{imoney_code}'")
def step_impl(context, corpuser_name, imoney_code):
	user_id = get_user_id_by_corpuser_id(context.client, context.client.cur_user_id)
	data = {
		"imoney_code": imoney_code,
		"_v": 2
	}
	resp = context.client.get('gplutus:imoney.balance', data)
	bdd_util.assert_api_call_success(resp)

	actual = resp.data['valid_balance']

	expected = json.loads(context.text)['balance']
	assert expected == actual, "e(%d) != a(%d)" % (expected, actual)

@Then(u"{corpuser_name}能获得平台的虚拟资产'{imoney_code}'")
def step_impl(context, corpuser_name, imoney_code):
	platform_user_id = get_platform_user_id(context.client)
	data = {
		"imoney_code": imoney_code,
		"user_ids": json.dumps([platform_user_id]),
		"_v": 2
	}
	resp = context.client.get('gplutus:imoney.users_balance', data)
	bdd_util.assert_api_call_success(resp)

	actual = resp.data[0]['valid_balance']

	expected = json.loads(context.text)['balance']
	assert expected == actual, "e(%d) != a(%d)" % (expected, actual)

