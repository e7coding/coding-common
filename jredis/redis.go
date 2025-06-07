// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Redis 服务器提供便捷的客户端接口。
// Redis 官方命令文档: https://redis.io/commands
// Redis 中文文档: http://redisdoc.com/

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/errs/jerr"
)

const (
	errorNilRedis = `the Redis object is nil`
)

const errorNilAdapter = `redis adapter is not set, missing configuration or adapter register? possible reference: https://github.com/gogf/gf/tree/master/contrib/nosql/redis`

// AdapterFunc 是创建 Redis 适配器的函数类型。
type AdapterFunc func(config *Config) Adapter

var (
	// defaultAdapterFunc 是用于创建默认 Redis 适配器的函数。
	defaultAdapterFunc AdapterFunc = func(config *Config) Adapter {
		return nil
	}
)

// New 创建并返回一个 Redis 客户端。
// 如果提供了配置参数，则使用该配置创建适配器；
// 否则尝试从全局配置中获取并创建适配器。
func New(config ...*Config) (*Redis, error) {
	var (
		usedConfig  *Config
		usedAdapter Adapter
	)
	if len(config) > 0 && config[0] != nil {
		usedConfig = config[0]
		usedAdapter = defaultAdapterFunc(config[0])
	} else if configFromGlobal, ok := GetConfig(); ok {
		usedConfig = configFromGlobal
		usedAdapter = defaultAdapterFunc(configFromGlobal)
	}
	if usedConfig == nil {
		return nil, jerr.WithMsgF("未找到用于创建 Redis 客户端的配置")
	}
	if usedAdapter == nil {
		return nil, jerr.WithMsgF(errorNilAdapter)
	}
	redis := &Redis{
		config:       usedConfig,
		localAdapter: usedAdapter,
	}
	return redis.initGroup(), nil
}

// NewWithAdapter 使用给定的适配器创建并返回一个 Redis 客户端。
func NewWithAdapter(adapter Adapter) (*Redis, error) {
	if adapter == nil {
		return nil, jerr.WithMsgF("adapter 不能为空")
	}
	redis := &Redis{localAdapter: adapter}
	return redis.initGroup(), nil
}

// RegisterAdapterFunc 注册用于创建 Redis 适配器的默认函数。
func RegisterAdapterFunc(adapterFunc AdapterFunc) {
	defaultAdapterFunc = adapterFunc
}

// Redis client.
type Redis struct {
	config *Config
	localAdapter
	localGroup
}

type (
	localGroup struct {
		localGroupGeneric
		localGroupHash
		localGroupList
		localGroupPubSub
		localGroupScript
		localGroupSet
		localGroupSortedSet
		localGroupString
	}
	localAdapter        = Adapter
	localGroupGeneric   = IGroupGeneric
	localGroupHash      = IGroupHash
	localGroupList      = IGroupList
	localGroupPubSub    = IGroupPubSub
	localGroupScript    = IGroupScript
	localGroupSet       = IGroupSet
	localGroupSortedSet = IGroupSortedSet
	localGroupString    = IGroupStr
)

// initGroup initializes the group object of redis.
func (r *Redis) initGroup() *Redis {
	r.localGroup = localGroup{
		localGroupGeneric:   r.localAdapter.GroupGeneric(),
		localGroupHash:      r.localAdapter.GroupHash(),
		localGroupList:      r.localAdapter.GroupList(),
		localGroupPubSub:    r.localAdapter.GroupPubSub(),
		localGroupScript:    r.localAdapter.GroupScript(),
		localGroupSet:       r.localAdapter.GroupSet(),
		localGroupSortedSet: r.localAdapter.SortedSet(),
		localGroupString:    r.localAdapter.GroupStr(),
	}
	return r
}

// SetAdapter changes the underlying adapter with custom adapter for current redis client.
func (r *Redis) SetAdapter(adapter Adapter) {
	if r == nil {
		panic(jerr.WithMsg(errorNilRedis))
	}
	r.localAdapter = adapter
}

// GetAdapter returns the adapter that is set in current redis client.
func (r *Redis) GetAdapter() Adapter {
	if r == nil {
		return nil
	}
	return r.localAdapter
}

// Conn retrieves and returns a connection object for continuous operations.
// Note that you should call Close function manually if you do not use this connection any further.
func (r *Redis) Conn() (Conn, error) {
	if r == nil {
		return nil, jerr.WithMsg(errorNilRedis)
	}
	if r.localAdapter == nil {
		return nil, jerr.WithMsg(errorNilAdapter)
	}
	return r.localAdapter.Conn()
}

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (r *Redis) Do(command string, args ...interface{}) (*jvar.Var, error) {
	if r == nil {
		return nil, jerr.WithMsg(errorNilRedis)
	}
	if r.localAdapter == nil {
		return nil, jerr.WithMsg(errorNilAdapter)
	}
	return r.localAdapter.Do(command, args...)
}

// Close closes current redis client, closes its connection pool and releases all its related resources.
func (r *Redis) Close() error {
	if r == nil || r.localAdapter == nil {
		return nil
	}
	return r.localAdapter.Close()
}
