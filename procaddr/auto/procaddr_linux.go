// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

// +build !egl

package auto

import "github.com/chsc/gogl2/procaddr"
import "github.com/chsc/gogl2/procaddr/glx"

func init() {
	procaddr.GetProcAddress = glx.GetProcAddress
}
