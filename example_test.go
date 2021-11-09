package mapinterface

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/mapitf"
)

func Test_BasicUsage1(t *testing.T) {
	jsonStr := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	var jsonMap map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &jsonMap)

	age, err := mapitf.From(jsonMap).Get("age").ToInt()
	if err != nil {
		return // Get Except return
	}
	fmt.Println("age:", age) // age = 47 int
}

func Test_BasicUsage2(t *testing.T) {
	score, err := mapitf.From(itfObj[4]).GetAny("num", 1002).Index(0).Get("math").ToInt()
	if err != nil {
		return // Get Except return
	}
	fmt.Println("score:", score) // score = 98 int
	t.Logf("score:%v", score)
}

func TestCommonUsage(t *testing.T) {
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}

	firstName, err := mapitf.From(mapList[0]).Get("name").Get("first").ToStr()
	if err != nil {
		t.Errorf("ByMapStrItf1:%v", err)
		return
	}
	t.Logf("ByMapStrItf1 result:%s", firstName)

	arrayInt, err := mapitf.From(mapList[1]).GetAny("a", "b", "c", "d", "e").ToInt()
	if err != nil {
		t.Errorf("ByListItf2:%v", err)
		return
	}
	t.Logf("ByListItf2 result:%v", arrayInt)

	toMap, err := mapitf.From(mapList[2]).Get("users").Index(1).Get("name").ToMap()
	if err != nil {
		t.Errorf("ByMapStrItf3:%v", err)
		return
	}
	t.Logf("ByMapStrItf3 result:%v", toMap)

	array, err := mapitf.From(mapList[3]).Get("vendor").Get("prices").ToList()
	if err != nil {
		t.Errorf("ByMapStrItf4:%v", err)
		return
	}
	t.Logf("ByMapStrItf4 result:%v", array)
}

func Test_ListJson(t *testing.T) {
	listItf, err := ListJson()
	if err != nil {
		t.Errorf("ListJson:%v", err)
	}
	toList, err := mapitf.From(listItf).Index(1).Get("100").ToList()
	if err != nil {
		t.Errorf("ListJson err:%v", err)
	}
	t.Logf("ListJson, ToList:%v", toList)

	toInt, err := mapitf.From(listItf).Index(1).Get("100").Index(1).ToInt()
	if err != nil {
		t.Errorf("ListJson err:%v", err)
	}
	t.Logf("ListJson, ToInt val:%v", toInt)

	toStr, err := mapitf.From(listItf).Index(2).Index(0).Get("name").ToStr()
	if err != nil {
		t.Errorf("ListJson err:%v", err)
	}
	t.Logf("ListJson, toStr val:%v", toStr)
}

func TestExceptUse(t *testing.T) {
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}

	// Index(0) - correct; Index(1) - err
	score, err := mapitf.From(mapList[4]).GetAny("name", 1002).Index(1).Get("math").ToInt64()
	if err != nil {
		t.Errorf("ByMapStrItf5:%v", err)
		return
	}
	t.Logf("ByMapStrItf5 result:%v", score)
}

func Test_GetAny(t *testing.T) {
	mi, err := pkg.JsonLoads(getAndCase)
	if err != nil {
		t.Errorf("JsonLoads err:%v", err)
		return
	}
	getInt, err := GetInt(mi, "vendor", "info", "age")
	t.Logf("int:%d,err:%v", getInt, err)
}

func GetInt(data map[string]interface{}, keys ...interface{}) (int64, error) {
	return mapitf.From(data).GetAny(keys...).ToInt64()
}

func Test_UpdateVal(t *testing.T) {
	mapList, err := JsonToMap()
	if err != nil {
		t.Errorf("JsonToMap:%v", err)
		return
	}
	t.Logf("before:%v", mapList[0])

	toMap, err := mapitf.From(mapList[0]).GetAny("name").ToMap()
	if err != nil {
		t.Errorf("ToMap:%v", err)
		return
	}
	toMap["first"] = "hello"
	t.Logf("after:%v", mapList[0])
}

func Test_Convert(t *testing.T) {
	strList := []string{"1", "2"}
	cvtList, err := mapitf.From(strList).ToListStr()
	if err != nil {
		t.Errorf("convert err:%v", err)
		return
	}
	t.Logf("convert result:%v", cvtList)

	intCase := 10
	toInt32, err := mapitf.From(intCase).ToInt32()
	if err != nil {
		t.Errorf("convert err:%v", err)
		return
	}
	t.Logf("Convert ToInt32:%v", toInt32)

	s, err := mapitf.From(itfObj[0]).GetAny(10011, "name").ToStr()
	if err != nil {
		t.Errorf("convert itfObj[0] err:%v", err)
		return
	}
	t.Logf("Convert itfObj[0] ToStr:%v", s)

	strF, err := mapitf.From(itfObj[1]).ToListStrF()
	if err != nil {
		t.Errorf("convert itfObj[1] err:%v", err)
		return
	}
	t.Logf("Convert itfObj[1] ToListStrF:%v", strF)

	lst, err := mapitf.From(itfObj[2]).ToListInt()
	if err != nil {
		t.Errorf("convert itfObj[2] err:%v", err)
		return
	}
	t.Logf("Convert itfObj[2] ToListInt:%v", lst)

	strF, err = mapitf.From(itfObj[3]).ToListStrF()
	if err != nil {
		t.Errorf("convert itfObj[3] err:%v", err)
		return
	}
	t.Logf("Convert itfObj[3] ToListStrF:%v", strF)
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
	if err != nil {
		t.Errorf("ForEach err:%v", err)
		return
	}
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
	if err != nil {
		t.Errorf("Foreach object err:%v", err)
		return
	}
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
	if err != nil {
		t.Errorf("Foreach Map to List err:%v", err)
		return
	}
	t.Logf("Foreach Map to List result:%v", nameList)

	_ = mapitf.From(mapList[3]).GetAny("vendor").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
		t.Logf("print,k:%v,v:%v", k, v)
		return nil, nil
	})
}
