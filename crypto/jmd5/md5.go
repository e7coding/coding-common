// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dmd5 提供 MD5 加密算法的实用接口
package jmd5

import (
	"crypto/md5"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"os"

	"github.com/e7coding/coding-common/jutil/jconv"
)

// Enc 对任意类型的数据使用 MD5 算法进行加密。
// 会使用 gconv 包将 data 转换为字节类型。
func Enc(data interface{}) (encrypt string, err error) {
	return EncBytes(jconv.Bytes(data))
}

// EncBytes 对字节数据使用 MD5 算法进行加密。
func EncBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write(data); err != nil {
		err = jerr.WithMsgErr(err, "hash.Write 失败")
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// EncStr 对字符串数据使用 MD5 算法进行加密。
func EncStr(data string) (encrypt string, err error) {
	return EncBytes([]byte(data))
}

// EncFile 对指定路径的文件内容使用 MD5 算法进行加密。
func EncFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err != nil {
		err = jerr.WithMsgErrF(err, `os.Open 失败，path="%s"`, path)
		return "", err
	}

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		err = jerr.WithMsgErr(err, "io.Copy 失败")
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
