// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dbase64 提供 BASE64 编码/解码 的便捷接口。
package jbase64

import (
	"encoding/base64"
	"github.com/e7coding/coding-common/errs/jerr"
	"os"
)

// Encode 将字节切片使用 BASE64 算法编码为字节切片。
func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

// EncodeStr 将字符串使用 BASE64 算法编码为字符串。
func EncodeStr(src string) string {
	return EncodeToStr([]byte(src))
}

// EncodeToStr 将字节切片使用 BASE64 算法编码为字符串。
func EncodeToStr(src []byte) string {
	return string(Encode(src))
}

// EncodeFile 从指定文件路径读取内容，并将其使用 BASE64 算法编码为字节切片。
func EncodeFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		err = jerr.WithMsgErrF(err, `os.ReadFile 读取文件 "%s" 失败`, path)
		return nil, err
	}
	return Encode(content), nil
}

// EncodeFileToStr 从指定文件路径读取内容，并将其使用 BASE64 算法编码为字符串。
func EncodeFileToStr(path string) (string, error) {
	content, err := EncodeFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Decode 将 BASE64 编码的字节切片解码为原始字节切片。
func Decode(data []byte) ([]byte, error) {
	var (
		src    = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err = base64.StdEncoding.Decode(src, data)
	)
	if err != nil {
		err = jerr.WithMsgErr(err, `base64.StdEncoding.Decode 解码失败`)
	}
	return src[:n], err
}

// DecodeStr 将 BASE64 编码的字符串解码为原始字节切片。
func DecodeStr(data string) ([]byte, error) {
	return Decode([]byte(data))
}

// DecodeToStr 将 BASE64 编码的字符串解码为字符串。
func DecodeToStr(data string) (string, error) {
	b, err := DecodeStr(data)
	return string(b), err
}
