// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"path/filepath"
)

type docRefEntry struct {
	XMLName    xml.Name     `xml:"refentry"`
	Id         string       `xml:"id,attr"`
	RefName    string       `xml:"refnamediv>refname"`
	RefPurpose string       `xml:"refnamediv>refpurpose"`
}

type docSvnIndex struct {
	XMLName    xml.Name     `xml:"svn"`
	FileRefs   []docFileRef `xml:"index>file"`
}

type docFileRef struct {
	Ref string `xml:"href,attr"`
}

type docFile struct {
	BaseName string
	FileName string
}

type docFunc struct {
	Purpose string
}

func writeKhronosDocCopyright(w io.Writer) {
	fmt.Fprintln(w, "// Copyright (c) 2010 Khronos Group.")
	fmt.Fprintln(w, "// This material may be distributed subject to the terms and conditions")
	fmt.Fprintln(w, "// set forth in the Open Publication License, v 1.0, 8 June 1999.")
	fmt.Fprintln(w, "// http://opencontent.org/openpub/.")
	fmt.Fprintln(w, "// ")
}

func writeSgiDocCopyright(w io.Writer) {
	fmt.Fprintln(w, "// Copyright (c) 1991-2006 Silicon Graphics, Inc.")
	fmt.Fprintln(w, "// This document is licensed under the SGI Free Software B License.")
	fmt.Fprintln(w, "// For details, see http://oss.sgi.com/projects/FreeB.")
	fmt.Fprintln(w, "//")
}

func writeFuncDocUrl(w io.Writer, majorVersion int, fName string) {
	manVer := "2"
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	fmt.Fprintf(w, "http://www.opengl.org/sdk/docs/man%s/xhtml/gl%s.xml", manVer, fName)
}

func makeExtensionSpecDocUrl(vendor, extension string) string {
	return fmt.Sprintf("https://www.opengl.org/registry/specs/%s/%s.txt", vendor, extension)
}

func makeGLDocUrl(majorVersion int) string {
	manVer := ""
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	return fmt.Sprintf("https://www.opengl.org/sdk/docs/man%s", manVer)
}

func makeFuncDocUrl(majorVersion int, fName string) string {
	manVer := "2"
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	return fmt.Sprintf("https://www.opengl.org/sdk/docs/man%s/xhtml/gl%s.xml", manVer, fName)
}


func readFileNonStrict(fileName string, data interface{}) error {
	reader, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer reader.Close()
	decoder := xml.NewDecoder(reader)
	decoder.Strict = false
	return decoder.Decode(data)
}

func parseDocIndex(fileName string) ([]docFile, error) {
	var di docSvnIndex
	if err := readFileNonStrict(fileName, &di); err != nil {
		fmt.Println("###", err)
		return nil, err
	}
	files := make([]docFile, 0, 256)
	for _, fr := range di.FileRefs {
		if strings.HasPrefix(fr.Ref, "glu") { // ignore
			continue
		}
		if strings.HasPrefix(fr.Ref, "glX") { // ignore
			continue
		}
		if strings.HasPrefix(fr.Ref, "gl") {
			fn := strings.TrimPrefix(strings.TrimSuffix(fr.Ref , ".xml"), "gl")
			files = append(files, docFile{BaseName: fn, FileName: fr.Ref})
		}
	}
	return files, nil
}

func parseDocFile(fileName string) (*docFunc, error) {
	var d docRefEntry
	if err := readFileNonStrict(fileName, &d); err != nil {
		return nil, err
	}
	return &docFunc{Purpose: d.RefPurpose}, nil
}

func downloadDocs(url, docCat, outDir string) error {
	complOutDir := filepath.Join(outDir, docCat)
	err := downloadFile(url, docCat, complOutDir, "index.xml")
	if err != nil {
		return err
	}
	files, err := parseDocIndex(filepath.Join(complOutDir, "index.xml"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = downloadFile(fmt.Sprintf("%s/%s", url, docCat), file.FileName, complOutDir, file.FileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFuncBaseName(name string, baseNames []docFile) string {
	for _, bn := range baseNames {
			if strings.HasPrefix(name, bn.BaseName) {
				return bn.BaseName
			}
	}
	return name
}
