// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package tracing provides some utility functions for tracing functionality.
package tracing

import (
	"github.com/e7coding/coding-common/container/jatomic"
	"math"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/e7coding/coding-common/encoding/jbinary"
	"github.com/e7coding/coding-common/jutil/jrand"
)

var (
	randomInitSequence = int32(jrand.Intn(math.MaxInt32))
	sequence           = jatomic.NewInt32(randomInitSequence)
)

// NewIDs creates and returns a new trace and span ID.
func NewIDs() (traceID trace.TraceID, spanID trace.SpanID) {
	return NewTraceID(), NewSpanID()
}

// NewTraceID creates and returns a trace ID.
func NewTraceID() (traceID trace.TraceID) {
	var (
		timestampNanoBytes = dbinary.EncodeInt64(time.Now().UnixNano())
		sequenceBytes      = dbinary.EncodeInt32(sequence.Add(1))
		randomBytes        = jrand.B(4)
	)
	copy(traceID[:], timestampNanoBytes)
	copy(traceID[8:], sequenceBytes)
	copy(traceID[12:], randomBytes)
	return
}

// NewSpanID creates and returns a span ID.
func NewSpanID() (spanID trace.SpanID) {
	copy(spanID[:], dbinary.EncodeInt64(time.Now().UnixNano()/1e3))
	copy(spanID[4:], jrand.B(4))
	return
}
