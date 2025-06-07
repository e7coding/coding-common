// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jcache

import (
	"time"

	"github.com/e7coding/coding-common/container/jvar"
)

// Adapter is the core adapter for cache features implements.
//
// Note that the implementer itself should guarantee the concurrent safety of these functions.
type Adapter interface {
	// Set sets cache with `key`-`value` pair, which is expired after `duration`.
	//
	// It does not expire if `duration` == 0.
	// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
	Set(key interface{}, value interface{}, duration time.Duration) error

	// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
	//
	// It does not expire if `duration` == 0.
	// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
	SetMap(data map[interface{}]interface{}, duration time.Duration) error

	// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
	// if `key` does not exist in the cache. It returns true the `key` does not exist in the
	// cache, and it sets `value` successfully to the cache, or else it returns false.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil.
	SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (ok bool, err error)

	// SetIfNotExistFunc sets `key` with result of function `f` and returns true
	// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
	//
	// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
	// result is nil.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil.
	SetIfNotExistFunc(key interface{}, f Func, duration time.Duration) (ok bool, err error)

	// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
	// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil.
	//
	// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
	// writing mutex lock for concurrent safety purpose.
	SetIfNotExistFuncLock(key interface{}, f Func, duration time.Duration) (ok bool, err error)

	// Get retrieves and returns the associated value of given `key`.
	// It returns nil if it does not exist, or its value is nil, or it's expired.
	// If you would like to check if the `key` exists in the cache, it's better using function Contains.
	Get(key interface{}) (*jvar.Var, error)

	// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
	// returns `value` if `key` does not exist in the cache. The key-value pair expires
	// after `duration`.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
	// if `value` is a function and the function result is nil.
	GetOrSet(key interface{}, value interface{}, duration time.Duration) (result *jvar.Var, err error)

	// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
	// function `f` and returns its result if `key` does not exist in the cache. The key-value
	// pair expires after `duration`.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
	// if `value` is a function and the function result is nil.
	GetOrSetFunc(key interface{}, f Func, duration time.Duration) (result *jvar.Var, err error)

	// GetOrSetFuncLock retrieves and returns the value of `key`, or sets `key` with result of
	// function `f` and returns its result if `key` does not exist in the cache. The key-value
	// pair expires after `duration`.
	//
	// It does not expire if `duration` == 0.
	// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
	// if `value` is a function and the function result is nil.
	//
	// Note that it differs from function `GetOrSetFunc` is that the function `f` is executed within
	// writing mutex lock for concurrent safety purpose.
	GetOrSetFuncLock(key interface{}, f Func, duration time.Duration) (result *jvar.Var, err error)

	// Contains checks and returns true if `key` exists in the cache, or else returns false.
	Contains(key interface{}) (bool, error)

	// Size returns the number of items in the cache.
	Size() (size int, err error)

	// Data returns a copy of all key-value pairs in the cache as map type.
	// Note that this function may lead lots of memory usage, you can implement this function
	// if necessary.
	Data() (data map[interface{}]interface{}, err error)

	// Keys returns all keys in the cache as slice.
	Keys() (keys []interface{}, err error)

	// Values returns all values in the cache as slice.
	Values() (values []interface{}, err error)

	// Update updates the value of `key` without changing its expiration and returns the old value.
	// The returned value `exist` is false if the `key` does not exist in the cache.
	//
	// It deletes the `key` if given `value` is nil.
	// It does nothing if `key` does not exist in the cache.
	Update(key interface{}, value interface{}) (oldValue *jvar.Var, exist bool, err error)

	// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
	//
	// It returns -1 and does nothing if the `key` does not exist in the cache.
	// It deletes the `key` if `duration` < 0.
	UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error)

	// GetExpire retrieves and returns the expiration of `key` in the cache.
	//
	// Note that,
	// It returns 0 if the `key` does not expire.
	// It returns -1 if the `key` does not exist in the cache.
	GetExpire(key interface{}) (time.Duration, error)

	// Remove deletes one or more keys from cache, and returns its value.
	// If multiple keys are given, it returns the value of the last deleted item.
	Remove(keys ...interface{}) (lastValue *jvar.Var, err error)

	// Clear clears all data of the cache.
	// Note that this function is sensitive and should be carefully used.
	Clear() error

	// Close closes the cache if necessary.
	Close() error
}
