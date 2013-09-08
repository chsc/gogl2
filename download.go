// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	khronosRegistryBaseURL = "https://cvs.khronos.org/svn/repos/ogl/trunk/doc/registry/public/api"
	openGLSpecFile         = "gl.xml"
	eglSpecFile            = "egl.xml"
	wglSpecFile            = "wgl.xml"
	glxSpecFile            = "glx.xml"
)

func makeURL(base, file string) string {
	return fmt.Sprintf("%s/%s", base, file)
}

func downloadFile(baseURL, fileName, outDir string) error {
	fullURL := makeURL(baseURL, fileName)
	fmt.Printf("Downloading %s ...\n", fullURL)
	r, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(absPath, fileName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func downloadOpenGLSpecs(baseURL, outDir string) error {
	err := downloadFile(baseURL, openGLSpecFile, outDir)
	if err != nil {
		return err
	}
	err = downloadFile(baseURL, wglSpecFile, outDir)
	if err != nil {
		return err
	}
	err = downloadFile(baseURL, glxSpecFile, outDir)
	if err != nil {
		return err
	}
	err = downloadFile(baseURL, eglSpecFile, outDir)
	if err != nil {
		return err
	}
	return nil
}
