// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package djson

import (
	"github.com/coding-common/encoding/dini"
	"github.com/coding-common/encoding/dproperties"
	"github.com/coding-common/encoding/dtoml"
	"github.com/coding-common/encoding/dxml"
	"github.com/coding-common/encoding/dyaml"
	"github.com/coding-common/internal/json"
)

// ========================================================================
// JSON
// ========================================================================

func (j *Json) ToJson() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return Encode(*(j.p))
}

func (j *Json) ToJsonString() (string, error) {
	b, e := j.ToJson()
	return string(b), e
}

func (j *Json) ToJsonIndent() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return json.MarshalIndent(*(j.p), "", "\t")
}

func (j *Json) ToJsonIndentString() (string, error) {
	b, e := j.ToJsonIndent()
	return string(b), e
}

// ========================================================================
// XML
// ========================================================================

func (j *Json) ToXml(rootTag ...string) ([]byte, error) {
	return dxml.Encode(j.Var().Map(), rootTag...)
}

func (j *Json) ToXmlString(rootTag ...string) (string, error) {
	b, e := j.ToXml(rootTag...)
	return string(b), e
}

func (j *Json) ToXmlIndent(rootTag ...string) ([]byte, error) {
	return dxml.EncodeWithIndent(j.Var().Map(), rootTag...)
}

func (j *Json) ToXmlIndentString(rootTag ...string) (string, error) {
	b, e := j.ToXmlIndent(rootTag...)
	return string(b), e
}

// ========================================================================
// YAML
// ========================================================================

func (j *Json) ToYaml() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return dyaml.Encode(*(j.p))
}

func (j *Json) ToYamlIndent(indent string) ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return dyaml.EncodeIndent(*(j.p), indent)
}

func (j *Json) ToYamlString() (string, error) {
	b, e := j.ToYaml()
	return string(b), e
}

// ========================================================================
// TOML
// ========================================================================

func (j *Json) ToToml() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return dtoml.Encode(*(j.p))
}

func (j *Json) ToTomlString() (string, error) {
	b, e := j.ToToml()
	return string(b), e
}

// ========================================================================
// INI
// ========================================================================

// ToIni json to ini
func (j *Json) ToIni() ([]byte, error) {
	return dini.Encode(j.Map())
}

// ToIniString ini to string
func (j *Json) ToIniString() (string, error) {
	b, e := j.ToIni()
	return string(b), e
}

// ========================================================================
// properties
// ========================================================================
// Toproperties json to properties
func (j *Json) ToProperties() ([]byte, error) {
	return dproperties.Encode(j.Map())
}

// ToPropertiesString properties to string
func (j *Json) ToPropertiesString() (string, error) {
	b, e := j.ToProperties()
	return string(b), e
}
