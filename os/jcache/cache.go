// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dcache provides kinds of cache management for process.
//
// It provides a concurrent-safe in-memory cache adapter for process in default.
package jcache

import (
	"context"
	"github.com/e7coding/coding-common/container/jvar"
	"time"
)

// Func is the cache function that calculates and returns the value.
type Func = func(ctx context.Context) (value interface{}, err error)

// DurationNoExpire represents the cache key-value pair that never expires.
const DurationNoExpire = time.Duration(0)

// Default cache object.
var defaultCache = New()

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	return defaultCache.Set(ctx, key, value, duration)
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	return defaultCache.SetMap(ctx, data, duration)
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExist(ctx, key, value, duration)
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func SetIfNotExistFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExistFunc(ctx, key, f, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func SetIfNotExistFuncLock(ctx context.Context, key interface{}, f Func, duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExistFuncLock(ctx, key, f, duration)
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func Get(ctx context.Context, key interface{}) (*jvar.Var, error) {
	return defaultCache.Get(ctx, key)
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (*jvar.Var, error) {
	return defaultCache.GetOrSet(ctx, key, value, duration)
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func GetOrSetFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) (*jvar.Var, error) {
	return defaultCache.GetOrSetFunc(ctx, key, f, duration)
}

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
func GetOrSetFuncLock(ctx context.Context, key interface{}, f Func, duration time.Duration) (*jvar.Var, error) {
	return defaultCache.GetOrSetFuncLock(ctx, key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func Contains(ctx context.Context, key interface{}) (bool, error) {
	return defaultCache.Contains(ctx, key)
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return defaultCache.GetExpire(ctx, key)
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func Remove(ctx context.Context, keys ...interface{}) (value *jvar.Var, err error) {
	return defaultCache.Remove(ctx, keys...)
}

// Removes deletes `keys` in the cache.
func Removes(ctx context.Context, keys []interface{}) error {
	return defaultCache.Removes(ctx, keys)
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func Update(ctx context.Context, key interface{}, value interface{}) (oldValue *jvar.Var, exist bool, err error) {
	return defaultCache.Update(ctx, key, value)
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return defaultCache.UpdateExpire(ctx, key, duration)
}

// Size returns the number of items in the cache.
func Size(ctx context.Context) (int, error) {
	return defaultCache.Size(ctx)
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func Data(ctx context.Context) (map[interface{}]interface{}, error) {
	return defaultCache.Data(ctx)
}

// Keys returns all keys in the cache as slice.
func Keys(ctx context.Context) ([]interface{}, error) {
	return defaultCache.Keys(ctx)
}

// KeyStrings returns all keys in the cache as string slice.
func KeyStrings(ctx context.Context) ([]string, error) {
	return defaultCache.KeyStrings(ctx)
}

// Values returns all values in the cache as slice.
func Values(ctx context.Context) ([]interface{}, error) {
	return defaultCache.Values(ctx)
}
