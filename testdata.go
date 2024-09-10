package mapinterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/runingriver/mapinterface/pkg"
	"strings"
)

var (
	listJson = `
	[
		{
			"key1": "val1"
		},
		{
			"100": [
				1,
				2,
				3
			]
		},
		[
			{
				"name": "zhangsan"
			},
			{
				"age": 15
			}
		]
	]`
	jsonStrList = []string{
		// 0---------------------------------------------------------
		`{"name":{"first":"Janet","last":"Prichard"},"age":47}`,
		// 1---------------------------------------------------------
		`
		{
			"a": {
				"b": {
					"c": {
						"d": {
							"e": 1
						}
					}
				}
			}
		}
		`,
		// 2---------------------------------------------------------
		`{
			"users": [
				{
					"id": 1,
					"name": {
						"first": "John",
						"last": "Ramboo"
					}
				},
				{
					"id": 2,
					"name": {
						"first": "Ethan",
						"last": "Hunt"
					}
				},
				{
					"id": 3,
					"name": {
						"first": "John",
						"last": "Doe"
					}
				}
			]
		}`,
		// 3---------------------------------------------------------
		`
		{
			"name": "computers",
			"description": "List of computer products",
			"vendor": {
				"name": "Star Trek",
				"email": "info@example.com",
				"website": "www.example.com",
				"items": [
					{
						"id": 1,
						"name": "MacBook Pro 13 inch retina",
						"price": 1350
					},
					{
						"id": 2,
						"name": "MacBook Pro 15 inch retina",
						"price": 1700
					},
					{
						"id": 3,
						"name": "Sony VAIO",
						"price": 1200
					},
					{
						"id": null,
						"name": "HP core i3 SSD",
						"price": 850
					}
				],
				"prices": [
					2400,
					2100,
					1200,
					400.87,
					89.9,
					150.1
				],
				"names": [
					"John Doe",
					"Jane Doe",
					"Tom",
					"Jerry",
					"Nicolas",
					"Abby"
				]
			}
		}
		`,
		// 4---------------------------------------------------------
		`
		{
			"name": {
				"1001": "Janet",
				"1002": [
					{
						"math": 98
					},
					{
						"geography": 88
					}
				],
				"type": "exam"
			},
			"age": 47
		}
		`,
	}
	itfObj = []interface{}{
		map[int64]interface{}{
			10010: map[string]interface{}{"name": "Janet", "age": 20},
			10011: map[string]interface{}{"name": "Tom", "age": 21},
			10012: map[string]interface{}{"name": "Kate", "age": 22},
		},
		[]string{"Janet", "Tom", "Kate"},
		[]uint64{10011, 10012, 10013},
		[]struct{ A int }{{A: 1}, {A: 2}},
		map[string]interface{}{
			"num": map[int]interface{}{
				1001: []map[string]interface{}{
					{"math": 66},
					{"geography": 70},
				},
				1002: []map[string]interface{}{
					{"math": 98},
					{"geography": 88},
				},
			},
		},
	}
	mapJsonStr = map[string]interface{}{
		"student": `{
			"name": "Jack",
			"age": 23
		}`,
		"int-map-itf": map[int]interface{}{
			10020: `{
				"name": "Jack",
				"age": 23
			}`,
		},
		"str-str": map[string]string{
			"info": "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}",
		},
		"str-itf-str": map[string]interface{}{
			"info": "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}",
		},
	}
	getAndCase = `
	{
		"name": "computers",
		"description": "List of computer products",
		"vendor": {
			"name": "Star Trek",
			"email": "info@example.com",
			"info": {
				"age": 25,
				"relation": {
					"zhangsan": "friend"
				}
			},
			"items": [
				{
					"id": 1,
					"name": "MacBook Pro 13 inch retina",
					"price": 1350
				},
				{
					"id": 2,
					"name": "MacBook Pro 15 inch retina",
					"price": 1700
				}
			],
			"prices": [
				2400,
				400.87,
				89.9,
				150.2
			],
			"names": [
				"John Doe",
				"Jane Doe",
				"Tom"
			]
		}
	}
	`
	intObjPrt  = MockInt(1)
	intObjPrt2 = MockInt(1)
	intPrt     = 1
	intPrt2    = 1
	UniqList   = map[string]interface{}{
		"UniqForInt":    []int{1, 1, 2, 2, 3},
		"UniqForIntPtr": []*int{&intPrt, &intPrt2, &intPrt, &intPrt2},
		"UniqForStrItf": []interface{}{"Jak", "Tom", "Kav", "Jak"},
		"UniqForStr":    []string{"Jak", "Tom", "Kav", "Jak"},
		"UniqForObj":    []MockInt{MockInt(1), MockInt(1), MockInt(2), MockInt(2), MockInt(3)},
		"UniqForObjPtr": []*MockInt{&intObjPrt, &intObjPrt2, &intObjPrt, &intObjPrt2},
	}
	CvtMap = map[string]interface{}{
		"StrToStrForStrItf": map[string]interface{}{
			"Kav": "101",
			"Tom": "102",
			"Jak": "103",
		},
		"StrToStrForItfItf": map[interface{}]interface{}{
			&MockObj{"Kav"}: "101",
			&MockObj{"Tom"}: "102",
			&MockObj{"Jak"}: "103",
		},
		"StrToStrForItfObj": map[string]MockObj{
			"101": {"Kav"},
			"102": {"Tom"},
			"103": {"Jak"},
		},
		"IntToIntForStrItf": map[string]interface{}{
			"101": 1,
			"102": 2,
			"103": 3,
		},
		"IntToIntForItfItf": map[interface{}]interface{}{
			"101": MockInt(1),
			"102": MockInt(2),
			"103": MockInt(3),
		},
		"IntToIntForItfObj": map[interface{}]MockInt{
			"101": MockInt(1),
			"102": MockInt(2),
			"103": MockInt(3),
		},
	}
	CvtList = map[string]interface{}{
		"ToListForRf":     []interface{}{&MockObj{"Kav"}, MockObj{"Tom"}, "Jak"},
		"ToListStrForItf": []interface{}{&MockObj{"Kav"}, MockObj{"Tom"}, "Jak"},
		"ToListStrForRf":  []*MockObj{{"Kav"}, {"Tom"}, {"Jak"}},
		"ToListIntForItf": []interface{}{MockInt(101), MockInt(102), MockInt(103)},
		"ToListIntForRf":  []MockInt{101, 102, 103},
	}

	mapStrToMap = map[int]map[string]string{
		1: {"a": "b"},
		2: {"c": "d"},
	}
	MapAny = map[string]interface{}{
		"map-int-map-str":     mapStrToMap,
		"ptr-map-int-map-str": &mapStrToMap,
		"map-str-ptr": map[string]*struct{ A string }{
			"1": {A: "1"},
			"2": {A: "2"},
		},
		"map-str-map-str": map[string]map[string]string{
			"1": {"a": "b"},
			"2": {"c": "d"},
		},
	}

	strCase           = "Jack"
	digitCase         = int32(111)
	OriginTypeChecker = map[string]interface{}{
		"str":       map[string]interface{}{"name": strCase},
		"str-ptr":   map[string]interface{}{"name": &strCase},
		"digit":     map[string]interface{}{"age": digitCase},
		"digit-ptr": map[string]interface{}{"age": &digitCase},
		"list":      []string{"1", "2"},
		"list-str":  &[]string{"2", "3"},
		"list-digit": map[string]interface{}{
			"score": []int32{1, 2, 3},
			"num":   []interface{}{int64(1), int8(2), uint(3), &digitCase}, // IsDigitList() == true
		},
		"map":         &map[string]interface{}{"name": "Tom"},
		"map-str-itf": map[string]interface{}{"map-str-itf-ptr": &map[string]interface{}{"name": "Tom"}},
	}

	MapInnerJsonStr = `{
    "users": [
        {
            "nam": "jack",
            "age": 23,
            "info": "{\"7351241250965703962\":\"item_id\",\"2329\":\"app_id\",\"23.29\":\"high\"}"
        },
        {
            "nam": "Tom",
            "age": 24,
            "info": "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}"
        }
    ]
	}`
	ListInnerJsonStr = `{
		"index": "[1,2,3.3,7351241250965703962]",
		"list_map": "[{\"1\":\"2\"}]"
	}`
	DataCase = `{
		"delivery_seconds": 259200,
		"ies_pricing_target": 3,
		"delivery_type": 1,
		"external_action": 96,
		"order_object_infos": [
			{
				"object_type": 1,
				"object_id": 7352741252632153353,
				"object_owner_id": 1372637839243246
			}
		],
		"product_id": 3576771487698444494,
		"payment_info": {
			"payment_type": 2,
			"pay_amount": 0,
			"balance": 100000
		},
		"cg_payment_info": {
			"url": "",
			"extra": "",
			"cash_serial": null
		},
		"extra": {
			"bound_user_ids": [
				1372637839243246
			],
			"item_id": 7352741252632153353,
			"order_package_id": 1775712180432919,
			"order_detail_type": 1,
			"is_async_order": true,
			"debug_info": {}
		},
		"estimate_profit": {
			"min_profit": 0,
			"max_profit": 0
		},
		"max_delivery_seconds": 259200
	}`
	SetAsMap = map[string]interface{}{
		"map-itf-str": map[string]interface{}{
			"coupon_info": "[\n\t{\n\t\t\"coupon_id\": 1792669145841964,\n\t\t\"coupon_discount\": 2000000,\n\t\t\"coupon_remain_discount\": 2000000,\n\t\t\"is_returned\": true\n\t}\n]",
		},
		"map-itf-list-str": map[string]interface{}{
			"coupon_list": []interface{}{"{\"coupon_id\": 1792669145841964}", "7401810649933939467", ""},
		},
		"map-str-str": map[string]interface{}{
			"key": "{\"key1\":\"{\\\"nested_key1\\\":\\\"nested_value1\\\"}\"}",
		},
		"map-str": "{\"key1\":\"{\\\"nested_key1\\\":\\\"nested_value1\\\"}\"}",
		"map-type-except": map[string]string{
			"key": "{\"nested_key1\":\"nested_value1\"}",
		},
		"map-itf-list-except": map[string]interface{}{
			"coupon_list": []string{"{\"coupon_id\": 1792669145841964}", "7401810649933939467", ""},
		},
	}
)

type MockInt int64

type MockObj struct {
	Name string
}

func (m *MockObj) String() string {
	return m.Name
}

func JsonToMap() (mapStrItf []map[string]interface{}, err error) {
	for _, s := range jsonStrList {
		mapStr, err := pkg.JsonLoadsMap(s)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("JsonLoadsMap err:%v", err))
		}
		mapStrItf = append(mapStrItf, mapStr)
	}
	return mapStrItf, nil
}

func ListJson() ([]interface{}, error) {
	var m []interface{}
	decoder := json.NewDecoder(strings.NewReader(listJson))
	decoder.UseNumber()
	err := decoder.Decode(&m)
	return m, err
}
