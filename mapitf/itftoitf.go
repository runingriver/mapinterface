package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapInterfaceToItf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapInterfaceToItf
}

type MapInterfaceToItfImpl struct {
	BaseItfImpl

	OriginMap map[interface{}]interface{}
}

func NewMapInterfaceToItfImpl(m interface{}) MapInterfaceToItf {
	if srcMap, ok := m.(map[interface{}]interface{}); ok {
		return &MapInterfaceToItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapInterfaceToItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapInterfaceToItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("MapInterfaceToItfImpl")},
	}
}

func (m *MapInterfaceToItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	m.IterLevel++
	if val, ok := m.OriginMap[key]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapInt32ItfImpl#Get(%v)#%d", key, m.IterLevel))
	}
	return m
}

func (m *MapInterfaceToItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapInterfaceToItfImpl#GetAny", "param used err maybe GetAny(keys...)")
		return m
	}

	for _, k := range keys {
		if m.ItfErr != nil {
			return m
		}

		if itf, err := m.GetByInterface(k); err == nil && m.ItfErr == nil {
			m.IterLevel++
			m.IterVal = itf
		}
	}

	return m
}

func (m *MapInterfaceToItfImpl) WithErr(err *itferr.MapItfError) MapInterfaceToItf {
	m.ItfErr = err
	return m
}
