// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"time"

	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/utils"
	"github.com/e7coding/coding-common/jutil/jconv/internal/localinterface"
	"github.com/e7coding/coding-common/os/jtime"
)

// Time converts `any` to time.Time.
func (c *Converter) Time(any interface{}, format ...string) (time.Time, error) {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v, nil
		}
	}
	t, err := c.GTime(any, format...)
	if err != nil {
		return time.Time{}, err
	}
	if t != nil {
		return t.Time, nil
	}
	return time.Time{}, nil
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func (c *Converter) Duration(any interface{}) (time.Duration, error) {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v, nil
	}
	s, err := c.String(any)
	if err != nil {
		return 0, err
	}
	if !utils.IsNumeric(s) {
		return jtime.ParseDuration(s)
	}
	i, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return time.Duration(i), nil
}

// GTime converts `any` to *time.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using jtime.NewFromTimeStamp if `any` is numeric,
// or using jtime.StrToTime if `any` is string.
func (c *Converter) GTime(any interface{}, format ...string) (*jtime.Time, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	if v, ok := any.(localinterface.IGTime); ok {
		return v.GTime(format...), nil
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*jtime.Time); ok {
			return v, nil
		}
		if t, ok := any.(time.Time); ok {
			return jtime.New(t), nil
		}
		if t, ok := any.(*time.Time); ok {
			return jtime.New(t), nil
		}
	}
	s, err := c.String(any)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return jtime.New(), nil
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := jtime.StrToTimeFormat(s, item)
			if err != nil {
				return nil, err
			}
			if t != nil {
				return t, nil
			}
		}
		return nil, nil
	}
	if utils.IsNumeric(s) {
		i, err := c.Int64(s)
		if err != nil {
			return nil, err
		}
		return jtime.NewFromTimeStamp(i), nil
	} else {
		return jtime.StrToTime(s)
	}
}
