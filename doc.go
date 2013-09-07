// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"strconv"
)

func MakeExtensionSpecDocUrl(vendor, extension string) string {
	return fmt.Sprintf("https://www.opengl.org/registry/specs/%s/%s.txt", vendor, extension)
}

func MakeGLDocUrl(majorVersion int) string {
	manVer := ""
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	return fmt.Sprintf("https://www.opengl.org/sdk/docs/man%s", manVer)
}

func MakeFuncDocUrl(majorVersion int, fName string) string {
	manVer := ""
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	return fmt.Sprintf("https://www.opengl.org/sdk/docs/man%s/xhtml/gl%s.xml", manVer, fName)
}
