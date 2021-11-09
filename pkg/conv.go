package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/runingriver/mapinterface/itferr"
)

func JsonLoads(jsonStr string) (map[string]interface{}, error) {
	var m map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	if err := decoder.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

func JsonLoadsList(jsonStr string) ([]interface{}, error) {
	var m []interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	if err := decoder.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

func JsonDumps(m interface{}) (string, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(m); err != nil {
		return "", err
	}
	return strings.TrimSpace(b.String()), nil
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
	case int8, int16, int, int32, int64:
		result = fmt.Sprintf("%d", vv)
	case uint8, uint16, uint, uint32, uint64:
		result = fmt.Sprintf("%d", vv)
	case float32, float64:
		result = fmt.Sprintf("%f", vv)
	case bool:
		result = strconv.FormatBool(vv)
	case []byte:
		result = ByteToStr(vv)
	default:
		dumps, _ := JsonDumps(vv)
		result = strings.TrimSpace(dumps)
	}
	return result
}

func ToInt64(v interface{}) (int64, error) {
	v = interpret(v)
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
		if strings.ContainsAny(vv, ".e") {
			result, convErr := strconv.ParseFloat(vv, 64)
			if convErr != nil {
				return 0, itferr.NewBaseTypeConvErr("ToInt64#ParseFloat()", "", convErr)
			}
			return int64(result), nil
		} else {
			result, convErr := strconv.ParseInt(vv, 10, 64)
			if convErr != nil {
				return 0, itferr.NewBaseTypeConvErr("ToInt64#ParseInt()", "", convErr)
			}
			return result, nil
		}
	case []byte:
		sv := string(vv)
		if strings.ContainsAny(sv, ".e") {
			result, convErr := strconv.ParseFloat(sv, 64)
			if convErr != nil {
				return 0, fmt.Errorf("ToInt64 convert err:%v", convErr)
			}
			return int64(result), nil
		}
		result, convErr := strconv.ParseInt(sv, 10, 64)
		if convErr != nil {
			return 0, fmt.Errorf("ToInt64 convert err:%v", convErr)
		}
		return result, nil
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
	return 0, itferr.NewBaseTypeConvErr("ToInt64", "unknown type", nil)
}

func ToFloat64(v interface{}) (float64, error) {
	v = interpret(v)
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
	v = interpret(v)
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

func interpret(a interface{}) interface{} {
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
