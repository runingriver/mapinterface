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
- 日志定义与实现剥离,若希望打印mapitf过程的异常,使用前请先将日志注入
```go
mapitf.Config().SetLogger(logs.DefaultLogger())
```

# 进阶
1. 字符串场景的支持: 默认把所有的Json str当作map[string]interface{}看待.
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

3. ForEach
```go
nameList, err := mapitf.From(mapList[2]).GetAny("users").ForEach(func(i int, k, v interface{}) (key, val interface{}) {
    idNum, cvtErr := mapitf.From(v).Get("id").ToInt64() 
    if cvtErr != nil || idNum <= 1 {
        return nil, nil // 不满足条件,则不加入结果集
    }
    firstName, cvtErr := mapitf.From(v).GetAny("name", "first").ToStr()
    if cvtErr != nil {
        return nil, nil // 不满足条件,则不加入结果集
    }
    return nil, firstName // 结果集为List
}).ToListStrF()
```

4. Setter(赋值)
```go
mapJsonStr = map[string]interface{}{
    "str-str": map[string]string{
        "info": "{\"item_id\":7351241250965703963,\"app_id\":\"2324\"}",
    },
}
orgVal, err := mapitf.From(mapJsonStr).SetMap("set-map", "hello world")
orgVal, setErr := mapitf.From(mapJsonStr).Get("str-str").SetAsMap("info")
```


注:下面的方式也是支持的,但不推荐:
```go
jsonStr := []map[string]interface{}{
    {"first":"Janet","last":"Prichard"},
    {"first":"Jack","last":"Jam"},
}
toMap, err := mapitf.From(jsonStr).Index(1).ToMap()
if err != nil {
    return
}
toMap["age"] = 24
toMap["first"] = "Kite"
```
原理: `mapinterface`在层层递进时会hold当前层的value指针,所以可以通过指针进行赋值;

5. 外传:强大的类型转换
```go
mapStr, err := mapitf.From(val).ToMapStrToStr()
```

```go
type OjbImpStr struct {
	Name string
	Age  int
}

func (o *OjbImpStr) String() string {
	return fmt.Sprintf("Name:%s,Age:%d", o.Name, o.Age)
}

ois := OjbImpStr{
   Name: "Jack",
   Age:  24,
}

oiStr, err := mapitf.From(ois).ToStr() 
// oiStr: Name:Jack,Age:24
```

```go
type MockObject int64 // 定义一个类型

mo := MockObject(10)
intStr, err := mapitf.From(mo).ToStr() // 转换为:10
```

6. Map内部Json字符串支持
- 查找时针对map内部val是json字符串的情况,我们也支持像普通map一样进行操作,内部实现原理是:把字符串反转成map再进行相关操作.
    - 注:不会把原来的map从str改成map类型;如果需要请用SetAsMap()
- ToMap系列,ToList系列的类型转换也支持将json str转成对应的类型;
```go
ListInnerJsonStr = map[string]interface{}{
    "index": "[1,2,3.3,7351241250965703962]",
    "list_map": "[{\"1\":\"2\"}]"
}

// 如下"Index(0)"取json字符串中index为0的内容再对他进行For循环
mapInt, err := mapitf.From(ListInnerJsonStr).Get("list_map").Index(0).ForEach(func(i int, k, v interface{}) (key, val interface{}) {
    return k, v
}).ToMapIntToInt()

listMap, err := mapitf.From(ListInnerJsonStr).Get("list_map").ToListMap() // 此时listMap是一个map对象
```


7. 支持ForEach能力(p0), 预案如下:  -- **v1.0.14已支持**
```go
// python => k = [get_predict(v) for k,v in dsl['predict']['risk'].iterms() if v['risk_level'] > 20]

mapitf.From(dsl).Get("predict").Get("risk").ForEach(operationFunc).ToListStr()

operationFunc = func (i int, k, v interface) (key, val interface) {
    if mapitf.From(v).Get('risk_level').ToInt() > 20 {
        return nil, nil
    }
    result := somePkg.getPredict(v)
    return nil, result
}
```

```go
// python => k = {k:get_predict(v) for k,v in dsl['predict']['risk'].iterms() if v['risk_level'] > 20}

mapitf.From(dsl).Get("predict").Get("risk").ForEach(operationFunc).ToMap()

operationFunc = func (i int, k, v interface) (key, val interface) {
    if mapitf.From(v).Get('risk_level').ToInt() > 20 {
        return nil, nil
    }
    result := somePkg.getPredict(v)
    return k, result
}
```

# 规划
1. 支持条件获取(p2), 预案如下:
```go
// select * from x where (id=1 or id>100) and name rlike "hu%"
mapitf.From().Get().Where("id", eq, 1).OrWhere("id", gt, 100).Where("name", startWith, "hu").ToList()
```

# 更新说明
1. 20240413:重大更新:
    - 支持对对象进行赋值,间Set相关方法
    - 支持ToStruct把结果转成struct
    - ToMap,ToList系列支持将json str转成Map或List
    - 支持非链式调用
    - 整体项目实现优化,删除基础类型map,统一由MapAny承担. 单测从打日志改为assert.