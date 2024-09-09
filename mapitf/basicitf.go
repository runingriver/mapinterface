package mapitf

import (
	"context"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"reflect"
)

type BasicItf interface {
	api.MapInterface

	WithIterChain(iterChain *IterChain) BasicItf
}

type BasicItfImpl struct {
	BaseItfImpl

	OriginItf interface{}
}

// NewBasicItfImpl 非map,list类型的结构,通常用于对对象做类型转换;
func NewBasicItfImpl(ctx context.Context, m interface{}) BasicItf {
	return &BasicItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       ctx,
			Class:     "BasicItf",
			IterChain: NewLinkedList(m),
			IterVal:   m,
			ItfErr:    nil,
		},
		OriginItf: m,
	}
}

func (m *BasicItfImpl) Index(index int) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicItfImpl#Index(%d)", index), "un-supported func")
	return m
}

func (m *BasicItfImpl) Get(key interface{}) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicItfImpl#Get(%v)", key), "un-supported func")
	return m
}

func (m *BasicItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicItfImpl#GetAny(%+v)", keys), "un-supported func")
	return m
}

func (m *BasicItfImpl) WithIterChain(iterChain *IterChain) BasicItf {
	if iterChain == nil {
		return m
	}
	if e := iterChain.Back(); e != nil {
		ic := e.Value.(*IterCtx)
		// 当当前Iter值是json字符串转成的对象时,替换成具体对象
		valType := reflect.TypeOf(ic.Val).Kind()
		if valType == reflect.String && reflect.TypeOf(m.IterVal).Kind() != valType {
			iterChain.ReplaceBack(m.IterVal)
		}
	}

	m.IterChain = iterChain
	return m
}

func (m *BasicItfImpl) New() api.MapInterface {
	return &BasicItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       m.Ctx,
			Class:     m.Class,
			ItfErr:    m.ItfErr,
			IterVal:   m.IterVal,
			IterChain: m.IterChain.Clone(),
		},
		OriginItf: m.OriginItf,
	}
}

func (m *BasicItfImpl) OrgVal() (interface{}, error) {
	return m.OriginItf, nil
}
