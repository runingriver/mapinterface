package api

// ForFunc 迭代函数
// i表示索引; k v表示迭代值, 如果循环的是list则k为nil;
// key val为返回值,若希望结果集为List则将key返回为nil,若希望为map则key val都不为空;若key val同时为nil则不将结果加入结果集
type ForFunc func(i int, k, v interface{}) (key, val interface{})

type MapInterface interface {
	ToBaseType
	ToMapType
	ToArrayType
	ToObjectType
	OriginTypeChecker
	SetValType

	// Get 获取map中对应类型的值,key必须和map[k]v中k同类型,否则请用GetAny
	Get(key interface{}) MapInterface

	// GetAny 按照keys的顺序往下找,keys类型可以不一样
	GetAny(keys ...interface{}) MapInterface

	// Valid当前获取路径下的值是否有效
	Valid() bool

	//Uniq 对list去重
	Uniq() MapInterface

	// Exist 当前key是否存在,存在则返回对应值+true,不存在返回nil,false.json str中存在也会返回true
	Exist(key interface{}) (interface{}, bool)

	// Index 返回当前List对应index的值
	Index(index int) MapInterface

	// ForEach迭代List或Map,不支持修改当前值
	// ForFunc 迭代函数,i表示索引; k v表示迭代值, 如果循环的是list则k为nil;
	// ForFunc 返回值:若希望结果集为List则将key返回为nil,若希望为map则key val都不为空;若key val同时为nil则不将结果加入结果集
	ForEach(forFunc ForFunc) MapInterface

	// New 用于分段调用,clone出一个新的当前现场
	New() MapInterface

	// PrintPath 打印迭代路径
	PrintPath() string
}

type ToBaseType interface {
	Val() (interface{}, error) // Val 直接返回当前节点的值
	ToStr() (string, error)
	ToByte() ([]byte, error)
	ToInt() (int, error)
	ToInt64() (int64, error)
	ToInt32() (int32, error)
	ToRune() (rune, error)
	ToUint() (uint, error)
	ToUint64() (uint64, error)
	ToUint32() (uint32, error)
	ToFloat32() (float32, error)
	ToFloat64() (float64, error)
	// ToBool 如:1转成ture,"true"成true
	ToBool() (bool, error)
}

type OriginTypeChecker interface {
	IsStr() (bool, error)
	IsDigit() (bool, error)

	// IsList 原始类型是否为一个数组,json list时为true
	IsList() (bool, error)
	// IsStrList 原始类型是否为一个字符串数组
	IsStrList() (bool, error)
	// IsDigitList 原始类型是否为一个数字型数组
	IsDigitList() (bool, error)

	// IsMap是否为map类型,json map时为true
	IsMap() (bool, error)
	// IsMap是否为map类型,json map时为true
	IsMapStrItf() (bool, error)
	//IsMapStrDigit() (bool, error)
	//IsMapDigitItf() (bool, error)
	//IsMapDigitStr() (bool, error)
}

type ToMapType interface {
	ToMap() (map[string]interface{}, error)
	ToMapInt() (map[int]interface{}, error)
	ToMapInt64() (map[int64]interface{}, error)
	ToMapInt32() (map[int32]interface{}, error)
	ToMapUint() (map[uint]interface{}, error)
	ToMapUint64() (map[uint64]interface{}, error)
	ToMapUint32() (map[uint32]interface{}, error)
	ToMapFloat32() (map[float32]interface{}, error)
	ToMapFloat64() (map[float64]interface{}, error)
	ToMapItf() (map[interface{}]interface{}, error)

	ToMapStrToStr() (map[string]string, error)
	ToMapIntToInt() (map[int]int, error)
	ToMapInt64ToInt64() (map[int64]int64, error)
	ToMapFloat64ToFloat64() (map[float64]float64, error)
	ToMapFloat32ToFloat32() (map[float32]float32, error)
}

type ToArrayType interface {
	ToList() ([]interface{}, error)
	ToListMap() ([]map[string]interface{}, error)
	ToListStr() ([]string, error)  // ToListStr to list string
	ToListStrF() ([]string, error) // ToListStrF 强制转换成[]string,只要是数组就能转成[]string
	ToListInt() ([]int, error)
	ToListInt32() ([]int32, error)
	ToListRune() ([]rune, error)
	ToListInt64() ([]int64, error)
	ToListUint() ([]uint, error)
	ToListUint64() ([]uint64, error)
	ToListUint32() ([]uint32, error)
	ToListFloat32() ([]float32, error)
	ToListFloat64() ([]float64, error)
	ToListBool() ([]bool, error)
}

type ToObjectType interface {
	ToStruct(out interface{}) (interface{}, error) // ToStruct 支持map,str,[]byte等对象转化为struct
}

type SetValType interface {
	// SetMap 设置key对应的值为val,另外,当key在json str中时,将该json序列化为map并赋值给上个节点.
	// orgVal 是开始传入的那个值,如果是str则会返回对应的map[string]interface{}
	SetMap(key interface{}, val interface{}) (orgVal interface{}, err error)
	// SetAsMap 指定key对应的值设置为map,仅当key对应的值是json str时有效
	// orgVal 是开始传入的那个值,如果是str则会返回对应的map[string]interface{}
	SetAsMap(key interface{}) (orgVal interface{}, err error)
	// SetList 设置指定index上的值为val
	//SetList(idx int, val interface{}) (orgVal interface{}, err error)
	// SetAllAsMap 递归map的每一个节点,将json str赋值为map
	SetAllAsMap() (orgVal interface{}, err error)
}
