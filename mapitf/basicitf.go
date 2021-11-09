package mapitf

import (
	"fmt"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type BasicItf interface {
	api.MapInterface
}

type BasicItfImpl struct {
	BaseItfImpl

	OriginItf interface{}
}

func NewBasicItfImpl(m interface{}) MapListItf {
	return &BasicItfImpl{
		BaseItfImpl: BaseItfImpl{
			IterVal:   m,
			IterLevel: 0,
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
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicItfImpl#Get(%d)", key), "un-supported func")
	return m
}

func (m *BasicItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("BasicItfImpl#GetAny(%d)", keys), "un-supported func")
	return m
}
