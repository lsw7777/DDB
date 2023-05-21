package indexmanager

import (
	"github.com/google/btree"
	catalogmanager "Distributed-MiniSQL/minisql/manager/catalogmanager"
)



type node struct {
	key   interface{} 
	value catalogmanager.Address
}

// btree存放的东西必须实现Less(),即Item接口
func (i *node) Less(b btree.Item) bool {
	switch i.key.(type) {
	case int:
		return i.key.(int) < b.(*node).key.(int)
	case float32:
		return i.key.(float32) < b.(*node).key.(float32)
	case float64:
		return i.key.(float64) < b.(*node).key.(float64)
	case string:
		return i.key.(string) < b.(*node).key.(string)
	}
	return false
}

type BPTree struct {
	tree *btree.BTree
}

func NewBPTree() *BPTree {
	return &BPTree{btree.New(16)}
}

func (b *BPTree) Insert(key interface{}, value catalogmanager.Address) { 
	b.tree.ReplaceOrInsert(&node{key: key, value: value})
}

func (b BPTree) FindEq(key interface{}) *catalogmanager.Address { 
	value := b.tree.Get(&node{key: key})
	if value == nil {
		return nil
	}
	f := value.(*node)
	return &f.value
}

func (b BPTree) FindNeq(key interface{}) []catalogmanager.Address {
	var value []catalogmanager.Address
	b.tree.Ascend(func(item btree.Item) bool {
		f := item.(*node)
		if f.key != key {
			value = append(value, f.value)
		}
		return true
	})
	return value
}

func (b BPTree) FindLess(key interface{}) []catalogmanager.Address {
	var value []catalogmanager.Address
	b.tree.AscendLessThan(&node{key: key}, func(item btree.Item) bool {
		f := item.(*node)
		value = append(value, f.value)
		return true
	})
	return value
}

func (b BPTree) FindLeq(key interface{}) []catalogmanager.Address {
	var value = b.FindLess(key)
	var value2 = b.FindEq(key)
	if value2 != nil {
		value = append(value, *value2)
	}
	return value
}

func (b BPTree) FindGeq(key interface{}) []catalogmanager.Address {
	var value []catalogmanager.Address
	b.tree.AscendGreaterOrEqual(&node{key: key}, func(item btree.Item) bool {
		f := item.(*node)
		value = append(value, f.value)
		return true
	})
	return value
}

func (b BPTree) FindGreater(key interface{}) []catalogmanager.Address {
	var value = b.FindGeq(key)
	if value != nil {
		if b.FindEq(key) != nil { //说明有相等的情况
			value = value[1:]
		}
		// if(value[0].Elements[0].Data == key){ //去除等于的情况
		// 	value = value[1:]
		// }
	}
	return value
}

func (b *BPTree) Delete(key interface{}) bool {
	var value = b.tree.Delete(&node{key: key})
	return value != nil
}

