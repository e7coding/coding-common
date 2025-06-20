// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package httputil provides HTTP functions for internal usage only.
package httputil

import (
	"github.com/e7coding/coding-common/encoding/jurl"
	"net/http"
	"strings"

	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/text/jstr"
)

const (
	fileUploadingKey = "@file:"
)

// BuildParams builds the request string for the http client. The `params` can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter `noUrlEncode` specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		return jconv.String(params)
	case []interface{}:
		if len(v) > 0 {
			params = v[0]
		} else {
			params = nil
		}
	}
	// Else converts it to map and does the url encoding.
	m, urlEncode := jconv.Map(params, jconv.MapOption{
		OmitEmpty: true,
	}), true
	if len(m) == 0 {
		return jconv.String(params)
	}
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	// If there's file uploading, it ignores the url encoding.
	if urlEncode {
		for k, v := range m {
			if jstr.Contains(k, fileUploadingKey) || jstr.Contains(jconv.String(v), fileUploadingKey) {
				urlEncode = false
				break
			}
		}
	}
	s := ""
	for k, v := range m {
		// Ignore nil attributes.
		if empty.IsNil(v) {
			continue
		}
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = jconv.String(v)
		if urlEncode {
			if strings.HasPrefix(s, fileUploadingKey) && len(s) > len(fileUploadingKey) {
				// No url encoding if uploading file.
			} else {
				s = jurl.Encode(s)
			}
		}
		encodedParamStr += k + "=" + s
	}
	return
}

// HeaderToMap coverts request headers to map.
func HeaderToMap(header http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range header {
		if len(v) > 1 {
			m[k] = v
		} else if len(v) == 1 {
			m[k] = v[0]
		}
	}
	return m
}
