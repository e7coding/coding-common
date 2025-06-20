// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gclient provides convenient http client functionalities.
package jclient

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"net/http"
	"os"
	"time"

	"github.com/e7coding/coding-common"

	"github.com/e7coding/coding-common/os/jfile"
)

// Client is the HTTP client for HTTP request management.
type Client struct {
	http.Client                         // Underlying HTTP Client.
	header            map[string]string // Custom header map.
	cookies           map[string]string // Custom cookie map.
	prefix            string            // Prefix for request.
	authUser          string            // HTTP basic authentication: user.
	authPass          string            // HTTP basic authentication: pass.
	retryCount        int               // Retry count when request fails.
	noUrlEncode       bool              // No url encoding for request parameters.
	retryInterval     time.Duration     // Retry interval when request fails.
	middlewareHandler []HandlerFunc     // Interceptor handlers
}

const (
	httpProtocolName          = `http`
	httpParamFileHolder       = `@file:`
	httpRegexParamJson        = `^[\w\[\]]+=.+`
	httpRegexHeaderRaw        = `^([\w\-]+):\s*(.+)`
	httpHeaderHost            = `Host`
	httpHeaderCookie          = `Cookie`
	httpHeaderUserAgent       = `User-Agent`
	httpHeaderContentType     = `Content-Type`
	httpHeaderContentTypeJson = `application/json`
	httpHeaderContentTypeXml  = `application/xml`
	httpHeaderContentTypeForm = `application/x-www-form-urlencoded`
)

var (
	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`GClient %s at %s`, gf.VERSION, hostname)
)

// New creates and returns a new HTTP client object.
func New() *Client {
	c := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header:  make(map[string]string),
		cookies: make(map[string]string),
	}
	c.header[httpHeaderUserAgent] = defaultClientAgent
	return c
}

// Clone deeply clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := New()
	*newClient = *c
	newClient.header = make(map[string]string, len(c.header))
	for k, v := range c.header {
		newClient.header[k] = v
	}
	newClient.cookies = make(map[string]string, len(c.cookies))
	for k, v := range c.cookies {
		newClient.cookies[k] = v
	}
	return newClient
}

// LoadKeyCrt creates and returns a TLS configuration object with given certificate and key files.
func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error) {
	crtPath, err := jfile.Search(crtFile)
	if err != nil {
		return nil, err
	}
	keyPath, err := jfile.Search(keyFile)
	if err != nil {
		return nil, err
	}
	crt, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		err = jerr.WithMsgErrF(err, `tls.LoadX509KeyPair failed for certFile "%s", keyFile "%s"`, crtPath, keyPath)
		return nil, err
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}
