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

const (
	PackageTypeUnknown PackageType = iota
	PackageTypeGL
	PackageTypeGLExt
	PackageTypeGLES
	PackageTypeWGL
	PackageTypeGLX
	PackageTypeEGL
)

type PackageType int

type Package struct {
	PackageType PackageType
	Name        string
	Version     Version
	TypeDefs    []TypeDef
	Enums       Enums
	Functions   Functions
}

type Packages []*Package

func (p *Package) writeHeader(w io.Writer) {
	fmt.Fprintln(w, "// GoGL2 - automatically generated OpenGL binding: http://github.com/chsc/gogl2")
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

func (p *Package) writeCGetProcAddressFunction(w io.Writer) {
	fmt.Fprintln(w, "// #ifdef _WIN32")
	fmt.Fprintln(w, "// static HMODULE opengl32 = NULL;")
	fmt.Fprintln(w, "// #endif\n// ")
	fmt.Fprintln(w, "// static void* goglGetProcAddress(const char* name) {")
	fmt.Fprintln(w, "// #ifdef __APPLE__")
	fmt.Fprintln(w, "// 	return dlsym(RTLD_DEFAULT, name);")
	fmt.Fprintln(w, "// #elif _WIN32")
	fmt.Fprintln(w, "// 	void* pf = wglGetProcAddress((LPCSTR)name);")
	fmt.Fprintln(w, "// 	if(pf) {")
	fmt.Fprintln(w, "// 		return pf;")
	fmt.Fprintln(w, "// 	}")
	fmt.Fprintln(w, "// 	if(opengl32 == NULL) {")
	fmt.Fprintln(w, "// 		opengl32 = LoadLibraryA(\"opengl32.dll\");")
	fmt.Fprintln(w, "// 	}")
	fmt.Fprintln(w, "// 	return GetProcAddress(opengl32, (LPCSTR)name);")
	fmt.Fprintln(w, "// #else")
	fmt.Fprintln(w, "// 	return glXGetProcAddress((const GLubyte*)name);")
	fmt.Fprintln(w, "// #endif")
	fmt.Fprintln(w, "// }")
	fmt.Fprintln(w, "//")
}

func (p *Package) writeCgoFlags(w io.Writer) {
	fmt.Fprintln(w, "// #cgo darwin  LDFLAGS: -framework OpenGL")
	fmt.Fprintln(w, "// #cgo linux   LDFLAGS: -lGL")
	fmt.Fprintln(w, "// #cgo windows LDFLAGS: -lopengl32")
	fmt.Fprintln(w, "//")
}

func (p *Package) writeCIncludes(w io.Writer) {
	fmt.Fprintln(w, "// #include <stdlib.h>")
	fmt.Fprintln(w, "// #if defined(__APPLE__)")
	fmt.Fprintln(w, "// #include <dlfcn.h>")
	fmt.Fprintln(w, "// #elif defined(_WIN32)")
	fmt.Fprintln(w, "// #define WIN32_LEAN_AND_MEAN 1")
	fmt.Fprintln(w, "// #include <windows.h>")
	fmt.Fprintln(w, "// #else")
	fmt.Fprintln(w, "// #include <X11/Xlib.h>")
	//fmt.Fprintln(w, "// #include <GL/glx.h>")
	fmt.Fprintln(w, "// #endif")
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

func (p *Package) writeCommands(dir string, useFuncPtrs bool) error {
	w, err := os.Create(filepath.Join(dir, "commands.go"))
	if err != nil {
		return err
	}
	defer w.Close()

	sf := p.Functions.Sort()

	p.writeHeader(w)
	p.writeCgoFlags(w)
	p.writeCIncludes(w)
	p.writeAPIDefinitions(w)
	p.writeCTypes(w)
	if useFuncPtrs {
		p.writeCGetProcAddressFunction(w)
		sf.WriteCFunctionPtrs(w)
		sf.WriteCBridgeDefinitions(w)
	} else {
		sf.WriteCDeclarations(w)
	}
	if useFuncPtrs {
		sf.WriteCInitProcAddresses(w)
	}
	fmt.Fprintln(w, "import \"C\"")
	fmt.Fprintln(w, "import \"errors\"")
	fmt.Fprintln(w, "")
	sf.WriteGoDefinitions(w, useFuncPtrs)
	sf.WriteGoInitPackage(w)
	p.writeFooter(w)

	return nil
}

func (p *Package) GeneratePackage() error {
	dir := ""
	switch p.PackageType {
	case PackageTypeGL:
		dir = filepath.Join("gl", p.Version.String(), p.Name)
	case PackageTypeGLExt:
		dir = filepath.Join("gl", p.Name)
	case PackageTypeGLES:
		dir = filepath.Join("gles", p.Version.String(), p.Name)
	case PackageTypeGLX:
		dir = filepath.Join("glx", p.Version.String(), p.Name)
	case PackageTypeWGL:
		dir = filepath.Join("wgl", p.Version.String(), p.Name)
	case PackageTypeEGL:
		dir = filepath.Join("egl", p.Version.String(), p.Name)
	default:
		return fmt.Errorf("Unknown package type %v", p.PackageType)
	}
		err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = p.writeEnums(dir)
	if err != nil {
		return err
	}
	err = p.writeCommands(dir, true)
	if err != nil {
		return err
	}
	return nil
}

func (ps Packages) GeneratePackages() error {
	for _, p := range ps {
		err := p.GeneratePackage()
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
