package jaspergo

import (
	"fmt"
	"io"
	"os"

	"github.com/tknie/log"
	"github.com/tknie/xmlquery"
)

func LoadJasperReportsDomFromFile(path string) (*xmlquery.Node, error) {
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
	// var lastNode *xmlquery.Node
	// for node := element.FirstChild; node != nil; node = node.NextSibling {
	// 	if node.Data[0] == 10 {
	// 		if node == element.FirstChild || lastNode == nil {
	// 			element.FirstChild = node.NextSibling
	// 			node.NextSibling.PrevSibling = nil
	// 			node.Parent = nil
	// 			node.PrevSibling = nil
	// 		} else {
	// 			lastNode.NextSibling = node.NextSibling
	// 			if node.NextSibling != nil {
	// 				node.NextSibling.PrevSibling = lastNode
	// 			}
	// 			node.Parent = nil
	// 			node.PrevSibling = nil
	// 		}
	// 	} else {
	// 		lastNode = node
	// 	}
	// }
	// element.LastChild = lastNode
}
