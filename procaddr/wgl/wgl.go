// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.
package wgl

// #cgo windows LDFLAGS: -lopengl32
// #define WIN32_LEAN_AND_MEAN 1
// #include <windows.h>
// static HMODULE ogl32dll = NULL;
// void* GoglGetProcAddress(const char* name) { 
// 	void* pf = wglGetProcAddress((LPCSTR)name);
// 	if(pf) {
// 		return pf;
// 	}
// 	if(ogl32dll == NULL) {
// 		ogl32dll = LoadLibraryA("opengl32.dll");
// 	}
// 	return GetProcAddress(ogl32dll, (LPCSTR)name);
// }
import "C"
import "unsafe"
import "github.com/chsc/gogl2/glt"

func GetProcAddress(name string) glt.Pointer {
	var cname *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return glt.Pointer(unsafe.Pointer(C.GoglGetProcAddress(cname)))
}

func init() {
	glt.GetProcAddress = GetProcAddress
}
