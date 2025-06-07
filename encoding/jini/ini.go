// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jini

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"strings"
)

// Decode 将 INI 格式的字节数据解析为 map。
func Decode(data []byte) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	fieldMap := make(map[string]interface{})
	reader := bufio.NewReader(bytes.NewReader(data))
	var (
		section     string
		lastSection string
		haveSection bool
		line        string
		err         error
	)

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, jerr.WithMsgErr(err, `读取 INI 行失败`)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		if idx1 := strings.Index(line, "["); idx1 >= 0 {
			if idx2 := strings.Index(line, "]"); idx2 > idx1 {
				section = line[idx1+1 : idx2]
				if section != lastSection {
					fieldMap = make(map[string]interface{})
					lastSection = section
				}
				haveSection = true
				continue
			}
		}
		if !haveSection || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		fieldMap[key] = val
		res[section] = fieldMap
	}

	if !haveSection {
		return nil, jerr.WithMsg("INI 格式解析失败：未找到任何节")
	}
	return res, nil
}

// Encode 将 map 转换为 INI 格式的字节数据。
func Encode(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	for section, v := range data {
		buf.WriteString(fmt.Sprintf("[%s]\n", section))
		if m, ok := v.(map[string]interface{}); ok {
			for k, val := range m {
				buf.WriteString(fmt.Sprintf("%s=%v\n", k, val))
			}
		}
	}
	return buf.Bytes(), nil
}

// ToJSON 将 INI 格式的字节数据转换为 JSON。
func ToJSON(data []byte) ([]byte, error) {
	iniMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(iniMap)
}
