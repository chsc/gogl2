// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"io"
	"sort"
)

type Parameter struct {
	Name string
	Type Type
}

type Function struct {
	Name       string
	Parameters []Parameter
	Return     Type
}

type Functions map[string]*Function
type SortedFunctions []*Function

func (sf SortedFunctions) Len() int {
	return len(sf)
}

func (sf SortedFunctions) Swap(i, j int) {
	sf[i], sf[j] = sf[j], sf[i]
}

func (sf SortedFunctions) Less(i, j int) bool {
	return sf[i].Name < sf[j].Name
}

func (f *Function) writeCParameters(w io.Writer) {
	for i := 0; i < len(f.Parameters); i++ {
		p := &f.Parameters[i]
		if i != 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%s %s", p.Type.CType(), RenameIfReservedCWord(p.Name))
	}
}

func (f *Function) WriteCFunctionPtr(w io.Writer) {
	ctype := f.Return.CType()
	fmt.Fprintf(w, "// %s (APIENTRYP pgl%s)(", ctype, f.Name)
	f.writeCParameters(w)
	fmt.Fprintln(w, ");")
}

func (f *Function) WriteCDeclaration(w io.Writer) {
	ctype := f.Return.CType()
	fmt.Fprintf(w, "// GLAPI %s APIENTRY gl%s(", ctype, f.Name)
	f.writeCParameters(w)
	fmt.Fprintln(w, ");")
}

func (f *Function) WriteCBridgeDefinition(w io.Writer) {
	ctype := f.Return.CType()
	fmt.Fprintf(w, "// %s gogl%s(", ctype, f.Name)
	f.writeCParameters(w)
	fmt.Fprintln(w, ") {")
	if f.Return.IsVoid() {
		fmt.Fprintf(w, "// 	")
	} else {
		fmt.Fprintf(w, "// 	return ")
	}
	fmt.Fprintf(w, "(*pgl%s)(", f.Name)
	for i, _ := range f.Parameters {
		p := &f.Parameters[i]
		if i != 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%s", RenameIfReservedCWord(p.Name))
	}
	fmt.Fprintln(w, ");")
	fmt.Fprintln(w, "// }")
}

func (f *Function) WriteCGetProcAddress(w io.Writer) {
	fmt.Fprintf(w, "// 	if((pgl%s = goglGetProcAddress(\"gl%s\")) == NULL) return 1;\n", f.Name, f.Name)
}

func (f *Function) WriteGoDefinition(w io.Writer, usePtr bool, d *Documentation, majorVersion int) {
	err := d.WriteGoCmdDoc(w, f.Name, majorVersion)
	if err != nil {
		//fmt.Printf("Unable to find function doc: %v\n", err)
	}
	fmt.Fprintf(w, "func %s(", f.Name)
	for i, _ := range f.Parameters {
		p := &f.Parameters[i]
		if i != 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%s %s", RenameIfReservedGoWord(p.Name), p.Type.GoType())
	}
	if f.Return.IsVoid() {
		fmt.Fprintln(w, ") {")
		if usePtr {
			fmt.Fprintf(w, "	C.gogl%s(", f.Name)
		} else {
			fmt.Fprintf(w, "	C.gl%s(", f.Name)
		}
	} else {
		ctype := f.Return.GoType()
		tconv := f.Return.GoConversion()
		fmt.Fprintf(w, ") %s {\n", ctype)
		if usePtr {
			fmt.Fprintf(w, "\treturn %s(C.gogl%s(", tconv, f.Name)
		} else {
			fmt.Fprintf(w, "\treturn %s(C.gl%s(", tconv, f.Name)
		}
	}
	for i, _ := range f.Parameters {
		p := &f.Parameters[i]
		if i != 0 {
			fmt.Fprintf(w, ", ")
		}
		tconv := p.Type.CgoConversion()
		fmt.Fprintf(w, "%s(%s)", tconv, RenameIfReservedGoWord(p.Name))
	}
	if f.Return.IsVoid() {
		fmt.Fprintln(w, ")")
	} else {
		fmt.Fprintln(w, "))")
	}
	fmt.Fprintln(w, "}")
}

func (fs Functions) Sort() SortedFunctions {
	sortedFunctions := make(SortedFunctions, 0, len(fs))
	for _, f := range fs {
		sortedFunctions = append(sortedFunctions, f)
	}
	sort.Sort(sortedFunctions)
	return sortedFunctions
}

func (sf SortedFunctions) WriteCFunctionPtrs(w io.Writer) {
	for _, f := range sf {
		f.WriteCFunctionPtr(w)
	}
	fmt.Fprintln(w, "// ")
}

func (sf SortedFunctions) WriteCDeclarations(w io.Writer) {
	for _, f := range sf {
		f.WriteCDeclaration(w)
	}
	fmt.Fprintln(w, "// ")
}

func (sf SortedFunctions) WriteCBridgeDefinitions(w io.Writer) {
	for _, f := range sf {
		f.WriteCBridgeDefinition(w)
	}
	fmt.Fprintln(w, "// ")
}

func (sf SortedFunctions) WriteCInitProcAddresses(w io.Writer) {
	fmt.Fprintln(w, "// int goglInit() {")
	for _, f := range sf {
		f.WriteCGetProcAddress(w)
	}
	fmt.Fprintln(w, "// \treturn 0;")
	fmt.Fprintln(w, "// }")
}

func (sf SortedFunctions) WriteGoDefinitions(w io.Writer, usePtr bool, d *Documentation, majorVersion int) {
	for _, f := range sf {
		f.WriteGoDefinition(w, usePtr, d, majorVersion)
	}
	fmt.Fprintln(w, "// ")
}

func (sf SortedFunctions) WriteGoInitPackage(w io.Writer) {
	fmt.Fprintln(w, "func Init() error {")
	fmt.Fprintln(w, "\tvar ret C.int")
	fmt.Fprintln(w, "\tif ret = C.goglInit(); ret != 0 {")
	fmt.Fprintln(w, "\t\treturn errors.New(\"unable to initialize OpenGL\")")
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "\treturn nil")
	fmt.Fprintln(w, "}")
}
