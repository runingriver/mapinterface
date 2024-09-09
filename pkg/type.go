package pkg

import "reflect"

type InterfaceType int

const (
	UnImplementInterface InterfaceType = 0
	ListInterface        InterfaceType = 1
	MapInterface         InterfaceType = 2
)

func GetInterfaceType(i interface{}) InterfaceType {
	k := reflect.ValueOf(i).Kind()
	switch k {
	case reflect.Map:
		return MapInterface
	case reflect.Array, reflect.Slice:
		return ListInterface
	}
	return UnImplementInterface
}

type ValType int

const (
	UnKnownType ValType = 0

	Str     ValType = 1
	Int     ValType = 2
	Int64   ValType = 3
	Int32   ValType = 4
	Uint    ValType = 5
	Uint64  ValType = 6
	Uint32  ValType = 7
	Float32 ValType = 8
	Float64 ValType = 9
	Bool    ValType = 10

	Map   ValType = 21
	Array ValType = 32
	Slice ValType = 23

	Struct ValType = 41
	Ptr    ValType = 42
	Func   ValType = 43

	MapStrItf     ValType = 61
	MapIntItf     ValType = 62
	MapInt32Itf   ValType = 63
	MapInt64Itf   ValType = 64
	MapUintItf    ValType = 65
	MapUint32Itf  ValType = 66
	MapUint64Itf  ValType = 67
	MapFloatItf   ValType = 68
	MapFloat64Itf ValType = 69
	MapStructItf  ValType = 70
	MapFuncItf    ValType = 71

	MapStrToStr         ValType = 72
	MapIntToInt         ValType = 73
	MapInt64ToInt64     ValType = 74
	MapFloat32ToFloat32 ValType = 75
	MapFloat64ToFloat64 ValType = 76

	ListItf     ValType = 81
	ListStr     ValType = 82
	ListInt     ValType = 83
	ListInt64   ValType = 84
	ListInt32   ValType = 85
	ListUint    ValType = 86
	ListUint32  ValType = 87
	ListUint64  ValType = 88
	ListFloat32 ValType = 89
	ListFloat64 ValType = 90
	ListBool    ValType = 91
)
