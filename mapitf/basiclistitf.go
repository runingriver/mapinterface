package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type BasicListItf interface {
	api.MapInterface
}

type BasicListItfImpl struct {
	BaseItfImpl

	OriginItf interface{}
}

func NewBasicListItfImpl(m interface{}) MapListItf {
	return &BasicListItfImpl{
		BaseItfImpl: BaseItfImpl{
			IterVal:   m,
			IterLevel: 0,
			ItfErr:    nil,
		},
		OriginItf: m,
	}
}

func (m *BasicListItfImpl) Index(index int) api.MapInterface {
	v := reflect.ValueOf(m.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() < index {
			m.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("BasicListItfImpl#Index(%d)#(%d)", index, v.Len()))
			return NewExceptItfImplErr(m.ItfErr)
		}
		return m.indexing(index)
	default:
		m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#Index(%d)", index), "un-supported func")
		return NewExceptItfImplErr(m.ItfErr)
	}
}

func (m *BasicListItfImpl) Get(key interface{}) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#Get(%d)", key), "un-supported func")
	return m
}

func (m *BasicListItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#GetAny(%d)", keys), "un-supported func")
	return m
}

func (m *BasicListItfImpl) indexing(index int) api.MapInterface {
	switch vv := m.IterVal.(type) {
	case []string:
		return doStrFromStr(vv[index])
	case []int8:
		return NewBasicItfImpl(vv[index])
	case []int16:
		return NewBasicItfImpl(vv[index])
	case []int32:
		return NewBasicItfImpl(vv[index])
	case []int:
		return NewBasicItfImpl(vv[index])
	case []int64:
		return NewBasicItfImpl(vv[index])
	case []uint8:
		return NewBasicItfImpl(vv[index])
	case []uint16:
		return NewBasicItfImpl(vv[index])
	case []uint32:
		return NewBasicItfImpl(vv[index])
	case []uint:
		return NewBasicItfImpl(vv[index])
	case []uint64:
		return NewBasicItfImpl(vv[index])
	case []float32:
		return NewBasicItfImpl(vv[index])
	case []float64:
		return NewBasicItfImpl(vv[index])
	case []bool:
		return NewBasicItfImpl(vv[index])
	}
	return NewExceptItfImpl()
}
