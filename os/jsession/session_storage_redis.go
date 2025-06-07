// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jsession

import (
	"github.com/e7coding/coding-common/jredis"
	"time"

	"github.com/e7coding/coding-common/container/jmap"
	"github.com/e7coding/coding-common/internal/intlog"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/os/jtimer"
)

// StorageRedis implements the Session Storage interface with redis.
type StorageRedis struct {
	StorageBase
	redis         *jredis.Redis       // Redis client for session storage.
	prefix        string              // Redis key prefix for session id.
	updatingIdMap *jmap.SafeStrIntMap // Updating TTL set for session id.
}

const (
	// DefaultStorageRedisLoopInterval is the interval updating TTL for session ids
	// in last duration.
	DefaultStorageRedisLoopInterval = 10 * time.Second
)

// NewStorageRedis creates and returns a redis storage object for session.
func NewStorageRedis(redis *jredis.Redis, prefix ...string) *StorageRedis {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedis{
		redis:         redis,
		updatingIdMap: jmap.NewSafeStrIntMap(),
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	// Batch updates the TTL for session ids timely.
	jtimer.AddSingleton(DefaultStorageRedisLoopInterval, func() {
		intlog.Print("StorageRedis.timer start")
		var (
			err        error
			sessionId  string
			ttlSeconds int
		)
		for {
			if sessionId, ttlSeconds = s.updatingIdMap.Pop(); sessionId == "" {
				break
			} else {
				if err = s.doUpdateExpireForSession(sessionId, ttlSeconds); err != nil {
					intlog.Errorf(`%+v`, err)
				}
			}
		}
		intlog.Print("StorageRedis.timer end")
	})
	return s
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedis) RemoveAll(sessionId string) error {
	_, err := s.redis.Del(s.sessionIdToRedisKey(sessionId))
	return err
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageRedis) GetSession(sessionId string, ttl time.Duration) (*jmap.StrAnyMap, error) {
	intlog.Printf("StorageRedis.GetSession: %s, %v", sessionId, ttl)
	r, err := s.redis.Get(s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	content := r.Bytes()
	if len(content) == 0 {
		return nil, nil
	}
	var m map[string]interface{}
	if err = json.UnmarshalUseNumber(content, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return jmap.NewStrAnyMapFrom(m), nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedis) SetSession(sessionId string, sessionData *jmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf("StorageRedis.SetSession: %s, %v, %v", sessionId, sessionData, ttl)
	content, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	err = s.redis.SetEX(s.sessionIdToRedisKey(sessionId), content, int64(ttl.Seconds()))
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedis) UpdateTTL(sessionId string, ttl time.Duration) error {
	intlog.Printf("StorageRedis.UpdateTTL: %s, %v", sessionId, ttl)
	if ttl >= DefaultStorageRedisLoopInterval {
		s.updatingIdMap.Put(sessionId, int(ttl.Seconds()))
	}
	return nil
}

// doUpdateExpireForSession updates the TTL for session id.
func (s *StorageRedis) doUpdateExpireForSession(sessionId string, ttlSeconds int) error {
	intlog.Printf("StorageRedis.doUpdateTTL: %s, %d", sessionId, ttlSeconds)
	_, err := s.redis.Expire(s.sessionIdToRedisKey(sessionId), int64(ttlSeconds))
	return err
}

// sessionIdToRedisKey converts and returns the redis key for given session id.
func (s *StorageRedis) sessionIdToRedisKey(sessionId string) string {
	return s.prefix + sessionId
}
