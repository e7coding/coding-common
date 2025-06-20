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
)

// StorageRedisHashTable implements the Session Storage interface with redis hash table.
type StorageRedisHashTable struct {
	StorageBase
	redis  *jredis.Redis // Redis client for session storage.
	prefix string        // Redis key prefix for session id.
}

// NewStorageRedisHashTable creates and returns a redis hash table storage object for session.
func NewStorageRedisHashTable(redis *jredis.Redis, prefix ...string) *StorageRedisHashTable {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedisHashTable{
		redis: redis,
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	return s
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageRedisHashTable) Get(sessionId string, key string) (value interface{}, err error) {
	v, err := s.redis.HGet(s.sessionIdToRedisKey(sessionId), key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	return v.String(), nil
}

// Data retrieves all key-value pairs as map from storage.
func (s *StorageRedisHashTable) Data(sessionId string) (data map[string]interface{}, err error) {
	m, err := s.redis.HGetAll(s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	return m.Map(), nil
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageRedisHashTable) GetSize(sessionId string) (size int, err error) {
	v, err := s.redis.HLen(s.sessionIdToRedisKey(sessionId))
	return int(v), err
}

// Set sets key-value session pair to the storage.
// The parameter `ttl` specifies the TTL for the session id (not for the key-value pair).
func (s *StorageRedisHashTable) Set(sessionId string, key string, value interface{}, ttl time.Duration) error {
	_, err := s.redis.HSet(s.sessionIdToRedisKey(sessionId), map[string]interface{}{
		key: value,
	})
	return err
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter `ttl` specifies the TTL for the session id(not for the key-value pair).
func (s *StorageRedisHashTable) SetMap(sessionId string, data map[string]interface{}, ttl time.Duration) error {
	err := s.redis.HMSet(s.sessionIdToRedisKey(sessionId), data)
	return err
}

// Remove deletes key with its value from storage.
func (s *StorageRedisHashTable) Remove(sessionId string, key string) error {
	_, err := s.redis.HDel(s.sessionIdToRedisKey(sessionId), key)
	return err
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedisHashTable) RemoveAll(sessionId string) error {
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
func (s *StorageRedisHashTable) GetSession(sessionId string, ttl time.Duration) (*jmap.StrAnyMap, error) {
	intlog.Printf("StorageRedisHashTable.GetSession: %s, %v", sessionId, ttl)
	v, err := s.redis.Exists(s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	if v > 0 {
		// It does not store the session data in memory, it so returns an empty map.
		// It retrieves session data items directly through redis server each time.
		return jmap.NewStrAnyMap(), nil
	}
	return nil, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedisHashTable) SetSession(sessionId string, sessionData *jmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf("StorageRedisHashTable.SetSession: %s, %v", sessionId, ttl)
	_, err := s.redis.Expire(s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedisHashTable) UpdateTTL(sessionId string, ttl time.Duration) error {
	intlog.Printf("StorageRedisHashTable.UpdateTTL: %s, %v", sessionId, ttl)
	_, err := s.redis.Expire(s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

// sessionIdToRedisKey converts and returns the redis key for given session id.
func (s *StorageRedisHashTable) sessionIdToRedisKey(sessionId string) string {
	return s.prefix + sessionId
}
