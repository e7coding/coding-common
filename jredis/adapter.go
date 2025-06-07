// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
)

// Adapter 定义通用的 Redis 操作接口，包括组操作和命令操作
type Adapter interface {
	AdapterGroup
	AdapterOpts
}

// AdapterGroup 定义 Redis 分组操作接口
type AdapterGroup interface {
	GroupGeneric
	GroupHash
	GroupList
	GroupPubSub
	GroupScript
	GroupSet
	GroupSortedSet
	GroupStr
}

type GroupGeneric interface {
	GroupGeneric() IGroupGeneric // 获取通用命令分组
}
type GroupHash interface {
	GroupHash() IGroupHash // 获取哈希命令分组
}
type GroupList interface {
	GroupList() IGroupList // 获取列表命令分组
}
type GroupPubSub interface {
	GroupPubSub() IGroupPubSub // 获取发布/订阅命令分组
}
type GroupScript interface {
	GroupScript() IGroupScript // 获取脚本命令分组
}
type GroupSet interface {
	GroupSet() IGroupSet // 获取集合命令分组
}
type GroupSortedSet interface {
	SortedSet() IGroupSortedSet // 获取有序集合命令分组
}
type GroupStr interface {
	GroupStr() IGroupStr // 获取字符串命令分组
}

// AdapterOpts 定义核心的 Redis 命令操作接口，可由自定义实现覆盖
type AdapterOpts interface {
	// Do 发送命令到服务器并返回结果，自动对结构体/切片/映射类型进行 JSON 编码
	Do(command string, args ...interface{}) (*jvar.Var, error)

	// Conn 获取一个可持续操作的连接，需手动调用 Close 归还或关闭
	Conn() (conn Conn, err error)

	// Close 关闭当前客户端，释放连接池及所有相关资源
	Close() error
}

// Conn 定义从通用客户端获取的连接接口
type Conn interface {
	ConnCmd

	// Do 在连接上发送命令并返回结果，自动对复杂类型进行 JSON 编码
	Do(command string, args ...interface{}) (*jvar.Var, error)

	// Close 将连接归还连接池或关闭连接
	Close() error
}

// ConnCmd 定义针对特定连接的命令操作接口
type ConnCmd interface {
	// Subscribe 订阅指定频道
	Subscribe(channel string, channels ...string) ([]*Subscription, error)

	// PSubscribe 按模式订阅频道
	// 支持通配符: ?, *, [字符集]
	PSubscribe(pattern string, patterns ...string) ([]*Subscription, error)

	// ReceiveMessage 接收一条发布/订阅消息
	ReceiveMessage() (*Message, error)

	// Receive 接收一条命令回复
	Receive() (*jvar.Var, error)
}
