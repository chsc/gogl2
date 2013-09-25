// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.
package darwin

// cgo darwin LDFLAGS: -framework OpenGL
// #include <dlfcn.h>
// void* GetProcAddress(const char* name) { 
// 	return dlsym(RTLD_DEFAULT, name);
// }
import "C"
import "unsafe"
import "github.com/chsc/gogl2/glt"

func GetProcAddress(name string) glt.Pointer {
	var n [64]byte
	glt.CopyString(n[:], name)
	return glt.Pointer(unsafe.Pointer(C.GetProcAddress((*C.char)(&n[0]))))
}

func init() {
	glt.GetProcAddress = GetProcAddress
}
