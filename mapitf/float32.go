package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
)

type MapFloat32Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapFloat32Itf
}

type MapFloat32ItfImpl struct {
	BaseItfImpl

	OriginMap map[float32]interface{}
}

func NewMapFloat32ItfImpl(m interface{}) MapFloat32Itf {
	if srcMap, ok := m.(map[float32]interface{}); ok {
		return &MapFloat32ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapFloat32Itf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapFloat32ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapFloat32ItfImpl")},
	}
}

func (m *MapFloat32ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToFloat32(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapFloat32ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[k]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapFloat32ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapFloat32ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if len(keys) == 1 && k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapFloat32ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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

func (m *MapFloat32ItfImpl) WithErr(err *itferr.MapItfError) MapFloat32Itf {
	m.ItfErr = err
	return m
}
