// Copyright 2013 The GoGL2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.mkd file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func generateGoPackages(specsDir string) {
	ps, err := ParseSpecFile(filepath.Join(specsDir, openGLSpecFile))
	if err != nil {
		fmt.Printf("Error: ", err)
	}
	ps.GeneratePackages()
}

func downloadSpec(name string, args []string) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	src := fs.String("src", "khronos", "Source URL or 'khronos'.")
	odir := fs.String("odir", "glspecs", "Output directory for spec files.")
	fs.Parse(args)
	fmt.Println("Downloading specs ...")
	switch *src {
	case "khronos":
		downloadOpenGLSpecs(khronosRegistryBaseURL, *odir)
	default:
		downloadOpenGLSpecs(*src, *odir)
	}
}

func downloadDoc(name string, args []string) {
	fmt.Println("Download docs not implemented.")
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.String("src", "", "Source URL.")
	fs.String("odir", "gldocs", "Output directory for doc files.")
	fs.Parse(args)
	fmt.Println("Downloading docs ...")
	// TODO: download docs
}

func generate(name string, args []string) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	sdir := fs.String("sdir", "glspecs", "OpenGL spec directory.")
	_ = fs.String("ddir", "gldocs", "Documentation directory (currently not used).")
	fs.Parse(args)
	fmt.Println("Generate Bindings ...")
	generateGoPackages(*sdir)
}

func printUsage(name string) {
	fmt.Printf("Usage:     %s command [arguments]\n", name)
	fmt.Println("Commands:")
	fmt.Println(" pullspec  Download spec files.")
	fmt.Println(" pulldoc   Download documentation files.")
	fmt.Println(" generate  Generate bindings.")
	fmt.Printf("Type %s <command> -help for a detailed command description.\n", name)
}

func main() {
	fmt.Println("GoGL2 - OpenGL binding generator for the Go programming language (http://golang.org).")
	fmt.Println("Copyright (c) 2013 by Christoph Schunk. All rights reserved. See LICENSE.mkd for more information.")

	name := os.Args[0]
	args := os.Args[1:]

	if len(args) < 1 {
		printUsage(name)
		os.Exit(-1)
	}

	command := args[0]

	switch command {
	case "pullspec":
		downloadSpec("pullspec", args[1:])
	case "pulldoc":
		downloadDoc("dldoc", args[1:])
	case "generate":
		generate("generate", args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command '%s'.", command)
		printUsage(name)
		os.Exit(-1)
	}
}

