package mapitf

import (
	"fmt"
	"reflect"

	"github.com/runingriver/mapinterface/api"

	"github.com/runingriver/mapinterface/itferr"
)

type MapListItf interface {
	api.MapInterface
}

type MapListItfImpl struct {
	BaseItfImpl

	OriginListItf []interface{}
}

func NewMapListItfImpl(m interface{}) MapListItf {
	if listItf, ok := m.([]interface{}); ok {
		return &MapListItfImpl{
			BaseItfImpl: BaseItfImpl{
				IterVal:   listItf,
				IterLevel: 0,
				ItfErr:    nil,
			},
			OriginListItf: listItf,
		}
	}
	return &MapListItfImpl{
		BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("MapListItfImpl")},
	}
}

func (m *MapListItfImpl) Index(index int) api.MapInterface {
	if len(m.OriginListItf) < index {
		m.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("MapListItf#Index(%d)#(%d)", index, len(m.OriginListItf)))
		return NewExceptItfImplErr(m.ItfErr)
	}
	v := reflect.ValueOf(m.OriginListItf[index])
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return From(m.OriginListItf[index])
	}

	return NewBasicItfImpl(m.OriginListItf[index])
}
