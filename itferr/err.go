package itferr

import (
	"errors"
	"fmt"
)

type MapItfErrorCode int

const (
	UnknownErr       MapItfErrorCode = -1
	InitParseFailed  MapItfErrorCode = 1001
	InitParamTypeErr MapItfErrorCode = 2001
	ExceptObject     MapItfErrorCode = 2002

	KeyTypeErr              MapItfErrorCode = 3001
	ValueTypeErr            MapItfErrorCode = 3002
	ValueConvertFailed      MapItfErrorCode = 3003
	BaseTypeConvertFailed   MapItfErrorCode = 3004
	KeyNotFound             MapItfErrorCode = 3005
	GetFuncTypeInconsistent MapItfErrorCode = 3006
	IllegalMapObject        MapItfErrorCode = 3007

	ListIndexIllegal MapItfErrorCode = 4001

	UnSupportInterfaceFunc MapItfErrorCode = 5001
	CurrentCannotUseIndex  MapItfErrorCode = 5002
	TypeMismatchErr        MapItfErrorCode = 5003
	FuncUsedErr            MapItfErrorCode = 5004
	UnrecoverableErr       MapItfErrorCode = 5005
)

func (err MapItfErrorCode) String() string {
	switch err {
	case InitParseFailed:
		return "Init_Parse_Failed"
	case InitParamTypeErr:
		return "InitParamTypeErr"
	case ExceptObject:
		return "ExceptObject"
	case KeyTypeErr:
		return "KeyTypeErr"
	case ValueTypeErr:
		return "ValueTypeErr"
	case ValueConvertFailed:
		return "ValueConvertFailed"
	case BaseTypeConvertFailed:
		return "BaseTypeConvertFailed"
	case KeyNotFound:
		return "KeyNotFound"
	case GetFuncTypeInconsistent:
		return "GetFuncTypeInconsistent"
	case ListIndexIllegal:
		return "ListIndexIllegal"
	case UnSupportInterfaceFunc:
		return "UnSupportInterfaceFunc"
	case CurrentCannotUseIndex:
		return "CurrentCannotUseIndex"
	case TypeMismatchErr:
		return "TypeMismatchErr"
	case FuncUsedErr:
		return "FuncUsedErr"
	case UnrecoverableErr:
		return "UnrecoverableErr"
	default:
		return "UnDefine_Err"
	}
}

type MapItfError struct {
	ErrCode  MapItfErrorCode
	Location string
	Err      error
	ErrMsg   string
}

// NewMapItfErr method,code必填,msg,err选填
func NewMapItfErr(locate string, code MapItfErrorCode, msg string, err error) *MapItfError {
	return &MapItfError{
		Location: locate,
		ErrCode:  code,
		Err:      err,
		ErrMsg:   msg,
	}
}

// NewMapItfErrX 精简版
func NewMapItfErrX(locate string, code MapItfErrorCode) *MapItfError {
	return &MapItfError{
		Location: locate,
		ErrCode:  code,
	}
}

func NewParamTypeErr(locate string) *MapItfError {
	return NewMapItfErr(locate, InitParamTypeErr, "", nil)
}

func NewListIndexIllegal(locate string) *MapItfError {
	return NewMapItfErr(locate, ListIndexIllegal, "", nil)
}

func NewKeyTypeErr(locate string) *MapItfError {
	return NewMapItfErr(locate, KeyTypeErr, "", nil)
}

func NewValueTypeErr(locate string) *MapItfError {
	return NewMapItfErr(locate, ValueTypeErr, "", nil)
}

func NewKeyNotFoundFailed(locate string) *MapItfError {
	return NewMapItfErr(locate, KeyNotFound, "", nil)
}

func NewConvFailed(locate string) *MapItfError {
	return NewMapItfErr(locate, ValueConvertFailed, "", nil)
}

func NewUnSupportInterfaceFunc(locate string) *MapItfError {
	return NewMapItfErr(locate, UnSupportInterfaceFunc, "", nil)
}

func NewGetFuncTypeInconsistent(locate string) *MapItfError {
	return NewMapItfErr(locate, GetFuncTypeInconsistent, "", nil)
}

func NewTypeMismatchErr(locate string) *MapItfError {
	return NewMapItfErr(locate, TypeMismatchErr, "", nil)
}

func NewFuncUsedErr(locate, msg string) *MapItfError {
	return NewMapItfErr(locate, FuncUsedErr, msg, nil)
}

func NewIllegalMapObject(locate string) *MapItfError {
	return NewMapItfErr(locate, IllegalMapObject, "", nil)
}

func NewCurrentCannotUseIndex(locate string) *MapItfError {
	return NewMapItfErr(locate, CurrentCannotUseIndex, "", nil)
}

func NewBaseTypeConvErr(locate string, msg string, err error) *MapItfError {
	return NewMapItfErr(locate, BaseTypeConvertFailed, msg, err)
}

func (mie MapItfError) String() string {
	if mie.Err == nil && mie.ErrMsg == "" {
		return fmt.Sprintf("%s{Location:%s,ErrCode:%d}", mie.ErrCode, mie.Location, mie.ErrCode)
	}
	if mie.ErrMsg != "" {
		return fmt.Sprintf("MapItfError{Location:%s,ErrCode:%d,ErrMsg:%s}", mie.Location, mie.ErrCode, mie.ErrMsg)
	}
	return fmt.Sprintf("MapItfError{Location:%s,ErrCode:%d,Err:%v,ErrMsg:%s}", mie.Location, mie.ErrCode, mie.Err, mie.ErrMsg)
}

func (mie MapItfError) Error() string {
	return mie.String()
}

func IsErrEqual(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return false
	}

	if errors.Is(err1, err2) {
		return true
	}

	err1Code := GetErrCode(err1)
	err2Code := GetErrCode(err2)
	if err1Code == err2Code && err1Code != -1 {
		return true
	}

	return false
}

func IsItfErr(err error) bool {
	if _, ok := err.(MapItfError); ok {
		return true
	}

	return false
}

func GetErrCode(err error) MapItfErrorCode {
	if err1, ok := err.(MapItfError); ok {
		return err1.ErrCode
	}

	return UnknownErr
}
