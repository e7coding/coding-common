// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dcompress 提供针对二进制/字节数据的 zlib 压缩与解压功能。
package jcompress

import (
	"bytes"
	"compress/zlib"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
)

// CompressZlib 使用 zlib 算法压缩 data。
// 如果 data 为 nil 或长度小于 13 字节，则直接返回原始 data。
func CompressZlib(data []byte) ([]byte, error) {
	if data == nil || len(data) < 13 {
		return data, nil
	}
	var (
		err    error
		buf    bytes.Buffer
		writer = zlib.NewWriter(&buf)
	)
	// 写入待压缩数据
	if _, err = writer.Write(data); err != nil {
		return nil, jerr.WithMsgErr(err, `zlib.Writer.Write 失败`)
	}
	// 关闭 writer 并刷新缓冲区
	if err = writer.Close(); err != nil {
		return buf.Bytes(), jerr.WithMsgErr(err, `zlib.Writer.Close 失败`)
	}
	return buf.Bytes(), nil
}

// DecompressZlib 使用 zlib 算法解压 data。
// 如果 data 为 nil 或长度小于 13 字节，则直接返回原始 data。
func DecompressZlib(data []byte) ([]byte, error) {
	if data == nil || len(data) < 13 {
		return data, nil
	}
	var (
		out         bytes.Buffer
		bytesReader = bytes.NewReader(data)
		reader, err = zlib.NewReader(bytesReader)
	)
	// 创建解压 reader
	if err != nil {
		return nil, jerr.WithMsgErr(err, `zlib.NewReader 失败`)
	}
	// 将解压后的数据拷贝到 out
	if _, err = io.Copy(&out, reader); err != nil {
		_ = reader.Close()
		return nil, jerr.WithMsgErr(err, `io.Copy 失败`)
	}
	// 关闭 reader
	if err = reader.Close(); err != nil {
		return out.Bytes(), jerr.WithMsgErr(err, `reader.Close 失败`)
	}
	return out.Bytes(), nil
}
