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

type GoHeaderHandler struct {
	Base
}

func (handler *GoHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		buffer.WriteString("// " + line + "\n")
	}

	buffer.WriteString("\n")

	return buffer.String()
}
