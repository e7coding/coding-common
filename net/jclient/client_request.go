// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jclient

import (
	"bytes"
	"context"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/e7coding/coding-common/encoding/jjson"

	"github.com/e7coding/coding-common/internal/httputil"
	"github.com/e7coding/coding-common/internal/utils"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/os/jfile"
	"github.com/e7coding/coding-common/text/jregex"
	"github.com/e7coding/coding-common/text/jstr"
)

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Get(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodGet, url, data...)
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Put(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPut, url, data...)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Post(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPost, url, data...)
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Delete(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodDelete, url, data...)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Head(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodHead, url, data...)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Patch(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPatch, url, data...)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Connect(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodConnect, url, data...)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Options(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodOptions, url, data...)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Trace(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodTrace, url, data...)
}

// PostForm is different from net/http.PostForm.
// It's a wrapper of Post method, which sets the Content-Type as "multipart/form-data;".
// and It will automatically set boundary characters for the request body and Content-Type.
//
// It's Seem like the following case:
//
// Content-Type: multipart/form-data; boundary=----Boundarye4Ghaog6giyQ9ncN
//
// And form data is like:
// ------Boundarye4Ghaog6giyQ9ncN
// Content-Disposition: form-data; name="checkType"
//
// none
//
// It's used for sending form data.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) PostForm(ctx context.Context, url string, data map[string]string) (resp *Response, err error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range data {
		err := w.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return c.ContentType(w.FormDataContentType()).Post(ctx, url, body)
}

// DoRequest sends request with given HTTP method and data and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
//
// Note that it uses "multipart/form-data" as its Content-Type if it contains file uploading,
// else it uses "application/x-www-form-urlencoded". It also automatically detects the post
// content for JSON format, and for that it automatically sets the Content-Type as
// "application/json".
func (c *Client) DoRequest(
	ctx context.Context, method, url string, data ...interface{},
) (resp *Response, err error) {
	req, err := c.prepareRequest(ctx, method, url, data...)
	if err != nil {
		return nil, err
	}

	// Client middleware.
	if len(c.middlewareHandler) > 0 {
		mdlHandlers := make([]HandlerFunc, 0, len(c.middlewareHandler)+1)
		mdlHandlers = append(mdlHandlers, c.middlewareHandler...)
		mdlHandlers = append(mdlHandlers, func(cli *Client, r *http.Request) (*Response, error) {
			return cli.callRequest(r)
		})
		ctx = context.WithValue(req.Context(), clientMiddlewareKey, &clientMiddleware{
			client:       c,
			handlers:     mdlHandlers,
			handlerIndex: -1,
		})
		req = req.WithContext(ctx)
		resp, err = c.Next(req)
	} else {
		resp, err = c.callRequest(req)
	}
	if resp != nil && resp.Response != nil {
		req.Response = resp.Response
	}
	return resp, err
}

// prepareRequest verifies request parameters, builds and returns http request.
func (c *Client) prepareRequest(ctx context.Context, method, url string, data ...interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	if len(c.prefix) > 0 {
		url = c.prefix + jstr.Trim(url)
	}
	if !jstr.ContainsI(url, httpProtocolName) {
		url = httpProtocolName + `://` + url
	}
	var (
		params             string
		allowFileUploading = true
	)
	if len(data) > 0 {
		switch c.header[httpHeaderContentType] {
		case httpHeaderContentTypeJson:
			switch data[0].(type) {
			case string, []byte:
				params = jconv.String(data[0])
			default:
				if b, err := jjson.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
			allowFileUploading = false

		case httpHeaderContentTypeXml:
			switch data[0].(type) {
			case string, []byte:
				params = jconv.String(data[0])
			default:
				if b, err := jjson.New(data[0]).ToXml(); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
			allowFileUploading = false

		default:
			params = httputil.BuildParams(data[0], c.noUrlEncode)
		}
	}
	if method == http.MethodGet {
		var bodyBuffer *bytes.Buffer
		if params != "" {
			switch c.header[httpHeaderContentType] {
			case
				httpHeaderContentTypeJson,
				httpHeaderContentTypeXml:
				bodyBuffer = bytes.NewBuffer([]byte(params))
			default:
				// It appends the parameters to the url
				// if http method is GET and Content-Type is not specified.
				if jstr.Contains(url, "?") {
					url = url + "&" + params
				} else {
					url = url + "?" + params
				}
				bodyBuffer = bytes.NewBuffer(nil)
			}
		} else {
			bodyBuffer = bytes.NewBuffer(nil)
		}
		if req, err = http.NewRequest(method, url, bodyBuffer); err != nil {
			err = jerr.WithMsgErrF(err, `http.NewRequest failed with method "%s" and URL "%s"`, method, url)
			return nil, err
		}
	} else {
		if allowFileUploading && strings.Contains(params, httpParamFileHolder) {
			// File uploading request.
			var (
				buffer          = bytes.NewBuffer(nil)
				writer          = multipart.NewWriter(buffer)
				isFileUploading = false
			)
			for _, item := range strings.Split(params, "&") {
				array := strings.Split(item, "=")
				if len(array) < 2 {
					continue
				}
				if len(array[1]) > 6 && strings.Compare(array[1][0:6], httpParamFileHolder) == 0 {
					path := array[1][6:]
					if !jfile.Exists(path) {
						return nil, jerr.WithMsgF(`"%s" does not exist`, path)
					}
					var (
						file          io.Writer
						formFileName  = jfile.Basename(path)
						formFieldName = array[0]
					)
					// it sets post content type as `application/octet-stream`
					if file, err = writer.CreateFormFile(formFieldName, formFileName); err != nil {
						return nil, jerr.WithMsgErrF(
							err, `CreateFormFile failed with "%s", "%s"`, formFieldName, formFileName,
						)
					}
					var f *os.File
					if f, err = jfile.Open(path); err != nil {
						return nil, err
					}
					if _, err = io.Copy(file, f); err != nil {
						_ = f.Close()
						return nil, jerr.WithMsgErrF(
							err, `io.Copy failed from "%s" to form "%s"`, path, formFieldName,
						)
					}
					if err = f.Close(); err != nil {
						return nil, jerr.WithMsgErrF(err, `close file descriptor failed for "%s"`, path)
					}
					isFileUploading = true
				} else {
					var (
						fieldName  = array[0]
						fieldValue = array[1]
					)
					if err = writer.WriteField(fieldName, fieldValue); err != nil {
						return nil, jerr.WithMsgErrF(
							err, `write form field failed with "%s", "%s"`, fieldName, fieldValue,
						)
					}
				}
			}
			// Close finishes the multipart message and writes the trailing
			// boundary end line to the output.
			if err = writer.Close(); err != nil {
				return nil, jerr.WithMsgErrF(err, `form writer close failed`)
			}

			if req, err = http.NewRequest(method, url, buffer); err != nil {
				return nil, jerr.WithMsgErrF(
					err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url,
				)
			}
			if isFileUploading {
				req.Header.Set(httpHeaderContentType, writer.FormDataContentType())
			}
		} else {
			// Normal request.
			paramBytes := []byte(params)
			if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
				err = jerr.WithMsgErrF(err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url)
				return nil, err
			}
			if v, ok := c.header[httpHeaderContentType]; ok {
				// Custom Content-Type.
				req.Header.Set(httpHeaderContentType, v)
			} else if len(paramBytes) > 0 {
				if (paramBytes[0] == '[' || paramBytes[0] == '{') && jjson.Valid(paramBytes) {
					// Auto-detecting and setting the post content format: JSON.
					req.Header.Set(httpHeaderContentType, httpHeaderContentTypeJson)
				} else if jregex.IsMatchString(httpRegexParamJson, params) {
					// If the parameters passed like "name=value", it then uses form type.
					req.Header.Set(httpHeaderContentType, httpHeaderContentTypeForm)
				}
			}
		}
	}

	// Context.
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}
	// It's necessary set the req.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not empty.
	if reqHeaderHost := req.Header.Get(httpHeaderHost); reqHeaderHost != "" {
		req.Host = reqHeaderHost
	}
	// Custom Cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set(httpHeaderCookie, headerCookie)
		}
	}
	// HTTP basic authentication.
	if len(c.authUser) > 0 {
		req.SetBasicAuth(c.authUser, c.authPass)
	}
	return req, nil
}

// callRequest sends request with give http.Request, and returns the responses object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) callRequest(req *http.Request) (resp *Response, err error) {
	resp = &Response{
		request: req,
	}
	// Dump feature.
	// The request body can be reused for dumping
	// raw HTTP request-response procedure.
	reqBodyContent, _ := io.ReadAll(req.Body)
	resp.requestBody = reqBodyContent
	for {
		req.Body = utils.NewReadCloser(reqBodyContent, false)
		if resp.Response, err = c.Do(req); err != nil {
			err = jerr.WithMsgErrF(err, `request failed`)
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				_ = resp.Response.Body.Close()
			}
			if c.retryCount > 0 {
				c.retryCount--
				time.Sleep(c.retryInterval)
			} else {
				// return resp, err
				break
			}
		} else {
			break
		}
	}
	return resp, err
}
