# -*- coding: utf-8 -*-
import json
import sys

from behave import *

from features.bdd import util as bdd_util
from features.steps import step_util
from features.bdd.client import RestClient
from features.steps.product import product_label_steps

def get_product_category_id_by_name(name):
	if name == '':
		return 0

	objs = bdd_util.exec_sql("select * from product_category where name = %s", [name])
	print objs, "objs"
	return objs[0]['id']

def get_product_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_product where name = %s", [name])
	return objs[0]['id']

def get_pool_product_id_by_name(name, corpuser_name=None):
	product_model = bdd_util.exec_sql("select * from product_product where name = %s", [name])[0]
	if corpuser_name:
		client = RestClient()
		corp_id = step_util.get_corp_id_for_corpuser(client, corpuser_name)
		pool_product_model = bdd_util.exec_sql("select * from product_pool_product where product_id = %s and corp_id=%s", [product_model['id'], corp_id])[0]
	else:
		pool_product_model = bdd_util.exec_sql("select * from product_pool_product where product_id = %s order by id asc", [product_model['id']])[0]
	return pool_product_model['id']

def get_product_sku_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_sku where name = %s and is_deleted = 0", [name])
	if len(objs) > 0:
		return objs[0]['id']
	else:
		return 0

def get_sku_name_from_display_name(sku_name):
	if sku_name == 'standard':
		return sku_name

	objs = bdd_util.exec_sql("select * from product_sku", [])
	for obj in objs:
		if sku_name in obj['code']:
			return obj['name']

def get_product_properties():
	objs = bdd_util.exec_sql("select * from product_sku_property where is_deleted = 0", [])
	return objs

def get_product_property_values_for(property_ids):
	property_ids_str = ','.join(map(lambda x: "{}".format(x), property_ids))
	objs = bdd_util.exec_sql("select * from product_sku_property_value where is_deleted = 0 and property_id in ({})".format(property_ids_str), [])
	return objs

def get_product_label_id_by_name(name):
	objs = bdd_util.exec_sql("select * from product_label where name = %s", [name])
	return objs[0]['id']

def __get_boolean(product, field, default=False):
	return 'true' if product.get(field, default) == True else 'false'

def __parse_sku_display_name(context, sku_display_name):
	"""
	解析sku display name
		1. 如果都是已存在的property value（黑色 M），生成
			1. 标准sku_name: 2:5_3:6
			2. model properties: [{property_id:2, property_value_id:5}, {property_id:3, property_value_id:6}]

		2. 如果是未存在的proeprty value (颜色:蓝色 尺寸:XS), 生成
			1. sku_name: -1:-1_-1:-1
			2. model properties: [{property_id:-1, property_value_id:-1}, ...]

		3. 对于混合的property value (黑色 尺寸:XL), 生成
			1. sku_name: 2:5_-1:-1
			2. model properties: [{property_id:2, property_value_id:5}, {property_id:-1, property_value_id:-1}]
	"""
	if sku_display_name == 'standard':
		return 'standard', []

	#获取properties
	properties = get_product_properties()
	name2property = {}
	for property in properties:
		name2property[property['name']] = property

	#获取property values
	property_ids = [property['id'] for property in properties]
	property_values = get_product_property_values_for(property_ids)
	name2value = {}
	for property_value in property_values:
		name = property_value['text']
		name2value[name] = property_value

	#从显示用的model display name(黑色 M)构造标准sku_name(2:5_3:6)
	normalized_sku_name_items = []
	#从(黑色 M)获得model properties
	properties = []
	for property_value_name in sku_display_name.split(' '):
		#property_value_name: '黒' 或是 'M'
		property_value = name2value.get(property_value_name)
		property_id = -1
		property_value_id = -1
		property_text = ''
		property_value_text = ''
		if property_value:
			property_id = property_value['property_id']
			property_value_id = property_value['id']
		else:
			items = property_value_name.split(':')
			property_text = items[0]
			property_value_text = items[1]
			if property_text in name2property:
				property_id = name2property[property_text]['id']
		normalized_sku_name_item = '%d:%d' % (property_id, property_value_id)
		normalized_sku_name_items.append(normalized_sku_name_item)
		properties.append({
			'property_id': property_id,
			'property_text': property_text,
			'property_value_id': property_value_id,
			'property_value_text': property_value_text
		})
	normalized_sku_name_items.sort(lambda x,y: cmp(int(x.split(':')[0]), int(y.split(':')[0])))
	normalized_sku_name = '_'.join(normalized_sku_name_items)

	return normalized_sku_name, properties


def __format_product_sku_info(context, product):
	#规格信息
	def should_create_property_and_value(sku_display_name):
		return ':' in sku_display_name

	skus_info = []
	if 'skus' in product:
		for sku_display_name, sku_info in product['skus'].items():
			data = {
				"price": int(round(sku_info.get('price', 1.0) * 100, 0)),
				"cost_price": int(round(sku_info.get('cost_price', 1.0) * 100, 0)),
				"stocks": sku_info.get('stocks', 10),
				"code": sku_info.get('code', 'code_%s' % sku_display_name)
			}

			normalized_sku_name, sku_properties = __parse_sku_display_name(context, sku_display_name)
			data['name'] = normalized_sku_name
			data['properties'] = sku_properties

			#在数据库中查询是否已存在该model，如果存在，则是更新；否则，为创建
			existed_product_sku_id = get_product_sku_id_by_name(normalized_sku_name)
			if existed_product_sku_id == 0:
				#不存在
				data['id'] = -1
			else:
				data['id'] = existed_product_sku_id

			skus_info.append(data)
	else:
		standard_sku = {
			"name": "standard",
			"price": int(round(product.get('price', 1.0) * 100, 0)),
			"cost_price": int(round(product.get('cost_price', 1.0) * 100, 0)),
			"stocks": 9999,
			"code": 'code',
			'properties': []
		}
		skus_info = [standard_sku]

	return skus_info

def __format_product_post_data(context, product):
	"""
	构造用于提交的product数据
	"""
	base_info = {
		'name': product['name'], #商品名
		'type': product.get('type', 'product'),
		'code': product.get('code', ''), #商品编码
		'category_id': get_product_category_id_by_name(product.get('category', '')),
		'promotion_title': product.get('promotion_title', ''), #促销标题
		'detail': product.get('detail', u'商品的详情') #详情
	}

	#规格信息
	skus_info = __format_product_sku_info(context, product)

	#图片信息
	media_info = {
		'images': product.get('images', [{'url': '/static/test_resource_img/default_product_img.jpg'}]),
		'thumbnail': ''
	}
	if len(media_info['images']) > 0:
		media_info['thumbnail'] = media_info['images'][0]['url']

	#虚拟资产信息
	imoney_codes = product.get('imoney_codes', [])

	#分组信息
	# categories = []
	# category_names = product.get('categories', None)
	# if category_names:
	# 	for category_name in category_names:
	# 		category_id = get_product_category_id_by_name(category_name)
	# 		categories.append(category_id)

	#物流信息
	postage = product.get('postage')
	if not postage:
		logistics_info = {
			'postage_type': 'unified',
			'unified_postage_money': 0
		}
	else:
		if type(postage) == str or type(postage) == unicode:
			logistics_info = {
				'postage_type': 'unified',
				'unified_postage_money': 0
			}
		else:
			logistics_info = {
				'postage_type': 'unified',
				'unified_postage_money': int(round(postage * 100.0, 0))
			}

	data = {
		'base_info': json.dumps(base_info),
		'skus_info': json.dumps(skus_info),
		'media_info': json.dumps(media_info),
		'imoney_codes': json.dumps(imoney_codes),
		'logistics_info': json.dumps(logistics_info)
	}

	return data

def __create_product(context, product):
	data = __format_product_post_data(context, product)

	response = context.client.put('product.product', data)
	bdd_util.assert_api_call_success(response)

def __get_product(context, corpuser_name, name):
	pool_product_id = get_pool_product_id_by_name(name, corpuser_name)
	#product_model = bdd_util.exec_sql("select * from product_product where name = %s", [name])[0]
	#pool_product_model = bdd_util.exec_sql("select * from product_pool_product where product_id = %s", [product_model['id']])[0]
	data = {
		"id": pool_product_id
	}

	response = context.client.get('product.product', data)
	bdd_util.assert_api_call_success(response)

	resp_data = response.data
	base_info = resp_data['base_info']
	product = {
		"id": resp_data['id'],
		"name": base_info['name'],
		"promotion_title": base_info['promotion_title'],
		"detail": base_info['detail'],
		"thumbnail": base_info['thumbnail'],
		"shelve_type": u"上架" if base_info['shelve_type'] == 'on_shelf' else u"下架"
	}

	#处理category
	category = resp_data['category']
	product['category'] = category['name'] if category else ''

	#处理媒体信息
	product['medias'] = resp_data['medias']

	#处理标签
	product['labels'] = [label['name'] for label in resp_data['labels']]

	#处理规格信息
	skus = resp_data['skus']
	name2sku = {}
	for sku in skus:
		sku['price'] = bdd_util.format_price(sku['price'])
		sku['cost_price'] = bdd_util.format_price(sku['cost_price'])
		name = sku['name']
		if name != 'standard':
			name = ' '.join([property_value['text'] for property_value in sku['property_values']])
		name2sku[name] = sku
	product['skus'] = name2sku
	return product


def __get_format_products(products):
	datas = []
	for product in products:
		data = {}
		data['name'] = product['base_info']['name']
		data['type'] = product['type']

		#分类
		if 'categories' in product:
			data['categories'] = ','.join([category['name'] for category in product['categories']])
		else:
			data['categories'] = ''

		#处理规格信息
		skus = product['skus']
		name2sku = {}
		for sku in skus:
			sku['price'] = bdd_util.format_price(sku['price'])
			sku['cost_price'] = bdd_util.format_price(sku['cost_price'])
			name = sku['name']
			if name != 'standard':
				name = ' '.join([property_value['text'] for property_value in sku['property_values']])
			name2sku[name] = sku
		data['skus'] = name2sku

		data['thumbnail'] = product['base_info']['thumbnail']
		data['status'] = product['status']

		if data['status'] == 'on_pool':
			data['status'] = u'未入库'
		elif data['status'] == 'off_shelf':
			data['status'] = u'已入库'
		elif data['status'] == 'on_shelf':
			data['status'] = u'已上架'

		datas.append(data)
	return datas


def __get_products(context, corp_name, type_name=u'在售', options={}):
	TYPE2URL = {
		u'待售': 'product.offshelf_products',
		u'在售': 'product.onshelf_products',
		u'分类': 'product.category_products',
		u'手机上分类': 'mall.category_products',
		u'平台': 'product.platform_pool_products',
		u'标签': 'product.labeled_products'
	}

	# if type_name == u'平台':
	# 	#print json.dumps(response.data, indent=2)
	# 	print products
	# 	raw_input()

	url = TYPE2URL[type_name]

	params = {}
	if type_name == u'分类' or type_name == u'手机上分类':
		params['category_id'] = options['category_id']
	if type_name == u'标签':
		params['label_id'] = options['label_id']
	response = context.client.get(url, params)
	bdd_util.assert_api_call_success(response)

	products = __get_format_products(response.data["products"])
	return products

@then(u"{user}能获得'{type_name}'商品列表")
def step_impl(context, user, type_name):
	actual = __get_products(context, user, type_name)
	#context.products = actual

	if context.table:
		expected = []
		for product in context.table:
			product = dict(product.as_dict())
			if 'barCode' in product:
				product['bar_code'] = product['barCode']
				del product['barCode']

			if 'categories' in product:
				product['categories'] = product['categories'].split(',')
				# 处理空字符串分割问题
				if product['categories'][0] == '':
					product['categories'] = []
			# 处理table中没有验证库存的行
			if 'stocks' in product and product['stocks'] == '':
				del product['stocks']
			# 处理table中没有验证条码的行
			if 'bar_code' in product and product['bar_code'] == '':
				del product['bar_code']

			expected.append(product)
	else:
		print 'load expected from context.text'
		expected = json.loads(context.text)

	bdd_util.assert_list(expected, actual)

@then(u"{user}能获得商品分类'{category_name}'下的商品列表")
def step_impl(context, user, category_name):
	category_id = get_product_category_id_by_name(category_name)
	actual = __get_products(context, user, u'分类', {
		'category_id': category_id
	})

	expected = json.loads(context.text)

	bdd_util.assert_list(expected, actual)

@then(u"{user}能获得商城中商品分类'{category_name}'下的商品列表")
def step_impl(context, user, category_name):
	category_id = get_product_category_id_by_name(category_name)
	actual = __get_products(context, user, u'手机上分类', {
		'category_id': category_id
	})

	expected = json.loads(context.text)

	bdd_util.assert_list(expected, actual)

@then(u"{user}能获得商品标签'{label_name}'下的商品列表")
def step_impl(context, user, label_name):
	label_id = product_label_steps.get_product_label_id_by_name(label_name)
	actual = __get_products(context, user, u'标签', {
		'label_id': label_id
	})

	expected = json.loads(context.text)

	bdd_util.assert_list(expected, actual)

# @When(u"{user}创建商品")
# def step_impl(context, user):
# 	products = json.loads(context.text)
# 	if isinstance(products, dict):
# 		products = [products]
# 	for product in products:
# 		#product['auto_on_shelf'] = True
# 		__create_product(context, product)
#
# 		#step = u'When {}将商品移动到\'在售\'货架\n"""\n["{}"]\n"""'.format(user, product['name'])
# 		#context.execute_steps(step)

@When(u"{user}添加商品")
def step_impl(context, user):
	products = json.loads(context.text)
	if isinstance(products, dict):
		products = [products]
	for product in products:
		__create_product(context, product)

		step = u'When {}将商品移动到\'在售\'货架\n"""\n["{}"]\n"""'.format(user, product['name'])
		context.execute_steps(step)

@When(u"{user}添加待审核商品")
def step_impl(context, user):
	products = json.loads(context.text)
	if isinstance(products, dict):
		products = [products]
	for product in products:
		__create_product(context, product)

@When(u"{user}添加虚拟商品")
def step_impl(context, user):
	products = json.loads(context.text)
	if isinstance(products, dict):
		products = [products]
	for product in products:
		data = {
			'name': product['name'],
			'skus': json.dumps([{
				'name': 'standard',
				'price': product['price']
			}])
		}
		response = context.client.put('product.virtual_product', data)
		bdd_util.assert_api_call_success(response)

@then(u"{user}能获取商品'{product_name}'")
def step_impl(context, user, product_name):
	expected = json.loads(context.text)

	actual = __get_product(context, user, product_name)
	bdd_util.assert_dict(expected, actual)

@When(u"{user}删除商品'{name}'")
def step_impl(context, user, name):
	id = get_product_id_by_name(name)
	resp = context.client.delete("product.product", {"id": id})
	bdd_util.assert_api_call_success(resp)

@when(u"{user}更新商品'{product_name}'")
def step_impl(context, user, product_name):
	product = __get_product(context, user, product_name)

	update_data = json.loads(context.text)
	for key, value in update_data.items():
		product[key] = value
	data = __format_product_post_data(context, product)
	data['id'] = product['id']

	response = context.client.post('product.product', data)
	bdd_util.assert_api_call_success(response)

@when(u"{user}为商品'{product_name}'增加标签")
def step_impl(context, user, product_name):
	product_id = get_pool_product_id_by_name(product_name)

	input_datas = json.loads(context.text)
	label_ids = []
	for label_name in input_datas:
		label_ids.append(get_product_label_id_by_name(label_name))

	data = {
		"product_id": product_id,
		"label_ids": json.dumps(label_ids)
	}
	response = context.client.put('product.labeled_product', data)
	bdd_util.assert_api_call_success(response)

@when(u"{user}为商品'{product_name}'删除标签'{label_name}'")
def step_impl(context, user, product_name, label_name):
	product_id = get_pool_product_id_by_name(product_name)
	label_id = get_product_label_id_by_name(label_name)

	data = {
		"product_id": product_id,
		"label_id": label_id
	}
	response = context.client.delete('product.labeled_product', data)
	bdd_util.assert_api_call_success(response)

@When(u"{user}启用商品'{name}'")
def step_impl(context, user, name):
	id = get_product_id_by_name(name)
	resp = context.client.delete("product.disabled_product", {"id": id})
	bdd_util.assert_api_call_success(resp)

@When(u"{user}禁用商品'{name}'")
def step_impl(context, user, name):
	id = get_product_id_by_name(name)
	resp = context.client.put("product.disabled_product", {"id": id})
	bdd_util.assert_api_call_success(resp)


@When(u"{user}商品排序")
def step_impl(context, user, name):
	id = get_product_id_by_name(name)
	action = json.loads(context.text)["action"]
	resp = context.client.post("product.product_index", {"id": id, "action": action})
	bdd_util.assert_api_call_success(resp)

@when(u"{user}将商品移动到'{shelf_name}'货架")
def step_impl(context, user, shelf_name):
	product_names = json.loads(context.text)
	pool_product_ids = []
	for product_name in product_names:
		pool_product_id = get_pool_product_id_by_name(product_name, user)
		pool_product_ids.append(pool_product_id)

	if shelf_name == u'在售':
		data = {
			'product_ids': json.dumps(pool_product_ids)
		}
		response = context.client.put('product.onshelf_products', data)
		bdd_util.assert_api_call_success(response)
	elif shelf_name == u'待售':
		data = {
			'product_ids': json.dumps(pool_product_ids)
		}
		response = context.client.put('product.offshelf_products', data)
		bdd_util.assert_api_call_success(response)


@When(u"{corp_user}修改商品的价格")
def step_impl(context, corp_user):
	product = json.loads(context.text)
	pool_product_id = get_pool_product_id_by_name(product['name'], corp_user)
	sku = get_sku_name_from_display_name(product['sku'])

	price_infos = None
	if 'price' in product:
		price_infos = [{
			'sku': sku,
			'price': product['price']
		}]

	supplier_price_infos = None
	if 'supplier_price' in product:
		supplier_price_infos = [{
			'sku': sku,
			'price': product['supplier_price']
		}]

	data = {
		"id": pool_product_id,
	}
	if price_infos:
		data['price_infos'] = json.dumps(price_infos)
	if supplier_price_infos:
		data['supplier_price_infos'] = json.dumps(supplier_price_infos)

	response = context.client.post('product.pool_product_price', data)
	bdd_util.assert_api_call_success(response)



@then(u"{user}能在手机上获得商品'{product_name}'")
def step_impl(context, user, product_name):
	expected = json.loads(context.text)

	actual = __get_product(context, context.corpuser_name, product_name)
	bdd_util.assert_dict(expected, actual)
