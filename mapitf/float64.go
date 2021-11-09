package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapFloat64Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapFloat64Itf
}

type MapFloat64ItfImpl struct {
	BaseItfImpl

	OriginMap map[float64]interface{}
}

func NewMapFloat64ItfImpl(m interface{}) MapFloat64Itf {
	if srcMap, ok := m.(map[float64]interface{}); ok {
		return &MapFloat64ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapFloat64Itf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapFloat64ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapFloat64ItfImpl")},
	}
}

func (m *MapFloat64ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToFloat64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapFloat64ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[k]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapFloat64ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapFloat64ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapFloat64ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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

func (m *MapFloat64ItfImpl) WithErr(err *itferr.MapItfError) MapFloat64Itf {
	m.ItfErr = err
	return m
}
