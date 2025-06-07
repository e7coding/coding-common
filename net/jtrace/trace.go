// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtrace provides convenience wrapping functionality for tracing feature using OpenTelemetry.
package jtrace

import (
	"context"
	"github.com/e7coding/coding-common/container/jmap"
	"github.com/e7coding/coding-common/errs/jerr"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/e7coding/coding-common/container/jvar"

	"github.com/e7coding/coding-common/net/jipv4"
	"github.com/e7coding/coding-common/net/jtrace/internal/provider"
	"github.com/e7coding/coding-common/text/jstr"
)

const (
	tracingCommonKeyIpIntranet        = `ip.intranet`
	tracingCommonKeyIpHostname        = `hostname`
	commandEnvKeyForMaxContentLogSize = "gf.gtrace.max.content.log.size" // To avoid too big tracing content.
	commandEnvKeyForTracingInternal   = "gf.gtrace.tracing.internal"     // For detailed controlling for tracing content.
)

var (
	intranetIps, _           = jipv4.GetIntranetIpArray()
	intranetIpStr            = strings.Join(intranetIps, ",")
	hostname, _              = os.Hostname()
	tracingInternal          = true       // tracingInternal enables tracing for internal type spans.
	tracingMaxContentLogSize = 512 * 1024 // Max log size for request and response body, especially for HTTP/RPC request.
	// defaultTextMapPropagator is the default propagator for context propagation between peers.
	defaultTextMapPropagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
)

func init() {
	// Default trace provider.
	otel.SetTracerProvider(provider.New())
	CheckSetDefaultTextMapPropagator()
}

// IsUsingDefaultProvider checks and return if currently using default trace provider.
func IsUsingDefaultProvider() bool {
	_, ok := otel.GetTracerProvider().(*provider.TracerProvider)
	return ok
}

// MaxContentLogSize returns the max log size for request and response body, especially for HTTP/RPC request.
func MaxContentLogSize() int {
	return tracingMaxContentLogSize
}

// CommonLabels returns common used attribute labels:
// ip.intranet, hostname.
func CommonLabels() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String(tracingCommonKeyIpHostname, hostname),
		attribute.String(tracingCommonKeyIpIntranet, intranetIpStr),
		semconv.HostName(hostname),
	}
}

// CheckSetDefaultTextMapPropagator sets the default TextMapPropagator if it is not set previously.
func CheckSetDefaultTextMapPropagator() {
	p := otel.GetTextMapPropagator()
	if len(p.Fields()) == 0 {
		otel.SetTextMapPropagator(GetDefaultTextMapPropagator())
	}
}

// GetDefaultTextMapPropagator returns the default propagator for context propagation between peers.
func GetDefaultTextMapPropagator() propagation.TextMapPropagator {
	return defaultTextMapPropagator
}

// GetTraceID retrieves and returns TraceId from context.
// It returns an empty string is tracing feature is not activated.
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	if traceID.IsValid() {
		return traceID.String()
	}
	return ""
}

// GetSpanID retrieves and returns SpanId from context.
// It returns an empty string is tracing feature is not activated.
func GetSpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanID := trace.SpanContextFromContext(ctx).SpanID()
	if spanID.IsValid() {
		return spanID.String()
	}
	return ""
}

// SetBaggageValue is a convenient function for adding one key-value pair to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func SetBaggageValue(ctx context.Context, key string, value interface{}) context.Context {
	return NewBaggage(ctx).SetValue(key, value)
}

// SetBaggageMap is a convenient function for adding map key-value pairs to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func SetBaggageMap(ctx context.Context, data map[string]interface{}) context.Context {
	return NewBaggage(ctx).SetMap(data)
}

// GetBaggageMap retrieves and returns the baggage values as map.
func GetBaggageMap(ctx context.Context) *jmap.StrAnyMap {
	return NewBaggage(ctx).GetMap()
}

// GetBaggageVar retrieves value and returns a *jvar.Var for specified key from baggage.
func GetBaggageVar(ctx context.Context, key string) *jvar.Var {
	return NewBaggage(ctx).GetVar(key)
}

// WithUUID injects custom trace id with UUID into context to propagate.
func WithUUID(ctx context.Context, uuid string) (context.Context, error) {
	return WithTraceID(ctx, jstr.Replace(uuid, "-", ""))
}

// WithTraceID injects custom trace id into context to propagate.
func WithTraceID(ctx context.Context, traceID string) (context.Context, error) {
	generatedTraceID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return ctx, jerr.WithMsgErrF(
			err,
			`invalid custom traceID "%s", a traceID string should be composed with [0-f] and fixed length 32`,
			traceID,
		)
	}
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasTraceID() {
		var span trace.Span
		ctx, span = NewSpan(ctx, "gtrace.WithTraceID")
		defer span.End()
		sc = trace.SpanContextFromContext(ctx)
	}
	ctx = trace.ContextWithRemoteSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    generatedTraceID,
		SpanID:     sc.SpanID(),
		TraceFlags: sc.TraceFlags(),
		TraceState: sc.TraceState(),
		Remote:     sc.IsRemote(),
	}))
	return ctx, nil
}
