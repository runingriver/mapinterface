package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/pkg"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type MapUint32Itf interface {
	api.MapInterface

	WithErr(err *itferr.MapItfError) MapUint32Itf
}

type MapUint32ItfImpl struct {
	BaseItfImpl

	OriginMap map[uint32]interface{}
}

func NewMapUint32ItfImpl(m interface{}) MapUint32Itf {
	if srcMap, ok := m.(map[uint32]interface{}); ok {
		return &MapUint32ItfImpl{
			BaseItfImpl: BaseItfImpl{
				Class:     "MapUin32tItf",
				IterVal:   srcMap,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginMap: srcMap,
		}
	}
	return &MapUint32ItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("NewMapUint32ItfImpl")},
	}
}

func (m *MapUint32ItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	k, err := pkg.ToInt64(key)
	if err != nil {
		m.ItfErr = itferr.NewGetFuncTypeInconsistent(fmt.Sprintf("MapUint32ItfImpl#Get(%+v)#err(%v)", key, err))
	}

	m.IterLevel++
	if val, ok := m.OriginMap[uint32(k)]; ok {
		m.IterVal = val
	} else {
		m.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("MapUint32ItfImpl#Get(%v)#%d", k, m.IterLevel))
	}

	return m
}

func (m *MapUint32ItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapUint32ItfImpl#GetAny", "param used err maybe GetAny(keys...)")
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
func (m *MapUint32ItfImpl) WithErr(err *itferr.MapItfError) MapUint32Itf {
	m.ItfErr = err
	return m
}
