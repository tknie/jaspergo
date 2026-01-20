package main

import (
	"flag"
	"fmt"

	"github.com/tknie/jaspergo"
)

const description = `This tool converts Jasperreport v6 JRXML files into Jasperreport v7 files.`

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
