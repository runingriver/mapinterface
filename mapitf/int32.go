package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapInt32Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapInt32Itf
}

type MapInt32ItfImpl struct {
	BaseItfImpl

	OriginMap map[int32]interface{}
}

func NewMapInt32ItfImpl(m interface{}) MapInt32Itf {
	if srcMap, ok := m.(map[int32]interface{}); ok {
		return &MapInt32ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapInt32Itf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapInt32ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapInt32ItfImpl")},
	}
}

func (m *MapInt32ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapInt32ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[int32(k)]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapInt32ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapInt32ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapInt32ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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

func (m *MapInt32ItfImpl) WithErr(err *itferr.MapItfError) MapInt32Itf {
	m.ItfErr = err
	return m
}
