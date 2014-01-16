// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
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

func (f *Function) WriteCFunctionPtrTypedef(w io.Writer) {
	ctype := f.Return.CType()
	fmt.Fprintf(w, "// typedef %s (APIENTRYP PGL%s)(", ctype, strings.ToUpper(f.Name))
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
	fmt.Fprintf(w, "// %s gogl%s(PGL%s glfptr", ctype, f.Name, strings.ToUpper(f.Name))
	if len(f.Parameters) != 0 {
		fmt.Fprintf(w, ", ")
	}
	f.writeCParameters(w)
	fmt.Fprintln(w, ") {")
	if f.Return.IsVoid() {
		fmt.Fprintf(w, "// 	")
	} else {
		fmt.Fprintf(w, "// 	return ")
	}
	fmt.Fprintf(w, "(*glfptr)(")
	for i, _ := range f.Parameters {
		p := &f.Parameters[i]
		if i != 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%s", p.Name)
	}
	fmt.Fprintln(w, ");")
	fmt.Fprintln(w, "// }")
}

func (f *Function) WriteGoFunctionPtr(w io.Writer) {
	fmt.Fprintf(w, "	pgl%s C.PGL%s\n", f.Name, strings.ToUpper(f.Name))
}

func (f *Function) WriteGoGetProcAddress(w io.Writer) {
	fmt.Fprintf(w, "	if pgl%s = (C.PGL%s)(unsafe.Pointer(procaddr.GetProcAddress(\"gl%s\"))); pgl%s == nil { return errors.New(\"gl%s\") }\n", f.Name, strings.ToUpper(f.Name), f.Name, f.Name, f.Name)
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
			fmt.Fprintf(w, "	C.gogl%s(pgl%s", f.Name, f.Name)
			if len(f.Parameters) != 0 {
				fmt.Fprintf(w, ", ")
			}
		} else {
			fmt.Fprintf(w, "	C.gl%s(", f.Name)
		}
	} else {
		ctype := f.Return.GoType()
		tconv := f.Return.GoConversion()
		fmt.Fprintf(w, ") %s {\n", ctype)
		if usePtr {
			fmt.Fprintf(w, "\treturn %s(C.gogl%s(pgl%s", tconv, f.Name, f.Name)
			if len(f.Parameters) != 0 {
				fmt.Fprintf(w, ", ")
			}
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

func (sf SortedFunctions) WriteCFunctionPtrTypedefs(w io.Writer) {
	for _, f := range sf {
		f.WriteCFunctionPtrTypedef(w)
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

func (sf SortedFunctions) WriteGoFunctionPtrs(w io.Writer) {
	fmt.Fprintln(w, "var (")
	for _, f := range sf {
		f.WriteGoFunctionPtr(w)
	}
	fmt.Fprintln(w, ")")
}

func (sf SortedFunctions) WriteGoInitPackage(w io.Writer) {
	fmt.Fprintln(w, "func Init() error {")
	for _, f := range sf {
		f.WriteGoGetProcAddress(w)
	}
	fmt.Fprintln(w, "	return nil")
	fmt.Fprintln(w, "}")
}

func (sf SortedFunctions) WriteGoDefinitions(w io.Writer, usePtr bool, d *Documentation, majorVersion int) {
	for _, f := range sf {
		f.WriteGoDefinition(w, usePtr, d, majorVersion)
	}
	fmt.Fprintln(w, "")
}
