package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapUintItf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapUintItf
}

type MapUintItfImpl struct {
	BaseItfImpl

	OriginMap map[uint]interface{}
}

func NewMapUintItfImpl(m interface{}) MapUintItf {
	if srcMap, ok := m.(map[uint]interface{}); ok {
		return &MapUintItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapUintItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapUintItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapUintItfImpl")},
	}
}

func (m *MapUintItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapUintItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[uint(k)]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapUintItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapUintItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapUintItfImpl#GetAny", "param used err maybe GetAny(keys...)")
		return m
	}

	for _, k := range keys {
		if m.ItfErr != nil {
			return m
		}

		// 如果出错已经在GetByInterface赋值了
		m.IterLevel++
		if itf, err := m.GetByInterface(k); err == nil && m.ItfErr == nil {
			m.IterVal = itf
		}
	}

	return m
}

func (m *MapUintItfImpl) WithErr(err *itferr.MapItfError) MapUintItf {
	m.ItfErr = err
	return m
}
