// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jcache

import (
	"github.com/e7coding/coding-common/jredis"
	"time"

	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// AdapterRedis is the dcache adapter implements using Redis server.
type AdapterRedis struct {
	redis *jredis.Redis
}

// NewAdapterRedis creates and returns a new memory cache object.
func NewAdapterRedis(redis *jredis.Redis) *AdapterRedis {
	return &AdapterRedis{
		redis: redis,
	}
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) Set(key interface{}, value interface{}, duration time.Duration) (err error) {
	redisKey := jconv.String(key)
	if value == nil || duration < 0 {
		_, err = c.redis.Del(redisKey)
	} else {
		if duration == 0 {
			_, err = c.redis.Set(redisKey, value)
		} else {
			_, err = c.redis.Set(redisKey, value, jredis.SetOption{TTLOption: jredis.TTLOption{PX: jconv.PtrInt64(duration.Milliseconds())}})
		}
	}
	return err
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetMap(data map[interface{}]interface{}, duration time.Duration) error {
	if len(data) == 0 {
		return nil
	}
	// DEL.
	if duration < 0 {
		var (
			index = 0
			keys  = make([]string, len(data))
		)
		for k := range data {
			keys[index] = jconv.String(k)
			index += 1
		}
		_, err := c.redis.Del(keys...)
		if err != nil {
			return err
		}
	}
	if duration == 0 {
		err := c.redis.MSet(jconv.Map(data))
		if err != nil {
			return err
		}
	}
	if duration > 0 {
		var err error
		for k, v := range data {
			if err = c.Set(k, v, duration); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	var (
		err      error
		redisKey = jconv.String(key)
	)
	// Execute the function and retrieve the result.
	f, ok := value.(Func)
	if !ok {
		// Compatible with raw function value.
		f, ok = value.(func() (value interface{}, err error))
	}
	if ok {
		if value, err = f(); err != nil {
			return false, err
		}
	}
	// DEL.
	if duration < 0 || value == nil {
		var delResult int64
		delResult, err = c.redis.Del(redisKey)
		if err != nil {
			return false, err
		}
		if delResult == 1 {
			return true, err
		}
		return false, err
	}
	ok, err = c.redis.SetNX(redisKey, value)
	if err != nil {
		return ok, err
	}
	if ok && duration > 0 {
		// Set the expiration.
		_, err = c.redis.PExpire(redisKey, duration.Milliseconds())
		if err != nil {
			return ok, err
		}
		return ok, err
	}
	return ok, err
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterRedis) SetIfNotExistFunc(key interface{}, f Func, duration time.Duration) (ok bool, err error) {
	value, err := f()
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(key, value, duration)
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterRedis) SetIfNotExistFuncLock(key interface{}, f Func, duration time.Duration) (ok bool, err error) {
	value, err := f()
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(key, value, duration)
}

// Get retrieves and returns the associated value of given <key>.
// It returns nil if it does not exist or its value is nil.
func (c *AdapterRedis) Get(key interface{}) (*jvar.Var, error) {
	return c.redis.Get(jconv.String(key))
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterRedis) GetOrSet(key interface{}, value interface{}, duration time.Duration) (result *jvar.Var, err error) {
	result, err = c.Get(key)
	if err != nil {
		return nil, err
	}
	if result.IsNil() {
		return jvar.New(value), c.Set(key, value, duration)
	}
	return
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterRedis) GetOrSetFunc(key interface{}, f Func, duration time.Duration) (result *jvar.Var, err error) {
	v, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		value, err := f()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return jvar.New(value), c.Set(key, value, duration)
	} else {
		return v, nil
	}
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
func (c *AdapterRedis) GetOrSetFuncLock(key interface{}, f Func, duration time.Duration) (result *jvar.Var, err error) {
	return c.GetOrSetFunc(key, f, duration)
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *AdapterRedis) Contains(key interface{}) (bool, error) {
	n, err := c.redis.Exists(jconv.String(key))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Size returns the number of items in the cache.
func (c *AdapterRedis) Size() (size int, err error) {
	n, err := c.redis.DBSize()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// Data returns a copy of all key-value pairs in the cache as map type.
// Note that this function may lead lots of memory usage, you can implement this function
// if necessary.
func (c *AdapterRedis) Data() (map[interface{}]interface{}, error) {
	// Keys.
	keys, err := c.redis.Keys("*")
	if err != nil {
		return nil, err
	}
	// Key-Value pairs.
	var m map[string]*jvar.Var
	m, err = c.redis.MGet(keys...)
	if err != nil {
		return nil, err
	}
	// Type converting.
	data := make(map[interface{}]interface{})
	for k, v := range m {
		data[k] = v.Val()
	}
	return data, nil
}

// Keys returns all keys in the cache as slice.
func (c *AdapterRedis) Keys() ([]interface{}, error) {
	keys, err := c.redis.Keys("*")
	if err != nil {
		return nil, err
	}
	return jconv.Interfaces(keys), nil
}

// Values returns all values in the cache as slice.
func (c *AdapterRedis) Values() ([]interface{}, error) {
	// Keys.
	keys, err := c.redis.Keys("*")
	if err != nil {
		return nil, err
	}
	// Key-Value pairs.
	var m map[string]*jvar.Var
	m, err = c.redis.MGet(keys...)
	if err != nil {
		return nil, err
	}
	// Values.
	var values []interface{}
	for _, key := range keys {
		if v := m[key]; !v.IsNil() {
			values = append(values, v.Val())
		}
	}
	return values, nil
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *AdapterRedis) Update(key interface{}, value interface{}) (oldValue *jvar.Var, exist bool, err error) {
	var (
		v        *jvar.Var
		oldPTTL  int64
		redisKey = jconv.String(key)
	)
	// TTL.
	oldPTTL, err = c.redis.PTTL(redisKey) // update ttl -> pttl(millisecond)
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		// It does not exist or expired.
		return
	}
	// Check existence.
	v, err = c.redis.Get(redisKey)
	if err != nil {
		return
	}
	oldValue = v
	// DEL.
	if value == nil {
		_, err = c.redis.Del(redisKey)
		if err != nil {
			return
		}
		return
	}
	// Update the value.
	if oldPTTL == -1 {
		_, err = c.redis.Set(redisKey, value)
	} else {
		// update SetEX -> SET PX Option(millisecond)
		// Starting with Redis version 2.6.12: Added the EX, PX, NX and XX options.
		_, err = c.redis.Set(redisKey, value, jredis.SetOption{TTLOption: jredis.TTLOption{PX: jconv.PtrInt64(oldPTTL)}})
	}
	return oldValue, true, err
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *AdapterRedis) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	var (
		v        *jvar.Var
		oldPTTL  int64
		redisKey = jconv.String(key)
	)
	// TTL.
	oldPTTL, err = c.redis.PTTL(redisKey)
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		// It does not exist or expired.
		oldPTTL = -1
		return
	}
	oldDuration = time.Duration(oldPTTL) * time.Millisecond
	// DEL.
	if duration < 0 {
		_, err = c.redis.Del(redisKey)
		return
	}
	// Update the expiration.
	if duration > 0 {
		_, err = c.redis.PExpire(redisKey, duration.Milliseconds())
	}
	// No expire.
	if duration == 0 {
		v, err = c.redis.Get(redisKey)
		if err != nil {
			return
		}
		_, err = c.redis.Set(redisKey, v.Val())
	}
	return
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (c *AdapterRedis) GetExpire(key interface{}) (time.Duration, error) {
	pttl, err := c.redis.PTTL(jconv.String(key))
	if err != nil {
		return 0, err
	}
	switch pttl {
	case -1:
		return 0, nil
	case -2, 0: // It does not exist or expired.
		return -1, nil
	default:
		return time.Duration(pttl) * time.Millisecond, nil
	}
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func (c *AdapterRedis) Remove(keys ...interface{}) (lastValue *jvar.Var, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	// Retrieves the last key value.
	if lastValue, err = c.redis.Get(jconv.String(keys[len(keys)-1])); err != nil {
		return nil, err
	}
	// Deletes all given keys.
	_, err = c.redis.Del(jconv.Strings(keys)...)
	return
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
// It uses `FLUSHDB` command in redis server, which might be disabled in server.
func (c *AdapterRedis) Clear() (err error) {
	// The "FLUSHDB" may not be available.
	err = c.redis.FlushDB()
	return
}

// Close closes the cache.
func (c *AdapterRedis) Close() error {
	// It does nothing.
	return nil
}
