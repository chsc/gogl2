// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"bytes"
	"fmt"
	"strings"
)

type Type struct {
	IsConst      bool
	PointerLevel int
	Name         string
}

type TypeDef struct {
	Name        string
	Comment     string
	Api         string
	CDefinition string
}

func (t Type) String() string {
	s := bytes.NewBufferString("")
	if t.IsConst {
		fmt.Fprint(s, "const ")
	}
	fmt.Fprint(s, t.Name)
	if t.PointerLevel == 1 {
		fmt.Fprint(s, " *")
	} else if t.PointerLevel == 2 {
		fmt.Fprint(s, " **")
	}
	return string(s.Bytes())
}

func (t Type) ptrStr() string {
	return strings.Repeat("*", t.PointerLevel)
}

func (t Type) IsVoid() bool {
	return (t.Name == "void" || t.Name == "GLvoid") && t.PointerLevel == 0
}

func (t Type) CType() string {
	if t.IsConst {
		return "const " + t.Name + t.ptrStr()
	}
	return t.Name + t.ptrStr()
}

func (t Type) GoType() string {
	switch t.Name {
	case "GLenum":
		return t.ptrStr() + "glt.Enum"
	case "GLbitfield":
		return t.ptrStr() + "glt.Bitfieled"
	case "GLboolean":
		if t.PointerLevel == 0 {
			return "bool"
		}
		return t.ptrStr() + "byte"
	case "GLint":
		return t.ptrStr() + "int32"
	case "GLuint":
		return t.ptrStr() + "uint32"
	case "GLint64", "GLint64EXT":
		return t.ptrStr() + "int64"
	case "GLuint64", "GLuint64EXT":
		return t.ptrStr() + "uint64"
	case "GLclampf", "GLfloat":
		return t.ptrStr() + "float32"
	case "GLclampd", "GLdouble":
		return t.ptrStr() + "float64"
	case "GLsizei":
		return t.ptrStr() + "uint32"
	case "GLbyte":
		return t.ptrStr() + "int8"
	case "GLfixed":
		return t.ptrStr() + "int32"
	case "void", "GLvoid":
		if t.PointerLevel == 1 {
			return "glt.Pointer"
		} else if t.PointerLevel == 2 {
			return "*glt.Pointer"
		}
		return ""
	case "GLintptr", "GLintptrARB":
		if t.PointerLevel == 0 {
			return "int"
		}
		return ""
	case "GLsizeiptrARB", "GLsizeiptr":
		if t.PointerLevel == 0 {
			return "int"
		}
		return ""
	case "GLcharARB", "GLchar":
		return t.ptrStr() + "int8"
	case "GLubyte":
		return t.ptrStr() + "uint8"
	case "GLshort":
		return t.ptrStr() + "int16"
	case "GLushort":
		return t.ptrStr() + "uint16"
	case "GLhandleARB":
		return t.ptrStr() + "glt.Pointer"
	case "GLhalfNV":
		return t.ptrStr() + "uint16"
	case "GLeglImageOES":
		return t.ptrStr() + "glt.Pointer"
	case "GLvdpauSurfaceARB":
		return t.ptrStr() + "glt.Pointer"
	case "GLsync":
		return t.ptrStr() + "glt.Pointer"
	}
	return ""
}

func (t Type) CgoConversion() string {
	switch t.Name {
	case "GLboolean":
		if t.PointerLevel == 0 {
			return "glt.GLBool"
		}
	case "void", "GLvoid":
		if t.PointerLevel > 0 {
			return "glt.Pointer"
		}
	case "GLintptr", "GLintptrARB":
		if t.PointerLevel == 0 {
			return "(int)"
		}
	case "GLsizeiptrARB", "GLsizeiptr":
		if t.PointerLevel == 0 {
			return "(int)"
		}
	}
	return fmt.Sprintf("(%sC.%s)", t.ptrStr(), t.Name)
}

func (t Type) GoConversion() string {
	switch t.Name {
	case "GLboolean":
		if t.PointerLevel == 0 {
			return "glt.GLBoolean"
		}
	case "void", "GLvoid":
		if t.PointerLevel > 0 {
			return "glt.Pointer"
		}
	case "GLintptr", "GLintptrARB":
		if t.PointerLevel == 0 {
			return "int"
		}
	case "GLsizeiptrARB", "GLsizeiptr":
		if t.PointerLevel == 0 {
			return "int"
		}
	}
	return fmt.Sprintf("(%sC.%s)", t.ptrStr(), t.Name)
}
