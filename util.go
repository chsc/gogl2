// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type ParamLenType int

const (
	ParamLenTypeUnknown ParamLenType = iota
	ParamLenTypeParamRef
	ParamLenTypeValue
	ParamLenTypeCompSize
)

type ParamLen struct {
	Type     ParamLenType
	ParamRef string
	Value    int
	Params   string
}

func TrimGLCmdPrefix(str string) string {
	if strings.HasPrefix(str, "gl") {
		return strings.TrimPrefix(str, "gl")
	}
	if strings.HasPrefix(str, "glx") {
		return strings.TrimPrefix(str, "glx")
	}
	if strings.HasPrefix(str, "wgl") {
		return strings.TrimPrefix(str, "wgl")
	}
	return str
}

func TrimGLEnumPrefix(str string) string {
	t := str
	p := ""
	if strings.HasPrefix(str, "GL_") {
		t = strings.TrimPrefix(t, "GL_")
		p = "GL_"
	} else if strings.HasPrefix(str, "GLX_") {
		t = strings.TrimPrefix(str, "GLX_")
		p = "GLX_"
	} else if strings.HasPrefix(str, "WGL_") {
		t = strings.TrimPrefix(str, "WGL_")
		p = "WGL_"
	}
	if strings.IndexAny(t, "0123456789") == 0 {
		return p + t
	}
	return t
}

func ParseLenString(lenStr string) ParamLen {
	if strings.HasPrefix(lenStr, "COMPSIZE") {
		p := strings.TrimSuffix(strings.TrimPrefix(lenStr, "COMPSIZE("), ")")
		return ParamLen{Type: ParamLenTypeCompSize, Params: p}
	}
	n, err := strconv.ParseInt(lenStr, 10, 32)
	if err == nil {
		return ParamLen{Type: ParamLenTypeValue, Value: (int)(n)}
	}
	return ParamLen{Type: ParamLenTypeParamRef, ParamRef: lenStr}
}

// Prevent name clashes.
func RenameIfReservedCWord(word string) string {
	switch word {
	case "near", "far":
		return fmt.Sprintf("%s", word)
	}
	return word
}

// Prevent name clashes.
func RenameIfReservedGoWord(word string) string {
	switch word {
	case "func", "type", "struct", "range", "map", "string":
		return fmt.Sprintf("gl%s", word)
	}
	return word
}

// Converts strings with underscores to Go-like names. e.g.: bla_blub_foo -> BlaBlubFoo
func CamelCase(n string) string {
	prev := '_'
	return strings.Map(
		func(r rune) rune {
			if r == '_' {
				prev = r
				return -1
			}
			if prev == '_' {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		n)
}
