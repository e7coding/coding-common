// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"fmt"
)

// IGroupPubSub manages redis pub/sub operations.
// Implements see jredis.GroupPubSub.
type IGroupPubSub interface {
	Publish(channel string, message interface{}) (int64, error)
	Subscribe(channel string, channels ...string) (Conn, []*Subscription, error)
	PSubscribe(pattern string, patterns ...string) (Conn, []*Subscription, error)
}

// Message received as result of a PUBLISH command issued by another client.
type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}

// Subscription received after a successful subscription to channel.
type Subscription struct {
	Kind    string // Can be "subscribe", "unsubscribe", "psubscribe" or "punsubscribe".
	Channel string // Channel name we have subscribed to.
	Count   int    // Number of channels we are currently subscribed to.
}

// String converts current object to a readable string.
func (m *Subscription) String() string {
	return fmt.Sprintf("%s: %s", m.Kind, m.Channel)
}
