// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"fmt"
	"io"
	"sort"
)

type Enum struct {
	Name  string
	Value string
	Group string
}

type Enums map[string]*Enum
type SortedEnums []*Enum

func (f SortedEnums) Len() int {
	return len(f)
}

func (f SortedEnums) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f SortedEnums) Less(i, j int) bool {
	return f[i].Name < f[j].Name
}

func (e Enum) cleanName() string {
	return e.Name
}

func (e Enum) WriteGoDefinition(w io.Writer) {
	fmt.Fprintf(w, "\t%s = %s\n", e.cleanName(), e.Value)
}

func (es Enums) Sort() SortedEnums {
	sortedEnums := make(SortedEnums, 0, len(es))
	for _, e := range es {
		sortedEnums = append(sortedEnums, e)
	}
	sort.Sort(sortedEnums)
	return sortedEnums
}

func (se SortedEnums) WriteGoDefinitions(w io.Writer) {
	fmt.Fprintf(w, "const (\n")
	for _, e := range se {
		e.WriteGoDefinition(w)
	}
	fmt.Fprintf(w, ")\n")
}
