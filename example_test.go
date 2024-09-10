package mapinterface

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/runingriver/mapinterface/mapitf"
	"github.com/runingriver/mapinterface/pkg"
	"github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicUsage1(t *testing.T) {
	jsonStr := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	var jsonMap map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &jsonMap)

	age, err := mapitf.From(jsonMap).Get("age").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 47, age)

	path := mapitf.From(jsonMap).Get("age").PrintPath()
	assert.Equal(t, "map[string]interface {} => age:float64", path)
}

func Test_BasicUsage2(t *testing.T) {
	score, err := mapitf.From(itfObj[4]).GetAny("num", 1002).Index(0).Get("math").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 98, score)
	t.Logf("score:%v,err:%v", score, err)
}

func TestCommonUsage(t *testing.T) {
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}

	firstName, err := mapitf.From(mapList[0]).Get("name").Get("first").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Janet", firstName)
	t.Logf("ByMapStrItf1 result:%s,err:%v", firstName, err)

	arrayInt, err := mapitf.From(mapList[1]).GetAny("a", "b", "c", "d", "e").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 1, arrayInt)
	t.Logf("ByListItf2 result:%v,err:%v", arrayInt, err)

	toMap, err := mapitf.From(mapList[2]).Get("users").Index(1).Get("name").ToMap()
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"first": "Ethan", "last": "Hunt"}, toMap)
	t.Logf("ByMapStrItf3 result:%v,err:%v", toMap, err)

	array, err := mapitf.From(mapList[3]).Get("vendor").Get("prices").ToList()
	assert.Nil(t, err)
	assert.Equal(t, 6, len(array))
	// []interface {}{"2400", "2100", "1200", "400.87", "89.9", "150.1"}
	t.Logf("ByMapStrItf4 result:%v,err:%v", array, err)
}

func TestExceptionUsage(t *testing.T) {
	// 特殊的数据结构&特殊的操作方式
	// list-json.map, list.json.list,map-json.list
	// Get.Get, GetAny.GetAny, GetAny.Get, Get.GetAny,Foreach.Get(),Foreach.Index
	holder := mapitf.From(MapInnerJsonStr).GetAny("users").Index(1).Get("info").Get("app_id")
	path := holder.New().PrintPath()
	age, err := holder.ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 2324, age)
	assert.Equal(t, "map[string]interface {} => users:[]interface {} => 1:map[string]interface {} => info:map[string]interface {} => app_id:string", path)

	holder = mapitf.From(ListInnerJsonStr).GetAny("list_map").Index(0).Get("1")
	path = holder.New().PrintPath()
	toInt, err := holder.ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 2, toInt)
	assert.Equal(t, "map[string]interface {} => list_map:[]interface {} => 0:map[string]interface {} => 1:string", path)

	holder = mapitf.From(MapInnerJsonStr).Get("users").Index(0).GetAny("info").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		return v, k
	}).Get("item_id")
	path = holder.New().PrintPath()
	itemId, err := holder.ToStr()

	assert.Nil(t, err)
	assert.Equal(t, "7351241250965703962", itemId)
	assert.Equal(t, "map[string]interface {} => users:[]interface {} => 0:map[string]interface {} => info:map[string]interface {} => item_id:string", path)

	holder = mapitf.From(MapInnerJsonStr).GetAny("users").Index(1).GetAny("info").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		return nil, k
	}).Index(1)
	path = holder.New().PrintPath()
	appId, err := holder.ToStr()

	assert.Nil(t, err)
	assert.Equal(t, "app_id", appId)
	assert.Equal(t, "map[string]interface {} => users:[]interface {} => 1:map[string]interface {} => info:map[string]interface {} => 1:[]interface {}", path)

	name, err := mapitf.From(OriginTypeChecker).GetAny("map").Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", name)

	name, err = mapitf.From(OriginTypeChecker).GetAny("map", "name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", name)

	name, err = mapitf.From(OriginTypeChecker).Get("map").GetAny("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", name)

	name, err = mapitf.From(OriginTypeChecker).GetAny("map").GetAny("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", name)

	toInt, err = mapitf.From(OriginTypeChecker).Get("list-str").Index(0).ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 2, toInt)
}

func Test_ListJson(t *testing.T) {
	listItf, err := ListJson()
	if err != nil {
		t.Errorf("ListJson:%v", err)
	}
	toList, err := mapitf.From(listItf).Index(1).Get("100").ToList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(toList))
	t.Logf("ListJson, ToList:%v,err:%v", toList, err)

	toInt, err := mapitf.From(listItf).Index(1).Get("100").Index(1).ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 2, toInt)
	t.Logf("ListJson, ToInt val:%v,err:%v", toInt, err)

	toStr, err := mapitf.From(listItf).Index(2).Index(0).Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "zhangsan", toStr)
	t.Logf("ListJson, toStr val:%v,err:%v", toStr, err)
}

func Test_ItfJsonStr(t *testing.T) {
	name, err := mapitf.From(mapJsonStr).Get("student").Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Jack", name)
	t.Logf("name:%s,err:%v", name, err)

	age, err := mapitf.From(mapJsonStr).GetAny("student", "age").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 23, age)
	t.Logf("age:%d,err:%v", age, err)

	name, err = mapitf.From(mapJsonStr).GetAny("int-map-itf", 10020, "name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Jack", name)
	t.Logf("name:%s,err:%v", name, err)

	info, err := mapitf.From(mapJsonStr).Get("str-str").Get("info").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}", info)
	t.Logf("info:%s,err:%v", info, err)

	info, err = mapitf.From(mapJsonStr).GetAny("str-str", "info").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}", info)
	t.Logf("info:%s,err:%v", info, err)
}

func TestExceptUse(t *testing.T) {
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}

	// Index(0) - correct; Index(1) - err
	score, err := mapitf.From(mapList[4]).GetAny("name", 1002).Index(0).Get("math").ToInt64()
	assert.Nil(t, err)
	assert.Equal(t, int64(98), score)
	t.Logf("ByMapStrItf5 result:%v,err:%v", score, err)

	score, err = mapitf.From(mapList[4]).GetAny("name", 1002).Index(1).Get("math").ToInt64()
	assert.NotNil(t, err)
	t.Logf("ByMapStrItf5 result:%v,err:%v", score, err)
}

func Test_GetAny(t *testing.T) {
	mi, err := pkg.JsonLoadsMap(getAndCase)
	if err != nil {
		t.Errorf("JsonLoadsMap err:%v", err)
		return
	}
	getInt, err := GetInt(mi, "vendor", "info", "age")
	assert.Nil(t, err)
	assert.Equal(t, int64(25), getInt)
	t.Logf("int:%d,err:%v", getInt, err)
}

func GetInt(data map[string]interface{}, keys ...interface{}) (int64, error) {
	return mapitf.From(data).GetAny(keys...).ToInt64()
}

func Test_UpdateVal(t *testing.T) {
	// test case 1
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}
	t.Logf("before:%v", mapList[0])

	toMap, err := mapitf.From(mapList[0]).GetAny("name").ToMap()
	assert.Nil(t, err)
	assert.NotEmpty(t, toMap)
	toMap["first"] = "Brant"
	assert.Equal(t, 2, len(toMap))
	t.Logf("after:%v", mapList[0])

	// test case 2
	listMap := []map[string]interface{}{
		{"first": "Janet", "last": "Prichard"},
		{"first": "Jack", "last": "Jam"},
	}
	t.Logf("before:%v", listMap)
	toMap, err = mapitf.From(listMap).Index(1).ToMap()
	assert.Nil(t, err)
	assert.NotEmpty(t, toMap)
	toMap["age"] = 24
	toMap["first"] = "Kite"
	assert.Equal(t, 3, len(toMap))
	t.Logf("after:%v", listMap)
}

func Test_Convert(t *testing.T) {
	strList := []string{"1", "2"}
	cvtList, err := mapitf.From(strList).ToListStr()
	assert.Nil(t, err)
	assert.Equal(t, strList, cvtList)
	t.Logf("convert result:%v,err:%v", cvtList, err)

	intCase := 10
	toInt32, err := mapitf.From(intCase).ToInt32()
	assert.Nil(t, err)
	assert.Equal(t, int32(10), toInt32)
	t.Logf("Convert ToInt32:%v,err:%v", toInt32, err)

	s, err := mapitf.From(itfObj[0]).GetAny(10011, "name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", s)
	t.Logf("Convert itfObj[0] ToStr:%v,err:%v", s, err)

	strF, err := mapitf.From(itfObj[1]).ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Janet", "Tom", "Kate"}, strF)
	t.Logf("Convert itfObj[1] ToListStrF:%v,err:%v", strF, err)

	lst, err := mapitf.From(itfObj[2]).ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, []int{10011, 10012, 10013}, lst)
	t.Logf("Convert itfObj[2] ToListInt:%v,err:%v", lst, err)

	strF, err = mapitf.From(itfObj[3]).ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"{\"A\":1}", "{\"A\":2}"}, strF)
	t.Logf("Convert itfObj[3] ToListStrF:%v,err:%v", strF, err)
}

func Test_Foreach(t *testing.T) {
	lstJson, err := ListJson()
	if err != nil {
		t.Errorf("load list json err:%v", err)
		return
	}
	filter, err := mapitf.From(lstJson).Index(1).Get("100").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		intV, cvtErr := pkg.ToInt64(v)
		if cvtErr != nil || intV <= 1 {
			return nil, nil
		}
		return i, intV * 10
	}).ToMapIntToInt()
	assert.Nil(t, err)
	assert.Equal(t, map[int]int{1: 20, 2: 30}, filter)
	t.Logf("ForEach result:%v", filter)

	each, err := mapitf.From(lstJson).Index(2).ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		toMap, cvtErr := mapitf.From(v).ToMap()
		if cvtErr != nil {
			return nil, nil
		}
		if name, ok := toMap["name"]; ok {
			return nil, name
		}
		if age, ok := toMap["age"]; ok {
			return nil, age
		}
		return nil, nil
	}).ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"zhangsan", "15"}, each)
	t.Logf("foreach object:%v", each)

	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}
	nameList, err := mapitf.From(mapList[2]).GetAny("users").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		idNum, cvtErr := mapitf.From(v).Get("id").ToInt64()
		if cvtErr != nil || idNum <= 1 {
			return nil, nil
		}
		firstName, cvtErr := mapitf.From(v).GetAny("name", "first").ToStr()
		if cvtErr != nil {
			return nil, nil
		}
		return nil, firstName
	}).ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Ethan", "John"}, nameList)
	t.Logf("Foreach Map to List result:%v", nameList)

	j := 0
	_ = mapitf.From(mapList[3]).GetAny("vendor").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		t.Logf("print,k:%v,v:%v", k, v)
		j++
		return nil, nil
	})
	assert.Equal(t, 6, j)
}

func Test_Uniq(t *testing.T) {
	toInt, err := mapitf.From(UniqList).Get("UniqForInt").Uniq().ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3}, toInt)
	t.Logf("UniqForInt ok:%v,err:%v", toInt, err)

	toInt, err = mapitf.From(UniqList).Get("UniqForIntPtr").Uniq().ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 1}, toInt)
	t.Logf("UniqForIntPtr ok:%v,err:%v", toInt, err)

	toStr, err := mapitf.From(UniqList).Get("UniqForStrItf").Uniq().ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Jak", "Tom", "Kav"}, toStr)
	t.Logf("UniqForStrItf ok:%v,err:%v", toStr, err)

	toStr, err = mapitf.From(UniqList).Get("UniqForStr").Uniq().ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Jak", "Tom", "Kav"}, toStr)
	t.Logf("UniqForStr ok:%v,err:%v", toStr, err)

	objList, err := mapitf.From(UniqList).Get("UniqForObj").Uniq().ToList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(objList))
	t.Logf("UniqForObj ok:%v,err:%v", objList, err)

	objListStr, err := mapitf.From(UniqList).Get("UniqForObjPtr").Uniq().ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "1"}, objListStr)
	t.Logf("UniqForObjPtr ok:%v,err:%v", objListStr, err)
}

func Test_CvtMap(t *testing.T) {
	m, err := mapitf.From(CvtMap).Get("StrToStrForStrItf").ToMapStrToStr()
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"Kav": "101",
		"Tom": "102",
		"Jak": "103",
	}, m)
	t.Logf("Cvt ToMapStrToStr of StrToStrForStrItf ok:%v,err:%v", m, err)

	m, err = mapitf.From(CvtMap).Get("StrToStrForItfItf").ToMapStrToStr()
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"Kav": "101",
		"Tom": "102",
		"Jak": "103",
	}, m)
	t.Logf("Cvt ToMapStrToStr of StrToStrForItfItf ok:%v,err:%v", m, err)

	m, err = mapitf.From(CvtMap).Get("StrToStrForItfObj").ToMapStrToStr()
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"101": "Kav",
		"102": "Tom",
		"103": "Jak",
	}, m)
	t.Logf("Cvt ToMapStrToStr of StrToStrForItfObj ok:%v,err:%v", m, err)

	// -------------------------------------------------------------------------------
	mapInt, err := mapitf.From(CvtMap).Get("IntToIntForStrItf").ToMapIntToInt()
	assert.Nil(t, err)
	assert.Equal(t, map[int]int{
		101: 1,
		102: 2,
		103: 3,
	}, mapInt)
	t.Logf("Cvt ToMapStrToStr of IntToIntForStrItf ok:%v,err:%v", mapInt, err)

	mapInt, err = mapitf.From(CvtMap).Get("IntToIntForItfItf").ToMapIntToInt()
	assert.Nil(t, err)
	assert.Equal(t, map[int]int{
		101: 1,
		102: 2,
		103: 3,
	}, mapInt)
	t.Logf("Cvt ToMapStrToStr of IntToIntForItfItf ok:%v,err:%v", mapInt, err)

	mapInt, err = mapitf.From(CvtMap).Get("IntToIntForItfObj").ToMapIntToInt()
	assert.Nil(t, err)
	assert.Equal(t, map[int]int{
		101: 1,
		102: 2,
		103: 3,
	}, mapInt)
	t.Logf("Cvt ToMapStrToStr of IntToIntForItfObj ok:%v,err:%v", mapInt, err)
}

func Test_CvtList(t *testing.T) {
	ctx := context.Background()
	toList, err := mapitf.Fr(ctx, CvtList).Get("ToListForRf").ToList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(toList))
	t.Logf("Cvt ToListForRf of ToList ok:%v,err:%v", toList, err)

	toListStr, err := mapitf.Fr(ctx, CvtList).Get("ToListStrForItf").ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Kav", "Tom", "Jak"}, toListStr)
	t.Logf("Cvt ToListStrForItf of ToListStrF ok:%v,err:%v", toListStr, err)

	toListStr, err = mapitf.Fr(ctx, CvtList).Get("ToListStrForRf").ToListStrF()
	assert.Nil(t, err)
	assert.Equal(t, []string{"Kav", "Tom", "Jak"}, toListStr)
	t.Logf("Cvt ToListStrForRf of ToListStrF ok:%v,err:%v", toListStr, err)

	toListInt, err := mapitf.Fr(ctx, CvtList).Get("ToListIntForItf").ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, []int{101, 102, 103}, toListInt)
	t.Logf("Cvt ToListIntForItf of ToListInt ok:%v,err:%v", toListInt, err)

	toListInt, err = mapitf.Fr(ctx, CvtList).Get("ToListIntForRf").ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, []int{101, 102, 103}, toListInt)
	t.Logf("Cvt ToListIntForRf of ToListInt ok:%v,err:%v", toListInt, err)
}

func Test_Val(t *testing.T) {
	jsonStr := `{"error_code":0,"error_message":"success","data":{"topic_id": 28074}}`
	topic, err := mapitf.From(jsonStr).Get("data").Val()
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"topic_id": json.Number("28074")}, topic)
	t.Logf("val:%v,err:%v", topic, err)
}

func Test_MapAny(t *testing.T) {
	// test case 1
	strVal, err := mapitf.From(MapAny).GetAny("ptr-map-int-map-str", 1, "a").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "b", strVal)
	t.Logf("Test_MapAny case1 strVal:%s,err:%v", strVal, err)

	// test case 2
	strVal, err = mapitf.From(MapAny["map-int-map-str"]).GetAny(1, "a").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "b", strVal)
	t.Logf("Test_MapAny case2 strVal:%s,err:%v", strVal, err)

	// test case 3
	strVal, err = mapitf.From(MapAny["ptr-map-int-map-str"]).GetAny(1, "a").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "b", strVal)
	t.Logf("Test_MapAny case3 strVal:%s,err:%v", strVal, err)

	// test case 4
	strVal, err = mapitf.From(MapAny["map-str-ptr"]).Get("1").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "{\"A\":\"1\"}", strVal)
	t.Logf("Test_MapAny case4 strVal:%s,err:%v", strVal, err)

	// test case 5
	strVal, err = mapitf.From(MapAny["map-str-map-str"]).GetAny("1", "a").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "b", strVal)
	t.Logf("Test_MapAny case4 strVal:%s,err:%v", strVal, err)
}

type OjbImpStr struct {
	Name string
	Age  int
}

func (o *OjbImpStr) String() string {
	return fmt.Sprintf("Name:%s,Age:%d", o.Name, o.Age)
}

type MockObject int64

func Test_Conv(t *testing.T) {
	toByte, err := mapitf.From(map[string]string{"foo": "bar"}).ToByte()
	assert.Nil(t, err)
	assert.Equal(t, len(toByte), 13)

	mo := MockObject(10)
	intStr, err := mapitf.From(mo).ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "10", intStr)
	t.Logf("ToStr result:%s,err:%v", intStr, err)

	ois := OjbImpStr{
		Name: "Jack",
		Age:  24,
	}
	oiStr, err := mapitf.From(ois).ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Name:Jack,Age:24", oiStr)
	t.Logf("ToStr result:%s,err:%v", oiStr, err)
}

func Test_OriginTypeChecker(t *testing.T) {
	// test case1
	isStr, err := mapitf.From(OriginTypeChecker).GetAny("str", "name").IsStr()
	assert.Nil(t, err)
	assert.Equal(t, true, isStr)
	t.Logf("IsStr result:%t,err:%v", isStr, err)

	// test case2
	isStr, err = mapitf.From(OriginTypeChecker).GetAny("str-ptr", "name").IsStr()
	assert.Nil(t, err)
	assert.Equal(t, true, isStr)
	t.Logf("IsStr result:%t,err:%v", isStr, err)

	// test case3
	isDigit, err := mapitf.From(OriginTypeChecker).GetAny("digit", "age").IsDigit()
	assert.Nil(t, err)
	assert.Equal(t, true, isDigit)
	t.Logf("IsDigit result:%t,err:%v", isDigit, err)

	// test case4
	isDigit, err = mapitf.From(OriginTypeChecker).GetAny("digit-ptr", "age").IsDigit()
	assert.Nil(t, err)
	assert.Equal(t, true, isDigit)
	t.Logf("IsDigit result:%t,err:%v", isDigit, err)

	// test case5
	isList, err := mapitf.From(OriginTypeChecker).GetAny("list").IsList()
	assert.Nil(t, err)
	assert.Equal(t, true, isList)
	t.Logf("IsList result:%t,err:%v", isList, err)

	// test case6
	isList, err = mapitf.From(OriginTypeChecker).GetAny("list-str").IsStrList()
	assert.Nil(t, err)
	assert.Equal(t, true, isList)
	t.Logf("IsList result:%t,err:%v", isList, err)

	// test case7
	isList, err = mapitf.From(OriginTypeChecker).GetAny("list-digit", "score").IsDigitList()
	assert.Nil(t, err)
	assert.Equal(t, true, isList)
	t.Logf("IsDigit result:%t,err:%v", isList, err)

	// test case8
	isList, err = mapitf.From(OriginTypeChecker).GetAny("list-digit", "num").IsDigitList()
	assert.Nil(t, err)
	assert.Equal(t, true, isList)
	t.Logf("IsDigitList result:%t,err:%v", isList == false, err)

	// test case9
	isMap, err := mapitf.From(OriginTypeChecker).GetAny("map").IsMap()
	assert.Nil(t, err)
	assert.Equal(t, true, isMap)
	t.Logf("IsMap result:%t,err:%v", isMap, err)

	// test case10
	isMap, err = mapitf.From(OriginTypeChecker).GetAny("map-str-itf", "map-str-itf-ptr").IsMapStrItf()
	assert.Nil(t, err)
	assert.Equal(t, true, isMap)
	t.Logf("IsMapStrItf result:%t,err:%v", isMap, err)
}

func Test_ListIndex(t *testing.T) {
	jsonStr := `{"name":["first","Janet","last","Prichard"],"age":47}`
	var jsonMap map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &jsonMap)

	val, err := mapitf.From(jsonMap).Get("name").Index(0).ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "first", val)
	fmt.Println("val:", val)

	jsonStr = `{"name":[],"age":47}`
	_ = json.Unmarshal([]byte(jsonStr), &jsonMap)

	val, err = mapitf.From(jsonMap).Get("name").Index(0).ToStr()
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	list, _ := mapitf.From(jsonMap).Get("name").ToList()
	val, err = mapitf.From(list).Index(0).ToStr()
	assert.NotNil(t, err)
	assert.Equal(t, "", val)
}

func Test_ToMapForJson(t *testing.T) {
	name, err := mapitf.From(MapInnerJsonStr).Get("users").Index(1).Get("nam").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Tom", name)

	mapRst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(1).Get("info").ToMap()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "2324", mapRst["app_id"])

	map64Rst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(0).Get("info").ToMapInt64()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "app_id", map64Rst[2329])

	mapItfRst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(0).Get("info").ToMapItf()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 3, len(mapItfRst))

	mapIntRst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(0).Get("info").ToMapInt()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "item_id", mapIntRst[7351241250965703962])

	mapf32Rst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(0).Get("info").ToMapFloat32()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "high", mapf32Rst[23.29])

	toMap, err := mapitf.From(MapInnerJsonStr).Get("users").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		name, err := mapitf.From(v).Get("nam").ToStr()
		if err != nil {
			return nil, nil
		}
		keyList, err := mapitf.From(v).Get("info").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
			return nil, pkg.ToStr(k)
		}).ToListStr()
		if err != nil {
			return nil, nil
		}
		return name, keyList
	}).ToMap()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(toMap))
	assert.Equal(t, 2, len(toMap["Tom"].([]string)))
}

func Test_ToListForJson(t *testing.T) {
	list, err := mapitf.From(ListInnerJsonStr).Get("index").ToList()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(list))

	listInt, err := mapitf.From(ListInnerJsonStr).Get("index").ToListInt()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(listInt))

	listMap, err := mapitf.From(ListInnerJsonStr).Get("list_map").ToListMap()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(listMap))
	assert.Equal(t, "2", listMap[0]["1"])

	listStr, err := mapitf.From(ListInnerJsonStr).Get("list_map").ToListStr()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(listStr))

	listI64, err := mapitf.From(ListInnerJsonStr).Get("index").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		intVal, CvtErr := pkg.ToInt64(v)
		if CvtErr != nil {
			return nil, nil
		}
		if intVal < 100 {
			return nil, nil
		}
		return nil, v
	}).ToListInt64()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(listI64))

	mapInt, err := mapitf.From(ListInnerJsonStr).Get("list_map").Index(0).ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		return k, v
	}).ToMapIntToInt()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(mapInt))
}

func Test_ToStruct(t *testing.T) {
	type Info struct {
		ItemId int64  `json:"item_id" mapstructure:"item_id"`
		AppId  string `json:"app_id" mapstructure:"app_id"`
	}
	convey.Convey("Test_ToStruct", t, func() {
		convey.Convey("str to struct", func() {
			var info Info
			infoRst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(1).Get("info").ToStruct(&info)
			assert.Nil(t, err)
			assert.Equal(t, &Info{
				ItemId: 7351241250965703963,
				AppId:  "2324",
			}, infoRst)
		})

		convey.Convey("map to struct", func() {
			mapRst, err := mapitf.From(MapInnerJsonStr).Get("users").Index(1).Get("info").ToMap()
			assert.Nil(t, err)

			var info Info
			infoRst, err := mapitf.From(mapRst).ToStruct(&info)
			assert.Nil(t, err)
			assert.Equal(t, &Info{
				ItemId: 7351241250965703963,
				AppId:  "2324",
			}, infoRst)
		})

		convey.Convey("object type equal", func() {
			mapStruct := map[string]Info{
				"Tom": {
					ItemId: 5739458739,
					AppId:  "222",
				},
			}
			var info Info
			infoRst, err := mapitf.From(mapStruct).Get("Tom").ToStruct(&info)
			assert.Nil(t, err)
			assert.Equal(t, Info{
				ItemId: 5739458739,
				AppId:  "222",
			}, infoRst)
		})
	})
}

func Test_SetMap(t *testing.T) {
	orgVal, err := mapitf.From(mapJsonStr).SetMap("set-map", "hello world")
	assert.Nil(t, err)
	orgMap := orgVal.(map[string]interface{})
	assert.Equal(t, "hello world", orgMap["set-map"])

	orgVal, err = mapitf.From(mapJsonStr).SetMap("student", map[string]interface{}{"name": "Jack", "age": 23})
	assert.Nil(t, err)
	val, err := mapitf.From(orgVal).GetAny("student").Val()
	assert.Nil(t, err)
	orgMap, ok := val.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Jack", orgMap["name"])

	orgVal, err = mapitf.From(mapJsonStr).Get("int-map-itf").Get(10020).SetMap("age", 24)
	assert.Nil(t, err)
	val, err = mapitf.From(orgVal).GetAny("int-map-itf", 10020, "age").Val()
	assert.Nil(t, err)
	age, ok := val.(int)
	assert.True(t, ok)
	assert.Equal(t, 24, age)
	t.Logf("%v", pkg.ToStr(mapJsonStr))
}

func Test_SetAsMap(t *testing.T) {
	orgVal, err := mapitf.From(MapInnerJsonStr).Get("users").Index(0).SetAsMap("info")
	assert.Nil(t, err)
	val, err := mapitf.From(orgVal).Get("users").Index(0).Get("info").Val()
	assert.Nil(t, err)
	_, ok := val.(map[string]interface{})
	assert.True(t, ok)

	orgVal, setErr := mapitf.From(mapJsonStr).Get("str-itf-str").SetAsMap("info")
	assert.Nil(t, setErr)
	val, setErr = mapitf.From(orgVal).GetAny("str-itf-str", "info").Val()
	assert.Nil(t, setErr)
	_, ok = val.(map[string]interface{})
	assert.True(t, ok)

	// exception: 1.不是当前值不是map; 2. 写入的对象不存在;
	_, err = mapitf.From(mapJsonStr).Get("str-str").SetAsMap("info")
	assert.NotNil(t, err)

	_, err = mapitf.From(mapJsonStr).GetAny("str-str").SetAsMap("age")
	assert.NotNil(t, err)

	_, err = mapitf.From(mapJsonStr).GetAny("int-map-itf", 10020, "name").SetAsMap("name")
	assert.NotNil(t, err)

	_, err = mapitf.From(DataCase).Get("extra").SetAsMap("cg_payment_info")
	assert.NotNil(t, err)
}

func Test_NewCodingMode(t *testing.T) {
	idxItf := mapitf.From(listJson).Index(2)

	str, err := idxItf.New().Index(0).Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "zhangsan", str)

	toInt, err := idxItf.New().Index(1).Get("age").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 15, toInt)

	path1 := idxItf.New().Index(0).Get("name").PrintPath()
	assert.Equal(t, "[]interface {} => 2:[]interface {} => 0:map[string]interface {} => name:string", path1)
	path2 := idxItf.New().Index(1).Get("age").PrintPath()
	assert.Equal(t, "[]interface {} => 2:[]interface {} => 1:map[string]interface {} => age:json.Number", path2)

	vendor := mapitf.From(jsonStrList).Index(3).Get("vendor")
	toMap, err := vendor.New().ToMap()
	assert.Nil(t, err)
	assert.Equal(t, 6, len(toMap))
	toMap["name"] = "Jack"
	toMap["age"] = 23

	toStr, err := vendor.New().Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Jack", toStr)

	age, err := vendor.New().Get("age").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 23, age)
	path := vendor.New().Get("age").PrintPath()
	assert.Equal(t, "[]string => 3:map[string]interface {} => vendor:map[string]interface {} => age:int", path)

	name, err := mapitf.From(itfObj).Index(0).Get(10010).Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Janet", name)

	normalMal := mapitf.From(itfObj).Index(0)
	name, err = normalMal.New().Get(10010).Get("name").ToStr()
	assert.Nil(t, err)
	assert.Equal(t, "Janet", name)

	age, err = normalMal.New().GetAny(10011, "age").ToInt()
	assert.Nil(t, err)
	assert.Equal(t, 21, age)
}

func TestIterPath(t *testing.T) {
	path := mapitf.From(listJson).Index(0).Get("key1").PrintPath()
	exceptPath := "[]interface {} => 0:map[string]interface {} => key1:string"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(listJson).Index(1).Get("100").Index(2).PrintPath()
	exceptPath = "[]interface {} => 1:map[string]interface {} => 100:[]interface {} => 2:json.Number"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(listJson).Index(2).Index(1).Get("age").PrintPath()
	exceptPath = "[]interface {} => 2:[]interface {} => 1:map[string]interface {} => age:json.Number"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(CvtList).Get("ToListIntForItf").Index(1).PrintPath()
	exceptPath = "map[string]interface {} => ToListIntForItf:[]interface {} => 1:mapinterface.MockInt"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(jsonStrList).Index(1).GetAny("a", "b", "c", "d", "e").PrintPath()
	exceptPath = "[]string => 1:map[string]interface {} => a:map[string]interface {} => b:map[string]interface {} => c:map[string]interface {} => d:map[string]interface {} => e:json.Number"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(jsonStrList).Index(3).GetAny("vendor", "items").Index(1).Get("name").PrintPath()
	exceptPath = "[]string => 3:map[string]interface {} => vendor:map[string]interface {} => items:[]interface {} => 1:map[string]interface {} => name:string"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(itfObj).Index(4).GetAny("num", 1001).Index(0).Get("math").PrintPath()
	exceptPath = "[]interface {} => 4:map[string]interface {} => num:map[int]interface {} => 1001:[]map[string]interface {} => 0:map[string]interface {} => math:int"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(MapInnerJsonStr).GetAny("users").Index(0).GetAny("info", "2329").PrintPath()
	exceptPath = "map[string]interface {} => users:[]interface {} => 0:map[string]interface {} => info:map[string]interface {} => 2329:string"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(MapInnerJsonStr).GetAny("users").Index(1).GetAny("info", "app_id").PrintPath()
	exceptPath = "map[string]interface {} => users:[]interface {} => 1:map[string]interface {} => info:map[string]interface {} => app_id:string"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(ListInnerJsonStr).GetAny("index").Index(1).PrintPath()
	exceptPath = "map[string]interface {} => index:[]interface {} => 1:json.Number"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)

	path = mapitf.From(ListInnerJsonStr).GetAny("list_map").Index(0).Get("1").PrintPath()
	exceptPath = "map[string]interface {} => list_map:[]interface {} => 0:map[string]interface {} => 1:string"
	assert.Equal(t, exceptPath, path)
	t.Logf("%s", path)
}

func Test_SetAllAsMap(t *testing.T) {
	val, err := mapitf.From(SetAsMap).Get("map-itf-str").SetAllAsMap()
	assert.Nil(t, err)
	couponVal, err := mapitf.From(val).GetAny("map-itf-str", "coupon_info").Val()
	assert.Nil(t, err)
	assert.IsType(t, []interface{}{}, couponVal)

	val, err = mapitf.From(SetAsMap).Get("map-itf-list-str").SetAllAsMap()
	assert.Nil(t, err)
	listVal, err := mapitf.From(val).GetAny("map-itf-list-str", "coupon_list").Val()
	assert.Nil(t, err)
	assert.Len(t, listVal, 3)
	assert.IsType(t, []interface{}{}, listVal)
	listItf := listVal.([]interface{})
	assert.IsType(t, map[string]interface{}{}, listItf[0])

	val, err = mapitf.From(SetAsMap).Get("map-str-str").SetAllAsMap()
	assert.Nil(t, err)
	key1, _ := mapitf.From(SetAsMap).GetAny("map-str-str", "key").Val()
	assert.IsType(t, map[string]interface{}{}, key1)
	key2, _ := mapitf.From(SetAsMap).GetAny("map-str-str", "key", "key1").Val()
	assert.IsType(t, map[string]interface{}{}, key2)

	_, err = mapitf.From(SetAsMap).Get("map-str").SetAllAsMap()
	assert.Nil(t, err)
	assert.IsType(t, map[string]interface{}{}, SetAsMap["map-str"])

	val, err = mapitf.From(SetAsMap).Get("map-type-except").SetAllAsMap()
	assert.Nil(t, err)
	val, _ = mapitf.From(SetAsMap).GetAny("map-type-except", "key").Val()
	assert.IsType(t, "", val)

	val, err = mapitf.From(SetAsMap).Get("map-itf-list-except").SetAllAsMap()
	assert.Nil(t, err)
	val, _ = mapitf.From(SetAsMap).GetAny("map-itf-list-except", "coupon_list").Val()
	assert.IsType(t, []string{}, val)
}

func TestName(t *testing.T) {
	js := map[string]interface{}{
		"key1": "{\"key1\":\"{\\\"nested_key1\\\":\\\"nested_value1\\\"}\"}",
	}
	marshalString, _ := sonic.MarshalString(js)
	t.Logf(marshalString)
}
