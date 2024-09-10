package mapitf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/conf"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/logx"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
	"strings"
)

type BaseItfImpl struct {
	Ctx   context.Context
	Class string

	ItfErr itferr.MapItfErr

	IterVal interface{}

	IterChain *IterChain
}

func (b *BaseItfImpl) Valid() bool {
	if b.ItfErr != nil || b.IterVal == nil {
		return false
	}

	return true
}

func (b *BaseItfImpl) Exist(key interface{}) (interface{}, bool) {
	if b.ItfErr != nil || b.IterVal == nil {
		return nil, false
	}

	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		jsonMap, err := pkg.JsonLoadsMap(js)
		if err != nil {
			return nil, false
		}
		if vv, ok := jsonMap[pkg.ToStr(key)]; ok {
			return vv, true
		}
		return nil, false
	}

	if _, val, err := b.GetByInterface(key); err == nil && val != nil {
		return val, true
	}

	return nil, false
}

// GetByInterface 如果key的类型和Map[k]v中k类型不一致会panic
func (b *BaseItfImpl) GetByInterface(key interface{}) (srcVal, itfVal interface{}, err itferr.MapItfErr) {
	defer func() {
		if err := recover(); err != nil {
			b.ItfErr = itferr.NewMapItfErrX(fmt.Sprintf("GetByInterface(%v)", key), itferr.UnrecoverablePanicErr)
		}
	}()

	if b.ItfErr != nil {
		return b.IterVal, nil, b.ItfErr
	}

	if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if strMapItf, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			val, ok := strMapItf[pkg.ToStr(key)]
			if !ok {
				b.ItfErr = itferr.NewKeyNotFoundFailed(fmt.Sprintf("GetByInterface#Get(%s)", key))
			}
			return strMapItf, val, b.ItfErr
		}
	}

	v := pkg.ReflectToVal(b.IterVal)
	if v.Kind() != reflect.Map {
		b.ItfErr = itferr.NewValueTypeErr(fmt.Sprintf("%s#GetByInterface(%+v)", b.Class, key))
		return b.IterVal, nil, b.ItfErr
	}

	keys := v.MapKeys()
	if len(keys) == 0 {
		b.ItfErr = itferr.NewMapItfErr(fmt.Sprintf("%s#GetByInterface(%s)#%+v", b.Class, key, b.IterVal), itferr.EmptyMapObject, "", nil)
		return b.IterVal, nil, b.ItfErr
	}

	keyV := reflect.ValueOf(key)
	if keys[0].Kind() != keyV.Kind() {
		for _, keyVal := range keys {
			if keyItf, itfError := b.toInterface(keyVal); itfError == nil {
				if pkg.ToStr(keyItf) == pkg.ToStr(key) {
					dstVal := v.MapIndex(keyVal)
					itfVal, b.ItfErr = b.toInterface(dstVal)
					return b.IterVal, itfVal, b.ItfErr
				}
			}
		}

		b.ItfErr = itferr.NewTypeMismatchErr(fmt.Sprintf("%s#GetByInterface(%s)%v:%v", b.Class, key, keys[0].Kind(), keyV.Kind()))
		return b.IterVal, nil, b.ItfErr
	}

	itfVal, b.ItfErr = b.toInterface(v.MapIndex(keyV))
	return b.IterVal, itfVal, b.ItfErr
}

func (b *BaseItfImpl) toInterface(dstVal reflect.Value) (interface{}, itferr.MapItfErr) {
	if !dstVal.IsValid() || !dstVal.CanInterface() {
		return nil, itferr.NewValueTypeErr(fmt.Sprintf("%s#GetByInterface():%v", b.Class, dstVal.Kind()))
	}

	return dstVal.Interface(), nil
}

func (b *BaseItfImpl) Val() (interface{}, error) {
	if b.ItfErr != nil {
		return "", b.ItfErr
	}
	return b.IterVal, nil
}

func (b *BaseItfImpl) ToStr() (string, error) {
	if b.ItfErr != nil {
		return "", b.ItfErr
	}
	return pkg.ToStr(b.IterVal), nil
}

func (b *BaseItfImpl) ToByte() ([]byte, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if s, ok := b.ToStr(); ok == nil {
		return pkg.StrToByte(s), nil
	}

	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToByte", b.Class))
}

func (b *BaseItfImpl) ToInt() (int, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return int(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt", b.Class))
}

func (b *BaseItfImpl) ToInt64() (int64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int64); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return v, nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt64", b.Class))
}

func (b *BaseItfImpl) ToInt32() (int32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(int32); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return int32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToInt32", b.Class))
}

func (b *BaseItfImpl) ToRune() (rune, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(rune); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return rune(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToRune", b.Class))
}

func (b *BaseItfImpl) ToUint() (uint, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint", b.Class))
}

func (b *BaseItfImpl) ToUint64() (uint64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint64); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint64(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint64", b.Class))
}

func (b *BaseItfImpl) ToUint32() (uint32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(uint32); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return uint32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToUint32", b.Class))
}

func (b *BaseItfImpl) ToFloat32() (float32, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(float32); ok {
		return v, nil
	}

	if v, ok := pkg.ToFloat64(b.IterVal); ok == nil {
		return float32(v), nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToFloat", b.Class))
}

func (b *BaseItfImpl) ToFloat64() (float64, error) {
	if b.ItfErr != nil {
		return 0, b.ItfErr
	}

	if v, ok := b.IterVal.(float64); ok {
		return v, nil
	}

	if v, ok := pkg.ToFloat64(b.IterVal); ok == nil {
		return v, nil
	}

	return 0, itferr.NewConvFailed(fmt.Sprintf("%s#ToFloat64", b.Class))
}

func (b *BaseItfImpl) ToBool() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}

	if v, ok := b.IterVal.(bool); ok {
		return v, nil
	}

	if v, ok := pkg.ToInt64(b.IterVal); ok == nil {
		return v == 1, nil
	}

	v := strings.TrimSpace(pkg.ToStr(b.IterVal))
	if v == "true" {
		return true, nil
	}
	if v == "false" {
		return false, nil
	}

	return false, itferr.NewConvFailed(fmt.Sprintf("%s#ToBool", b.Class))
}

// ToMapType ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) ToMap() (map[string]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	result, err := toMap(b.Ctx, b.IterVal)
	if err == nil {
		return result, nil
	}

	return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMap", b.Class), "", err)
}

func toMap(ctx context.Context, m interface{}) (map[string]interface{}, error) {
	if v, ok := m.(map[string]interface{}); ok {
		return v, nil
	}

	if ok, jsonStr := pkg.JsonChecker(m); ok {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			return rstMap, nil
		} else {
			return nil, fmt.Errorf("convert map type un-match:%w", err)
		}
	}

	if v, ok := m.(map[interface{}]interface{}); ok {
		result := make(map[string]interface{}, len(v))
		for k, v := range v {
			result[pkg.ToStr(k)] = v
		}
		return result, nil
	}

	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil, errors.New("type is not map[string]interface{}")
	}
	srcLen := v.Len()
	result := make(map[string]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}
		result[pkg.ToStr(rfK.Interface())] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, errors.New("convert to map exist failed convert item")
	}
	if len(result) != srcLen {
		logx.CtxWarn(ctx, "ToMap %d key value convert failed", srcLen-len(result))
	}

	return result, nil
}

func (b *BaseItfImpl) ToMapInt() (map[int]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[int]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[int(cnvK)] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt", b.Class), "map key not int type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt", b.Class))
	}
	srcLen := v.Len()
	result := make(map[int]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapInt convert err:%v", cnvKErr)
			continue
		}
		result[int(cnvK)] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapInt %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapInt64() (map[int64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int64]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt64", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[int64]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[cnvK] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt64", b.Class), "map key not int64 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64", b.Class))
	}
	srcLen := v.Len()
	result := make(map[int64]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapInt64 convert err:%v", cnvKErr)
			continue
		}
		result[cnvK] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapInt64 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapInt32() (map[int32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if v, ok := b.IterVal.(map[int32]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt32", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[int32]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[int32(cnvK)] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt32", b.Class), "map key not int32 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt32", b.Class))
	}
	srcLen := v.Len()
	result := make(map[int32]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapInt32 convert err:%v", cnvKErr)
			continue
		}
		result[int32(cnvK)] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt32", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapInt32 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapUint() (map[uint]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[uint]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[uint(cnvK)] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint", b.Class), "map key not uint type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint", b.Class))
	}
	srcLen := v.Len()
	result := make(map[uint]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapUint convert err:%v", cnvKErr)
			continue
		}
		result[uint(cnvK)] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapUint %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}
func (b *BaseItfImpl) ToMapUint64() (map[uint64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint64]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint64", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[uint64]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[uint64(cnvK)] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint64", b.Class), "map key not uint64 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint64", b.Class))
	}
	srcLen := v.Len()
	result := make(map[uint64]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapUint64 convert err:%v", cnvKErr)
			continue
		}
		result[uint64(cnvK)] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint64", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapUint64 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapUint32() (map[uint32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[uint32]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint32", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[uint32]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[uint32(cnvK)] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapUint32", b.Class), "map key not uint32 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint32", b.Class))
	}
	srcLen := v.Len()
	result := make(map[uint32]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapUint32 convert err:%v", cnvKErr)
			continue
		}
		result[uint32(cnvK)] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapUint32", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapUint32 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapFloat32() (map[float32]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float32]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat32", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[float32]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToFloat32(kk); cnvKErr == nil {
				result[cnvK] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat32", b.Class), "map key not float32 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32", b.Class))
	}
	srcLen := v.Len()
	result := make(map[float32]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToFloat32(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapFloat32 convert err:%v", cnvKErr)
			continue
		}
		result[cnvK] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapFloat32 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapFloat64() (map[float64]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float64]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat64", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[float64]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToFloat64(kk); cnvKErr == nil {
				result[cnvK] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat64", b.Class), "map key not float64 type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64", b.Class))
	}
	srcLen := v.Len()
	result := make(map[float64]interface{}, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToFloat64(rfK.Interface())
		if cnvKErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapFloat64 convert err:%v", cnvKErr)
			continue
		}
		result[cnvK] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapFloat64 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}
func (b *BaseItfImpl) ToMapItf() (map[interface{}]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapItf", b.Class), "not json map str", err)
		}
	}

	if cvtOk {
		result := make(map[interface{}]interface{}, len(mapStrItf))
		for kk, vv := range mapStrItf {
			if cnvK, cnvKErr := pkg.ToInt64(kk); cnvKErr == nil {
				result[cnvK] = vv
			} else if !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapItf", b.Class), "map key not interface{} type", cnvKErr)
			}
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapItf", b.Class))
	}
	srcLen := v.Len()
	result := make(map[interface{}]interface{}, v.Len())
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}
		result[rfK.Interface()] = mpV.Interface()
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapItf", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapItf %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapStrToStr() (map[string]string, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	// 1. 直接转换
	if v, ok := b.IterVal.(map[string]string); ok {
		return v, nil
	}

	// 2. 经验转换
	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		result := make(map[string]string, len(v))
		for k, v := range v {
			result[pkg.ToStr(k)] = pkg.ToStr(v)
		}
		return result, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapStrToStr", b.Class), "not json map str", err)
		}
	}
	if cvtOk {
		result := make(map[string]string, len(mapStrItf))
		for kk, vv := range mapStrItf {
			result[kk] = pkg.ToStr(vv)
		}
		return result, nil
	}

	// 3. 反射转换,不在入口判断是否map的原因:我们信任转换是可以成功的,从而减少反射判断
	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapStrToStr", b.Class))
	}
	srcLen := v.Len()
	result := make(map[string]string, v.Len())
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}
		result[pkg.ToStr(rfK.Interface())] = pkg.ToStr(mpV.Interface())
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapStrToStr", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapStrToStr %d key value convert failed", srcLen-len(result))
	}

	return result, nil
}

func (b *BaseItfImpl) ToMapIntToInt() (map[int]int, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[int]int); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapIntToInt", b.Class), "not json map str", err)
		}
	}
	if cvtOk {
		result := make(map[int]int, len(mapStrItf))
		for kk, vv := range mapStrItf {
			cnvK, cnvKErr := pkg.ToInt64(kk)
			cnvV, cnvVErr := pkg.ToInt64(vv)
			if (cnvKErr != nil || cnvVErr != nil) && !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapIntToInt", b.Class), "map key not int type", cnvKErr)
			}
			if cnvKErr == nil && cnvVErr == nil {
				result[int(cnvK)] = int(cnvV)
			}
		}
		return result, nil
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		result := make(map[int]int, len(v))
		for kk, vv := range v {
			cnvK, cnvKErr := pkg.ToInt64(kk)
			cnvV, cnvVErr := pkg.ToInt64(vv)
			if cnvKErr != nil || cnvVErr != nil {
				logx.CtxWarn(b.Ctx, "ToMapIntToInt convert err:%v:%v", cnvKErr, cnvVErr)
				continue
			}
			result[int(cnvK)] = int(cnvV)
		}
		if len(result) != len(v) && !conf.CONF.SkipCvtFailForToMapType {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapIntToInt", b.Class))
		}
		if len(result) != len(v) {
			logx.CtxWarn(b.Ctx, "ToMapIntToInt %d key value convert failed", len(v)-len(result))
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapIntToInt", b.Class))
	}
	srcLen := v.Len()
	result := make(map[int]int, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		cnvV, cnvVErr := pkg.ToInt64(mpV.Interface())
		if cnvKErr != nil || cnvVErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapIntToInt convert err:%v:%v", cnvKErr, cnvVErr)
			continue
		}
		result[int(cnvK)] = int(cnvV)
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapIntToInt", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapIntToInt %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapInt64ToInt64() (map[int64]int64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[int64]int64); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class), "not json map str", err)
		}
	}
	if cvtOk {
		result := make(map[int64]int64, len(mapStrItf))
		for kk, vv := range mapStrItf {
			cnvK, cnvKErr := pkg.ToInt64(kk)
			cnvV, cnvVErr := pkg.ToInt64(vv)
			if (cnvKErr != nil || cnvVErr != nil) && !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class), "map key not int64 type", cnvKErr)
			}
			if cnvKErr == nil && cnvVErr == nil {
				result[cnvK] = cnvV
			}
		}
		return result, nil
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		result := make(map[int64]int64, len(v))
		for kk, vv := range v {
			cnvK, cnvKErr := pkg.ToInt64(kk)
			cnvV, cnvVErr := pkg.ToInt64(vv)
			if cnvKErr != nil || cnvVErr != nil {
				logx.CtxWarn(b.Ctx, "ToMapInt64ToInt64 convert err:%v:%v", cnvKErr, cnvVErr)
				continue
			}
			result[cnvK] = cnvV
		}
		if len(result) != len(v) && !conf.CONF.SkipCvtFailForToMapType {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class))
		}
		if len(result) != len(v) {
			logx.CtxWarn(b.Ctx, "ToMapInt64ToInt64 %d key value convert failed", len(v)-len(result))
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class))
	}
	srcLen := v.Len()
	result := make(map[int64]int64, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToInt64(rfK.Interface())
		cnvV, cnvVErr := pkg.ToInt64(mpV.Interface())
		if cnvKErr != nil || cnvVErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapInt64ToInt64 convert err:%v:%v", cnvKErr, cnvVErr)
			continue
		}
		result[cnvK] = cnvV
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapInt64ToInt64", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapInt64ToInt64 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapFloat64ToFloat64() (map[float64]float64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float64]float64); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class), "not json map str", err)
		}
	}
	if cvtOk {
		result := make(map[float64]float64, len(mapStrItf))
		for kk, vv := range mapStrItf {
			cnvK, cnvKErr := pkg.ToFloat64(kk)
			cnvV, cnvVErr := pkg.ToFloat64(vv)
			if (cnvKErr != nil || cnvVErr != nil) && !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class), "map key not float64 type", cnvKErr)
			}
			if cnvKErr == nil && cnvVErr == nil {
				result[cnvK] = cnvV
			}
		}
		return result, nil
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		result := make(map[float64]float64, len(v))
		for kk, vv := range v {
			cnvK, cnvKErr := pkg.ToFloat64(kk)
			cnvV, cnvVErr := pkg.ToFloat64(vv)
			if cnvKErr != nil || cnvVErr != nil {
				logx.CtxWarn(b.Ctx, "ToMapFloat64ToFloat64 convert err:%v:%v", cnvKErr, cnvVErr)
				continue
			}
			result[cnvK] = cnvV
		}
		if len(result) != len(v) && !conf.CONF.SkipCvtFailForToMapType {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class))
		}
		if len(result) != len(v) {
			logx.CtxWarn(b.Ctx, "ToMapFloat64ToFloat64 %d key value convert failed", len(v)-len(result))
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class))
	}
	srcLen := v.Len()
	result := make(map[float64]float64, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToFloat64(rfK.Interface())
		cnvV, cnvVErr := pkg.ToFloat64(mpV.Interface())
		if cnvKErr != nil || cnvVErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapFloat64ToFloat64 convert err:%v:%v", cnvKErr, cnvVErr)
			continue
		}
		result[cnvK] = cnvV
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat64ToFloat64", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapFloat64ToFloat64 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

func (b *BaseItfImpl) ToMapFloat32ToFloat32() (map[float32]float32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.(map[float32]float32); ok {
		return v, nil
	}

	var mapStrItf map[string]interface{}
	cvtOk := false
	if v, ok := b.IterVal.(map[string]interface{}); ok {
		mapStrItf, cvtOk = v, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstMap, err := pkg.JsonLoadsMap(jsonStr); err == nil {
			mapStrItf, cvtOk = rstMap, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class), "not json map str", err)
		}
	}
	if cvtOk {
		result := make(map[float32]float32, len(mapStrItf))
		for kk, vv := range mapStrItf {
			cnvK, cnvKErr := pkg.ToFloat32(kk)
			cnvV, cnvVErr := pkg.ToFloat32(vv)
			if (cnvKErr != nil || cnvVErr != nil) && !conf.CONF.SkipCvtFailForToMapType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class), "map key not float32 type", cnvKErr)
			}
			if cnvKErr == nil && cnvVErr == nil {
				result[cnvK] = cnvV
			}
		}
		return result, nil
	}

	if v, ok := b.IterVal.(map[interface{}]interface{}); ok {
		result := make(map[float32]float32, len(v))
		for kk, vv := range v {
			cnvK, cnvKErr := pkg.ToFloat32(kk)
			cnvV, cnvVErr := pkg.ToFloat32(vv)
			if cnvKErr != nil || cnvVErr != nil {
				logx.CtxWarn(b.Ctx, "ToMapFloat32ToFloat32 convert err:%v:%v", cnvKErr, cnvVErr)
				continue
			}
			result[cnvK] = cnvV
		}
		if len(result) != len(v) && !conf.CONF.SkipCvtFailForToMapType {
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class))
		}
		if len(result) != len(v) {
			logx.CtxWarn(b.Ctx, "ToMapFloat32ToFloat32 %d key value convert failed", len(v)-len(result))
		}
		return result, nil
	}

	v := reflect.ValueOf(b.IterVal)
	if v.Kind() != reflect.Map {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class))
	}
	srcLen := v.Len()
	result := make(map[float32]float32, srcLen)
	for _, rfK := range v.MapKeys() {
		if !rfK.IsValid() || !rfK.CanInterface() {
			continue
		}
		mpV := v.MapIndex(rfK)
		if !mpV.IsValid() || !mpV.CanInterface() {
			continue
		}

		cnvK, cnvKErr := pkg.ToFloat32(rfK.Interface())
		cnvV, cnvVErr := pkg.ToFloat32(mpV.Interface())
		if cnvKErr != nil || cnvVErr != nil {
			logx.CtxWarn(b.Ctx, "ToMapFloat32ToFloat32 convert err:%v:%v", cnvKErr, cnvVErr)
			continue
		}
		result[cnvK] = cnvV
	}

	if len(result) != srcLen && !conf.CONF.SkipCvtFailForToMapType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToMapFloat32ToFloat32", b.Class))
	}
	if len(result) != srcLen {
		logx.CtxWarn(b.Ctx, "ToMapFloat32ToFloat32 %d key value convert failed", srcLen-len(result))
	}
	return result, nil
}

// ToArrayType ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) ToList() ([]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]interface{}); ok {
		return v, nil
	}

	if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			return rstList, nil
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToList", b.Class), "not json list str", err)
		}
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToList", b.Class))
	}

	resultList := make([]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToList", b.Class))
		}
		resultList = append(resultList, ele.Interface())
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToList %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListMap() ([]map[string]interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]map[string]interface{}); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListMap", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		result := make([]map[string]interface{}, 0, len(listItf))
		for _, v := range listItf {
			if vv, err := toMap(b.Ctx, v); err == nil {
				result = append(result, vv)
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListMap", b.Class), "convert list val to map err", err)
			}
		}
		return result, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListMap", b.Class))
	}

	resultList := make([]map[string]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListMap", b.Class))
		}
		m, err := toMap(b.Ctx, ele.Interface())
		if err != nil || !conf.CONF.SkipCvtFailForToArrayType {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListMap", b.Class), "convert to map err", err)
		}
		if m != nil {
			resultList = append(resultList, m)
		}
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListMap %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListStr() ([]string, error) {
	return b.toListStr()
}

func (b *BaseItfImpl) ToListStrF() ([]string, error) {
	return b.toListStr(true)
}

// ToListStr to list string, force表示是否强转,只要是数组就能转成[]string
func (b *BaseItfImpl) toListStr(force ...bool) ([]string, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]string); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListStr", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		result := make([]string, 0, len(listItf))
		for _, v := range listItf {
			result = append(result, pkg.ToStr(v))
		}
		return result, nil
	}

	// 不强制转换-则报错返回; 强制转换且b.IterVal不是list类型-则报错返回
	rv := reflect.ValueOf(b.IterVal)
	if !(len(force) > 0 && force[0]) || (rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice) {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
	}

	resultList := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
		}
		resultList = append(resultList, pkg.ToStr(ele.Interface()))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListStr %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListInt() ([]int, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListStr", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]int, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, int(vint64))
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListInt", b.Class), "list val not int type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt", b.Class))
	}

	resultList := make([]int, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListStr", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, int(iv))
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListInt %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListInt32() ([]int32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int32); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListInt32", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]int32, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, int32(vint64))
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListInt32", b.Class), "list val not int32 type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
	}

	resultList := make([]int32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, int32(iv))
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt32", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListInt32 %s convert failed", rv.Len()-len(resultList))
	}

	return resultList, nil
}

func (b *BaseItfImpl) ToListRune() ([]rune, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]rune); ok {
		return v, nil
	}

	return b.ToListInt32()
}

func (b *BaseItfImpl) ToListInt64() ([]int64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]int64); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListInt64", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]int64, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, vint64)
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListInt64", b.Class), "list val not int64 type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
	}

	resultList := make([]int64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, iv)
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListInt64", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListInt64 %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListUint() ([]uint, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]uint, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, uint(vint64))
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint", b.Class), "list val not uint type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
	}

	resultList := make([]uint, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint(iv))
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUInt", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListUInt %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListUint64() ([]uint64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint64); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint64", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]uint64, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, uint64(vint64))
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint64", b.Class), "list val not uint64 type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
	}

	resultList := make([]uint64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint64(iv))
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint64", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListUint64 %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListUint32() ([]uint32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]uint32); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint32", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listInt := make([]uint32, 0, len(listItf))
		for _, v := range listItf {
			if vint64, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listInt = append(listInt, uint32(vint64))
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListUint32", b.Class), "list val not uint32 type", cnvKErr)
			}
		}
		return listInt, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
	}

	resultList := make([]uint32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
		}
		if iv, ok := pkg.ToInt64(ele.Interface()); ok == nil {
			resultList = append(resultList, uint32(iv))
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListUint32", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListUint32 %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListFloat32() ([]float32, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]float32); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListFloat32", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listf32 := make([]float32, 0, len(listItf))
		for _, v := range listItf {
			if vf32, cnvKErr := pkg.ToFloat32(v); cnvKErr == nil {
				listf32 = append(listf32, vf32)
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListFloat32", b.Class), "list val not float32 type", cnvKErr)
			}
		}
		return listf32, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
	}

	resultList := make([]float32, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
		}
		if iv, ok := pkg.ToFloat32(ele.Interface()); ok == nil {
			resultList = append(resultList, iv)
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat32", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListFloat32 %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListFloat64() ([]float64, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]float64); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListFloat64", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listf64 := make([]float64, 0, len(listItf))
		for _, v := range listItf {
			if vf64, cnvKErr := pkg.ToFloat64(v); cnvKErr == nil {
				listf64 = append(listf64, vf64)
			} else if !conf.CONF.SkipCvtFailForToArrayType {
				return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListFloat64", b.Class), "list val not float64 type", cnvKErr)
			}
		}
		return listf64, nil
	}

	rv := reflect.ValueOf(b.IterVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
	}

	resultList := make([]float64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ele := rv.Index(i)
		if !(ele.IsValid() && ele.CanInterface()) {
			if conf.CONF.SkipCvtFailForToArrayType {
				continue
			}
			return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
		}
		if iv, ok := pkg.ToFloat64(ele.Interface()); ok == nil {
			resultList = append(resultList, iv)
		}
	}

	if len(resultList) != rv.Len() && !conf.CONF.SkipCvtFailForToArrayType {
		return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListFloat64", b.Class))
	}
	if rv.Len() != len(resultList) {
		logx.CtxWarn(b.Ctx, "ToListFloat64 %s convert failed", rv.Len()-len(resultList))
	}
	return resultList, nil
}

func (b *BaseItfImpl) ToListBool() ([]bool, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	if v, ok := b.IterVal.([]bool); ok {
		return v, nil
	}

	var listItf []interface{}
	cvtOk := false
	if vList, ok := b.IterVal.([]interface{}); ok {
		listItf, cvtOk = vList, true
	} else if isJson, jsonStr := pkg.JsonChecker(b.IterVal); isJson {
		if rstList, err := pkg.JsonLoadsList(jsonStr); err == nil {
			listItf, cvtOk = rstList, true
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToListBool", b.Class), "not json list str", err)
		}
	}
	if cvtOk {
		listb := make([]bool, 0, len(listItf))
		for _, v := range listItf {
			if bl, ok := v.(bool); ok {
				listb = append(listb, bl)
			} else if i, cnvKErr := pkg.ToInt64(v); cnvKErr == nil {
				listb = append(listb, i != 0)
			}
		}
		return listb, nil
	}

	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToListBool", b.Class))
}

// ToObjectType ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) ToStruct(stc interface{}) (interface{}, error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}

	// 1. 字符串转换
	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if obj, err := pkg.JsonLoadsObj(js, stc); err == nil {
			return obj, nil
		} else {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToStruct", b.Class), "json un-match struct", err)
		}
	}

	// 2. 类型原本就一致
	rfSrc, rfDst := pkg.Interpret(b.IterVal), pkg.Interpret(stc)
	if reflect.TypeOf(rfSrc).Kind() == reflect.TypeOf(rfDst).Kind() {
		return b.IterVal, nil
	}

	//3. map转换到obj
	if reflect.TypeOf(rfSrc).Kind() == reflect.Map {
		dstStruct, err := pkg.MapToStruct(b.IterVal, stc)
		if err != nil {
			return nil, itferr.NewConvFailedX(fmt.Sprintf("%s#ToStruct", b.Class), "map cannot cvt to struct", err)
		}
		return dstStruct, nil
	}

	return nil, itferr.NewConvFailed(fmt.Sprintf("%s#ToStruct", b.Class))
}

func (b *BaseItfImpl) Uniq() api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	b.IterVal = pkg.UniqList(b.IterVal)
	return b
}

// Interface ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) Get(key interface{}) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	b.ItfErr = itferr.NewUnSupportInterfaceFunc(fmt.Sprintf("BaseItfImpl#Get(%v)", key))
	return b
}

func (b *BaseItfImpl) GetAny(keys ...interface{}) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	b.ItfErr = itferr.NewUnSupportInterfaceFunc(fmt.Sprintf("BaseItfImpl#GetAny(%v)", keys))
	return b
}

func (b *BaseItfImpl) Index(index int) api.MapInterface {
	if b.ItfErr != nil {
		return b
	}

	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		listItf, err := pkg.JsonLoadsList(js)
		if err != nil {
			b.ItfErr = itferr.NewConvFailedX(fmt.Sprintf("BaseItfImpl#Index(%d)", index), "not json list str", err)
			return b
		}

		b.IterChain.ReplaceBack(listItf)
		b.IterVal = listItf[index]
		b.IterChain.PushBackByIdx(index, b.IterVal)

		v := reflect.ValueOf(b.IterVal)
		switch v.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.String, reflect.Interface:
			return FrWithChain(b.Ctx, b.IterVal, b.IterChain)
		}
		return b
	}

	v := pkg.ReflectToVal(b.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() <= index {
			b.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("Index(%d)#(%d)", index, v.Len()))
			return b
		}
	default:
		b.ItfErr = itferr.NewCurrentCannotUseIndex(fmt.Sprintf("BaseItfImpl#Index(%v)", index))
		return b
	}

	// 尝试类型转换,尝试转换两种常用类型
	switch vv := b.IterVal.(type) {
	case []interface{}:
		b.IterVal = vv[index]
		b.IterChain.PushBackByIdx(index, b.IterVal)

		rfVV := pkg.ReflectToVal(b.IterVal)
		switch rfVV.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.String, reflect.Interface:
			return FrWithChain(b.Ctx, b.IterVal, b.IterChain)
		}
		return b
	case []map[string]interface{}:
		b.IterVal = vv[index]
		b.IterChain.PushBackByIdx(index, b.IterVal)
		return FrWithChain(b.Ctx, b.IterVal, b.IterChain)
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		indexV := v.Index(index)
		if !indexV.IsValid() || !indexV.CanInterface() {
			b.ItfErr = itferr.NewListIndexIllegal(fmt.Sprintf("Index(%d)", index))
			return b
		}

		interfaceV := indexV.Interface()

		b.IterVal = interfaceV
		b.IterChain.PushBackByIdx(index, b.IterVal)

		rfVV := pkg.ReflectToVal(b.IterVal)
		switch rfVV.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.String, reflect.Interface:
			return FrWithChain(b.Ctx, b.IterVal, b.IterChain)
		}
		return b
	}

	b.ItfErr = itferr.NewCurrentCannotUseIndex(fmt.Sprintf("BaseItfImpl#Index(%v)", index))
	return b
}

func (b *BaseItfImpl) ForEach(forFunc api.ForFunc) api.MapInterface {
	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if mapObj, err := pkg.JsonLoadsMap(js); err == nil {
			b.IterVal = mapObj
			b.IterChain.ReplaceBack(mapObj)
			resultList := make([]interface{}, 0, len(mapObj))
			resultMap := make(map[interface{}]interface{}, len(mapObj))
			i := 0
			for kk, vv := range mapObj {
				key, val := forFunc(i, kk, vv)
				i++

				if key == nil && val == nil {
					continue
				}
				if key != nil {
					resultMap[key] = val
				} else {
					resultList = append(resultList, val)
				}
			}
			return NewForeachItfImpl(b.Ctx, resultList, resultMap).WithIterChain(b.IterChain)
		}

		if listObj, err := pkg.JsonLoadsList(js); err == nil {
			resultList := make([]interface{}, 0, len(listObj))
			resultMap := make(map[interface{}]interface{}, len(listObj))
			for i, vv := range listObj {
				key, val := forFunc(i, nil, vv)
				if key == nil && val == nil {
					continue
				}
				if key != nil {
					resultMap[key] = val
				} else {
					resultList = append(resultList, val)
				}
			}
			return NewForeachItfImpl(b.Ctx, resultList, resultMap).WithIterChain(b.IterChain)
		}

		return NewExceptItfImpl()
	}

	v := pkg.ReflectToVal(b.IterVal)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		resultList := make([]interface{}, 0, v.Len())
		resultMap := make(map[interface{}]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			idxV := v.Index(i)
			if !idxV.IsValid() || !idxV.CanInterface() {
				continue
			}
			key, val := forFunc(i, nil, idxV.Interface())
			if key == nil && val == nil {
				continue
			}
			if key != nil {
				resultMap[key] = val
			} else {
				resultList = append(resultList, val)
			}
		}
		return NewForeachItfImpl(b.Ctx, resultList, resultMap).WithIterChain(b.IterChain)
	case reflect.Map:
		resultList := make([]interface{}, 0, v.Len())
		resultMap := make(map[interface{}]interface{}, v.Len())
		for i, rfK := range v.MapKeys() {
			if !rfK.IsValid() || !rfK.CanInterface() {
				continue
			}
			mpV := v.MapIndex(rfK)
			if !mpV.IsValid() || !mpV.CanInterface() {
				continue
			}
			key, val := forFunc(i, rfK.Interface(), mpV.Interface())
			if key == nil && val == nil {
				continue
			}
			if key != nil {
				resultMap[key] = val
			} else {
				resultList = append(resultList, val)
			}
		}
		return NewForeachItfImpl(b.Ctx, resultList, resultMap).WithIterChain(b.IterChain)
	}
	return NewExceptItfImpl()
}

// OriginTypeChecker ----------------------------------------------------------------------------------------
func (b *BaseItfImpl) IsStr() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}

	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsStr", itferr.ValueTypeErr, "", nil)
	}

	if _, ok := b.IterVal.(string); ok {
		return true, nil
	}

	if v := reflect.ValueOf(b.IterVal); v.Kind() == reflect.Ptr && !v.IsNil() {
		if v.Elem().Kind() == reflect.String {
			return true, nil
		}
	}

	return false, nil
}

func (b *BaseItfImpl) IsDigit() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsDigit", itferr.ValueTypeErr, "", nil)
	}

	val := b.IterVal
	if v := reflect.ValueOf(b.IterVal); v.Kind() == reflect.Ptr && !v.IsNil() {
		val = v.Elem().Interface()
	}

	switch val.(type) {
	case json.Number, int8, int16, int32, int, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true, nil
	}

	return false, nil
}

func (b *BaseItfImpl) IsList() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsList", itferr.ValueTypeErr, "", nil)
	}

	v := pkg.ReflectToVal(b.IterVal)
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		return true, nil
	}

	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if _, err := pkg.JsonLoadsList(js); err == nil {
			return true, nil
		}
	}

	return false, nil
}

func (b *BaseItfImpl) IsStrList() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsStrList", itferr.ValueTypeErr, "", nil)
	}

	if _, ok := b.IterVal.([]string); ok {
		return true, nil
	}
	if _, ok := b.IterVal.([]*string); ok {
		return true, nil
	}

	if v := reflect.ValueOf(b.IterVal); v.Kind() == reflect.Ptr && !v.IsNil() {
		vv := v.Elem().Interface()
		if _, ok := vv.([]string); ok {
			return true, nil
		}
		if _, ok := b.IterVal.([]*string); ok {
			return true, nil
		}
	}

	return false, nil
}

func (b *BaseItfImpl) IsDigitList() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsDigitList", itferr.ValueTypeErr, "", nil)
	}

	switch b.IterVal.(type) {
	case []int8, []int16, []int32, []int, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64:
		return true, nil
	case []*int8, []*int16, []*int32, []*int, []*int64, []*uint, []*uint8, []*uint16, []*uint32, []*uint64, []*float32, []*float64:
		return true, nil
	}

	rfV := pkg.ReflectToVal(b.IterVal)
	if rfV.Kind() != reflect.Slice && rfV.Kind() != reflect.Array {
		return false, nil
	}
	if listInt, err := b.ToListInt64(); err == nil {
		return len(listInt) == rfV.Len(), nil
	}

	return false, nil
}

func (b *BaseItfImpl) IsMap() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsMap", itferr.ValueTypeErr, "", nil)
	}

	v := pkg.ReflectToVal(b.IterVal)
	if v.Kind() == reflect.Map {
		return true, nil
	}

	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if _, err := pkg.JsonLoadsMap(js); err == nil {
			return true, nil
		}
	}

	return false, nil
}

func (b *BaseItfImpl) IsMapStrItf() (bool, error) {
	if b.ItfErr != nil {
		return false, b.ItfErr
	}
	if b.IterVal == nil {
		return false, itferr.NewMapItfErr("IsMapStrItf", itferr.ValueTypeErr, "", nil)
	}

	if _, ok := b.IterVal.(map[string]interface{}); ok {
		return true, nil
	}

	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if _, err := pkg.JsonLoadsMap(js); err == nil {
			return true, nil
		}
		return false, nil
	}

	if v := reflect.ValueOf(b.IterVal); v.Kind() == reflect.Ptr && !v.IsNil() {
		if _, ok := v.Elem().Interface().(map[string]interface{}); ok {
			return true, nil
		}
	}

	return false, nil
}

func (b *BaseItfImpl) SetMap(key interface{}, val interface{}) (orgVal interface{}, err error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if b.IterVal == nil || b.IterChain.Back() == nil {
		return nil, itferr.NewSetValueErr(fmt.Sprintf("%s#SetMap", b.Class), "be set val illegal", nil)
	}
	// 1. 当前IterVal是Map
	switch vMap := b.IterVal.(type) {
	case map[string]interface{}:
		vMap[pkg.ToStr(key)] = val
		return b.OrgVal()
	case map[interface{}]interface{}:
		vMap[key] = val
		return b.OrgVal()
	}

	rfV := pkg.ReflectToVal(b.IterVal)
	if rfV.Kind() == reflect.Map {
		if rfV.Type().Key().Kind() != reflect.TypeOf(key).Kind() {
			return nil, itferr.NewSetValueErr(fmt.Sprintf("%s#SetMap", b.Class), "param-key with map-key not-match", nil)
		}
		valKind := rfV.Type().Elem().Kind()
		if valKind != reflect.TypeOf(val).Kind() && valKind != reflect.Interface {
			return nil, itferr.NewSetValueErr(fmt.Sprintf("%s#SetMap", b.Class), "param-val type illegal cannot set val", nil)
		}
		rfV.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		return b.OrgVal()
	}

	// 2. 当前IterVal是Json-Str-Map
	if isJson, js := pkg.JsonChecker(b.IterVal); isJson {
		if result, err := pkg.JsonLoadsMap(js); err == nil {
			// 把key上一层的值改为map
			setErr := b.IterChain.SetBackPreVal(result)
			if err != nil {
				b.ItfErr = itferr.NewSetValueErr(fmt.Sprintf("%s#SetMap", b.Class), "SetBackPreVal err", setErr)
				return nil, b.ItfErr
			}
			result[pkg.ToStr(key)] = val
			return b.OrgVal()
		}
	}

	return nil, itferr.NewUnSupportSetValErr(fmt.Sprintf("%s#SetMap", b.Class), "val is not map or json map", nil)
}

func (b *BaseItfImpl) SetAsMap(key interface{}) (orgVal interface{}, err error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	rfIterVal := pkg.ReflectToVal(b.IterVal)
	if rfIterVal.Kind() != reflect.Map {
		return nil, itferr.NewUnSupportSetValErr(fmt.Sprintf("%s#SetAsMap", b.Class), "val is not map", nil)
	}

	rfIterValType := rfIterVal.Type().Elem().Kind()
	if rfIterValType != reflect.Interface {
		return nil, itferr.NewUnSupportSetValErr(fmt.Sprintf("%s#SetAsMap", b.Class), "map val type un-match", nil)
	}

	rfKey := reflect.ValueOf(key)
	mapVal, itfErr := b.toInterface(rfIterVal.MapIndex(rfKey))
	if itfErr != nil {
		b.ItfErr = itfErr
		return nil, b.ItfErr
	}

	if reflect.TypeOf(mapVal).Kind() != reflect.String {
		return nil, itferr.NewUnSupportSetValErr(fmt.Sprintf("%s#SetAsMap", b.Class), "val is not json-str", nil)
	}

	if jsonMap, err := pkg.JsonLoadsMap(pkg.ToStr(mapVal)); err == nil {
		rfIterVal.SetMapIndex(rfKey, reflect.ValueOf(jsonMap))
		return b.OrgVal()
	}

	return nil, itferr.NewUnSupportSetValErr(fmt.Sprintf("%s#SetAsMap", b.Class), "val is not json string", nil)
}

func (b *BaseItfImpl) SetAllAsMap() (orgVal interface{}, err error) {
	if b.ItfErr != nil {
		return nil, b.ItfErr
	}
	if ok, val := b.deepSetMap(b.IterVal); ok {
		// 此时b.IterVal为json str,val为map,递归map中的内容对json-str进行处理
		b.deepSetMap(val)
		// 把key上一层的值改为map
		setErr := b.IterChain.SetBackPreVal(val)
		if err != nil {
			b.ItfErr = itferr.NewSetValueErr(fmt.Sprintf("%s#SetMap", b.Class), "SetBackPreVal err", setErr)
			return nil, b.ItfErr
		}
	}
	return b.OrgVal()
}

// deepSetMap 深度优先递归遍历,不限制递归深度
func (b *BaseItfImpl) deepSetMap(itVal interface{}) (bool, interface{}) {
	rfIterVal := pkg.ReflectToVal(itVal)
	if !rfIterVal.IsValid() || !rfIterVal.CanInterface() {
		return false, itVal
	}

redo:
	switch rfIterVal.Kind() {
	case reflect.Map:
		for _, rfK := range rfIterVal.MapKeys() {
			mpV := rfIterVal.MapIndex(rfK)
			if !mpV.IsValid() || !mpV.CanInterface() {
				continue
			}
			if ok, val := b.deepSetMap(mpV.Interface()); ok {
				// 为确保赋值成功,rfIterVal map的类型一定是map[x]interface{}
				if rfIterVal.Type().Elem().Kind() == reflect.Interface {
					rfIterVal.SetMapIndex(rfK, reflect.ValueOf(val))
					goto redo
				}
			}
		}
	case reflect.String:
		// 如果是json,序列化为map,赋值;并将该map纳入递归中.
		if ok, s := pkg.JsonChecker(rfIterVal.Interface()); ok {
			// json load后必须是map[string]interface{}, []interface{},才去赋值
			if val, err := pkg.JsonLoadsMap(s); err == nil {
				return true, val
			}
			if val, err := pkg.JsonLoadsList(s); err == nil {
				return true, val
			}
			return false, itVal
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < rfIterVal.Len(); i++ {
			idxV := rfIterVal.Index(i)
			if !idxV.IsValid() || !idxV.CanInterface() {
				continue
			}
			if ok, val := b.deepSetMap(idxV.Interface()); ok {
				// 为确保能赋值成功,rfIterVal list类型必须是[]interface{}
				if _, ok := itVal.([]interface{}); !ok {
					return false, itVal
				}
				idxV.Set(reflect.ValueOf(val))
				goto redo
			}
		}
	default:
		return false, itVal
	}
	return false, itVal
}

func (b *BaseItfImpl) New() api.MapInterface {
	return &BaseItfImpl{
		Ctx:       b.Ctx,
		Class:     b.Class,
		ItfErr:    b.ItfErr,
		IterVal:   b.IterVal,
		IterChain: b.IterChain.Clone(),
	}
}

func (b *BaseItfImpl) PrintPath() string {
	if b.IterChain.Len() == 0 {
		return ""
	}

	resultStr, isFirst := make([]string, 0, b.IterChain.Len()), true
	for e := b.IterChain.Front(); e != nil; e = e.Next() {
		iterCtx, ok := e.Value.(*IterCtx)
		if !ok {
			logx.CtxWarn(b.Ctx, "[%s#PrintPath] Value cannot cvt to IterCtx", b.Class)
			continue
		}
		if isFirst {
			resultStr = append(resultStr, fmt.Sprintf("%T", iterCtx.Val))
			isFirst = false
			continue
		}
		if iterCtx.Key == nil {
			resultStr = append(resultStr, fmt.Sprintf("%d:%T", iterCtx.Idx, iterCtx.Val))
			continue
		}
		resultStr = append(resultStr, fmt.Sprintf("%s:%T", pkg.ToStr(iterCtx.Key), iterCtx.Val))
	}
	return strings.Join(resultStr, " => ")
}

func (b *BaseItfImpl) OrgVal() (interface{}, error) {
	//if b.ItfErr != nil {
	//	return nil, b.ItfErr
	//}
	return b.IterChain.HeadVal(), b.ItfErr
}
