// Copyright (c) 2014-2019, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"strings"
)

type RawHeader struct {
	Content string
	Lines   []string
}

func NewRawHeader(content string) *RawHeader {
	ret := &RawHeader{Content: content}
	ret.Lines = strings.Split(content, "\n")

	return ret
}

type HeaderHandler interface {
	GetExt() string
	Execute(rh *RawHeader) string
}

var HeaderHandlers = []HeaderHandler{
	&GoHeaderHandler{Base{Ext: ".go"}},
	&JSHeaderHandler{Base{Ext: ".js"}},
	&CSSHeaderHandler{Base{Ext: ".css"}},
}

func GetHandler(ext string) HeaderHandler {
	for _, handler := range HeaderHandlers {
		if ext == handler.GetExt() {
			return handler
		}
	}

	return nil
}

//////// Base Handler ////////

type Base struct {
	Ext string
}

func (base *Base) GetExt() string {
	return base.Ext
}

//////// Handlers ////////

//// Go ////
type GoHeaderHandler struct {
	Base
}

func (handler *GoHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString("//\n")
		} else {
			buffer.WriteString("// " + line + "\n")
		}
	}

	return buffer.String()
}

//// JavaScript ////
type JSHeaderHandler struct {
	Base
}

func (handler *JSHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	buffer.WriteString("/*\n")
	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString(" *\n")
		} else {
			buffer.WriteString(" * " + line + "\n")
		}
	}
	buffer.WriteString(" */\n")

	return buffer.String()
}

//// CSS ////
type CSSHeaderHandler struct {
	Base
}

func (handler *CSSHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	buffer.WriteString("/*\n")
	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString(" *\n")
		} else {
			buffer.WriteString(" * " + line + "\n")
		}
	}
	buffer.WriteString(" */\n")

	return buffer.String()
}
