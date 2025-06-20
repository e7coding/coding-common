// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jjson

import (
	"bytes"
	"github.com/e7coding/coding-common/errs/jerr"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// Valid checks whether `data` is a valid JSON data type.
// The parameter `data` specifies the json format data, which can be either
// bytes or string type.
func Valid(data interface{}) bool {
	return json.Valid(jconv.Bytes(data))
}

// Marshal is alias of Encode in order to fit the habit of jjson.Marshal/Unmarshal functions.
func Marshal(v interface{}) (marshaledBytes []byte, err error) {
	return Encode(v)
}

// MarshalIndent is alias of jjson.MarshalIndent in order to fit the habit of jjson.MarshalIndent function.
func MarshalIndent(v interface{}, prefix, indent string) (marshaledBytes []byte, err error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal is alias of DecodeTo in order to fit the habit of jjson.Marshal/Unmarshal functions.
func Unmarshal(data []byte, v interface{}) (err error) {
	return DecodeTo(data, v)
}

// Encode encodes any golang variable `value` to JSON bytes.
func Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// MustEncode performs as Encode, but it panics if any error occurs.
func MustEncode(value interface{}) []byte {
	b, err := Encode(value)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeString encodes any golang variable `value` to JSON string.
func EncodeString(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	return string(b), err
}

// MustEncodeString encodes any golang variable `value` to JSON string.
// It panics if any error occurs.
func MustEncodeString(value interface{}) string {
	return string(MustEncode(value))
}

// Decode decodes json format `data` to golang variable.
// The parameter `data` can be either bytes or string type.
func Decode(data interface{}, options ...Options) (interface{}, error) {
	var value interface{}
	if err := DecodeTo(jconv.Bytes(data), &value, options...); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

// DecodeTo decodes json format `data` to specified golang variable `v`.
// The parameter `data` can be either bytes or string type.
// The parameter `v` should be a pointer type.
func DecodeTo(data interface{}, v interface{}, options ...Options) (err error) {
	decoder := json.NewDecoder(bytes.NewReader(jconv.Bytes(data)))
	if len(options) > 0 {
		// The StrNumber option is for certain situations, not for all.
		// For example, it causes converting issue for other data formats, for example: yaml.
		if options[0].StrNumber {
			decoder.UseNumber()
		}
	}
	if err = decoder.Decode(v); err != nil {
		err = jerr.WithMsgErr(err, `json Decode failed`)
	}
	return
}

// DecodeToJson codes json format `data` to a Json object.
// The parameter `data` can be either bytes or string type.
func DecodeToJson(data interface{}, options ...Options) (*Json, error) {
	if v, err := Decode(jconv.Bytes(data), options...); err != nil {
		return nil, err
	} else {
		if len(options) > 0 {
			return New(v, options[0].Safe), nil
		}
		return New(v), nil
	}
}
