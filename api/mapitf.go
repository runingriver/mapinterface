package api

// ForFunc 迭代函数
// i表示索引; k v表示迭代值,如果是list则k为nil;key val返回值,如果希望返回list则key为nil
// key val同时为nil则不将结果加入结果集
type ForFunc func(i int, k, v interface{}) (key, val interface{})

type MapInterface interface {
	ToBaseType
	ToMapType
	ToArrayType

	// Get 获取map中对应类型的值,key必须和map[k]v中k同类型,否则请用GetAny
	Get(key interface{}) MapInterface

	// GetAny 按照keys的顺序往下找,keys类型可以不一样
	GetAny(keys ...interface{}) MapInterface

	// Valid当前获取路径下的值是否有效
	Valid() bool

	// Exist 当前key是否存在,存在则返回对应值+true,不存在返回nil,false
	Exist(key interface{}) (interface{}, bool)

	// Index 返回当前List对应index的值
	Index(index int) MapInterface

	// ForEach迭代List或Map,不支持修改当前值
	ForEach(forFunc ForFunc) MapInterface
}

type ToBaseType interface {
	Val() (interface{}, error)
	ToStr() (string, error)
	ToByte() ([]byte, error)
	ToInt() (int, error)
	ToInt64() (int64, error)
	ToInt32() (int32, error)
	ToUint() (uint, error)
	ToUint64() (uint64, error)
	ToUint32() (uint32, error)
	ToFloat32() (float32, error)
	ToFloat64() (float64, error)
	// ToBool 如:1转成ture,"true"成true
	ToBool() (bool, error)
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
	ToListStr() ([]string, error)  // ToListStr to list string
	ToListStrF() ([]string, error) // ToListStrF 强制转换成[]string,只要是数组就能转成[]string
	ToListInt() ([]int, error)
	ToListInt32() ([]int32, error)
	ToListInt64() ([]int64, error)
	ToListUint() ([]uint, error)
	ToListUint64() ([]uint64, error)
	ToListUint32() ([]uint32, error)
	ToListFloat32() ([]float32, error)
	ToListFloat64() ([]float64, error)
	ToListBool() ([]bool, error)
}
