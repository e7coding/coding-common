// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jcompress

import (
	"bytes"
	"compress/gzip"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"os"

	"github.com/e7coding/coding-common/os/jfile"
)

// Gzip 使用 gzip 算法压缩 data。
// 可选参数 level 指定压缩级别（1~9，层级越高压缩率越好）。
func Gzip(data []byte, level ...int) ([]byte, error) {
	var (
		writer *gzip.Writer
		buf    bytes.Buffer
		err    error
	)
	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(&buf, level[0])
		if err != nil {
			return nil, jerr.WithMsgErrF(err, `gzip.NewWriterLevel 失败，level="%d"`, level[0])
		}
	} else {
		writer = gzip.NewWriter(&buf)
	}
	if _, err = writer.Write(data); err != nil {
		return nil, jerr.WithMsgErr(err, `writer.Write 失败`)
	}
	if err = writer.Close(); err != nil {
		return nil, jerr.WithMsgErr(err, `writer.Close 失败`)
	}
	return buf.Bytes(), nil
}

// GzipFile 将 srcFilePath 指定的文件压缩为 gzip 格式，输出到 dstFilePath。
func GzipFile(srcFilePath, dstFilePath string, level ...int) error {
	dstFile, err := jfile.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	return GzipPathWriter(srcFilePath, dstFile, level...)
}

// GzipPathWriter 将 filePath 指定的文件或目录内容通过 gzip 算法写入到 writer。
// 可选参数 level 指定压缩级别。
func GzipPathWriter(filePath string, writer io.Writer, level ...int) error {
	var (
		gzipWriter *gzip.Writer
		err        error
	)

	srcFile, err := jfile.Open(filePath)
	if err != nil {
		return err
	}

	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	if len(level) > 0 {
		gzipWriter, err = gzip.NewWriterLevel(writer, level[0])
		if err != nil {
			return jerr.WithMsgErr(err, `gzip.NewWriterLevel 失败`)
		}
	} else {
		gzipWriter = gzip.NewWriter(writer)
	}

	defer func(gzipWriter *gzip.Writer) {
		_ = gzipWriter.Close()
	}(gzipWriter)

	if _, err = io.Copy(gzipWriter, srcFile); err != nil {
		return jerr.WithMsgErr(err, `io.Copy 失败`)
	}
	return nil
}

// UnGzip 解压缩 gzip 格式的 data，并返回原始字节数据。
func UnGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, jerr.WithMsgErr(err, `gzip.NewReader 失败`)
	}

	defer func(reader *gzip.Reader) {
		_ = reader.Close()
	}(reader)

	if _, err = io.Copy(&buf, reader); err != nil {
		return nil, jerr.WithMsgErr(err, `io.Copy 失败`)
	}
	return buf.Bytes(), nil
}

// UnGzipFile 将 srcFilePath 指定的 gzip 文件解压，输出到 dstFilePath。
func UnGzipFile(srcFilePath, dstFilePath string) error {
	srcFile, err := jfile.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	dstFile, err := jfile.Create(dstFilePath)
	if err != nil {
		return err
	}

	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	reader, err := gzip.NewReader(srcFile)
	if err != nil {
		return jerr.WithMsgErr(err, `gzip.NewReader 失败`)
	}

	defer func(reader *gzip.Reader) {
		_ = reader.Close()
	}(reader)

	if _, err = io.Copy(dstFile, reader); err != nil {
		return jerr.WithMsgErr(err, `io.Copy 失败`)
	}
	return nil
}
