// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"strings"
	"unicode"
)

func TrimGLCmdPrefix(str string) string {
	return strings.TrimPrefix(str, "gl")
}

func TrimGLEnumPrefix(str string) string {
	t := strings.TrimPrefix(str, "GL_")
	if strings.IndexAny(t, "0123456789") == 0 {
		return fmt.Sprintf("GL_%s", t) // keep it if we have a number
	}
	return t
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
