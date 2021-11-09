package mapitf

import (
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type OneLevelMapItf interface {
	api.MapInterface
}

type OneLevelMapImpl struct {
	BaseItfImpl

	OriginMap interface{}
}

func NewOneLevelMap(m interface{}) *OneLevelMapImpl {
	return &OneLevelMapImpl{
		OriginMap: m,
	}
}

func (o *OneLevelMapImpl) Get(key interface{}) api.MapInterface {
	if o.ItfErr != nil {
		return o
	}

	if itf, err := o.GetByInterface(key); err != nil {
		o.IterLevel++
		o.IterVal = itf
	}
	return o
}

func (o *OneLevelMapImpl) GetAny(keys ...interface{}) api.MapInterface {
	if len(keys) == 0 {
		return o
	}

	k := reflect.ValueOf(keys[0])
	if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
		o.ItfErr = itferr.NewFuncUsedErr("OneLevelMapImpl#GetAny", "param used err maybe GetAny(keys...)")
		return o
	}

	for _, k := range keys {
		if o.ItfErr != nil {
			return o
		}

		if itf, err := o.GetByInterface(k); err == nil && o.ItfErr == nil {
			o.IterLevel++
			o.IterVal = itf
		}
	}

	return o
}
