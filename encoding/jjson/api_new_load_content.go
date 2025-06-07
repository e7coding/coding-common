// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jjson

import (
	"bytes"
	"github.com/e7coding/coding-common/errs/jerr"

	"github.com/e7coding/coding-common/encoding/jini"
	"github.com/e7coding/coding-common/encoding/jproperties"
	"github.com/e7coding/coding-common/encoding/jtoml"
	"github.com/e7coding/coding-common/encoding/jxml"
	"github.com/e7coding/coding-common/encoding/jyaml"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/text/jregex"
	"github.com/e7coding/coding-common/text/jstr"
)

// LoadWithOptions creates a Json object from given JSON format content and options.
func LoadWithOptions(data []byte, options Options) (*Json, error) {
	return loadContentWithOptions(data, options)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeJson,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeXml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeIni,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeYaml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeToml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadProperties creates a Json object from given TOML format content.
func LoadProperties(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeProperties,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadContent creates a Json object from given content, it checks the data type of `content`
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data []byte, safe ...bool) (*Json, error) {
	return LoadContentType("", data, safe...)
}

// LoadContentType creates a Json object from given type and content,
// supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContentType(dataType ContentType, data []byte, safe ...bool) (*Json, error) {
	if len(data) == 0 {
		return New(nil, safe...), nil
	}
	var options = Options{
		Type:      dataType,
		StrNumber: true,
	}
	if len(safe) > 0 && safe[0] {
		options.Safe = true
	}
	return loadContentWithOptions(data, options)
}

// IsValidDataType checks and returns whether given `dataType` a valid data type for loading.
func IsValidDataType(dataType ContentType) bool {
	if dataType == "" {
		return false
	}
	if dataType[0] == '.' {
		dataType = dataType[1:]
	}
	switch dataType {
	case
		ContentTypeJson,
		ContentTypeJs,
		ContentTypeXml,
		ContentTypeYaml,
		ContentTypeYml,
		ContentTypeToml,
		ContentTypeIni,
		ContentTypeProperties:
		return true
	}
	return false
}

func trimBOM(data []byte) []byte {
	if len(data) < 3 {
		return data
	}
	if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	return data
}

// loadContentWithOptions creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func loadContentWithOptions(data []byte, options Options) (*Json, error) {
	var (
		err    error
		result interface{}
	)
	data = trimBOM(data)
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type, err = checkDataType(data)
		if err != nil {
			return nil, err
		}
	}
	options.Type = ContentType(jstr.TrimLeft(
		string(options.Type), "."),
	)
	switch options.Type {
	case ContentTypeJson, ContentTypeJs:

	case ContentTypeXml:
		data, err = jxml.ToJson(data)

	case ContentTypeYaml, ContentTypeYml:
		data, err = jyaml.ToJson(data)

	case ContentTypeToml:
		data, err = jtoml.ToJson(data)

	case ContentTypeIni:
		data, err = jini.ToJSON(data)

	case ContentTypeProperties:
		data, err = jproperties.ToJson(data)

	default:
		err = jerr.WithMsgF(
			`unsupported type "%s" for loading`,
			options.Type,
		)
	}
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	if options.StrNumber {
		decoder.UseNumber()
	}
	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, jerr.WithMsgF(`json decoding failed for content: %s`, data)
	}
	return NewWithOptions(result, options), nil
}

// checkDataType automatically checks and returns the data type for `content`.
// Note that it uses regular expression for loose checking, you can use LoadXXX/LoadContentType
// functions to load the content for certain content type.
// TODO it is not graceful here automatic judging the data type.
// TODO it might be removed in the future, which lets the user explicitly specify the data type not automatic checking.
func checkDataType(data []byte) (ContentType, error) {
	switch {
	case json.Valid(data):
		return ContentTypeJson, nil

	case isXmlContent(data):
		return ContentTypeXml, nil

	case isYamlContent(data):
		return ContentTypeYaml, nil

	case isTomlContent(data):
		return ContentTypeToml, nil

	case isIniContent(data):
		// Must contain "[xxx]" section.
		return ContentTypeIni, nil

	case isPropertyContent(data):
		return ContentTypeProperties, nil

	default:
		return "", jerr.WithMsg(
			`unable auto check the data format type`,
		)
	}
}

func isXmlContent(data []byte) bool {
	return jregex.IsMatch(`^\s*<.+>[\S\s]+<.+>\s*$`, data)
}

func isYamlContent(data []byte) bool {
	return !jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""`, data) &&
		!jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''`, data) &&
		((jregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`, data) ||
			jregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`, data)) ||
			(jregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*".+"`, data) ||
				jregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*\w+`, data)))
}

func isTomlContent(data []byte) bool {
	return !jregex.IsMatch(`^[\s\t\n\r]*;.+`, data) &&
		!jregex.IsMatch(`[\s\t\n\r]+;.+`, data) &&
		!jregex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, data) &&
		(jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data) ||
			jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data))
}

func isIniContent(data []byte) bool {
	return jregex.IsMatch(`\[[\w\.]+\]`, data) &&
		(jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data) ||
			jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data))
}

func isPropertyContent(data []byte) bool {
	return jregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data)
}
