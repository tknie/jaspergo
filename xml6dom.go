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
package jaspergo

import (
	"fmt"
	"io"
	"os"

	"github.com/tknie/log"
	"github.com/tknie/xmlquery"
)

func loadJasperReportsDomFromFile(path string) (*xmlquery.Node, error) {
	fmt.Println("Load and parsing JRXML ...", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := io.Reader(f)
	// Convert UTF-16 XML to UTF-8
	doc, err := xmlquery.Parse(r)
	return doc, err
}

func cleanEmptyNodes(m *xmlquery.Node) {
	for _, nodeName := range []string{"band", "jasperReport", "element", "style", "field",
		"group", "variable", "box", "group", "groupHeader", "groupFooter", "pageFooter",
		"columnHeader", "columnFooter", "pageHeader", "detail"} {
		xmlquery.FindEach(m, "//"+nodeName, func(i int, element *xmlquery.Node) {
			log.Log.Debugf("Clean empty node: %s", element.Data)
			cleanAllEmptySubNodes(element)
		})
	}
}

func cleanAllEmptySubNodes(element *xmlquery.Node) {
	log.Log.Debugf("Clean all empty sub nodes %s", element.Data)
}
