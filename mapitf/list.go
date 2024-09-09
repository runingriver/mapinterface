package mapitf

import (
	"context"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

type MapListItf interface {
	api.MapInterface

	WithIterChain(iterChain *IterChain) MapListItf
}

type MapListItfImpl struct {
	BaseItfImpl

	OriginVal interface{}
}

func NewMapListItfImpl(ctx context.Context, m interface{}) MapListItf {
	if isJson, js := pkg.JsonChecker(m); isJson {
		listItf, err := pkg.JsonLoadsList(js)
		if err != nil {
			return &MapListItfImpl{BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("MapListItfImpl")}}
		}
		m = listItf
	} else if rv := reflect.ValueOf(m); rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return &MapListItfImpl{BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("MapListItfImpl")}}
	}

	return &MapListItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       ctx,
			Class:     "MapListItf",
			IterChain: NewLinkedList(m),
			IterVal:   m,
			ItfErr:    nil,
		},
		OriginVal: m,
	}
}

func (m *MapListItfImpl) Index(index int) api.MapInterface {
	rv := reflect.ValueOf(m.IterVal)
	if rv.Len() <= index {
		m.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("MapListItf#Index(%d)#(%d)", index, rv.Len()))
		return NewExceptItfImplErr(m.ItfErr)
	}
	iv := rv.Index(index)
	if !iv.IsValid() || !iv.CanInterface() {
		return &MapListItfImpl{
			BaseItfImpl: BaseItfImpl{ItfErr: itferr.NewParamTypeErr("Index")},
		}
	}

	interfaceV := iv.Interface()

	m.IterVal = interfaceV
	m.IterChain.PushBackByIdx(index, m.IterVal)

	rfVV := pkg.ReflectToVal(m.IterVal)
	switch rfVV.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.String, reflect.Interface:
		return FrWithChain(m.Ctx, m.IterVal, m.IterChain)
	}

	return NewBasicItfImpl(m.Ctx, m.IterVal).WithIterChain(m.IterChain)
}

func (m *MapListItfImpl) WithIterChain(iterChain *IterChain) MapListItf {
	if iterChain == nil {
		return m
	}
	if e := iterChain.Back(); e != nil {
		ic := e.Value.(*IterCtx)
		valType := reflect.TypeOf(ic.Val).Kind()
		if valType == reflect.String && reflect.TypeOf(m.IterVal).Kind() != valType {
			iterChain.ReplaceBack(m.IterVal)
		}
	}

	m.IterChain = iterChain
	return m
}

func (m *MapListItfImpl) New() api.MapInterface {
	return &MapListItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       m.Ctx,
			Class:     m.Class,
			ItfErr:    m.ItfErr,
			IterVal:   m.IterVal,
			IterChain: m.IterChain.Clone(),
		},
		OriginVal: m.OriginVal,
	}
}

func (m *MapListItfImpl) OrgVal() (interface{}, error) {
	return m.OriginVal, nil
}
