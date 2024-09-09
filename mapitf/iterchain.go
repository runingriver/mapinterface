package mapitf

import (
	"container/list"
	"github.com/jinzhu/copier"
	"github.com/runingriver/mapinterface/itferr"
	"github.com/runingriver/mapinterface/pkg"
	"reflect"
)

// IterChain 记录迭代路径
type IterChain struct {
	*list.List
}

type IterCtx struct {
	Idx int // 当是list取值时,对应的是索引

	Key interface{} // 取到当前值的Key
	Val interface{} // 当前的值

}

func NewEnterIterCtx(val interface{}) *IterCtx {
	return &IterCtx{
		Val: val,
	}
}

func NewIterCtx(key, val interface{}, idx int) *IterCtx {
	return &IterCtx{
		Val: val,
		Key: key,
		Idx: idx,
	}
}

func NewLinkedList(val interface{}) *IterChain {
	ic := &IterChain{list.New()}
	ic.List.PushBack(NewEnterIterCtx(val)) // linked list's head
	return ic
}

func (i *IterChain) PushBack(val interface{}) {
	i.List.PushBack(NewIterCtx(nil, val, 0))
}

func (i *IterChain) PushBackByKey(key, val interface{}) {
	i.List.PushBack(NewIterCtx(key, val, 0))
}

func (i *IterChain) PushBackByIdx(idx int, val interface{}) {
	i.List.PushBack(NewIterCtx(nil, val, idx))
}

func (i *IterChain) ReplaceBack(val interface{}) {
	if e := i.List.Back(); e != nil {
		i.List.Remove(e)
		ic := e.Value.(*IterCtx)
		i.List.PushBack(NewIterCtx(ic.Key, val, ic.Idx))
		return
	}

	i.List.PushBack(NewEnterIterCtx(val))
}

// SetBackPreVal 当当前Iter值为json-map时,把json-map序列化为map[string]interface{}并赋给上一层的map上
func (i *IterChain) SetBackPreVal(val interface{}) error {
	backElement := i.List.Back()
	if backElement == nil {
		return itferr.NewMapItfErrX("IterChain#SetPreIterVal", itferr.IterChainIsEmpty)
	}
	preElement := backElement.Prev()
	if preElement == nil {
		return itferr.NewMapItfErrX("IterChain#SetPreIterVal", itferr.IterChainPreElementIsNil)
	}

	backIterCtx := backElement.Value.(*IterCtx)
	preIterCtx := preElement.Value.(*IterCtx)
	rfV := pkg.ReflectToVal(preIterCtx.Val)
	if rfV.Kind() != reflect.Map {
		return itferr.NewMapItfErrX("IterChain#SetPreIterVal", itferr.ValueTypeErr)
	}
	rfV.SetMapIndex(reflect.ValueOf(backIterCtx.Key), reflect.ValueOf(val))
	return nil
}

func (i *IterChain) Clone() *IterChain {
	ic := &IterChain{list.New()}
	for e := i.List.Front(); e != nil; e = e.Next() {
		iterCtx := e.Value.(*IterCtx)
		var newIterCtx IterCtx
		_ = copier.Copy(&newIterCtx, iterCtx)
		ic.List.PushBack(&newIterCtx)
	}
	return ic
}

func (i *IterChain) HeadVal() interface{} {
	if e := i.List.Front(); e != nil {
		return e.Value.(*IterCtx).Val
	}
	return nil
}
