// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package djson 提供 JSON/XML/INI/YAML/TOML 等格式的数据处理接口。
package jjson

import (
	"github.com/e7coding/coding-common/errs/jerr"
	"reflect"
	"strconv"
	"strings"

	"github.com/e7coding/coding-common/internal/reflection"
	"github.com/e7coding/coding-common/internal/rwmutex"
	"github.com/e7coding/coding-common/internal/utils"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/text/jstr"
)

type ContentType string

const (
	ContentTypeJson       ContentType = `json`
	ContentTypeJs         ContentType = `js`
	ContentTypeXml        ContentType = `xml`
	ContentTypeIni        ContentType = `ini`
	ContentTypeYaml       ContentType = `yaml`
	ContentTypeYml        ContentType = `yml`
	ContentTypeToml       ContentType = `toml`
	ContentTypeProperties ContentType = `properties`
)

const (
	defaultSplitChar = '.' // 层级数据访问的分隔符
)

// Json 是定制化的 JSON 结构体。
type Json struct {
	mu rwmutex.RWMutex
	p  *interface{} // 层级数据访问的根指针，默认指向数据根节点。
	c  byte         // 分隔符（默认 '.'）。
	vc bool         // 激进检查（默认 false），当层级键包含分隔符时启用。
}

// Options 定义 Json 对象创建/加载时的可选项。
type Options struct {
	Safe      bool        // 是否并发安全，主要用于 Json 对象创建时。
	Tags      string      // 解码时使用的自定义优先标签，如 "json,yaml,MyTag"，用于结构体解析。
	Type      ContentType // 数据类型，如 json、xml、yaml、toml、ini。
	StrNumber bool        // 数字是否以字符串形式解析到 interface{} 而非 float64。
}

// InterfaceFunc 用于类型断言，提供 Interfaces() 方法。
type InterfaceFunc interface {
	Interfaces() []interface{}
}

// IMapStrAny 支持将结构体转换为 map[string]interface{} 的接口。
type IMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// IVal 用于获取底层 interface{} 值的接口。
type IVal interface {
	Val() interface{}
}

// setValue 根据路径 pattern 设置或删除值 value。
// 备注：
// 1. value 为 nil 且 removed 为 true 表示删除该值；
// 2. 内部逻辑涉及层级搜索、节点创建和数据赋值。
func (j *Json) setValue(pattern string, value interface{}, removed bool) error {
	var (
		err    error
		array  = strings.Split(pattern, string(j.c))
		length = len(array)
	)
	if value, err = j.convertValue(value); err != nil {
		return err
	}
	// 初始化检查
	if *j.p == nil {
		if jstr.IsNumeric(array[0]) {
			*j.p = make([]interface{}, 0)
		} else {
			*j.p = make(map[string]interface{})
		}
	}
	var (
		pparent *interface{} = nil // 父节点指针
		pointer *interface{} = j.p // 当前节点指针
	)
	j.mu.Lock()
	defer j.mu.Unlock()
	for i := 0; i < length; i++ {
		switch (*pointer).(type) {
		case map[string]interface{}:
			// 处理 map 节点
			if i == length-1 {
				if removed && value == nil {
					// 从 map 中删除项
					delete((*pointer).(map[string]interface{}), array[i])
				} else {
					if (*pointer).(map[string]interface{}) == nil {
						*pointer = map[string]interface{}{}
					}
					(*pointer).(map[string]interface{})[array[i]] = value
				}
			} else {
				// map 中不存在该键
				if v, ok := (*pointer).(map[string]interface{})[array[i]]; !ok {
					if removed && value == nil {
						goto done
					}
					// 创建新节点
					if jstr.IsNumeric(array[i+1]) {
						// 创建数组节点
						n, _ := strconv.Atoi(array[i+1])
						var v interface{} = make([]interface{}, n+1)
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					} else {
						// 创建 map 节点
						var v interface{} = make(map[string]interface{})
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					}
				} else {
					pparent = pointer
					pointer = &v
				}
			}

		case []interface{}:
			// 处理 slice 节点
			if !jstr.IsNumeric(array[i]) {
				// 字符串键
				if i == length-1 {
					*pointer = map[string]interface{}{array[i]: value}
				} else {
					var v interface{} = make(map[string]interface{})
					*pointer = v
					pparent = pointer
					pointer = &v
				}
				continue
			}
			// 数字索引
			valueNum, err := strconv.Atoi(array[i])
			if err != nil {
				err = jerr.WithMsgErrF(err, `strconv.Atoi 失败："%s"`, array[i])
				return err
			}

			if i == length-1 {
				// 叶子节点
				if len((*pointer).([]interface{})) > valueNum {
					if removed && value == nil {
						// 删除元素
						if pparent == nil {
							*pointer = append((*pointer).([]interface{})[:valueNum], (*pointer).([]interface{})[valueNum+1:]...)
						} else {
							j.setPointerWithValue(pparent, array[i-1], append((*pointer).([]interface{})[:valueNum], (*pointer).([]interface{})[valueNum+1:]...))
						}
					} else {
						(*pointer).([]interface{})[valueNum] = value
					}
				} else {
					if removed && value == nil {
						goto done
					}
					if pparent == nil {
						// 根节点
						j.setPointerWithValue(pointer, array[i], value)
					} else {
						// 非根节点
						s := make([]interface{}, valueNum+1)
						copy(s, (*pointer).([]interface{}))
						s[valueNum] = value
						j.setPointerWithValue(pparent, array[i-1], s)
					}
				}
			} else {
				// 分支节点
				if jstr.IsNumeric(array[i+1]) {
					n, _ := strconv.Atoi(array[i+1])
					pSlice := (*pointer).([]interface{})
					if len(pSlice) > valueNum {
						item := pSlice[valueNum]
						if s, ok := item.([]interface{}); ok {
							for i := 0; i < n-len(s); i++ {
								s = append(s, nil)
							}
							pparent = pointer
							pointer = &pSlice[valueNum]
						} else {
							if removed && value == nil {
								goto done
							}
							var v interface{} = make([]interface{}, n+1)
							pparent = j.setPointerWithValue(pointer, array[i], v)
							pointer = &v
						}
					} else {
						if removed && value == nil {
							goto done
						}
						var v interface{} = make([]interface{}, n+1)
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					}
				} else {
					pSlice := (*pointer).([]interface{})
					if len(pSlice) > valueNum {
						pparent = pointer
						pointer = &(*pointer).([]interface{})[valueNum]
					} else {
						s := make([]interface{}, valueNum+1)
						copy(s, pSlice)
						s[valueNum] = make(map[string]interface{})
						if pparent != nil {
							// i > 0
							j.setPointerWithValue(pparent, array[i-1], s)
							pparent = pointer
							pointer = &s[valueNum]
						} else {
							// i = 0
							var v interface{} = s
							*pointer = v
							pparent = pointer
							pointer = &s[valueNum]
						}
					}
				}
			}

		default:
			// pointer 指向的值非引用类型，通过父节点修改
			if removed && value == nil {
				goto done
			}
			if jstr.IsNumeric(array[i]) {
				n, _ := strconv.Atoi(array[i])
				s := make([]interface{}, n+1)
				if i == length-1 {
					s[n] = value
				}
				if pparent != nil {
					pparent = j.setPointerWithValue(pparent, array[i-1], s)
				} else {
					*pointer = s
					pparent = pointer
				}
			} else {
				var v1, v2 interface{}
				if i == length-1 {
					v1 = map[string]interface{}{array[i]: value}
				} else {
					v1 = map[string]interface{}{array[i]: nil}
				}
				if pparent != nil {
					pparent = j.setPointerWithValue(pparent, array[i-1], v1)
				} else {
					*pointer = v1
					pparent = pointer
				}
				v2 = v1.(map[string]interface{})[array[i]]
				pointer = &v2
			}
		}
	}
done:
	return nil
}

// convertValue 将 value 转换为 map[string]interface{} 或 []interface{}，以支持层级访问。
func (j *Json) convertValue(value interface{}) (convertedValue interface{}, err error) {
	if value == nil {
		return
	}

	switch value.(type) {
	case map[string]interface{}:
		return value, nil

	case []interface{}:
		return value, nil

	default:
		var reflectInfo = reflection.OriginValueAndKind(value)
		switch reflectInfo.OriginKind {
		case reflect.Array:
			return jconv.Interfaces(value), nil

		case reflect.Slice:
			return jconv.Interfaces(value), nil

		case reflect.Map:
			return jconv.Map(value), nil

		case reflect.Struct:
			if v, ok := value.(IMapStrAny); ok {
				convertedValue = v.MapStrAny()
			}
			if utils.IsNil(convertedValue) {
				if v, ok := value.(InterfaceFunc); ok {
					convertedValue = v.Interfaces()
				}
			}
			if utils.IsNil(convertedValue) {
				convertedValue = jconv.Map(value)
			}
			if utils.IsNil(convertedValue) {
				err = jerr.WithMsgF(`不支持的类型 "%s"`, reflect.TypeOf(value))
			}
			return

		default:
			return value, nil
		}
	}
}

// setPointerWithValue 在 pointer 上设置 key:value，可用于 map 键或 slice 索引，返回新节点指针。
func (j *Json) setPointerWithValue(pointer *interface{}, key string, value interface{}) *interface{} {
	switch (*pointer).(type) {
	case map[string]interface{}:
		(*pointer).(map[string]interface{})[key] = value
		return &value
	case []interface{}:
		n, _ := strconv.Atoi(key)
		if len((*pointer).([]interface{})) > n {
			(*pointer).([]interface{})[n] = value
			return &(*pointer).([]interface{})[n]
		} else {
			s := make([]interface{}, n+1)
			copy(s, (*pointer).([]interface{}))
			s[n] = value
			*pointer = s
			return &s[n]
		}
	default:
		*pointer = value
	}
	return pointer
}

// getPointerByPattern 根据路径 pattern 返回值指针。
func (j *Json) getPointerByPattern(pattern string) *interface{} {
	if j.p == nil {
		return nil
	}
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	} else {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}
}

// getPointerByPatternWithViolenceCheck 带激进检查的路径访问。
func (j *Json) getPointerByPatternWithViolenceCheck(pattern string) *interface{} {
	if !j.vc {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}

	// pattern 为空时返回 nil
	if pattern == "" {
		return nil
	}
	// pattern 为 "." 时返回根节点
	if pattern == "." {
		return j.p
	}

	var (
		index   = len(pattern)
		start   = 0
		length  = 0
		pointer = j.p
	)
	if index == 0 {
		return pointer
	}
	for {
		if r := j.checkPatternByPointer(pattern[start:index], pointer); r != nil {
			if length += index - start; start > 0 {
				length += 1
			}
			start = index + 1
			index = len(pattern)
			if length == len(pattern) {
				return r
			} else {
				pointer = r
			}
		} else {
			// 获取下一分隔符位置
			index = strings.LastIndexByte(pattern[start:index], j.c)
			if index != -1 && length > 0 {
				index += length + 1
			}
		}
		if start >= index {
			break
		}
	}
	return nil
}

// getPointerByPatternWithoutViolenceCheck 普通路径访问，无激进检查。
func (j *Json) getPointerByPatternWithoutViolenceCheck(pattern string) *interface{} {
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	}

	// pattern 为空时返回 nil
	if pattern == "" {
		return nil
	}
	// pattern 为 "." 时返回根节点
	if pattern == "." {
		return j.p
	}

	pointer := j.p
	if len(pattern) == 0 {
		return pointer
	}
	array := strings.Split(pattern, string(j.c))
	for k, v := range array {
		if r := j.checkPatternByPointer(v, pointer); r != nil {
			if k == len(array)-1 {
				return r
			} else {
				pointer = r
			}
		} else {
			break
		}
	}
	return nil
}

// checkPatternByPointer 在 pointer 指向的 map 或 slice 上按 key/index 查找，并返回指针。
func (j *Json) checkPatternByPointer(key string, pointer *interface{}) *interface{} {
	switch (*pointer).(type) {
	case map[string]interface{}:
		if v, ok := (*pointer).(map[string]interface{})[key]; ok {
			return &v
		}
	case []interface{}:
		if jstr.IsNumeric(key) {
			n, err := strconv.Atoi(key)
			if err == nil && len((*pointer).([]interface{})) > n {
				return &(*pointer).([]interface{})[n]
			}
		}
	}
	return nil
}
