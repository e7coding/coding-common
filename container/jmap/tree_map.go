package jmap

import (
	"github.com/e7coding/coding-common/container/jtree"
)

type TreeMap = jtree.RedBlackTree

func NewTreeMap(comparator func(v1, v2 interface{}) int) *TreeMap {
	return jtree.NewRedBlackTree(comparator)
}

func NewTreeMapFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}) *TreeMap {
	return jtree.NewRedBlackTreeFrom(comparator, data)
}
