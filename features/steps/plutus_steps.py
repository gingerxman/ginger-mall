# -*- coding:utf-8 -*-
#pylint: disable=E0602,E0102

import json

from behave import *
from features.bdd import util as bdd_util
from features.bdd import client as bdd_client
from features.bdd.client import RestClient
from features.steps import step_util

@given(u"系统配置虚拟资产")
def step_impl(context):
	rest_client = bdd_client.login('backend', "ginger", password=None, context=context)

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

		resp = rest_client.put('ginger-finance:imoney.imoney', data)
		bdd_util.assert_api_call_success(resp)

@when(u"{user}充值'{amount}'个'{imoney_code}'")
def step_impl(context, user, amount, imoney_code):
	user_id = step_util.get_user_id_by_name(context.client, user)

	data = {
		"source_user_id": 0,#client.cur_user_id,
		"dest_user_id": user_id,
		"imoney_code": imoney_code,
		"amount": int(amount) * 100,
		"bid": "bdd"
	}
	resp = context.client.put('ginger-finance:imoney.transfer', data)
	bdd_util.assert_api_call_success(resp)


@Then(u"{user_name}能获得虚拟资产'{imoney_code}'")
def step_impl(context, user_name, imoney_code):
	data = {
		"imoney_code": imoney_code
	}
	resp = context.client.get('ginger-finance:imoney.balance', data)
	bdd_util.assert_api_call_success(resp)

	actual = bdd_util.format_price(resp.data)

	expected = json.loads(context.text)['balance']
	assert expected == actual, "e(%d) != a(%d)" % (expected, actual)

@Then(u"{user_name}能获得公司的虚拟资产'{imoney_code}'")
def step_impl(context, user_name, imoney_code):
	data = {
		"imoney_code": imoney_code,
		'view_corp_account': True
	}
	resp = context.client.get('ginger-finance:imoney.balance', data)
	bdd_util.assert_api_call_success(resp)

	actual = bdd_util.format_price(resp.data)

	expected = json.loads(context.text)['balance']
	assert expected == actual, "e(%d) != a(%d)" % (expected, actual)


