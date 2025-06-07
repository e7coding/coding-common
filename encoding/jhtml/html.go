// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dhtml 提供处理 HTML 内容的实用 API。
package jhtml

import (
	"github.com/e7coding/coding-common/errs/jerr"
	"html"
	"reflect"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

// StripTags 从字符串中移除所有 HTML 标签，返回纯文本内容。
func StripTags(s string) string {
	return strip.StripTags(s)
}

// Entities 将字符串中的所有特殊字符转换为 HTML 实体。
func Entities(s string) string {
	return html.EscapeString(s)
}

// EntitiesDecode 将 HTML 实体还原为原始字符。
func EntitiesDecode(s string) string {
	return html.UnescapeString(s)
}

// SpecialChars 将字符串中的常见特殊字符（&, <, >, ", '）转换为 HTML 实体。
func SpecialChars(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&#34;",
		"'", "&#39;",
	).Replace(s)
}

// SpecialCharsDecode 将字符串中的 HTML 实体（&, <, >, ", '）还原为原始字符。
func SpecialCharsDecode(s string) string {
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	).Replace(s)
}

// SpecialCharsMapOrStruct 自动对 map 或 指向 struct 的字段值进行 htmlspecialchars 编码。
// 只会对 string 类型的值进行编码；
// 如果传入 struct，必须传入指针类型，否则返回错误。
func SpecialCharsMapOrStruct(mapOrStruct interface{}) error {
	var (
		reflectValue = reflect.ValueOf(mapOrStruct)
		reflectKind  = reflectValue.Kind()
		originalKind = reflectKind
	)
	for reflectValue.IsValid() && (reflectKind == reflect.Ptr || reflectKind == reflect.Interface) {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}

	switch reflectKind {
	case reflect.Map:
		for _, key := range reflectValue.MapKeys() {
			val := reflectValue.MapIndex(key)
			switch val.Kind() {
			case reflect.String:
				reflectValue.SetMapIndex(key, reflect.ValueOf(SpecialChars(val.String())))
			case reflect.Interface:
				if val.Elem().Kind() == reflect.String {
					reflectValue.SetMapIndex(
						key,
						reflect.ValueOf(SpecialChars(val.Elem().String())),
					)
				}
			}
		}

	case reflect.Struct:
		if originalKind != reflect.Ptr {
			return jerr.WithMsgF(
				`参数类型 "%s" 无效，应为指向 struct 的指针`,
				reflect.TypeOf(mapOrStruct).String(),
			)
		}
		for i := 0; i < reflectValue.NumField(); i++ {
			field := reflectValue.Field(i)
			if field.Kind() == reflect.String {
				field.Set(reflect.ValueOf(SpecialChars(field.String())))
			}
		}

	default:
		return jerr.WithMsgF(
			`参数类型 "%s" 无效`,
			reflect.TypeOf(mapOrStruct).String(),
		)
	}
	return nil
}
