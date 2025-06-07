// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jjson

import (
	"reflect"

	"github.com/e7coding/coding-common/internal/reflection"
	"github.com/e7coding/coding-common/internal/rwmutex"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// New creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context,
// which is false in default.
func New(data interface{}, safe ...bool) *Json {
	return NewWithTag(data, string(ContentTypeJson), safe...)
}

// NewWithTag creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `tags` specifies priority tags for struct conversion to map, multiple tags joined
// with char ','.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func NewWithTag(data interface{}, tags string, safe ...bool) *Json {
	option := Options{
		Tags: tags,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return NewWithOptions(data, option)
}

// NewWithOptions creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
func NewWithOptions(data interface{}, options Options) *Json {
	var j *Json
	switch result := data.(type) {
	case []byte:
		if r, err := loadContentWithOptions(result, options); err == nil {
			j = r
			break
		}
		j = &Json{
			p:  &data,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	case string:
		if r, err := loadContentWithOptions([]byte(result), options); err == nil {
			j = r
			break
		}
		j = &Json{
			p:  &data,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	default:
		var (
			pointedData interface{}
			reflectInfo = reflection.OriginValueAndKind(data)
		)
		switch reflectInfo.OriginKind {
		case reflect.Slice, reflect.Array:
			pointedData = jconv.Interfaces(data)

		case reflect.Map:
			pointedData = jconv.Map(data)

		case reflect.Struct:
			if v, ok := data.(IVal); ok {
				return NewWithOptions(v.Val(), options)
			}
			pointedData = jconv.Map(data)

		default:
			pointedData = data
		}
		j = &Json{
			p:  &pointedData,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	}
	j.mu = rwmutex.Create(options.Safe)
	return j
}
