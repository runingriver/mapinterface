package mapinterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/runingriver/mapinterface/pkg"
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
		`{"name":{"first":"Janet","last":"Prichard"},"age":47}`,
		// 0---------------------------------------------------------
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
		// 1---------------------------------------------------------
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
		// 2---------------------------------------------------------
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
		// 3---------------------------------------------------------
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
)

func JsonToMap() (mapStrItf []map[string]interface{}, err error) {
	for _, s := range jsonStrList {
		mapStr, err := pkg.JsonLoads(s)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("JsonLoads err:%v", err))
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
