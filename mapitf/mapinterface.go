package mapitf

import (
	"reflect"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/pkg"
)

// -------------------Enter/入口-------------------------------------
func From(itf interface{}) api.MapInterface {
	switch vv := itf.(type) {
	case string:
		return doStrFromStr(vv)
	case []byte:
		return doStrFromStr(pkg.ByteToStr(vv))
	}

	v := reflect.ValueOf(itf)
	switch v.Kind() {
	case reflect.Map:
		return doForMap(itf)
	case reflect.Slice, reflect.Array:
		return doForList(itf)
	}

	return NewBasicItfImpl(itf)
}

func doStrFromStr(vv string) api.MapInterface {
	if mapStrItf, err := pkg.JsonLoads(vv); err == nil {
		return NewItfTypeImpl(mapStrItf).ByMapStr()
	}

	if listItf, err := pkg.JsonLoadsList(vv); err == nil {
		return NewItfTypeImpl(listItf).ByList()
	}
	return NewBasicItfImpl(vv)
}

func doForMap(itf interface{}) api.MapInterface {
	switch vv := itf.(type) {
	case map[string]interface{}:
		return NewItfTypeImpl(vv).ByMapStr()
	case map[int]interface{}:
		return NewItfTypeImpl(vv).ByMapInt()
	case map[int32]interface{}:
		return NewItfTypeImpl(vv).ByMapInt32()
	case map[int64]interface{}:
		return NewItfTypeImpl(vv).ByMapInt64()
	case map[uint]interface{}:
		return NewItfTypeImpl(vv).ByMapUint()
	case map[uint32]interface{}:
		return NewItfTypeImpl(vv).ByMapUint32()
	case map[uint64]interface{}:
		return NewItfTypeImpl(vv).ByMapUint64()
	case map[float32]interface{}:
		return NewItfTypeImpl(vv).ByMapFloat32()
	case map[float64]interface{}:
		return NewItfTypeImpl(vv).ByMapFloat64()
	case map[interface{}]interface{}:
		return NewItfTypeImpl(vv).ByMapItfToItf()
	}

	return NewOneLevelMap(itf)
}

func doForList(itf interface{}) api.MapInterface {
	switch vv := itf.(type) {
	case []interface{}:
		return NewItfTypeImpl(vv).ByList()
	}

	return NewBasicListItfImpl(itf)
}
