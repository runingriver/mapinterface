package itferr

import (
	"errors"
	"fmt"
)

type MapItfErr interface {
	error
	Is(error) bool
	Wrap(error) MapItfErr
	Code() MapItfErrorCode
	String() string
	IsErrEqual(err error) bool
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

func NewConvFailedX(locate string, msg string, err error) *MapItfError {
	return NewMapItfErr(locate, ValueConvertFailed, msg, err)
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

func NewUnSupportSetValErr(locate string, msg string, err error) *MapItfError {
	return NewMapItfErr(locate, UnSupportSetValTypeErr, msg, err)
}

func NewSetValueErr(locate string, msg string, err error) *MapItfError {
	return NewMapItfErr(locate, SetValueErr, msg, err)
}

func (mie *MapItfError) String() string {
	if mie.Err == nil && mie.ErrMsg == "" {
		return fmt.Sprintf("MapItfError{Location:%s,ErrCode:%d(%s)}", mie.Location, mie.ErrCode, mie.ErrCode.String())
	}
	if mie.ErrMsg != "" {
		return fmt.Sprintf("MapItfError{Location:%s,ErrCode:%d(%s),ErrMsg:%s}", mie.Location, mie.ErrCode, mie.ErrCode.String(), mie.ErrMsg)
	}
	return fmt.Sprintf("MapItfError{Location:%s,ErrCode:%d(%s),Err:%s,ErrMsg:%s}", mie.Location, mie.ErrCode, mie.ErrCode.String(), mie.Err.Error(), mie.ErrMsg)
}

func (mie *MapItfError) Error() string {
	return mie.String()
}

func (mie *MapItfError) Wrap(err error) MapItfErr {
	mie.Err = err
	return mie
}

func (mie *MapItfError) Is(err error) bool {
	return IsItfErr(err)
}

func (mie *MapItfError) Code() MapItfErrorCode {
	return mie.ErrCode
}

func (mie *MapItfError) IsErrEqual(err2 error) bool {
	return IsErrEqual(mie, err2)
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
	if _, ok := err.(*MapItfError); ok {
		return true
	}
	return false
}

func GetErrCode(err error) MapItfErrorCode {
	if err1, ok := err.(*MapItfError); ok {
		return err1.ErrCode
	}
	return UnknownErr
}
