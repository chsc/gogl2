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

type DocFunc struct {
	BaseName string
	Purpose  string
	// TODO: more?
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
	fmt.Fprintf(w, "https://www.opengl.org/sdk/docs/man%s/xhtml/gl%s.xml", manVer, fName)
}

func makeExtensionSpecDocUrl(vendor, extension string) string {
	return fmt.Sprintf("https://www.opengl.org/registry/specs/%s/%s.txt", vendor, extension)
}

func makeGLDocUrl(majorVersion int) string {
	manVer := "2"
	if majorVersion >= 3 {
		manVer = strconv.Itoa(majorVersion)
	}
	return fmt.Sprintf("https://www.opengl.org/sdk/docs/man%s", manVer)
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

func DownloadDocs(url, docCat, outDir string) error {
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

func parseDocFile(fileName string) (*DocFunc, error) {
	var d docRefEntry
	if err := readFileNonStrict(fileName, &d); err != nil {
		return nil, err
	}
	return &DocFunc{Purpose: d.RefPurpose}, nil
}

func parseDocs(docCat, dir string) ([]*DocFunc, error) {
	complOutDir := filepath.Join(dir, docCat)
	files, err := parseDocIndex(filepath.Join(complOutDir, "index.xml"))
	if err != nil {
		return nil, err
	}
	docFuncs := make([]*DocFunc, 0, 256)
	for _, file := range files {
		df, err := parseDocFile(filepath.Join(complOutDir, file.FileName))
		if err != nil {
			return nil, err
		}
		df.BaseName = file.BaseName		
		docFuncs = append(docFuncs, df)
	}
	return docFuncs, nil
}

func ParseAllDocs(dir string) ([]*DocFunc, error) {
	df2, err := parseDocs("man2", dir)
	if err != nil {
		return nil, err
	}
	for _, ff := range df2 {
		fmt.Println("DOC", ff)
	}
	//TODO: distinguish between versions
/*	df3, err := parseDocs("man3", dir)
	if err != nil {
		return nil, err
	}
	df4, err := parseDocs("man4", dir)
	if err != nil {
		return nil, err
	}*/	
	return df2, err
}

func GetFuncDoc(name string, funcDocs []*DocFunc) (*DocFunc, error) {
//	fmt.Println("DOC", name)
	for _, fd := range funcDocs {
			if strings.HasPrefix(name, fd.BaseName) {
				return fd, nil
			}
	}
	return nil, fmt.Errorf("unable to find doc for function %s", name)
}
