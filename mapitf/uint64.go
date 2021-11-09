package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapUint64Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapUint64Itf
}

type MapUint64ItfImpl struct {
	BaseItfImpl

	OriginMap map[uint64]interface{}
}

func NewMapUint64ItfImpl(m interface{}) MapUint64Itf {
	if srcMap, ok := m.(map[uint64]interface{}); ok {
		return &MapUint64ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapUin64tItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapUint64ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapUint64ItfImpl")},
	}
}

func (m *MapUint64ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapUint64ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[uint64(k)]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapUint64ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapUint64ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapUint64ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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
func (m *MapUint64ItfImpl) WithErr(err *itferr.MapItfError) MapUint64Itf {
	m.ItfErr = err
	return m
}
