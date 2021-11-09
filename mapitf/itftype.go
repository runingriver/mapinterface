package mapitf

import (
	"github.com/runingriver/mapinterface/itferr"
)

type ItfType interface {
	ByMapStr() MapStrItf
	ByMapInt() MapIntItf
	ByMapInt32() MapInt32Itf
	ByMapInt64() MapInt64Itf
	ByMapUint() MapUintItf
	ByMapUint32() MapUint32Itf
	ByMapUint64() MapUint64Itf
	ByMapFloat32() MapFloat32Itf
	ByMapFloat64() MapFloat64Itf
	ByMapItfToItf() MapInterfaceToItf
	ByList() MapListItf
}

type ItfTypeImpl struct {
	Itf interface{}

	Err *itferr.MapItfError
}

func NewItfTypeImpl(itf interface{}) ItfType {
	return &ItfTypeImpl{Itf: itf}
}

func NewItfTypeImplErr(itf interface{}, err *itferr.MapItfError) ItfType {
	return &ItfTypeImpl{
		Itf: itf,
		Err: err,
	}
}

func (i *ItfTypeImpl) ByMapStr() MapStrItf {
	return NewMapStrItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapInt() MapIntItf {
	return NewMapIntItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapInt32() MapInt32Itf {
	return NewMapInt32ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapInt64() MapInt64Itf {
	return NewMapInt64ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapUint() MapUintItf {
	return NewMapUintItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapUint32() MapUint32Itf {
	return NewMapUint32ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapUint64() MapUint64Itf {
	return NewMapUint64ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapFloat32() MapFloat32Itf {
	return NewMapFloat32ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapFloat64() MapFloat64Itf {
	return NewMapFloat64ItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByMapItfToItf() MapInterfaceToItf {
	return NewMapInterfaceToItfImpl(i.Itf).WithErr(i.Err)
}

func (i *ItfTypeImpl) ByList() MapListItf {
	return NewMapListItfImpl(i.Itf)
}
