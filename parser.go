// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type SpecRegistry struct {
	XMLName  xml.Name        `xml:"registry"`
	Comment  string          `xml:"comment"`
	Types    []SpecType      `xml:"types>type"`
	Groups   []SpecGroup     `xml:"groups>group"`
	Enums    []SpecEnumToken `xml:"enums"`
	Commands []SpecCommand   `xml:"commands>command"`
	Features []SpecFeature   `xml:"feature"`
	// TODO: extensions
}

type SpecType struct {
	Name     string `xml:"name,attr"`
	Comment  string `xml:"comment,attr"`
	Requires string `xml:"requires,attr"`
	Api      string `xml:"api,attr"`
	Inner    []byte `xml:",innerxml"`
}

type SpecGroup struct {
	Name    string      `xml:"name,attr"`
	Comment string      `xml:"comment,attr"`
	Enums   []SpecGEnum `xml:"enum"`
}

type SpecGEnum struct {
	Name string `xml:"name,attr"`
}

type SpecEnumToken struct {
	Namespace string      `xml:"namespace,attr"`
	Group     string      `xml:"group,attr"`
	Type      string      `xml:"type,attr"`
	Comment   string      `xml:"comment,attr"`
	Enums     []SpecTEnum `xml:"enum"`
	// TODO: vendor, start, end attr
}

type SpecTEnum struct {
	Value string `xml:"value,attr"`
	Name  string `xml:"name,attr"`
}

type SpecCommand struct {
	Proto  SpecProto   `xml:"proto"`
	Params []SpecParam `xml:"param"`
}

type SpecSignature []byte

type SpecProto struct {
	Inner SpecSignature `xml:",innerxml"`
}

type SpecParam struct {
	Group string        `xml:"group,attr"`
	Len   string        `xml:"len,attr"`
	Inner SpecSignature `xml:",innerxml"`
}

type SpecFeature struct {
	Api      string        `xml:"api,attr"`
	Name     string        `xml:"name,attr"`
	Number   string        `xml:"number,attr"`
	Requires []SpecRequire `xml:"require"`
	Removes  []SpecRemove  `xml:"remove"`
}

type SpecRequire struct {
	Comment  string           `xml:"comment,attr"`
	Enums    []SpecEnumRef    `xml:"enum"`
	Commands []SpecCommandRef `xml:"command"`
}

type SpecRemove struct {
	Comment  string           `xml:"comment,attr"`
	Enums    []SpecEnumRef    `xml:"enum"`
	Commands []SpecCommandRef `xml:"command"`
}

type SpecEnumRef struct {
	Name string `xml:"name,attr"`
}

type SpecCommandRef struct {
	Name string `xml:"name,attr"`
}

func (st *SpecType) Parse() (TypeDef, error) {
	typed := TypeDef{Name: st.Name, Comment: st.Comment, Api: st.Api, CDefinition: ""}
	readName := false
	decoder := xml.NewDecoder(bytes.NewBuffer(st.Inner))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return typed, err
		}
		switch t := token.(type) {
		case xml.CharData:
		//	fmt.Println("- ", (string)(t))
			typed.CDefinition += (string)(t)
			if readName {
				typed.Name = (string)(t)
			}
		case xml.StartElement:
			if t.Name.Local == "name" {
				readName = true
			} else if t.Name.Local == "apientry" {
				typed.CDefinition += "APIENTRY"
			} else {
				return typed, fmt.Errorf("Wrong start element: %s", t.Name.Local)
			}
		case xml.EndElement:
			if t.Name.Local == "name" {
				readName = false
			} else if t.Name.Local == "apientry" {
			} else {
				return typed, fmt.Errorf("Wrong start element: %s", t.Name.Local)
			}
		}
	}
	//fmt.Println("-", typed)
	return typed, nil
}

func (r SpecRegistry) ParseTypedefs() ([]TypeDef, error) {
	tdefs := make([]TypeDef, 0, len(r.Types))
	for _, s := range r.Types {
		td, err := s.Parse()
		if err != nil {
			return nil, err
		}
		tdefs = append(tdefs, td)
	}
	return tdefs, nil
}

func (si SpecSignature) Parse() (string, Type, error) {
	name := ""
	ctype := Type{}
	readName := false
	readType := false
	decoder := xml.NewDecoder(bytes.NewBuffer(si))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return name, ctype, err
		}
		switch t := token.(type) {
		case xml.CharData:
			s := strings.Trim((string)(t), " ")
			//fmt.Println(" char", s)
			if readName {
				name = TrimGLCmdPrefix(s)
			} else if readType {
				ctype.Name = s
			} else if s == "" {
				// skip
			} else if s == "void" {
				ctype.Name = "void"
			} else if s == "void *" {
				ctype.Name = "void"
				ctype.PointerLevel = 1
			} else if s == "const void *" {
				ctype.Name = "void"
				ctype.IsConst = true
				ctype.PointerLevel = 1
			} else if s == "*" {
				ctype.PointerLevel = 1
			} else if s == "**" {
				ctype.PointerLevel = 2
			} else if s == "*const*" {
				ctype.PointerLevel = 2
			} else if s == "const" {
				ctype.IsConst = true
			} else {
				return name, ctype, fmt.Errorf("Unknown %s", s)
			}
		case xml.StartElement:
			//fmt.Println(" se", t.Name.Local)
			if t.Name.Local == "ptype" {
				readType = true
			} else if t.Name.Local == "name" {
				readName = true
			} else {
				return name, ctype, fmt.Errorf("Wrong start element: %s", t.Name.Local)
			}
		case xml.EndElement:
			//fmt.Println(" ee", t.Name.Local)
			if t.Name.Local == "ptype" {
				readType = false
			} else if t.Name.Local == "name" {
				readName = false
			} else {
				return name, ctype, fmt.Errorf("Wrong end element: %s", t.Name.Local)
			}
		}
	}
	return name, ctype, nil
}

func readSpecFile(file string) (*SpecRegistry, error) {
	var reg SpecRegistry
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := xml.NewDecoder(f)
	err = d.Decode(&reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func commandsToFunctions(commands []SpecCommand) Functions {
	functions := make(Functions)
	for _, c := range commands {
		cname, ct, err := c.Proto.Inner.Parse()
		if err != nil {
			fmt.Printf("Unable to parse proto signature '%s': %s\n", string(c.Proto.Inner), err)
		} else {
			parameters := make([]Parameter, 0, 4)
			for _, p := range c.Params {
				pname, pt, err := p.Inner.Parse()
				if err != nil {
					fmt.Printf("Unable to parse parameter signature '%s' of function '%s': %s\n", (string)(p.Inner), cname, err)
				} else {
					parameters = append(parameters, Parameter{Name: pname, Type: pt})
				}
			}
			//fmt.Println(cname)
			functions[cname] = &Function{Name: cname, Parameters: parameters, Return: ct}
		}
	}
	return functions
}

func findEnum(enumName string, est []SpecEnumToken) (string, string) {
	for _, es := range est {
		//fmt.Println(es.Type)
		for _, e := range es.Enums {
			if e.Name == enumName {
				return e.Value, es.Group
			}
		}
	}
	return "", ""
}

func addEnums(ps []*Package, ver Version, enumNames []SpecEnumRef, et []SpecEnumToken) {
	for _, pc := range ps {
		if pc.Version.Compare(ver) < 0 {
			continue
		}
		//fmt.Println("adding enums to package", pc.Name, ver)
		for _, en := range enumNames {
			val, grp := findEnum(en.Name, et)
			if val == "" {
				fmt.Println("Not found:", en.Name)
			}
			//fmt.Println("adding", en)
			pc.Enums[en.Name] = &Enum{Name: strings.TrimPrefix(en.Name, "GL_"), Value: val, Group: grp}
		}
	}
}

func removeEnums(ps Packages, ver Version, enumNames []SpecEnumRef) {
	for _, pc := range ps {
		if pc.Version.Compare(ver) < 0 {
			continue
		}
		//fmt.Println("removing enums from package", pc.Name, ver)
		for _, en := range enumNames {
			if _, ok := pc.Enums[en.Name]; ok {
				delete(pc.Enums, en.Name)
			}
		}
	}
}

func addCommands(ps Packages, ver Version, cmdNames []SpecCommandRef, functions Functions) {
	for _, pc := range ps {
		if pc.Version.Compare(ver) < 0 {
			continue
		}
		//fmt.Println("adding enums to package", pc.Name, ver)
		for _, cn := range cmdNames {
			fname := strings.TrimPrefix(cn.Name, "gl")
			f, ok := functions[fname]
			if !ok {
				fmt.Println("add cmd: Cmd not found:", fname)
			} else {
				//fmt.Println("adding", cn)
				pc.Functions[fname] = f
			}
		}
	}
}

func removeCommands(ps Packages, ver Version, cmdNames []SpecCommandRef) {
	for _, pc := range ps {
		if pc.Version.Compare(ver) < 0 {
			continue
		}
		fmt.Println("removing cmds from package", pc.Name, ver)
		for _, cn := range cmdNames {
			fname := strings.TrimPrefix(cn.Name, "gl")
			if _, ok := pc.Functions[fname]; !ok {
				fmt.Println("Remove cmd: Cmd not found", fname)
			} else {
				delete(pc.Functions, fname)
			}
		}
	}
}

func ParseSpecFile(file string) (Packages, error) {
	pacs := make(Packages, 0)

	reg, err := readSpecFile(file)
	if err != nil {
		return nil, err
	}

	functions := commandsToFunctions(reg.Commands)
	tds, err := reg.ParseTypedefs()
	if err != nil {
			return nil, err
	}
	
	for _, ft := range reg.Features {
		if ft.Api != "gl" || ft.Number != "2.1" { // TODO: only for testing
			continue
		}

		version, err := ParseVersion(ft.Number)
		if err != nil { // currently only for gl
			return nil, err
		}

		ptype := PackageTypeUnknown
		pname := ""
		switch ft.Api {
		case "gl":
			ptype = PackageTypeGL
			pname = "gl"
		case "gles1", "gles2":
			ptype = PackageTypeGLES
			pname = "gles"
		default:
			return nil, fmt.Errorf("Unknown api: '%s'", ft.Api)
		}

		p := &Package{PackageType: ptype, Name: pname, Version: version, TypeDefs: tds, Enums: make(Enums), Functions: make(Functions)}
		pacs = append(pacs, p)
	}

	for _, f := range reg.Features {
		if f.Api != "gl" || f.Number != "2.1" { // TODO: only for testing
			continue
		}

		fmt.Println("Feature:", f.Api, f.Name, f.Number)

		version, err := ParseVersion(f.Number)
		if err != nil {
			return nil, err
		}

		for _, r := range f.Requires {
			addEnums(pacs, version, r.Enums, reg.Enums)
		}
		for _, d := range f.Removes {
			removeEnums(pacs, version, d.Enums)
		}

		for _, r := range f.Requires {
			addCommands(pacs, version, r.Commands, functions)
		}
		for _, d := range f.Removes {
			removeCommands(pacs, version, d.Commands)
		}

	}

	return pacs, nil
}
