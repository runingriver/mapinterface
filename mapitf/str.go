package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
)

type MapStrItf interface {
	api.MapInterface

	GetByPath(keys ...string) MapStrItf

	WithErr(err *itferr.MapItfError) MapStrItf
}

// MapStrItfImpl 典型实现,整个通过接口交互,对外封闭,对内开放
type MapStrItfImpl struct {
	BaseItfImpl

	OriginMap map[string]interface{}
}

func NewMapStrItfImpl(m interface{}) MapStrItf {
	if srcMap, ok := m.(map[string]interface{}); ok {
		return &MapStrItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapStrItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapStrItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapStrInterface")},
	}
}

func (m *MapStrItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, ok := key.(string)
	if !ok {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapStrItfImpl#Get(%+v)", key))
	}

	iterMap := m.toMapStrInterface(m.IterVal)
	if iterMap == nil {
		m.ItfErr = itferr.NewConvFailed(fmt.Sprintf("MapStrItfImpl#Get(%s)", k))
		return m
	}

	m.IterLevel++
	if val, ok := iterMap[k]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapStrItfImpl#Get(%s)", k))
	}
	return m
}

// Get 类似于py中json_obj[key1][key2][[key3] 或 json_obj.get(key1, {}).get(key2, [])
func (m *MapStrItfImpl) GetByPath(keys ...string) MapStrItf {
	for _, key := range keys {
		if m.ItfErr != nil {
			return m
		}

		m.IterLevel++

		// 尝试转换成功map[string]interface
		if iterMap := m.toMapStrInterface(m.IterVal); iterMap != nil {
			if val, ok := iterMap[key]; ok {
				m.IterVal = val
			} else {
				m.ItfErr = itferr.NewKeyNotFoundFailed("MapStrItfImpl#GetByPath")
			}
			continue
		}

		// 尝试用反射获取
		vv := reflect.ValueOf(m.IterVal)
		if vv.Kind() != reflect.Map {
			m.ItfErr = itferr.NewMapItfErr("MapStrItfImpl#GetByPath", itferr.ExceptObject, "interface is not map", nil)
			continue
		}

		dstVal := vv.MapIndex(reflect.ValueOf(key))
		if dstVal.Kind() == reflect.Ptr || !dstVal.IsValid() || !dstVal.CanInterface() {
			m.ItfErr = itferr.NewValueTypeErr(fmt.Sprintf("MapStrItfImpl(%s)#GetByPath", key))
			continue
		}

		m.IterVal = dstVal.Interface()
	}

	return m
}

func (m *MapStrItfImpl) toMapStrInterface(v interface{}) map[string]interface{} {
	if mapStr, ok := v.(map[string]interface{}); ok {
		return mapStr
	}

	if mapItf, ok := v.(map[interface{}]interface{}); ok {
		result := make(map[string]interface{}, len(mapItf))
		for k, v := range mapItf {
			kStr := pkg.ToStr(k)
			result[kStr] = v
		}
		return result
	}

	return nil
}

func (m *MapStrItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	if keyStr, ok := m.isAllStrParam(keys...); ok {
		return m.GetByPath(keyStr...)
	}

	if k, ok := keys[0].([]string); ok {
		return m.GetByPath(k...)
	}

	// 检查一下用户容易填错参数
	k := reflect.ValueOf(keys[0])
	if len(keys) == 1 && k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapStrItfImpl#GetAny", "param used err maybe GetAny(keys...)")
		return m
	}

	for _, key := range keys {
		if m.ItfErr != nil {
			return m
		}

		// 如果出错已经在GetByInterface赋值了
		if itf, err := m.GetByInterface(key); err == nil && m.ItfErr == nil {
			m.IterLevel++
			m.IterVal = itf
		}
	}
	return m
}

func (m *MapStrItfImpl) isAllStrParam(keys ...interface{}) ([]string, bool) {
	keyStrList := make([]string, 0, len(keys))
	for _, key := range keys {
		if s, ok := key.(string); ok {
			keyStrList = append(keyStrList, s)
		} else {
			return keyStrList, false
		}
	}
	return keyStrList, true
}

func (m *MapStrItfImpl) WithErr(err *itferr.MapItfError) MapStrItf {
	m.ItfErr = err
	return m
}
