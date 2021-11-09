package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapInt64Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapInt64Itf
}

type MapInt64ItfImpl struct {
	BaseItfImpl

	OriginMap map[int64]interface{}
}

func NewMapInt64ItfImpl(m interface{}) MapInt64Itf {
	if srcMap, ok := m.(map[int64]interface{}); ok {
		return &MapInt64ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapInt64Itf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapInt64ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapInt64ItfImpl")},
	}
}

func (m *MapInt64ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapInt64ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[k]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapInt64ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapInt64ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapInt64ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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

func (m *MapInt64ItfImpl) WithErr(err *itferr.MapItfError) MapInt64Itf {
	m.ItfErr = err
	return m
}
