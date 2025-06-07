// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jcompress

import (
	"archive/zip"
	"bytes"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/e7coding/coding-common/internal/intlog"
	"github.com/e7coding/coding-common/os/jfile"
	"github.com/e7coding/coding-common/text/jstr"
)

// ZipPath 将 fileOrFolderPaths（可逗号分隔多路径）压缩为 dstFilePath。
// 可选 prefix 指定压缩包内文件名前缀。
func ZipPath(fileOrFolderPaths, dstFilePath string, prefix ...string) error {
	writer, err := os.Create(dstFilePath)
	if err != nil {
		return jerr.WithMsgErrF(err, `os.Create 失败，文件名="%s"`, dstFilePath)
	}

	defer func(writer *os.File) {
		_ = writer.Close()
	}(writer)

	zw := zip.NewWriter(writer)

	defer func(zw *zip.Writer) {
		_ = zw.Close()
	}(zw)

	for _, path := range strings.Split(fileOrFolderPaths, ",") {
		path = strings.TrimSpace(path)
		if err = doZipPathWriter(path, jfile.RealPath(dstFilePath), zw, prefix...); err != nil {
			return err
		}
	}
	return nil
}

// ZipPathWriter 将 fileOrFolderPaths（可逗号分隔多路径）写入任意 writer。
// 可选 prefix 指定压缩包内文件名前缀。
func ZipPathWriter(fileOrFolderPaths string, writer io.Writer, prefix ...string) error {
	zw := zip.NewWriter(writer)

	defer func(zw *zip.Writer) {
		_ = zw.Close()
	}(zw)

	for _, path := range strings.Split(fileOrFolderPaths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, "", zw, prefix...); err != nil {
			return err
		}
	}
	return nil
}

// ZipPathContent 将 fileOrFolderPaths（可逗号分隔多路径）压缩并返回字节内容。
// 可选 prefix 指定压缩包内文件名前缀。
func ZipPathContent(fileOrFolderPaths string, prefix ...string) ([]byte, error) {
	var buf bytes.Buffer
	if err := ZipPathWriter(fileOrFolderPaths, &buf, prefix...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// doZipPathWriter 遍历 fileOrFolderPath（文件或文件夹），写入 zipWriter。
// exclude 可指定跳过的文件路径，prefix 为压缩包内文件名前缀。
func doZipPathWriter(fileOrFolderPath, exclude string, zw *zip.Writer, prefix ...string) error {
	fileOrFolderPath, err := jfile.Search(fileOrFolderPath)
	if err != nil {
		return err
	}

	var files []string
	if jfile.IsDir(fileOrFolderPath) {
		files, err = jfile.ScanDir(fileOrFolderPath, "*", true)
		if err != nil {
			return err
		}
	} else {
		files = []string{fileOrFolderPath}
	}

	headerPrefix := ""
	if len(prefix) > 0 && prefix[0] != "" {
		headerPrefix = strings.TrimRight(prefix[0], "\\/")
	}
	if jfile.IsDir(fileOrFolderPath) {
		if headerPrefix != "" {
			headerPrefix += "/"
		} else {
			headerPrefix = jfile.Basename(fileOrFolderPath)
		}
	}
	headerPrefix = strings.ReplaceAll(headerPrefix, "//", "/")

	for _, file := range files {
		if exclude == file {
			intlog.Printf(`跳过文件: %s`, file)
			continue
		}
		relDir := jfile.Dir(file[len(fileOrFolderPath):])
		if relDir == "." {
			relDir = ""
		}
		if err = zipFile(file, headerPrefix+relDir, zw); err != nil {
			return err
		}
	}
	return nil
}

// UnZipFile 将 zip 文件解压到 dstFolderPath，可选 zippedPrefix 指定只解压前缀匹配的条目。
func UnZipFile(zippedFilePath, dstFolderPath string, zippedPrefix ...string) error {
	r, err := zip.OpenReader(zippedFilePath)
	if err != nil {
		return jerr.WithMsgErrF(err, `zip.OpenReader 失败，文件名="%s"`, zippedFilePath)
	}
	defer func(r *zip.ReadCloser) {
		_ = r.Close()
	}(r)

	return unZipFileWithReader(&r.Reader, dstFolderPath, zippedPrefix...)
}

// UnZipContent 将 zip 字节内容解压到 dstFolderPath，可选 zippedPrefix 指定只解压前缀匹配的条目。
func UnZipContent(zippedContent []byte, dstFolderPath string, zippedPrefix ...string) error {
	reader, err := zip.NewReader(bytes.NewReader(zippedContent), int64(len(zippedContent)))
	if err != nil {
		return jerr.WithMsgErr(err, `zip.NewReader 失败`)
	}
	return unZipFileWithReader(reader, dstFolderPath, zippedPrefix...)
}

func unZipFileWithReader(reader *zip.Reader, dstFolderPath string, zippedPrefix ...string) error {
	prefix := ""
	if len(zippedPrefix) > 0 {
		prefix = jstr.Replace(zippedPrefix[0], `\`, `/`)
	}
	if err := os.MkdirAll(dstFolderPath, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		name := strings.Trim(jstr.Replace(file.Name, `\`, `/`), "/")
		if prefix != "" {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = name[len(prefix):]
		}
		dstPath := filepath.Join(dstFolderPath, name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(dstPath, file.Mode()); err != nil {
				return jerr.WithMsgErrF(err, `创建目录失败 "%s"`, dstPath)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return jerr.WithMsgErrF(err, `创建目录失败 "%s"`, filepath.Dir(dstPath))
		}
		fr, err := file.Open()
		if err != nil {
			return jerr.WithMsgErr(err, `file.Open 失败`)
		}
		if err := doCopyForUnZipFileWithReader(file, fr, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func doCopyForUnZipFileWithReader(file *zip.File, fr io.ReadCloser, dstPath string) error {
	defer func(fr io.ReadCloser) {
		_ = fr.Close()
	}(fr)
	fw, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
	if err != nil {
		return jerr.WithMsgErrF(err, `os.OpenFile 失败，文件名="%s"`, dstPath)
	}
	defer func(fw *os.File) {
		_ = fw.Close()
	}(fw)

	if _, err = io.Copy(fw, fr); err != nil {
		return jerr.WithMsgErrF(err, `io.Copy 失败 从 "%s" 到 "%s"`, file.Name, dstPath)
	}
	return nil
}

func zipFile(filePath, prefix string, zw *zip.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return jerr.WithMsgErrF(err, `os.Open 失败，文件名="%s"`, filePath)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	info, err := f.Stat()
	if err != nil {
		return jerr.WithMsgErrF(err, `file.Stat 失败，文件名="%s"`, filePath)
	}
	header, err := createFileHeader(info, prefix)
	if err != nil {
		return err
	}
	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return jerr.WithMsgErrF(err, `CreateHeader 失败，header="%#v"`, header)
	}
	if !info.IsDir() {
		if _, err = io.Copy(writer, f); err != nil {
			return jerr.WithMsgErrF(err, `io.Copy 失败 从 "%s" 到 "%s"`, filePath, header.Name)
		}
	}
	return nil
}

func createFileHeader(info os.FileInfo, prefix string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		err = jerr.WithMsgErrF(err, `zip.FileInfoHeader failed for info "%#v"`, info)
		return nil, err
	}

	if len(prefix) > 0 {
		prefix = strings.ReplaceAll(prefix, `\`, `/`)
		prefix = strings.TrimRight(prefix, `/`)
		header.Name = prefix + `/` + header.Name
	}
	return header, nil
}
