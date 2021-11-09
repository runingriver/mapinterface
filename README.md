# mapinterface - A data parser lib for Go with pythonic grammar sugar and as concern as possible for high performance

mapinterface 旨在消灭对map/list解析而产生的层层断言,冗余代码从而提升代码的可读性

# 快速开始
1. `mapinterface` 对map[type]interface{},[]interface进行路径查找,并获取对应的值.

2. `mapinterface` Focus的问题: 一个大的json序列为map-interface后, 需要层层断言解析, 导致代码可读性低

3. `mapinterface` 期望提供类似于python操作dict的语法糖
 - `json.get('key1).get('key2')[1]` 或 `json['key1']['key2'][1]`
 - `[get_predict(v) for k,v in dsl['predict']['risk'].iterms() if v['risk_level'] > 20]`

4. 消灭箭头代码

6. 方便的类型转换,数据过滤功能

## 安装

```shell
$ go get -u github.com/runingriver/mapinterface
```
调试: `replace github.com/runingriver/mapinterface => ../../runingriver/mapinterface`

## Example
1. 获取map[string]interface{}中某个key对应的值
```go
jsonStr := '{"name":{"first":"Janet","last":"Prichard"},"age":47}'
var jsonMap map[string]interface{}
_ = json.Unmarshal([]byte(jsonStr), &jsonMap)
// 语法: mapitf.From().Get().GetAny().Index().GetAny().ToStr()
age, err := mapitf.From(jsonMap).Get("name").Get("first").ToStr()
if err != nil {
    return // Get Except return
}
fmt.Println("age:", age) // Janet
```

2. 获取map[string]interface{}中某个key对应的值
```go
jsonMap := map[string]interface{}{
    "num": map[int]interface{}{
        1002: []map[string]interface{}{
            {"math": 98},
            {"geography": 88},
        },
    },
}
score, err := mapitf.From(jsonMap).GetAny("num", 1002).Index(0).Get("math").ToInt()
if err != nil {
    return // Get Except return
}
fmt.Println("score:", score) // score = 98 int
```
3. 其他操作
```go
// 判断节点是否合法
isValid := mapitf.From(jsonMap).Get("name").Valid()
// 节点是否存在,如果存在返回该节点的值
if v,ok := mapitf.From(jsonMap).Get("name").Exist("first"); ok {
	// do something...
}
```

## 注意
1. 性能考虑, 不支持以下使用方式
```go
nameNode := mapitf.From(jsonMap).Get("name")
first := nameNode.Get("first")
second := nameNode.Get("first")
// 因为,如果支持以上方式,会导致过多的内存拷贝,所以建议如下使用:
first := mapitf.From(jsonMap).Get("name").Get("first")
second := mapitf.From(jsonMap).Get("name").Get("first")
```

# 进阶
1. 字符串场景的支持
```go
jsonStr := '{"name":{"first":"Janet","last":"Prichard"},"age":47}'
age, err := mapitf.From(jsonStr).Get("name").Get("first").ToStr()
if err != nil {
    return // Get Except return
}
fmt.Println("age:", age) // string ==> Janet
```

2. List Json支持
```go
jsonStr := '[{"first":"Janet","last":"Prichard"},{"first":"Jack","last":"Jam"}]'
age, err := mapitf.From(jsonStr).Index(1).Get("first").ToStr()
if err != nil {
    return // Get Except return
}
fmt.Println("age:", age) // string ==> Jack
```

# 规划
1. 支持Foreach能力(p0), 预案如下:  -- **已支持**
```go
// python => k = [get_predict(v) for k,v in dsl['predict']['risk'].iterms() if v['risk_level'] > 20]

mapitf.From(dsl).Get("predict").Get("risk").ForEach(operationFunc).ToListStr()

operationFunc = func (i int, k, v interface) (key, val interface, skip bool) {
	if mapitf.From(v).Get('risk_level').ToInt() > 20 {
		return nil, nil, true
    }   
    result := somePkg.getPredict(v)
	return nil, result, false
}
```

```go
// python => k = {k:get_predict(v) for k,v in dsl['predict']['risk'].iterms() if v['risk_level'] > 20}

mapitf.From(dsl).Get("predict").Get("risk").ForEach(operationFunc).ToMap()

operationFunc = func (i int, k, v interface) (key, val interface, skip bool) {
	if mapitf.From(v).Get('risk_level').ToInt() > 20 {
		return nil, nil, true
    }  
    result := somePkg.getPredict(v)
	return k, result, false
}
```
**进度:** v1.0.14版已支持

2. 支持条件获取(p2), 预案如下:
```go
// select * from x where (id=1 or id>100) and name rlike "hu%"
mapitf.From().Get().Where("id", eq, 1).OrWhere("id", gt, 100).Where("name", startWith, "hu").ToList()
```
