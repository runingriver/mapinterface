package mapitf

import (
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
)

type ExceptItf interface {
	api.MapInterface
}

// MapStrItfImpl 典型实现,整个通过接口交互,对外封闭,对内开放
type ExceptItfImpl struct {
	BaseItfImpl
}

func NewExceptItfImplErr(err itferr.MapItfErr) ExceptItf {
	return &ExceptItfImpl{BaseItfImpl{ItfErr: err}}
}

func NewExceptItfImpl() ExceptItf {
	return &ExceptItfImpl{
		BaseItfImpl{
			ItfErr: itferr.NewMapItfErr("ExceptItf", itferr.ExceptObject, "Unsupported Object", nil),
		},
	}
}

func (e *ExceptItfImpl) Get(key interface{}) api.MapInterface {
	e.ItfErr = itferr.NewUnSupportInterfaceFunc(fmt.Sprintf("ExceptItfImpl#Get(%v)", key))
	return e
}
