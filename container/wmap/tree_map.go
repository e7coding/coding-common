package wmap

import (
	"github.com/coding-common/container/tree"
)

type TreeMap = tree.RedBlackTree

func NewTreeMap(comparator func(v1, v2 interface{}) int) *TreeMap {
	return tree.NewRedBlackTree(comparator)
}

func NewTreeMapFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}) *TreeMap {
	return tree.NewRedBlackTreeFrom(comparator, data)
}
