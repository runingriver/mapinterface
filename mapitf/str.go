package mapitf

import (
	"context"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

type MapStrItf interface {
	api.MapInterface

	GetByPath(keys ...string) MapStrItf

	WithErr(err itferr.MapItfErr) MapStrItf

	WithIterChain(iterChain *IterChain) MapStrItf
}

// MapStrItfImpl 典型实现,整个通过接口交互,对外封闭,对内开放
type MapStrItfImpl struct {
	BaseItfImpl

	OriginVal map[string]interface{}
}

func NewMapStrItfImpl(ctx context.Context, m interface{}) MapStrItf {
	if srcMap, ok := m.(map[string]interface{}); ok {
		return &MapStrItfImpl{
			BaseItfImpl: BaseItfImpl{
				Ctx:       ctx,
				Class:     "MapStrItf",
				IterChain: NewLinkedList(srcMap),
				IterVal:   srcMap,
				ItfErr:    nil,
			},
			OriginVal: srcMap,
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

	if k, ok := key.(string); ok {
		return m.GetOne(k)
	}

	return m.GetAny(key)
}

// Get 类似于py中json_obj[key1][key2][[key3] 或 json_obj.get(key1, {}).get(key2, [])
func (m *MapStrItfImpl) GetByPath(keys ...string) (msi MapStrItf) {
	defer func() {
		if err := recover(); err != nil {
			m.ItfErr = itferr.NewMapItfErrX("MapStrItfImpl#GetByPath", itferr.UnrecoverablePanicErr)
			msi = m
		}
	}()
	for _, key := range keys {
		if m.ItfErr != nil {
			return m
		}

		m.GetOne(key)
	}

	return m
}

func (m *MapStrItfImpl) GetOne(key string) api.MapInterface {
	// 尝试转换成功map[string]interface
	if iterMap := m.toMapStrInterface(m.IterVal); iterMap != nil {
		if isStr, _ := pkg.IsStrType(m.IterVal); isStr {
			m.IterChain.ReplaceBack(iterMap)
		}
		if val, ok := iterMap[key]; ok {
			m.IterVal = val
			m.IterChain.PushBackByKey(key, val)
		} else {
			m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapStrItfImpl#Get(%s)", key))
		}
		return m
	}

	// 尝试用反射获取
	vv := pkg.ReflectToVal(m.IterVal)
	if vv.Kind() != reflect.Map {
		m.ItfErr = itferr.NewMapItfErr("MapStrItfImpl#GetByPath", itferr.ExceptObject, "interface is not map", nil)
		return m
	}

	dstVal := vv.MapIndex(reflect.ValueOf(key))
	if dstVal.Kind() == reflect.Ptr || !dstVal.IsValid() || !dstVal.CanInterface() {
		m.ItfErr = itferr.NewValueTypeErr(fmt.Sprintf("MapStrItfImpl(%s)#GetByPath", key))
		return m
	}

	m.IterVal = dstVal.Interface()
	m.IterChain.PushBackByKey(key, m.IterVal)
	return m
}

func (m *MapStrItfImpl) toMapStrInterface(v interface{}) map[string]interface{} {
	if mapStr, ok := v.(map[string]interface{}); ok {
		return mapStr
	}

	if isJson, js := pkg.JsonChecker(v); isJson {
		if result, err := pkg.JsonLoadsMap(js); err == nil {
			return result
		}
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
		if srcVal, itf, err := m.GetByInterface(key); err == nil && m.ItfErr == nil {
			if isStr, _ := pkg.IsStrType(m.IterVal); isStr {
				m.IterChain.ReplaceBack(srcVal)
			}
			m.IterVal = itf
			m.IterChain.PushBackByKey(key, itf)
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

func (m *MapStrItfImpl) New() api.MapInterface {
	return &MapStrItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       m.Ctx,
			Class:     m.Class,
			ItfErr:    m.ItfErr,
			IterVal:   m.IterVal,
			IterChain: m.IterChain.Clone(),
		},
		OriginVal: m.OriginVal,
	}
}

func (m *MapStrItfImpl) WithErr(err itferr.MapItfErr) MapStrItf {
	m.ItfErr = err
	return m
}

func (m *MapStrItfImpl) WithIterChain(iterChain *IterChain) MapStrItf {
	if iterChain == nil {
		return m
	}
	if e := iterChain.Back(); e != nil {
		ic := e.Value.(*IterCtx)
		valType := reflect.TypeOf(ic.Val).Kind()
		if valType == reflect.String && reflect.TypeOf(m.IterVal).Kind() != valType {
			iterChain.ReplaceBack(m.IterVal)
		}
	}

	m.IterChain = iterChain
	return m
}

func (m *MapStrItfImpl) OrgVal() (interface{}, error) {
	return m.OriginVal, nil
}
