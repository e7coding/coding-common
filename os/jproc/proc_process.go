// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jproc

import (
	"context"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/e7coding/coding-common"

	"github.com/e7coding/coding-common/internal/intlog"
	"github.com/e7coding/coding-common/net/jtrace"
	"github.com/e7coding/coding-common/os/jenv"
	"github.com/e7coding/coding-common/text/jstr"
)

// Process is the struct for a single process.
type Process struct {
	exec.Cmd
	Manager *Manager
	PPid    int
}

// NewProcess creates and returns a new Process.
func NewProcess(path string, args []string, environment ...[]string) *Process {
	eenv := os.Environ()
	if len(environment) > 0 {
		eenv = append(eenv, environment[0]...)
	}
	process := &Process{
		Manager: nil,
		PPid:    os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Env:        eenv,
			ExtraFiles: make([]*os.File, 0),
		},
	}
	process.Dir, _ = os.Getwd()
	if len(args) > 0 {
		// Exclude of current binary path.
		start := 0
		if strings.EqualFold(path, args[0]) {
			start = 1
		}
		process.Args = append(process.Args, args[start:]...)
	}
	return process
}

// NewProcessCmd creates and returns a process with given command and optional environment variable array.
func NewProcessCmd(cmd string, environment ...[]string) *Process {
	return NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...), environment...)
}

// Start starts executing the process in non-blocking way.
// It returns the pid if success, or else it returns an error.
func (p *Process) Start(ctx context.Context) (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	// OpenTelemetry for command.
	var (
		span trace.Span
		tr   = otel.GetTracerProvider().Tracer(
			tracingInstrumentName,
			trace.WithInstrumentationVersion(gf.VERSION),
		)
	)
	ctx, span = tr.Start(
		otel.GetTextMapPropagator().Extract(
			ctx,
			propagation.MapCarrier(jenv.Map()),
		),
		jstr.Join(os.Args, " "),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()
	span.SetAttributes(jtrace.CommonLabels()...)

	// OpenTelemetry propagation.
	tracingEnv := tracingEnvFromCtx(ctx)
	if len(tracingEnv) > 0 {
		p.Env = append(p.Env, tracingEnv...)
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", envKeyPPid, p.PPid))
	p.Env = jenv.Filter(p.Env)

	// On Windows, this works and doesn't work on other platforms
	if runtime.GOOS == "windows" {
		joinProcessArgs(p)
	}

	if err := p.Cmd.Start(); err == nil {
		if p.Manager != nil {
			p.Manager.processes.Put(p.Process.Pid, p)
		}
		return p.Process.Pid, nil
	} else {
		return 0, err
	}
}

// Run executes the process in blocking way.
func (p *Process) Run(ctx context.Context) error {
	if _, err := p.Start(ctx); err == nil {
		return p.Wait()
	} else {
		return err
	}
}

// Pid retrieves and returns the PID for the process.
func (p *Process) Pid() int {
	if p.Process != nil {
		return p.Process.Pid
	}
	return 0
}

// Send sends custom data to the process.
func (p *Process) Send(data []byte) error {
	if p.Process != nil {
		return Send(p.Process.Pid, data)
	}
	return jerr.WithMsg("invalid process")
}

// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (p *Process) Release() error {
	return p.Process.Release()
}

// Kill causes the Process to exit immediately.
func (p *Process) Kill() (err error) {
	err = p.Process.Kill()
	if err != nil {
		err = jerr.WithMsgErrF(err, `kill process failed for pid "%d"`, p.Process.Pid)
		return err
	}
	if p.Manager != nil {
		p.Manager.processes.Delete(p.Pid())
	}
	if runtime.GOOS != "windows" {
		if err = p.Process.Release(); err != nil {
			intlog.Errorf(`%+v`, err)
		}
	}
	// It ignores this error, just log it.
	_, err = p.Process.Wait()
	intlog.Errorf(`%+v`, err)
	return nil
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
	return p.Process.Signal(sig)
}
