// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
)

// SortedSetWriter 只包含对有序集合的写入/修改操作
type SortedSetWriter interface {
	ZAdd(key string, option *ZAddOption, member ZAddMember, members ...ZAddMember) (*jvar.Var, error)
	ZIncrBy(key string, increment float64, member interface{}) (float64, error)
	ZRem(key string, member interface{}, members ...interface{}) (int64, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRemRangeByScore(key string, min, max string) (int64, error)
	ZRemRangeByLex(key string, min, max string) (int64, error)
}

// SortedSetReader 只包含对有序集合的只读/查询操作
type SortedSetReader interface {
	ZScore(key string, member interface{}) (float64, error)
	ZCard(key string) (int64, error)
	ZCount(key string, min, max string) (int64, error)
	ZRange(key string, start, stop int64, option ...ZRangeOption) (

		jvar.Vars, error)
	ZRevRange(key string, start, stop int64, option ...ZRevRangeOption) (*jvar.Var, error)
	ZRank(key string, member interface{}) (int64, error)
	ZRevRank(key string, member interface{}) (int64, error)
	ZLexCount(key, min, max string) (int64, error)
}

// IGroupSortedSet 聚合了读写接口，向后兼容
type IGroupSortedSet interface {
	SortedSetWriter
	SortedSetReader
}

// ZAddOption provides options for function ZAdd.
type ZAddOption struct {
	XX bool // Only update elements that already exist. Don't add new elements.
	NX bool // Only add new elements. Don't update already existing elements.
	// Only update existing elements if the new score is less than the current score.
	// This flag doesn't prevent adding new elements.
	LT bool

	// Only update existing elements if the new score is greater than the current score.
	// This flag doesn't prevent adding new elements.
	GT bool

	// Modify the return value from the number of new elements added, to the total number of elements changed (CH is an abbreviation of changed).
	// Changed elements are new elements added and elements already existing for which the score was updated.
	// So elements specified in the command line having the same score as they had in the past are not counted.
	// Note: normally the return value of ZAdd only counts the number of new elements added.
	CH bool

	// When this option is specified ZAdd acts like ZIncrBy. Only one score-element pair can be specified in this mode.
	INCR bool
}

// ZAddMember is element struct for set.
type ZAddMember struct {
	Score  float64
	Member interface{}
}

// ZRangeOption provides extra option for ZRange function.
type ZRangeOption struct {
	ByScore bool
	ByLex   bool
	// The optional REV argument reverses the ordering, so elements are ordered from highest to lowest score,
	// and score ties are resolved by reverse lexicographical ordering.
	Rev   bool
	Limit *ZRangeOptionLimit
	// The optional WithScores argument supplements the command's reply with the scores of elements returned.
	WithScores bool
}

// ZRangeOptionLimit provides LIMIT argument for ZRange function.
// The optional LIMIT argument can be used to obtain a sub-range from the matching elements
// (similar to SELECT LIMIT offset, count in SQL). A negative `Count` returns all elements from the `Offset`.
type ZRangeOptionLimit struct {
	Offset *int
	Count  *int
}

// ZRevRangeOption provides options for function ZRevRange.
type ZRevRangeOption struct {
	WithScores bool
}
