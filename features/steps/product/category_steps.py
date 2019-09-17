# -*- coding: utf-8 -*-
import json

from behave import *

from features.bdd import util as bdd_util

def get_product_category_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_category where name = %s", [name])
	return objs[0]['id']

@Then(u"{user}能看到'{category_name}'分类下的商品分类列表")
def step_impl(context, user, category_name):
	expected = json.loads(context.text)
	father_id = get_product_category_id_by_name(category_name) if category_name != u'一级' else 0
	resp = context.client.get("ginger-mall:product.sub_categories", {
		"father_id": father_id
	})
	actual = [category for category in resp.data["categories"] if category['is_enabled']]

	bdd_util.assert_api_call_success(resp)
	bdd_util.assert_list(expected, actual)

@Then(u"{user}能在手机上看到'{category_name}'分类下的商品分类列表")
def step_impl(context, user, category_name):
	expected = json.loads(context.text)
	father_id = get_product_category_id_by_name(category_name) if category_name != u'一级' else 0
	resp = context.client.get("ginger-mall:mall.sub_categories", {
		"father_id": father_id
	})
	actual = [category for category in resp.data["categories"] if category['is_enabled']]

	bdd_util.assert_api_call_success(resp)
	bdd_util.assert_list(expected, actual)

@When(u"{user}创建商品分类")
def step_impl(context, user):
	def create_product_category(category, parent_category_name=None):
		for category_name, sub_categories in category.items():
			post_data = {
				"name": category_name,
				"father_id": get_product_category_id_by_name(parent_category_name) if parent_category_name else 0
			}
			resp = context.client.put("ginger-mall:product.category", post_data)
			bdd_util.assert_api_call_success(resp)

			if len(sub_categories) > 0:
				for sub_category in sub_categories:
					create_product_category(sub_category, category_name)

	datas = json.loads(context.text)
	for data in datas:
		create_product_category(data)

@When(u"{user}删除商品分类'{name}'")
def step_impl(context, user, name):
	id = get_product_category_id_by_name(name)
	resp = context.client.delete("ginger-mall:product.category", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}修改商品分类'{name}'的信息")
def step_impl(context, user, name):
	params = json.loads(context.text)
	id = get_product_category_id_by_name(name)
	params['id'] = id

	resp = context.client.post("ginger-mall:product.category", params)
	bdd_util.assert_api_call_success(resp)

@When(u"{user}启用商品分类'{name}'")
def step_impl(context, user, name):
	id = get_product_category_id_by_name(name)
	resp = context.client.delete("ginger-mall:product.disabled_category", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}禁用商品分类'{name}'")
def step_impl(context, user, name):
	id = get_product_category_id_by_name(name)
	resp = context.client.put("ginger-mall:product.disabled_category", {"id": id})
	bdd_util.assert_api_call_success(resp)


@When(u"{user}商品分类排序")
def step_impl(context, user, name):
	id = get_product_category_id_by_name(name)
	action = json.loads(context.text)["action"]
	resp = context.client.post("ginger-mall:product.category_display_index", {"id": id, "action": action})
	bdd_util.assert_api_call_success(resp)
