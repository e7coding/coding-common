// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package intlog provides internal logging for GoFrame development usage only.
package intlog

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/e7coding/coding-common/internal/utils"
)

// Print prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Print(ctx context.Context, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprint(v...))
}

// Printf prints `v` with format `format` using fmt.Printf.
// The parameter `v` can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprintf(format, v...))
}

// Error prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Error(ctx context.Context, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprint(v...))
}

// Errorf prints `v` with format `format` using fmt.Printf.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprintf(format, v...))
}

// PrintFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func PrintFunc(ctx context.Context, f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(ctx, s)
}

// ErrorFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func ErrorFunc(ctx context.Context, f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(ctx, s)
}

func doPrint(ctx context.Context, content string) {
	if !utils.IsDebugEnabled() {
		return
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
	buffer.WriteString(" [INTE] ")
	buffer.WriteString(" ")
	if s := traceIdStr(ctx); s != "" {
		buffer.WriteString(s + " ")
	}
	buffer.WriteString(content)
	buffer.WriteString("\n")
	fmt.Print(buffer.String())
}

// traceIdStr retrieves and returns the trace id string for logging output.
func traceIdStr(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if traceId := spanCtx.TraceID(); traceId.IsValid() {
		return "{" + traceId.String() + "}"
	}
	return ""
}
