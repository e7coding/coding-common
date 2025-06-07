// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dsha1 提供 SHA1 加密算法的实用接口
package jsha1

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"os"

	"github.com/e7coding/coding-common/jutil/jconv"
)

// Enc 对任意类型的数据使用 SHA1 算法进行加密，返回十六进制字符串。
func Enc(v interface{}) string {
	r := sha1.Sum(jconv.Bytes(v))
	return hex.EncodeToString(r[:])
}

// EncFile 对指定路径的文件内容使用 SHA1 算法进行加密，返回十六进制字符串。
func EncFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err != nil {
		err = jerr.WithMsgErrF(err, `os.Open 失败，path="%s"`, path)
		return "", err
	}

	h := sha1.New()
	if _, err = io.Copy(h, f); err != nil {
		err = jerr.WithMsgErr(err, "io.Copy 失败")
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
