package mapitf

import (
	"context"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
)

type DataType string

const (
	EmptyDataType = "empty_data"
	ListDataType  = "list"
	MapDataType   = "map"
)

type ForeachItf interface {
	api.MapInterface

	WithIterChain(iterChain *IterChain) ForeachItf
}

type ForeachItfImpl struct {
	BaseItfImpl

	DataType DataType
	ListItf  []interface{}
	MapItf   map[interface{}]interface{}
}

// NewForeachItfImpl 提供对map,list进行for循环的支持,返回的对象可以根据返回值定义为list或map.
func NewForeachItfImpl(ctx context.Context, listItf []interface{}, mapItf map[interface{}]interface{}) ForeachItf {
	fi := &ForeachItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:     ctx,
			Class:   "ForeachItf",
			IterVal: nil,
			ItfErr:  nil,
		},
		ListItf: listItf,
		MapItf:  mapItf,
	}
	if len(listItf) != 0 {
		fi.BaseItfImpl.IterVal = listItf
		fi.BaseItfImpl.IterChain = NewLinkedList(listItf)
		fi.DataType = ListDataType
	} else if len(mapItf) != 0 {
		fi.BaseItfImpl.IterVal = mapItf
		fi.BaseItfImpl.IterChain = NewLinkedList(mapItf)
		fi.DataType = MapDataType
	} else {
		fi.DataType = EmptyDataType
	}
	return fi
}

func (m *ForeachItfImpl) WithIterChain(iterChain *IterChain) ForeachItf {
	if iterChain == nil {
		return m
	}
	m.IterChain = iterChain
	return m
}

func (m *ForeachItfImpl) Index(index int) api.MapInterface {
	if m.DataType != ListDataType {
		m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("ForeachItfImpl#Index(%d)", index), "un-supported func")
		return m
	}

	if len(m.ListItf) > index {
		m.IterChain.PushBackByIdx(index, m.ListItf)
		return FrWithChain(m.Ctx, m.ListItf[index], m.IterChain)
	}

	m.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("ForeachItfImpl#Index(%d)#(%d)", index, len(m.ListItf)))
	return m
}

func (m *ForeachItfImpl) Get(key interface{}) api.MapInterface {
	if m.DataType != MapDataType {
		m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("ForeachItfImpl#Get(%v)", key), "un-supported func")
		return m
	}
	return FrWithChain(m.Ctx, m.MapItf, m.IterChain).Get(key)
}

func (m *ForeachItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if m.DataType != MapDataType {
		m.ItfErr = itferr.NewFuncUsedErr(fmt.Sprintf("ForeachItfImpl#GetAny(%v)", keys), "un-supported func")
		return m
	}
	return FrWithChain(m.Ctx, m.MapItf, m.IterChain).GetAny(keys...)
}

func (m *ForeachItfImpl) ToMap() (map[string]interface{}, error) {
	cvtResult := make(map[string]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		cvtResult[pkg.ToStr(k)] = v
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapInt() (map[int]interface{}, error) {
	cvtResult := make(map[int]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[int(iv)] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapInt64() (map[int64]interface{}, error) {
	cvtResult := make(map[int64]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[iv] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapInt32() (map[int32]interface{}, error) {
	cvtResult := make(map[int32]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[int32(iv)] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapUint() (map[uint]interface{}, error) {
	cvtResult := make(map[uint]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[uint(iv)] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapUint64() (map[uint64]interface{}, error) {
	cvtResult := make(map[uint64]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[uint64(iv)] = v
		}
	}
	return cvtResult, nil
}
func (m *ForeachItfImpl) ToMapUint32() (map[uint32]interface{}, error) {
	cvtResult := make(map[uint32]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToInt64(k); ok == nil {
			cvtResult[uint32(iv)] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapFloat32() (map[float32]interface{}, error) {
	cvtResult := make(map[float32]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToFloat32(k); ok == nil {
			cvtResult[iv] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapFloat64() (map[float64]interface{}, error) {
	cvtResult := make(map[float64]interface{}, len(m.MapItf))
	for k, v := range m.MapItf {
		if iv, ok := pkg.ToFloat64(k); ok == nil {
			cvtResult[iv] = v
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapItf() (map[interface{}]interface{}, error) {
	return m.MapItf, nil
}

func (m *ForeachItfImpl) ToMapStrToStr() (map[string]string, error) {
	cvtResult := make(map[string]string, len(m.MapItf))
	for k, v := range m.MapItf {
		cvtResult[pkg.ToStr(k)] = pkg.ToStr(v)
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapIntToInt() (map[int]int, error) {
	cvtResult := make(map[int]int, len(m.MapItf))
	for k, v := range m.MapItf {
		ik, ok1 := pkg.ToInt64(k)
		iv, ok2 := pkg.ToInt64(v)
		if ok1 == nil && ok2 == nil {
			cvtResult[int(ik)] = int(iv)
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapInt64ToInt64() (map[int64]int64, error) {
	cvtResult := make(map[int64]int64, len(m.MapItf))
	for k, v := range m.MapItf {
		ik, ok1 := pkg.ToInt64(k)
		iv, ok2 := pkg.ToInt64(v)
		if ok1 == nil && ok2 == nil {
			cvtResult[ik] = iv
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapFloat64ToFloat64() (map[float64]float64, error) {
	cvtResult := make(map[float64]float64, len(m.MapItf))
	for k, v := range m.MapItf {
		ik, ok1 := pkg.ToFloat64(k)
		iv, ok2 := pkg.ToFloat64(v)
		if ok1 == nil && ok2 == nil {
			cvtResult[ik] = iv
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToMapFloat32ToFloat32() (map[float32]float32, error) {
	cvtResult := make(map[float32]float32, len(m.MapItf))
	for k, v := range m.MapItf {
		ik, ok1 := pkg.ToFloat32(k)
		iv, ok2 := pkg.ToFloat32(v)
		if ok1 == nil && ok2 == nil {
			cvtResult[ik] = iv
		}
	}
	return cvtResult, nil
}

func (m *ForeachItfImpl) ToList() ([]interface{}, error) {
	return m.ListItf, nil
}
func (m *ForeachItfImpl) ToListMap() ([]map[string]interface{}, error) {
	return m.BaseItfImpl.ToListMap()
}
func (m *ForeachItfImpl) ToListStr() ([]string, error) {
	return m.BaseItfImpl.ToListStr()
}
func (m *ForeachItfImpl) ToListStrF() ([]string, error) {
	return m.BaseItfImpl.ToListStrF()
}
func (m *ForeachItfImpl) ToListInt() ([]int, error) {
	return m.BaseItfImpl.ToListInt()
}
func (m *ForeachItfImpl) ToListInt32() ([]int32, error) {
	return m.BaseItfImpl.ToListInt32()
}

func (m *ForeachItfImpl) ToListRune() ([]int32, error) {
	return m.BaseItfImpl.ToListRune()
}

func (m *ForeachItfImpl) ToListInt64() ([]int64, error) {
	return m.BaseItfImpl.ToListInt64()
}
func (m *ForeachItfImpl) ToListUint() ([]uint, error) {
	return m.BaseItfImpl.ToListUint()
}
func (m *ForeachItfImpl) ToListUint64() ([]uint64, error) {
	return m.BaseItfImpl.ToListUint64()
}
func (m *ForeachItfImpl) ToListUint32() ([]uint32, error) {
	return m.BaseItfImpl.ToListUint32()
}
func (m *ForeachItfImpl) ToListFloat32() ([]float32, error) {
	return m.BaseItfImpl.ToListFloat32()
}
func (m *ForeachItfImpl) ToListFloat64() ([]float64, error) {
	return m.BaseItfImpl.ToListFloat64()
}
func (m *ForeachItfImpl) ToListBool() ([]bool, error) {
	return m.BaseItfImpl.ToListBool()
}

func (m *ForeachItfImpl) ToStruct(out interface{}) (interface{}, error) {
	return m.BaseItfImpl.ToStruct(out)
}

func (m *ForeachItfImpl) IsStr() (bool, error) {
	return m.BaseItfImpl.IsStr()
}
func (m *ForeachItfImpl) IsDigit() (bool, error) {
	return m.BaseItfImpl.IsDigit()
}
func (m *ForeachItfImpl) IsList() (bool, error) {
	return m.BaseItfImpl.IsList()
}
func (m *ForeachItfImpl) IsStrList() (bool, error) {
	return m.BaseItfImpl.IsStrList()
}
func (m *ForeachItfImpl) IsDigitList() (bool, error) {
	return m.BaseItfImpl.IsDigitList()
}
func (m *ForeachItfImpl) IsMap() (bool, error) {
	return m.BaseItfImpl.IsMap()
}
func (m *ForeachItfImpl) IsMapStrItf() (bool, error) {
	return m.BaseItfImpl.IsMapStrItf()
}

func (m *ForeachItfImpl) Uniq() api.MapInterface {
	return m.BaseItfImpl.Uniq()
}

func (m *ForeachItfImpl) ForEach(forFunc api.ForFunc) api.MapInterface {
	m.ItfErr = itferr.NewMapItfErr("ForeachItfImpl#ForEach", itferr.UnSupportInterfaceFunc, "un-support ForEach to ForEach", nil)
	return m
}

func (m *ForeachItfImpl) SetMap(key interface{}, val interface{}) (orgVal interface{}, err error) {
	m.ItfErr = itferr.NewMapItfErr("ForeachItfImpl#SetMap", itferr.UnSupportInterfaceFunc, "un-support do ForEach then Set", nil)
	return nil, m.ItfErr
}

func (m *ForeachItfImpl) SetAsMap(key interface{}) (orgVal interface{}, err error) {
	m.ItfErr = itferr.NewMapItfErr("ForeachItfImpl#SetAsMap", itferr.UnSupportInterfaceFunc, "un-support do ForEach then Set", nil)
	return nil, m.ItfErr
}

func (m *ForeachItfImpl) New() api.MapInterface {
	return &ForeachItfImpl{
		BaseItfImpl: BaseItfImpl{
			Ctx:       m.Ctx,
			Class:     m.Class,
			ItfErr:    m.ItfErr,
			IterVal:   m.IterVal,
			IterChain: m.IterChain.Clone(),
		},
		DataType: m.DataType,
		ListItf:  m.ListItf,
		MapItf:   m.MapItf,
	}
}

func (m *ForeachItfImpl) PrintPath() string {
	return m.BaseItfImpl.PrintPath()
}
