// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jtimer

import (
	"time"

	"github.com/e7coding/coding-common/container/jatomic"
)

// New creates and returns a Timer.
func New(options ...TimerOptions) *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: jatomic.NewInt(StatusRunning),
		ticks:  jatomic.NewInt64(),
	}
	if len(options) > 0 {
		t.options = options[0]
		if t.options.Interval == 0 {
			t.options.Interval = defaultInterval
		}
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// Add adds a timing job to the timer, which runs in interval of `interval`.
func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Interval:    interval,
		Job:         job,
		IsSingleton: false,
		Times:       -1,
		Status:      StatusReady,
	})
}

// AddEntry adds a timing job to the timer with detailed parameters.
//
// The parameter `interval` specifies the running interval of the job.
//
// The parameter `singleton` specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when it's a singleton mode job.
//
// The parameter `times` specifies limit for the job running times, which means the job
// exits if its run times exceeds the `times`.
//
// The parameter `status` specifies the job status when it's firstly added to the timer.
func (t *Timer) AddEntry(interval time.Duration, job JobFunc, isSingleton bool, times int, status int) *Entry {
	return t.createEntry(createEntryInput{
		Interval:    interval,
		Job:         job,
		IsSingleton: isSingleton,
		Times:       times,
		Status:      status,
	})
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       -1,
		Status:      StatusReady,
	})
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       1,
		Status:      StatusReady,
	})
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       times,
		Status:      StatusReady,
	})
}

// DelayAdd adds a timing job after delay of `delay` duration.
// Also see Add.
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

// DelayAddEntry adds a timing job after delay of `delay` duration.
// Also see AddEntry.
func (t *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddEntry(interval, job, isSingleton, times, status)
	})
}

// DelayAddSingleton adds a timing job after delay of `delay` duration.
// Also see AddSingleton.
func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddSingleton(interval, job)
	})
}

// DelayAddOnce adds a timing job after delay of `delay` duration.
// Also see AddOnce.
func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddOnce(interval, job)
	})
}

// DelayAddTimes adds a timing job after delay of `delay` duration.
// Also see AddTimes.
func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddTimes(interval, times, job)
	})
}

// Start starts the timer.
func (t *Timer) Start() {
	t.status.Store(StatusRunning)
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.status.Store(StatusStopped)
}

// Close closes the timer.
func (t *Timer) Close() {
	t.status.Store(StatusClosed)
}

type createEntryInput struct {
	Interval    time.Duration
	Job         JobFunc
	IsSingleton bool
	Times       int
	Status      int
}

// createEntry creates and adds a timing job to the timer.
func (t *Timer) createEntry(in createEntryInput) *Entry {
	var (
		infinite  = false
		nextTicks int64
	)
	if in.Times <= 0 {
		infinite = true
	}
	var (
		intervalTicksOfJob = int64(in.Interval / t.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	if t.options.Quick {
		// If the quick mode is enabled, which means it will be run right now.
		// Don't need to wait for the first interval.
		nextTicks = t.ticks.Load()
	} else {
		nextTicks = t.ticks.Load() + intervalTicksOfJob
	}
	var (
		entry = &Entry{
			job:         in.Job,
			timer:       t,
			ticks:       intervalTicksOfJob,
			times:       jatomic.NewInt(in.Times),
			status:      jatomic.NewInt(in.Status),
			isSingleton: jatomic.NewBool(in.IsSingleton),
			nextTicks:   jatomic.NewInt64(nextTicks),
			infinite:    jatomic.NewBool(infinite),
		}
	)
	t.queue.Push(entry, nextTicks)
	return entry
}
