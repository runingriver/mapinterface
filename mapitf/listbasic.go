package mapitf

import (
	"context"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/logx"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

type BasicListItf interface {
	api.MapInterface

	WithIterChain(iterChain *IterChain) BasicListItf
}

type BasicListItfImpl struct {
	BaseItfImpl

	OriginVal interface{}
}

// NewBasicListItfImpl 表示一个基础的list类型,val是基础类型,如:[]int.
func NewBasicListItfImpl(ctx context.Context, m interface{}) BasicListItf {
	return &BasicListItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       ctx,
			Class:     "BasicListItf",
			IterChain: NewLinkedList(m),
			IterVal:   m,
			ItfErr:    nil,
		},
		OriginVal: m,
	}
}

func (m *BasicListItfImpl) Index(index int) api.MapInterface {
	if isJson, js := pkg.JsonChecker(m.IterVal); isJson {
		listItf, err := pkg.JsonLoadsList(js)
		if err != nil {
			m.ItfErr = itferr.NewConvFailedX(fmt.Sprintf("BasicListItfImpl#Index(%d)", index), "not json list str", err)
			return NewExceptItfImplErr(m.ItfErr)
		}
		m.IterVal = listItf
		m.IterChain.ReplaceBack(m.IterVal)
		return m.indexing(index)
	}

	v := pkg.ReflectToVal(m.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() <= index {
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

func (m *BasicListItfImpl) SetMap(key interface{}, val interface{}) (orgVal interface{}, err error) {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#Get(%d)", key), "un-supported func")
	return nil, m.ItfErr
}

func (m *BasicListItfImpl) SetAsMap(key interface{}) (orgVal interface{}, err error) {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#Get(%d)", key), "un-supported func")
	return nil, m.ItfErr
}
func (m *BasicListItfImpl) Exist(key interface{}) (interface{}, bool) {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicListItfImpl#Get(%d)", key), "un-supported func")
	return nil, false
}
func (m *BasicListItfImpl) indexing(index int) api.MapInterface {
	switch vv := m.IterVal.(type) {
	case []string:
		if isJson, jsonStr := pkg.JsonChecker(vv[index]); isJson {
			var err error
			if mapStrItf, err := pkg.JsonLoadsMap(jsonStr); err == nil {
				m.IterChain.PushBackByIdx(index, mapStrItf)
				return NewMapStrItfImpl(m.Ctx, mapStrItf).WithErr(m.ItfErr).WithIterChain(m.IterChain)
			}

			if listItf, err := pkg.JsonLoadsList(jsonStr); err == nil {
				m.IterChain.PushBackByIdx(index, listItf)
				return NewMapListItfImpl(m.Ctx, listItf).WithIterChain(m.IterChain)
			}
			logx.CtxWarn(m.Ctx, "[BasicListItfImpl] index string err of un-loads json,err:%v", err)
		}

		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []interface{}:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []int8:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []int16:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []int32:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []int:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []int64:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []uint8:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []uint16:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []uint32:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []uint:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []uint64:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []float32:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []float64:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	case []bool:
		m.IterChain.PushBackByIdx(index, vv[index])
		return NewBasicItfImpl(m.Ctx, vv[index]).WithIterChain(m.IterChain)
	}
	return NewExceptItfImpl()
}

func (m *BasicListItfImpl) WithIterChain(iterChain *IterChain) BasicListItf {
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

func (m *BasicListItfImpl) New() api.MapInterface {
	return &BasicListItfImpl{
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

func (m *BasicListItfImpl) OrgVal() (interface{}, error) {
	return m.OriginVal, nil
}
