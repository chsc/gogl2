// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Name        string
	Api         string
	Version     Version
	TypeDefs    []TypeDef
	Enums       Enums
	Functions   Functions
}

type Packages []*Package

func (p *Package) writeHeader(w io.Writer) {
	fmt.Fprintln(w, "// GoGL2 - automatically generated OpenGL binding: http://github.com/chsc/gogl2")
	fmt.Fprintln(w, "//")
	writeKhronosDocCopyright(w)
	writeSgiDocCopyright(w)
	fmt.Fprintf(w, "package %s\n\n", p.Name)
}

func (p *Package) writeFooter(w io.Writer) {
	fmt.Fprintf(w, "// package %s EOF\n", p.Name)
}

func (p *Package) writeCTypes(w io.Writer) {
	for _, t := range p.TypeDefs {
		//fmt.Println(t.Api)
		if t.Name == "khrplatform" {
			continue
		}
		if t.Api == "" {
			if len(t.Comment) > 0 {
				fmt.Fprintln(w, "// /*", t.Comment, "*/")
			}
			fmt.Fprintln(w, "// ", strings.Replace(t.CDefinition, "\n", "\n// ", -1))
		}
	}
	fmt.Fprintln(w, "// ")
}

func (p *Package) writeCgoFlags(w io.Writer) {
	fmt.Fprintln(w, "// #cgo darwin  LDFLAGS: -framework OpenGL")
	fmt.Fprintln(w, "// #cgo linux   LDFLAGS: -lGL")
	fmt.Fprintln(w, "// #cgo windows LDFLAGS: -lopengl32")
	fmt.Fprintln(w, "//")
}

func (p *Package) writeAPIDefinitions(w io.Writer) {
	fmt.Fprintln(w, "// #ifndef APIENTRY")
	fmt.Fprintln(w, "// #define APIENTRY")
	fmt.Fprintln(w, "// #endif")
	fmt.Fprintln(w, "// #ifndef APIENTRYP")
	fmt.Fprintln(w, "// #define APIENTRYP APIENTRY *")
	fmt.Fprintln(w, "// #endif")
	fmt.Fprintln(w, "// #ifndef GLAPI")
	fmt.Fprintln(w, "// #define GLAPI extern")
	fmt.Fprintln(w, "// #endif")
	fmt.Fprintln(w, "//")
}

func (p *Package) writeConvFunctions(w io.Writer) {
	fmt.Fprintln(w, "func GLBoolean(b C.GLboolean) bool {")
	fmt.Fprintln(w, "	return b == TRUE")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "func GoBoolean(b bool) C.GLboolean {")
	fmt.Fprintln(w, "	if b { return TRUE }")
	fmt.Fprintln(w, "	return FALSE")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "func cgoPtr1(p *glt.Pointer) *unsafe.Pointer {")
	fmt.Fprintln(w, " return (*unsafe.Pointer)(unsafe.Pointer(p))")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "func cgoChar2(p **int8) **C.GLchar {")
	fmt.Fprintln(w, " return (**C.GLchar)(unsafe.Pointer(p))")
	fmt.Fprintln(w, "}")

}

func (p *Package) writeEnums(dir string) error {
	w, err := os.Create(filepath.Join(dir, "enums.go"))
	if err != nil {
		return err
	}
	defer w.Close()
	p.writeHeader(w)
	p.Enums.Sort().WriteGoDefinitions(w)
	p.writeFooter(w)
	return nil
}

func (p *Package) writeCommands(dir string, useFuncPtrs bool, d *Documentation) error {
	w, err := os.Create(filepath.Join(dir, "commands.go"))
	if err != nil {
		return err
	}
	defer w.Close()

	sf := p.Functions.Sort()

	p.writeHeader(w)
	p.writeCgoFlags(w)
	p.writeAPIDefinitions(w)
	p.writeCTypes(w)
	if useFuncPtrs {
		sf.WriteCFunctionPtrTypedefs(w)
		sf.WriteCBridgeDefinitions(w)
	} else {
		sf.WriteCDeclarations(w)
	}
	fmt.Fprintln(w, "import \"C\"")
	fmt.Fprintln(w, "import \"errors\"")
	fmt.Fprintln(w, "import \"github.com/chsc/gogl2/glt\"")
	fmt.Fprintln(w, "import \"github.com/chsc/gogl2/procaddr\"")
	fmt.Fprintln(w, "import \"unsafe\"")
	fmt.Fprintln(w, "")
	sf.WriteGoFunctionPtrs(w)
	p.writeConvFunctions(w)
	sf.WriteGoDefinitions(w, useFuncPtrs, d, p.Version.Major)
	sf.WriteGoInitPackage(w)
	p.writeFooter(w)

	return nil
}

func (p *Package) GeneratePackage(d *Documentation) error {
	fmt.Println("Generating package", p.Name, p.Version)
	dir := ""
	usePtr := true
	dir = filepath.Join(p.Api, p.Version.String(), p.Name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = p.writeEnums(dir)
	if err != nil {
		return err
	}
	err = p.writeCommands(dir, usePtr, d)
	if err != nil {
		return err
	}
	return nil
}

func (ps Packages) GeneratePackages(df *Documentation) error {
	for _, p := range ps {
		err := p.GeneratePackage(df)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	fmt.Fprintln(w, "// typedef unsigned int GLenum;")
	fmt.Fprintln(w, "// typedef unsigned char GLboolean;")
	fmt.Fprintln(w, "// typedef unsigned int GLbitfield;")
	fmt.Fprintln(w, "// typedef signed char GLbyte;")
	fmt.Fprintln(w, "// typedef short GLshort;")
	fmt.Fprintln(w, "// typedef int GLint;")
	fmt.Fprintln(w, "// typedef int GLsizei;")
	fmt.Fprintln(w, "// typedef unsigned char GLubyte;")
	fmt.Fprintln(w, "// typedef unsigned short GLushort;")
	fmt.Fprintln(w, "// typedef unsigned int GLuint;")
	fmt.Fprintln(w, "// typedef unsigned short GLhalf;")
	fmt.Fprintln(w, "// typedef float GLfloat;")
	fmt.Fprintln(w, "// typedef float GLclampf;")
	fmt.Fprintln(w, "// typedef double GLdouble;")
	fmt.Fprintln(w, "// typedef double GLclampd;")
	fmt.Fprintln(w, "// typedef void GLvoid;\n// ")	*/

/*	fmt.Fprintln(w, "type (\n")
	fmt.Fprintln(w, "	Enum     C.GLenum\n")
	fmt.Fprintln(w, "	Boolean  C.GLboolean\n")
	fmt.Fprintln(w, "	Bitfield C.GLbitfield\n")
	fmt.Fprintln(w, "	Byte     C.GLbyte\n")
	fmt.Fprintln(w, "	Short    C.GLshort\n")
	fmt.Fprintln(w, "	Int      C.GLint\n")
	fmt.Fprintln(w, "	Sizei    C.GLsizei\n")
	fmt.Fprintln(w, "	Ubyte    C.GLubyte\n")
	fmt.Fprintln(w, "	Ushort   C.GLushort\n")
	fmt.Fprintln(w, "	Uint     C.GLuint\n")
	fmt.Fprintln(w, "	Half     C.GLhalf\n")
	fmt.Fprintln(w, "	Float    C.GLfloat\n")
	fmt.Fprintln(w, "	Clampf   C.GLclampf\n")
	fmt.Fprintln(w, "	Double   C.GLdouble\n")
	fmt.Fprintln(w, "	Clampd   C.GLclampd\n")
	fmt.Fprintln(w, "	Char     C.GLchar\n")
	fmt.Fprintln(w, "	Pointer  unsafe.Pointer\n")
	fmt.Fprintln(w, "	Sync     C.GLsync\n")
	fmt.Fprintln(w, "	Int64    C.GLint64\n")
	fmt.Fprintln(w, "	Uint64   C.GLuint64\n")
	fmt.Fprintln(w, "	Intptr   C.GLintptr\n")
	fmt.Fprintln(w, "	Sizeiptr C.GLsizeiptr\n")
	fmt.Fprintln(w, ")\n\n")*/
