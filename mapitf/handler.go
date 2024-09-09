package mapitf

import (
	"context"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

// -------------------Enter/入口-------------------------------------
// From 等价于Fr,入口方法
func From(itf interface{}) api.MapInterface {
	return Fr(context.TODO(), itf)
}

// Fr 等价于From,携带Context的版本
func Fr(ctx context.Context, itf interface{}) api.MapInterface {
	return FrWithChain(ctx, itf, nil)
}

func FrWithChain(ctx context.Context, itf interface{}, iterChain *IterChain) api.MapInterface {
	itf = pkg.Interpret(itf)

	switch vv := itf.(type) {
	case string:
		return doStrFromStr(ctx, vv, iterChain)
	case []byte:
		return doStrFromStr(ctx, pkg.ByteToStr(vv), iterChain)
	}

	v := reflect.ValueOf(itf)
	switch v.Kind() {
	case reflect.Map:
		return doForMap(ctx, itf, iterChain)
	case reflect.Slice, reflect.Array:
		return doForList(ctx, itf, iterChain)
	}

	return NewBasicItfImpl(ctx, itf).WithIterChain(iterChain)
}

func doStrFromStr(ctx context.Context, vv string, iterChain *IterChain) api.MapInterface {
	if mapStrItf, err := pkg.JsonLoadsMap(vv); err == nil {
		return NewMapStrItfImpl(ctx, mapStrItf).WithIterChain(iterChain)
	}

	if listItf, err := pkg.JsonLoadsList(vv); err == nil {
		return NewMapListItfImpl(ctx, listItf).WithIterChain(iterChain)
	}
	return NewBasicItfImpl(ctx, vv).WithIterChain(iterChain)
}

func doForMap(ctx context.Context, itf interface{}, iterChain *IterChain) api.MapInterface {
	switch vv := itf.(type) {
	case map[string]interface{}:
		return NewMapStrItfImpl(ctx, vv).WithIterChain(iterChain)
	default:
		return NewMapAnyToItfImpl(ctx, vv).WithIterChain(iterChain)
	}
}

func doForList(ctx context.Context, itf interface{}, iterChain *IterChain) api.MapInterface {
	v := pkg.ReflectToVal(itf)

	tt := v.Type().Elem().Kind()
	switch tt {
	case reflect.Interface, reflect.Map, reflect.String, reflect.Array, reflect.Slice, reflect.Ptr:
		return NewMapListItfImpl(ctx, itf).WithIterChain(iterChain)
	}

	return NewBasicListItfImpl(ctx, itf).WithIterChain(iterChain)
}
