// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtrace

import (
	"context"
	"github.com/coding-common/container/wmap"

	"go.opentelemetry.io/otel/baggage"

	"github.com/coding-common/container/dvar"
	"github.com/coding-common/util/dconv"
)

// Baggage holds the data through all tracing spans.
type Baggage struct {
	ctx context.Context
}

// NewBaggage creates and returns a new Baggage object from given tracing context.
func NewBaggage(ctx context.Context) *Baggage {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Baggage{
		ctx: ctx,
	}
}

// Ctx returns the context that Baggage holds.
func (b *Baggage) Ctx() context.Context {
	return b.ctx
}

// SetValue is a convenient function for adding one key-value pair to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func (b *Baggage) SetValue(key string, value interface{}) context.Context {
	member, _ := baggage.NewMember(key, dconv.String(value))
	bag, _ := baggage.New(member)
	b.ctx = baggage.ContextWithBaggage(b.ctx, bag)
	return b.ctx
}

// SetMap is a convenient function for adding map key-value pairs to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func (b *Baggage) SetMap(data map[string]interface{}) context.Context {
	members := make([]baggage.Member, 0)
	for k, v := range data {
		member, _ := baggage.NewMember(k, dconv.String(v))
		members = append(members, member)
	}
	bag, _ := baggage.New(members...)
	b.ctx = baggage.ContextWithBaggage(b.ctx, bag)
	return b.ctx
}

// GetMap retrieves and returns the baggage values as map.
func (b *Baggage) GetMap() *wmap.StrAnyMap {
	m := wmap.NewStrAnyMap()
	members := baggage.FromContext(b.ctx).Members()
	for i := range members {
		m.Put(members[i].Key(), members[i].Value())
	}
	return m
}

// GetVar retrieves value and returns a *dvar.Var for specified key from baggage.
func (b *Baggage) GetVar(key string) *dvar.Var {
	value := baggage.FromContext(b.ctx).Member(key).Value()
	return dvar.New(value)
}
