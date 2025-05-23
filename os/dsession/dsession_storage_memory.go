// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dsession

import (
	"context"
	"time"

	"github.com/coding-common/container/dvar"
	"github.com/coding-common/container/wmap"
	"github.com/coding-common/os/dcache"
)

// StorageMemory implements the Session Storage interface with memory.
type StorageMemory struct {
	StorageBase
	// cache is the memory data cache for session TTL,
	// which is available only if the Storage does not store any session data in synchronizing.
	// Please refer to the implements of StorageFile, StorageMemory and StorageRedis.
	//
	// Its value is type of `*gmap.StrAnyMap`.
	cache *dcache.Cache
}

// NewStorageMemory creates and returns a file storage object for session.
func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		cache: dcache.New(),
	}
}

// RemoveAll deletes session from storage.
func (s *StorageMemory) RemoveAll(ctx context.Context, sessionId string) error {
	_, err := s.cache.Remove(ctx, sessionId)
	return err
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageMemory) GetSession(ctx context.Context, sessionId string, ttl time.Duration) (*wmap.SafeStrAnyMap, error) {
	// Retrieve memory session data from manager.
	var (
		v   *dvar.Var
		err error
	)
	v, err = s.cache.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if v != nil {
		return v.Val().(*wmap.SafeStrAnyMap), nil
	}
	return wmap.NewSafeStrAnyMap(), nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageMemory) SetSession(ctx context.Context, sessionId string, sessionData *wmap.StrAnyMap, ttl time.Duration) error {
	return s.cache.Set(ctx, sessionId, sessionData, ttl)
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageMemory) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	_, err := s.cache.UpdateExpire(ctx, sessionId, ttl)
	return err
}
