/*
* Copyright 2026 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */
package main

import (
	"flag"
	"fmt"

	"github.com/tknie/jaspergo"
)

const description = `This tool converts Jasperreport v6 JRXML files into Jasperreport v7 files.
It can convert either a single file or all files in a directory.
The converted files are stored in the destination directory preserving the
directory structure of the source directory.
`

func main() {
	jaspergo.InitLogLevelWithFile("convert.log")

	source := ""
	sourceFile := ""
	destination := ""
	flag.StringVar(&source, "s", "", "Source directory containing Jasperreport v6 files")
	flag.StringVar(&sourceFile, "f", "", "Source file containing a Jasperreport v6 file")
	flag.StringVar(&destination, "d", "", "Destination directory putting Jasperreport v7 converted files")
	flag.Usage = func() {
		fmt.Println(description)
		fmt.Println("Default flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if source == "" && sourceFile == "" {
		fmt.Println("No source directory or file given")
		flag.Usage()
		return
	}
	if destination == "" {
		fmt.Println("No destination directory given")
		flag.Usage()
		return
	}

	if source != "" {
		err := jaspergo.ConvertDirectoryToPath(source, destination)
		if err != nil {
			fmt.Println("Convert of directory got error:", err)
		}
	}
	if sourceFile != "" {
		err := jaspergo.ConvertFileToPath(sourceFile, destination)
		if err != nil {
			fmt.Println("Convert of directory got error:", err)
		}
	}
}
