// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package glx

// #cgo linux LDFLAGS: -lGL
// #include <stdlib.h>
// #include <GL/glx.h>
import "C"
import "unsafe"

func GetProcAddress(name string) unsafe.Pointer {
	var cname *C.GLubyte = (*C.GLubyte)(C.CString(name))
	defer C.free(unsafe.Pointer(cname))
	return unsafe.Pointer(C.glXGetProcAddress(cname))
}
