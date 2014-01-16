// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package egl

// #cgo linux LDFLAGS: -lEGL
// #include <stdlib.h>
// #include <EGL/egl.h>
import "C"
import "unsafe"

func GetProcAddress(name string) unsafe.Pointer {
	var cname *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return unsafe.Pointer(C.eglGetProcAddress(cname))
}
