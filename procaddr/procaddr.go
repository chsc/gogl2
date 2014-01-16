// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package procaddr

import "unsafe"

type GetProcAddressFunc func(name string) unsafe.Pointer

// Global function for loading OpenGL procedure addresses.
var GetProcAddress GetProcAddressFunc
