// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package genv provides operations for environment variables of system.
package jenv

import (
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"os"
	"strings"

	"github.com/e7coding/coding-common/container/jvar"
)

// All returns a copy of strings representing the environment,
// in the form "key=value".
func All() []string {
	return os.Environ()
}

// Map returns a copy of strings representing the environment as a map.
func Map() map[string]string {
	return MapFromEnv(os.Environ())
}

// Get creates and returns a Var with the value of the environment variable
// named by the `key`. It uses the given `def` if the variable does not exist
// in the environment.
func Get(key string, def ...interface{}) *jvar.Var {
	v, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return jvar.New(def[0])
		}
		return nil
	}
	return jvar.New(v)
}

// Set sets the value of the environment variable named by the `key`.
// It returns an error, if any.
func Set(key, value string) (err error) {
	err = os.Setenv(key, value)
	if err != nil {
		err = jerr.WithMsgErrF(err, `set environment key-value failed with key "%s", value "%s"`, key, value)
	}
	return
}

// SetMap sets the environment variables using map.
func SetMap(m map[string]string) (err error) {
	for k, v := range m {
		if err = Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Contains checks whether the environment variable named `key` exists.
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Remove deletes one or more environment variables.
func Remove(key ...string) (err error) {
	for _, v := range key {
		if err = os.Unsetenv(v); err != nil {
			err = jerr.WithMsgErrF(err, `delete environment key failed with key "%s"`, v)
			return err
		}
	}
	return nil
}

// Build builds a map to an environment variable slice.
func Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}

// MapFromEnv converts environment variables from slice to map.
func MapFromEnv(envs []string) map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range envs {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

// MapToEnv converts environment variables from map to slice.
func MapToEnv(m map[string]string) []string {
	envs := make([]string, 0)
	for k, v := range m {
		envs = append(envs, fmt.Sprintf(`%s=%s`, k, v))
	}
	return envs
}

// Filter filters repeated items from given environment variables.
func Filter(envs []string) []string {
	return MapToEnv(MapFromEnv(envs))
}
