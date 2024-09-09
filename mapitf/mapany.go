package mapitf

import (
	"context"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

type MapAnyToItf interface {
	api.MapInterface

	WithErr(err itferr.MapItfErr) MapAnyToItf
	WithIterChain(iterChain *IterChain) MapAnyToItf
}

type MapAnyToItfImpl struct {
	BaseItfImpl

	OriginVal interface{}
}

// NewMapAnyToItfImpl 支持任意map的取数操作
func NewMapAnyToItfImpl(ctx context.Context, m interface{}) MapAnyToItf {
	return &MapAnyToItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       ctx,
			Class:     "MapAnyToItf",
			IterChain: NewLinkedList(m),
			IterVal:   m,
			ItfErr:    nil,
		},
		OriginVal: m,
	}
}

func (m *MapAnyToItfImpl) Get(key interface{}) api.MapInterface {
	if m.ItfErr != nil {
		return m
	}

	if srcVal, itf, err := m.GetByInterface(key); err == nil && m.ItfErr == nil {
		if isStr, _ := pkg.IsStrType(m.IterVal); isStr {
			m.IterChain.ReplaceBack(srcVal)
		}

		m.IterVal = itf
		m.IterChain.PushBackByKey(key, itf)
	}

	return m
}

func (m *MapAnyToItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return m
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		m.ItfErr = itferr.NewFuncUsedErr("MapAnyToItfImpl#GetAny", "param used err maybe GetAny(keys...)")
		return m
	}

	for _, key := range keys {
		if m.ItfErr != nil {
			return m
		}

		if srcVal, itf, err := m.GetByInterface(key); err == nil && m.ItfErr == nil {
			if isStr, _ := pkg.IsStrType(m.IterVal); isStr {
				m.IterChain.ReplaceBack(srcVal)
			}

			m.IterVal = itf
			m.IterChain.PushBackByKey(key, itf)
		}
	}

	return m
}

func (m *MapAnyToItfImpl) New() api.MapInterface {
	return &MapAnyToItfImpl{
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

func (m *MapAnyToItfImpl) WithErr(err itferr.MapItfErr) MapAnyToItf {
	m.ItfErr = err
	return m
}

func (m *MapAnyToItfImpl) WithIterChain(iterChain *IterChain) MapAnyToItf {
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

func (m *MapAnyToItfImpl) OrgVal() (interface{}, error) {
	return m.OriginVal, nil
}
