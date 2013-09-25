// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.
package glt

import (
	"reflect"
	"fmt"
)

type Pointer uintptr
type GetProcAddressFunc func(name string) Pointer

var GetProcAddress GetProcAddressFunc 

func Ptr(data interface{}) Pointer {
	if data == nil {
		return Pointer(0)
	}
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr: // for pointers: *byte, *int, ...
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return Pointer(e.UnsafeAddr())
		}
	case reflect.Uintptr:
		return Pointer(v.Pointer())
	case reflect.Slice: // for slices and arrays: []int, []float32, ...
		return Pointer(v.Index(0).UnsafeAddr())
	case reflect.Array:
		e := v.Index(0)
		fmt.Println(e, e.Kind())	
		return Pointer(v.UnsafeAddr())
	}
	panic(fmt.Sprintf("unknown type: %s: must be a pointer or a slice.", v.Type()))
}

func (p Pointer) Offset(o uintptr) Pointer {
	return Pointer(uintptr(p) + uintptr(o))
}

func CopyString(dest []byte, str string) {
	for i := 0; i < len(str); i++ {
		dest[i] = str[i]
	}
	dest[len(str)] = 0
}

/*
//Go bool to GL boolean.
func GLBool(b bool) Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

// GL boolean to Go bool.
func GoBool(b Boolean) bool {
	return b == TRUE
}

// Go string to GL string.
func GLString(str string) *Char {
	return (*Char)(C.CString(str))
}

// Allocates a GL string.
func GLStringAlloc(length Sizei) *Char {
	return (*Char)(C.malloc(C.size_t(length)))
}

// Frees GL string.
func GLStringFree(str *Char) {
	C.free(unsafe.Pointer(str))
}

// GL string (GLchar*) to Go string.
func GoString(str *Char) string {
	return C.GoString((*C.char)(str))
}

// GL string (GLubyte*) to Go string.
func GoStringUb(str *Ubyte) string {
	return C.GoString((*C.char)(unsafe.Pointer(str)))
}

// GL string (GLchar*) with length to Go string.
func GoStringN(str *Char, length Sizei) string {
	return C.GoStringN((*C.char)(str), C.int(length))
}

// Converts a list of Go strings to a slice of GL strings.
// Usefull for ShaderSource().
func GLStringArray(strs ...string) []*Char {
	strSlice := make([]*Char, len(strs))
	for i, s := range strs {
		strSlice[i] = (*Char)(C.CString(s))
	}
	return strSlice
}

// Free GL string slice allocated by GLStringArray().
func GLStringArrayFree(strs []*Char) {
	for _, s := range strs {
		C.free(unsafe.Pointer(s))
	}
}

// Add offset to a pointer. Usefull for VertexAttribPointer, TexCoordPointer, NormalPointer, ...
func Offset(p Pointer, o uintptr) Pointer {
	return Pointer(uintptr(p) + o)
}
*/
