package mapitf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
)

type BaseItfImpl struct {
	Class     string
	IterVal   interface{}
	IterLevel int
	ItfErr    *itferr.MapItfError
}

func (b *BaseItfImpl) Valid() bool {
	if b.ItfErr != nil || b.IterVal == nil {
		return false
	}

	return true
}

func (b *BaseItfImpl) Exist(key interface{}) (interface{}, bool) {
	if b.ItfErr != nil || b.IterVal == nil {
		return nil, false
	}
	if val, err := b.GetByInterface(key); err == nil && val != nil {
		return val, true
	}

	return nil, false
}

// GetByInterface 如果key的类型和Map[k]v中k类型不一致会panic
func (b *BaseItfImpl) GetByInterface(key interface{}) (itfVal interface{}, err *itferr.MapItfError) {
	defer func() {
		if err := recover(); err != nil {
			b.ItfErr = itferr.NewMapItfErrX(fmt.Sprintf("GetByInterface(%v)", key), itferr.UnrecoverableErr)
		}
	}()

	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		b.ItfErr = itferr.NewKeyTypeErr(fmt.Sprintf("GetByInterface(%s)", key))
		return nil, b.ItfErr
	}

	keys := v.MapKeys()
	if len(keys) == 0 {
		b.ItfErr = itferr.NewIllegalMapObject(fmt.Sprintf("GetByInterface(%s)#%+v", key, b.IterVal))
		return nil, b.ItfErr
	}

	keyV := reflect.ValueOf(key)
	if keys[0].Kind() != keyV.Kind() {
		for _, keyVal := range keys {
			if keyItf, itfError := b.toInterface(keyVal); itfError == nil {
				if pkg.ToStr(keyItf) == pkg.ToStr(key) {
					dstVal := v.MapIndex(keyVal)
					return b.toInterface(dstVal)
				}
			}
		}

		b.ItfErr = itferr.NewTypeMismatchErr(fmt.Sprintf("GetByInterface(%s)%v:%v", key, keys[0].Kind(), keyV.Kind()))
		return nil, b.ItfErr
	}

	itfVal, b.ItfErr = b.toInterface(v.MapIndex(keyV))
	return itfVal, b.ItfErr
}

func (b *BaseItfImpl) toInterface(dstVal reflect.Value) (interface{}, *itferr.MapItfError) {
	if dstVal.Kind() == reflect.Ptr || !dstVal.IsValid() || !dstVal.CanInterface() {
		return nil, itferr.NewValueTypeErr("GetByInterface#MapIndex")
	}

	return dstVal.Interface(), nil
}

func (b *BaseItfImpl) Val() (interface{}, error) {
	return b.IterVal, b.ItfErr
}

func (b *BaseItfImpl) ToStr() (string, error) {
	if b.ItfErr != nil {
		return "", b.ItfErr
	}
	return pkg.ToStr(b.IterVal), nil
}

func (b *BaseItfImpl) ToByte() ([]byte, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if s, ok := b.ToStr(); ok == nil {
		return pkg.StrToByte(s), nil
	}

	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToByte", b.Class))
}

func (b *BaseItfImpl) ToInt() (int, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return int(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt", b.Class))
}

func (b *BaseItfImpl) ToInt64() (int64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int64); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return v, nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt64", b.Class))
}

func (b *BaseItfImpl) ToInt32() (int32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int32); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return int32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt32", b.Class))
}

func (b *BaseItfImpl) ToUint() (uint, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint", b.Class))
}

func (b *BaseItfImpl) ToUint64() (uint64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint64); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint64(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint64", b.Class))
}

func (b *BaseItfImpl) ToUint32() (uint32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint32); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint32", b.Class))
}

func (b *BaseItfImpl) ToFloat32() (float32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(float32); ok {
		return v, nil
	}

	if v, ok := pkg.ToFloat64(b.IterVal); ok == nil {
		return float32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToFloat", b.Class))
}

func (b *BaseItfImpl) ToFloat64() (float64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(float64); ok {
		return v, nil
	}

	if v, ok := pkg.ToFloat64(b.IterVal); ok == nil {
		return v, nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToFloat64", b.Class))
}

func (b *BaseItfImpl) ToBool() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}

	if v, ok := b.IterVal.(bool); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return v == 1, nil
	}

	v := strings.TrimSpace(pkg.ToStr(b.IterVal))
	if v == "true" {
		return true, nil
	}
	if v == "false" {
		return false, nil
	}

	return false, itferr.NewConvFailed(fmt.Sprintf("%s#ToBool", b.Class))
}

// ToMapType ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) ToMap() (map[string]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[string]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMap", b.Class))

}

func (b *BaseItfImpl) ToMapInt() (map[int]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt", b.Class))
}

func (b *BaseItfImpl) ToMapInt64() (map[int64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int64]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64", b.Class))
}

func (b *BaseItfImpl) ToMapInt32() (map[int32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int32]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt32", b.Class))
}

func (b *BaseItfImpl) ToMapUint() (map[uint]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint", b.Class))
}
func (b *BaseItfImpl) ToMapUint64() (map[uint64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint64]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint64", b.Class))
}

func (b *BaseItfImpl) ToMapUint32() (map[uint32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint32]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64", b.Class))
}

func (b *BaseItfImpl) ToMapFloat32() (map[float32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float32]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32", b.Class))
}

func (b *BaseItfImpl) ToMapFloat64() (map[float64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float64]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64", b.Class))
}
func (b *BaseItfImpl) ToMapItf() (map[interface{}]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapItf", b.Class))
}

func (b *BaseItfImpl) ToMapStrToStr() (map[string]string, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[string]string); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapStrToStr", b.Class))
}

func (b *BaseItfImpl) ToMapIntToInt() (map[int]int, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[int]int); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapIntToInt", b.Class))
}

func (b *BaseItfImpl) ToMapInt64ToInt64() (map[int64]int64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[int64]int64); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class))
}

func (b *BaseItfImpl) ToMapFloat64ToFloat64() (map[float64]float64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float64]float64); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class))
}

func (b *BaseItfImpl) ToMapFloat32ToFloat32() (map[float32]float32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float32]float32); ok {
		return v, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class))
}

// ToArrayType ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) ToList() ([]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]interface{}); ok {
		return v, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToList", b.Class))
	}

	resultList := make([]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToList", b.Class))
		}
		resultList = append(resultList, ele.Interface())
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListStr() ([]string, error) {
	return b.toListStr()
}

func (b *BaseItfImpl) ToListStrF() ([]string, error) {
	return b.toListStr(true)
}

// ToListStr to list string, force表示是否强转,只要是数组就能转成[]string
func (b *BaseItfImpl) toListStr(force ...bool) ([]string, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]string); ok {
		return v, nil
	}
	if vList, ok := b.IterVal.([]interface{}); ok {
		result := make([]string, 0, len(vList))
		for _, v := range vList {
			result = append(result, pkg.ToStr(v))
		}
		return result, nil
	}

	// 不强制转换-则报错返回; 强制转换且b.IterVal不是list类型-则报错返回
	rv := reflect.ValueOf(b.IterVal)
	if !(len(force) > 0 && force[0]) || (rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice) {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
	}

	resultList := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
		}
		resultList = append(resultList, pkg.ToStr(ele.Interface()))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListInt() ([]int, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]int, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, int(vint64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt", b.Class))
	}

	resultList := make([]int, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, int(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListInt32() ([]int32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int32); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]int32, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, int32(vint64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
	}

	resultList := make([]int32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, int32(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListInt64() ([]int64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int64); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]int64, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, vint64)
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
	}

	resultList := make([]int64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, int64(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListUint() ([]uint, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]uint, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, uint(vint64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
	}

	resultList := make([]uint, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListUint64() ([]uint64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint64); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]uint64, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, uint64(vint64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
	}

	resultList := make([]uint64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint64(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListUint32() ([]uint32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint32); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]uint32, 0, len(vs))
		for _, v := range vs {
			if vint64, ok := pkg.ToInt64(v); ok == nil {
				listInt = append(listInt, uint32(vint64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
	}

	resultList := make([]uint32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint32(iv))
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListFloat32() ([]float32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]float32); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]float32, 0, len(vs))
		for _, v := range vs {
			if f64, ok := pkg.ToFloat64(v); ok == nil {
				listInt = append(listInt, float32(f64))
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
	}

	resultList := make([]float32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
		}
		if iv, ok := pkg.ToFloat32(ele.Interface()); ok == nil {
			resultList = append(resultList, iv)
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListFloat64() ([]float64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]float64); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]float64, 0, len(vs))
		for _, v := range vs {
			if f64, ok := pkg.ToFloat64(v); ok == nil {
				listInt = append(listInt, f64)
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
	}

	resultList := make([]float64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
		}
		if iv, ok := pkg.ToFloat64(ele.Interface()); ok == nil {
			resultList = append(resultList, iv)
		}
	}

	if len(resultList) == 0 {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListBool() ([]bool, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]bool); ok {
		return v, nil
	}

	if vs, ok := b.IterVal.([]interface{}); ok {
		listInt := make([]bool, 0, len(vs))
		for _, v := range vs {
			if bl, ok := v.(bool); ok {
				listInt = append(listInt, bl)
			} else {
				return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
			}
		}
		return listInt, nil
	}
	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListBool", b.Class))
}

// Interface ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) Get(key interface{}) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	b.ItfErr = itferr.NewUnSupportInterfaceFunc(fmt.Sprintf("BaseItfImpl#Get(%v)", key))
	return b
}

func (b *BaseItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	b.ItfErr = itferr.NewUnSupportInterfaceFunc(fmt.Sprintf("BaseItfImpl#GetAny(%v)", keys))
	return b
}

func (b *BaseItfImpl) Index(index int) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	v := reflect.ValueOf(b.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() < index {
			b.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("Index(%d)#(%d)", index, v.Len()))
			return b
		}
	default:
		b.ItfErr = itferr.NewCurrentCannotUseIndex(fmt.Sprintf("BaseItfImpl#Index(%v)", index))
		return b
	}

	switch vv := b.IterVal.(type) {
	case []interface{}:
		b.IterLevel++
		b.IterVal = vv[index]

		v := reflect.ValueOf(vv[index])
		switch v.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			return From(vv[index])
		}
		return b
	case []map[string]interface{}:
		b.IterLevel++
		b.IterVal = vv[index]
		return From(vv[index])
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		indexV := v.Index(index)
		if !indexV.IsValid() || !indexV.CanInterface() {
			b.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("Index(%d)", index))
			return b
		}

		interfaceV := indexV.Interface()
		b.IterLevel++
		b.IterVal = interfaceV

		switch indexV.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			return From(interfaceV)
		}
		return b
	}

	b.ItfErr = itferr.NewCurrentCannotUseIndex(fmt.Sprintf("BaseItfImpl#Index(%v)", index))
	return b
}

func (b *BaseItfImpl) ForEach(forFunc api.ForFunc) api.MapInterface {
	v := reflect.ValueOf(b.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		resultList := make([]interface{}, 0, v.Len())
		resultMap := make(map[interface{}]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			idxV := v.Index(i)
			if !idxV.IsValid() || !idxV.CanInterface() {
				continue
			}
			key, val := forFunc(i, nil, idxV.Interface())
			if key == nil && val == nil {
				continue
			}
			if key != nil {
				resultMap[key] = val
			} else {
				resultList = append(resultList, val)
			}
		}
		return NewForeachItfImpl(resultList, resultMap)
	case reflect.Map:
		resultList := make([]interface{}, 0, v.Len())
		resultMap := make(map[interface{}]interface{}, v.Len())
		for i, rfK := range v.MapKeys() {
			if !rfK.IsValid() || !rfK.CanInterface() {
				continue
			}
			mpV := v.MapIndex(rfK)
			if !mpV.IsValid() || !mpV.CanInterface() {
				continue
			}
			key, val := forFunc(i, rfK.Interface(), mpV.Interface())
			if key == nil && val == nil {
				continue
			}
			if key != nil {
				resultMap[key] = val
			} else {
				resultList = append(resultList, val)
			}
		}
		return NewForeachItfImpl(resultList, resultMap)
	}
	return NewExceptItfImpl()
}
