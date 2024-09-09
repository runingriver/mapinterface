package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/mitchellh/mapstructure"
	"github.com/runingriver/mapinterface/conf"
	"github.com/runingriver/mapinterface/itferr"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// JsonLoadsMap
// Tip: map类型的json key必须是字符串;
func JsonLoadsMap(jsonStr string) (map[string]interface{}, error) {
	var m map[string]interface{}
	dc := decoder.NewDecoder(jsonStr)
	dc.UseNumber()

	if err := dc.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

// JsonLoadsList
func JsonLoadsList(jsonStr string) ([]interface{}, error) {
	var m []interface{}
	decoder := decoder.NewDecoder(jsonStr)
	decoder.UseNumber()

	if err := decoder.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

// JsonLoadsObj 注:o必须是对象的指针
func JsonLoadsObj(jsonStr string, o interface{}) (interface{}, error) {
	if t := reflect.TypeOf(o); t.Kind() != reflect.Ptr {
		return nil, errors.New("param o must be ptr")
	}
	decoder := decoder.NewDecoder(jsonStr)
	decoder.UseNumber()

	if err := decoder.Decode(o); err == nil {
		return o, nil
	} else {
		return o, err
	}
}

// JsonUnmarshalObj 注:o必须是对象的指针
func JsonUnmarshalObj(jsonStr string, o interface{}) (interface{}, error) {
	if t := reflect.TypeOf(o); t.Kind() != reflect.Ptr {
		return nil, errors.New("param o must be ptr")
	}
	err := sonic.Unmarshal(StrToByte(jsonStr), o)
	return o, err
}

func JsonDumps(m interface{}) (string, error) {
	if str, err := sonic.MarshalString(m); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(str), nil
	}
}

// StrChecker 判断v是不是一个字符串类型,首先判断是否为字符串,其次判断是否为json
func JsonChecker(v interface{}) (bool, string) {
	if isStr, s := IsStrType(v); isStr {
		if ok := sonic.Valid(StrToByte(s)); ok { // 底层调用的json.Valid
			return true, s
		}
		return false, s
	}

	return false, ""
}

func MapToStruct(inputMap interface{}, outputStruct interface{}) (interface{}, error) {
	if t := reflect.TypeOf(inputMap); t.Kind() != reflect.Map {
		return nil, errors.New("param inputMap must be map")
	}
	if t := reflect.TypeOf(outputStruct); t.Kind() != reflect.Ptr {
		return nil, errors.New("param outputStruct must be ptr")
	}
	err := mapstructure.WeakDecode(inputMap, outputStruct)
	return outputStruct, err
}

func IsStrType(v interface{}) (bool, string) {
	v = Interpret(v)
	if v == nil {
		return false, ""
	}

	switch vv := v.(type) {
	case string:
		return true, vv
	case []byte:
		return true, ByteToStr(vv)
	default:
		return false, ""
	}
}

func ToStr(v interface{}) string {
	result := ""
	if v == nil {
		return result
	}
	switch vv := v.(type) {
	case json.Number:
		result = vv.String()
	case string:
		result = vv
	case int:
		return strconv.FormatInt(int64(vv), 10)
	case int8:
		return strconv.FormatInt(int64(vv), 10)
	case int16:
		return strconv.FormatInt(int64(vv), 10)
	case int32:
		return strconv.FormatInt(int64(vv), 10)
	case int64:
		return strconv.FormatInt(vv, 10)
	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case uint64:
		return strconv.FormatUint(vv, 10)
	case float32:
		return strconv.FormatFloat(float64(vv), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(vv, 'G', -1, 64)
	case bool:
		result = strconv.FormatBool(vv)
	case []byte:
		result = ByteToStr(vv)
	case nil:
		return ``
	case error:
		return vv.Error()
	default:
		if conf.CONF.CvtStrUseStringMethod {
			if f, ok := v.(fmt.Stringer); ok {
				return f.String()
			}

			if callRst, ok := CallMethod(v, "String"); ok {
				if result, ok = callRst.(string); ok {
					return result
				}
			}
		}

		result, _ = JsonDumps(vv)
	}
	return result
}

func CallMethod(v interface{}, methodName string) (interface{}, bool) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(v)

	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		// 如果v是值类型,需要组装出一个指针类型的变量赋值给ptr
		ptr = reflect.New(reflect.TypeOf(v))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// 检查methodName是否存在于值类型的方法中
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// 检查methodName是否存在于指针类型的方法中
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call([]reflect.Value{})[0].Interface(), true
	}

	// 无对应的方法实现
	return "", false
}

// ToInt64 如果v是float或uint,会进行强转,uint超过int64的最大数值可能丢失精度
func ToInt64(v interface{}) (int64, error) {
	v = Interpret(v)
	switch vv := v.(type) {
	case json.Number:
		if strings.ContainsAny(vv.String(), ".e") {
			result, convErr := vv.Float64()
			if convErr != nil {
				return 0, itferr.NewBaseTypeConvErr("ToInt64#vv.Float64()", "", convErr)
			}
			return int64(result), nil
		} else {
			result, convErr := vv.Int64()
			if convErr != nil {
				return 0, itferr.NewBaseTypeConvErr("ToInt64#vv.Int64()", "", convErr)
			}
			return result, nil
		}
	case string:
		return strToInt64(vv)
	case []byte:
		return strToInt64(string(vv))
	case int64:
		return vv, nil
	case int:
		return int64(vv), nil
	case int8:
		return int64(vv), nil
	case int16:
		return int64(vv), nil
	case int32:
		return int64(vv), nil
	case uint:
		return int64(vv), nil
	case uint8:
		return int64(vv), nil
	case uint16:
		return int64(vv), nil
	case uint32:
		return int64(vv), nil
	case uint64:
		return int64(vv), nil
	case float32:
		return int64(vv), nil
	case float64:
		return int64(vv), nil
	case uintptr:
		return int64(vv), nil
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil
	}

	// 再次尝试转换,与上面有重叠,这里主要解决类似:type MockInt int64的转换
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float()), nil
	}

	return 0, itferr.NewBaseTypeConvErr("ToInt64", "unknown type", nil)
}

func strToInt64(vv string) (int64, error) {
	if strings.ContainsAny(vv, ".e") {
		result, convErr := strconv.ParseFloat(vv, 64)
		if convErr != nil {
			return 0, itferr.NewBaseTypeConvErr("ToInt64#ParseFloat()", "", convErr)
		}
		return int64(result), nil
	}
	result, convErr := strconv.ParseInt(vv, 10, 64)
	if convErr != nil {
		return 0, itferr.NewBaseTypeConvErr("ToInt64#ParseInt()", "", convErr)
	}
	return result, nil
}

// ToFloat64 如果v是int64或uint64,会进行强转,数字超过float64最大数值可能丢失精度
func ToFloat64(v interface{}) (float64, error) {
	v = Interpret(v)
	switch vv := v.(type) {
	case json.Number:
		if strings.ContainsAny(vv.String(), ".e") {
			if result, convErr := vv.Float64(); convErr != nil { // will shadow origin result
				return 0, itferr.NewBaseTypeConvErr("ToFloat64#Float64()", "", convErr)
			} else {
				return result, nil
			}
		} else {
			if result, convErr := vv.Int64(); convErr != nil {
				return 0, itferr.NewBaseTypeConvErr("ToFloat64#Int64()", "", convErr)
			} else {
				return float64(result), nil
			}
		}
	case string:
		result, convErr := strconv.ParseFloat(vv, 64)
		if convErr != nil {
			return 0, itferr.NewBaseTypeConvErr("ToFloat64#ParseFloat()", "", convErr)
		}
		return result, nil
	case []byte:
		result, convErr := strconv.ParseFloat(string(vv), 64)
		if convErr != nil {
			return 0, fmt.Errorf("ToFloat64 convert err:%v", convErr)
		}
		return result, nil
	case int64:
		return float64(vv), nil
	case int:
		return float64(vv), nil
	case int8:
		return float64(vv), nil
	case int16:
		return float64(vv), nil
	case int32:
		return float64(vv), nil
	case uint:
		return float64(vv), nil
	case uint8:
		return float64(vv), nil
	case uint16:
		return float64(vv), nil
	case uint32:
		return float64(vv), nil
	case uint64:
		return float64(vv), nil
	case float32:
		return float64(vv), nil
	case float64:
		return vv, nil
	case uintptr:
		return float64(vv), nil
	case bool:
		if vv {
			return 1, nil
		}
		return 0, nil
	}

	// 再次尝试转换,与上面有重叠,这里主要解决类似:type MockFloat float64的转换
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		return float64(rv.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	}

	return 0, itferr.NewBaseTypeConvErr("ToFloat64", "unknown type", nil)
}

func ToFloat32(v interface{}) (float32, error) {
	k, err := ToFloat64(v)
	if err != nil {
		return 0, err
	}
	return float32(k), nil
}

// StrToByte 高效转换,避免内存拷贝
func StrToByte(s string) (b []byte) {
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return
}

// ByteToStr 高效转换,避免内存拷贝,
func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func IsBaseType(v interface{}) bool {
	v = Interpret(v)
	if v == nil {
		return false
	}
	switch v.(type) {
	case json.Number, string, []byte, bool:
		return true
	case int8, int16, int, int32, int64:
		return true
	case uint8, uint16, uint, uint32, uint64:
		return true
	case float32, float64:
		return true
	}
	return false
}

func Interpret(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func ReflectToVal(v interface{}) reflect.Value {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return vv
	}

	for vv.Kind() == reflect.Ptr && !vv.IsNil() {
		vv = vv.Elem()
	}
	return vv
}

func UniqList(data interface{}) interface{} {
	rfV := reflect.ValueOf(data)
	if rfV.Kind() != reflect.Array && rfV.Kind() != reflect.Slice {
		return data
	}

	// How To Reduce the duplicate code
	switch vv := data.(type) {
	case []string:
		uniqMap := make(map[string]bool, len(vv))
		uniqList := make([]string, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []int:
		uniqMap := make(map[int]bool, len(vv))
		uniqList := make([]int, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []int32:
		uniqMap := make(map[int32]bool, len(vv))
		uniqList := make([]int32, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []int8:
		uniqMap := make(map[int8]bool, len(vv))
		uniqList := make([]int8, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []int16:
		uniqMap := make(map[int16]bool, len(vv))
		uniqList := make([]int16, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []int64:
		uniqMap := make(map[int64]bool, len(vv))
		uniqList := make([]int64, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uint:
		uniqMap := make(map[uint]bool, len(vv))
		uniqList := make([]uint, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uint32:
		uniqMap := make(map[uint32]bool, len(vv))
		uniqList := make([]uint32, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uint8:
		uniqMap := make(map[uint8]bool, len(vv))
		uniqList := make([]uint8, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uint16:
		uniqMap := make(map[uint16]bool, len(vv))
		uniqList := make([]uint16, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uint64:
		uniqMap := make(map[uint64]bool, len(vv))
		uniqList := make([]uint64, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []float64:
		uniqMap := make(map[float64]bool, len(vv))
		uniqList := make([]float64, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []float32:
		uniqMap := make(map[float32]bool, len(vv))
		uniqList := make([]float32, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []interface{}:
		uniqMap := make(map[interface{}]bool, len(vv))
		uniqList := make([]interface{}, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	case []uintptr:
		uniqMap := make(map[uintptr]bool, len(vv))
		uniqList := make([]uintptr, 0, len(vv))
		for _, item := range vv {
			if _, ok := uniqMap[item]; !ok {
				uniqMap[item] = true
				uniqList = append(uniqList, item)
			}
		}
		return uniqList
	default:
		return UniqByReflect(vv)
	}
}

func UniqByReflect(data interface{}) interface{} {
	inArr := reflect.ValueOf(data)
	if inArr.Kind() != reflect.Slice && inArr.Kind() != reflect.Array {
		return data
	}

	existMap := make(map[interface{}]bool)
	outArr := reflect.MakeSlice(inArr.Type(), 0, inArr.Len())

	for i := 0; i < inArr.Len(); i++ {
		iVal := inArr.Index(i)

		if _, ok := existMap[iVal.Interface()]; !ok {
			outArr = reflect.Append(outArr, inArr.Index(i))
			existMap[iVal.Interface()] = true
		}
	}

	return outArr.Interface()
}
