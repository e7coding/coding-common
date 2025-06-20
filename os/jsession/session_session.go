// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jsession

import (
	"errors"
	"github.com/e7coding/coding-common/errs/jerr"
	"time"

	"github.com/e7coding/coding-common/container/jmap"
	"github.com/e7coding/coding-common/container/jvar"

	"github.com/e7coding/coding-common/internal/intlog"
)

// Session struct for storing single session data, which is bound to a single request.
// The Session struct is the interface with user, but the Storage is the underlying adapter designed interface
// for functionality implements.
type Session struct {
	id      string          // Session id. It retrieves the session if id is custom specified.
	data    *jmap.StrAnyMap // Current Session data, which is retrieved from Storage.
	dirty   bool            // Used to mark session is modified.
	start   bool            // Used to mark session is started.
	manager *Manager        // Parent session Manager.

	// idFunc is a callback function used for creating custom session id.
	// This is called if session id is empty ever when session starts.
	idFunc func(ttl time.Duration) (id string)
}

// init does the lazy initialization for session, which retrieves the session if session id is specified,
// or else it creates a new empty session.
func (s *Session) init() error {
	if s.start {
		return nil
	}
	var err error
	// Session retrieving.
	if s.id != "" {
		// Retrieve stored session data from storage.
		if s.manager.storage != nil {
			s.data, err = s.manager.storage.GetSession(s.id, s.manager.GetTTL())
			if err != nil {
				intlog.Errorf(`session restoring failed for id "%s": %+v`, s.id, err)
				return err
			}
		}
	}
	// Session id creation.
	if s.id == "" {
		if s.idFunc != nil {
			// Use custom session id creating function.
			s.id = s.idFunc(s.manager.ttl)
		} else {
			// Use default session id creating function of storage.
			s.id, err = s.manager.storage.New(s.manager.ttl)
			if err != nil && !errors.Is(err, ErrorDisabled) {
				intlog.Errorf("create session id failed: %+v", err)
				return err
			}
			// If session storage does not implements id generating functionality,
			// it then uses default session id creating function.
			if s.id == "" {
				s.id = NewSessionId()
			}
		}
	}
	if s.data == nil {
		s.data = jmap.NewStrAnyMap()
	}
	s.start = true
	return nil
}

// Close closes current session and updates its ttl in the session manager.
// If this session is dirty, it also exports it to storage.
//
// NOTE that this function must be called ever after a session request done.
func (s *Session) Close() error {
	if s.manager.storage == nil {
		return nil
	}
	if s.start && s.id != "" {
		size := s.data.Len()
		if s.dirty {
			err := s.manager.storage.SetSession(s.id, s.data, s.manager.ttl)
			if err != nil {
				return err
			}
		} else if size > 0 {
			err := s.manager.storage.UpdateTTL(s.id, s.manager.ttl)
			if err != nil && !errors.Is(err, ErrorDisabled) {
				return err
			}
		}
	}
	return nil
}

// Set sets key-value pair to this session.
func (s *Session) Set(key string, value interface{}) (err error) {
	if err = s.init(); err != nil {
		return err
	}
	if err = s.manager.storage.Set(s.id, key, value, s.manager.ttl); err != nil {
		if !errors.Is(err, ErrorDisabled) {
			return err
		}
		s.data.Put(key, value)
	}
	s.dirty = true
	return nil
}

// SetMap batch sets the session using map.
func (s *Session) SetMap(data map[string]interface{}) (err error) {
	if err = s.init(); err != nil {
		return err
	}
	if err = s.manager.storage.SetMap(s.id, data, s.manager.ttl); err != nil {
		s.data.PutAll(data)
	}
	s.dirty = true
	return nil
}

// Remove removes key along with its value from this session.
func (s *Session) Remove(keys ...string) (err error) {
	if s.id == "" {
		return nil
	}
	if err = s.init(); err != nil {
		return err
	}
	for _, key := range keys {
		if err = s.manager.storage.Remove(s.id, key); err != nil {
			s.data.Del(key)
		}
	}
	s.dirty = true
	return nil
}

// RemoveAll deletes all key-value pairs from this session.
func (s *Session) RemoveAll() (err error) {
	if s.id == "" {
		return nil
	}
	if err = s.init(); err != nil {
		return err
	}
	if err = s.manager.storage.RemoveAll(s.id); err != nil {
		if !errors.Is(err, ErrorDisabled) {
			return err
		}
	}
	// Remove data from memory.
	if s.data != nil {
		s.data.Empty()
	}
	s.dirty = true
	return nil
}

// Id returns the session id for this session.
// It creates and returns a new session id if the session id is not passed in initialization.
func (s *Session) Id() (id string, err error) {
	if err = s.init(); err != nil {
		return "", err
	}
	return s.id, nil
}

// SetId sets custom session before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetId(id string) error {
	if s.start {
		return jerr.WithMsg("session already started")
	}
	s.id = id
	return nil
}

// SetIdFunc sets custom session id creating function before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetIdFunc(f func(ttl time.Duration) string) error {
	if s.start {
		return jerr.WithMsg("session already started")
	}
	s.idFunc = f
	return nil
}

// Data returns all data as map.
// Note that it's using value copy internally for concurrent-safe purpose.
func (s *Session) Data() (sessionData map[string]interface{}, err error) {
	if s.id == "" {
		return map[string]interface{}{}, nil
	}
	if err = s.init(); err != nil {
		return nil, err
	}
	sessionData, err = s.manager.storage.Data(s.id)
	if err != nil && !errors.Is(err, ErrorDisabled) {
		intlog.Errorf(`%+v`, err)
	}
	if sessionData != nil {
		return sessionData, nil
	}
	return s.data.ToMap(), nil
}

// Size returns the size of the session.
func (s *Session) Size() (size int, err error) {
	if s.id == "" {
		return 0, nil
	}
	if err = s.init(); err != nil {
		return 0, err
	}
	size, err = s.manager.storage.GetSize(s.id)
	if err != nil && !errors.Is(err, ErrorDisabled) {
		intlog.Errorf(`%+v`, err)
	}
	if size > 0 {
		return size, nil
	}
	return s.data.Len(), nil
}

// Contains checks whether key exist in the session.
func (s *Session) Contains(key string) (ok bool, err error) {
	if s.id == "" {
		return false, nil
	}
	if err = s.init(); err != nil {
		return false, err
	}
	v, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return !v.IsNil(), nil
}

// IsDirty checks whether there's any data changes in the session.
func (s *Session) IsDirty() bool {
	return s.dirty
}

// Get retrieves session value with given key.
// It returns `def` if the key does not exist in the session if `def` is given,
// or else it returns nil.
func (s *Session) Get(key string, def ...interface{}) (value *jvar.Var, err error) {
	if s.id == "" {
		return nil, nil
	}
	if err = s.init(); err != nil {
		return nil, err
	}
	v, err := s.manager.storage.Get(s.id, key)
	if err != nil && !errors.Is(err, ErrorDisabled) {
		intlog.Errorf(`%+v`, err)
		return nil, err
	}
	if v != nil {
		return jvar.New(v), nil
	}
	if v = s.data.Get(key); v != nil {
		return jvar.New(v), nil
	}
	if len(def) > 0 {
		return jvar.New(def[0]), nil
	}
	return nil, nil
}

// RegenerateId regenerates a new session id for current session.
// It keeps the session data and updates the session id with a new one.
// This is commonly used to prevent session fixation attacks and increase security.
//
// The parameter `deleteOld` specifies whether to delete the old session data:
// - If true: the old session data will be deleted immediately
// - If false: the old session data will be kept and expire according to its TTL
func (s *Session) RegenerateId(deleteOld bool) (newId string, err error) {
	if err = s.init(); err != nil {
		return "", err
	}

	// Generate new session id
	if s.idFunc != nil {
		newId = s.idFunc(s.manager.ttl)
	} else {
		newId, err = s.manager.storage.New(s.manager.ttl)
		if err != nil && !errors.Is(err, ErrorDisabled) {
			return "", err
		}
		if newId == "" {
			newId = NewSessionId()
		}
	}

	// If using storage, need to copy data to new id
	if s.manager.storage != nil {
		if err = s.manager.storage.SetSession(newId, s.data, s.manager.ttl); err != nil {
			if !errors.Is(err, ErrorDisabled) {
				return "", err
			}
		}
		// Delete old session data if requested
		if deleteOld {
			if err = s.manager.storage.RemoveAll(s.id); err != nil {
				if !errors.Is(err, ErrorDisabled) {
					return "", err
				}
			}
		}
	}

	// Update session id
	s.id = newId
	s.dirty = true
	return newId, nil
}

// MustRegenerateId performs as function RegenerateId, but it panics if any error occurs.
func (s *Session) MustRegenerateId(deleteOld bool) string {
	newId, err := s.RegenerateId(deleteOld)
	if err != nil {
		panic(err)
	}
	return newId
}
