package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapIntItf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapIntItf
}

type MapIntItfImpl struct {
	BaseItfImpl

	OriginMap map[int]interface{}
}

func NewMapIntItfImpl(m interface{}) MapIntItf {
	if srcMap, ok := m.(map[int]interface{}); ok {
		return &MapIntItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapIntItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapIntItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapIntItfImpl")},
	}
}

func (m *MapIntItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapIntItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[int(k)]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapIntItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapIntItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapIntItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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

func (m *MapIntItfImpl) WithErr(err *itferr.MapItfError) MapIntItf {
	m.ItfErr = err
	return m
}
