package jaspergo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/tknie/log"

	"github.com/tknie/xmlquery"
)

var re = regexp.MustCompile(`^is([A-Z])`)
var docNode *xmlquery.Node

func ConvertFileToPath(file, destination string) error {
	source := filepath.Dir(file)
	fileName := filepath.Base(file)
	return convertFile(source, fileName, destination)
}

func convertFile(path, fileName, destination string) error {
	m, err := LoadJasperReportsDomFromFile(path + string(os.PathSeparator) + fileName)
	if err != nil {
		fmt.Printf("Error loading '%s': %v\n", fileName, err)
		return err
	}
	if m == nil {
		return fmt.Errorf("DOM read and parsing error")
	}
	docNode = m
	xmlquery.FindEach(m, "//style", func(i int, style *xmlquery.Node) {
		changeBooleanPrefix(style)

		for i, attr := range style.Attr {
			if attr.Name.Local == "fontSize" {
				f, _ := strconv.ParseFloat(attr.Value, 64)
				style.Attr[i].Value = fmt.Sprintf("%.1f", f)
			}
		}
	})
	jasperReprt := m.SelectElement("jasperReport")
	rs := []string{"xmlns", "xsi", "schemaLocation", "name"}
	attr := jasperReprt.Attr
	name := fileName
	for _, a := range attr {
		if a.Name.Local == "name" {
			name = a.Value
			break
		}
	}
	jasperReprt.Attr = []xmlquery.Attr{}
	jasperReprt.SetAttr("name", name)
	jasperReprt.SetAttr("language", "java")
	for _, a := range attr {
		if !slices.Contains(rs, a.Name.Local) {
			jasperReprt.Attr = append(jasperReprt.Attr, a)
		}
	}
	changeBooleanPrefix(jasperReprt)
	xmlquery.FindEach(m, "//crosstab", workConvertElements)
	xmlquery.FindEach(m, "//textField", workConvertElements)
	xmlquery.FindEach(m, "//staticText", workConvertElements)
	xmlquery.FindEach(m, "//subreport", workConvertElements)
	xmlquery.FindEach(m, "//line", workConvertElements)
	xmlquery.FindEach(m, "//crosstabRowHeader", removeCrosstabRowHeaderNode)
	xmlquery.FindEach(m, "//crosstabTotalRowHeader", removeCrosstabTotalRowHeader)

	xmlquery.FindEach(m, "//componentElement", workConvertElements)
	xmlquery.FindEach(m, "//crosstabCell", func(i int, crosstabCell *xmlquery.Node) {
		crosstabCell.Data = "cell"
		for node := crosstabCell.FirstChild; node != nil; node = node.NextSibling {
			if node.Data == "cellContents" {
				node.Data = "contents"
			}
		}
	})

	xmlquery.FindEach(m, "//crosstabColumnHeader", func(i int, crosstabColumnHeader *xmlquery.Node) {
		crosstabColumnHeader.Data = "header"
		cellContents := crosstabColumnHeader.SelectElement("cellContents")
		if cellContents != nil {
			crosstabColumnHeader.Attr = cellContents.Attr
			xmlquery.RemoveFromTree(cellContents)
		}
	})
	xmlquery.FindEach(m, "//crosstabTotalColumnHeader", func(i int, crosstabColumnHeader *xmlquery.Node) {
		crosstabColumnHeader.Data = "totalHeader"
		cellContents := crosstabColumnHeader.SelectElement("cellContents")
		if cellContents != nil {
			crosstabColumnHeader.Attr = cellContents.Attr
			xmlquery.RemoveFromTree(cellContents)
		}
	})

	xmlquery.FindEach(m, "//subDataset", func(i int, subDataset *xmlquery.Node) {
		subDataset.Data = "dataset"
	})

	expressionList := []string{"groupExpression", "bucketExpression", "variableExpression",
		"datasetParameterExpression", "measureExpression"}
	for _, e := range expressionList {
		xmlquery.FindEach(m, "//"+e, func(i int, expression *xmlquery.Node) {
			expression.Data = "expression"
		})
	}
	xmlquery.FindEach(m, "//group", func(i int, group *xmlquery.Node) {
		changeBooleanPrefix(group)
		cleanAllEmptySubNodes(group)
	})
	xmlquery.FindEach(m, "//datasetParameter", func(i int, parameter *xmlquery.Node) {
		parameter.Data = "parameter"
	})
	xmlquery.FindEach(m, "//pageHeader", removeBandNode)
	xmlquery.FindEach(m, "//title", removeBandNode)
	xmlquery.FindEach(m, "//summary", removeBandNode)
	xmlquery.FindEach(m, "//lastPageFooter", removeBandNode)
	xmlquery.FindEach(m, "//pageFooter", removeBandNode)
	xmlquery.FindEach(m, "//columnHeader", removeBandNode)
	xmlquery.FindEach(m, "//columnFooter", removeBandNode)
	// xmlquery.FindEach(m, "//groupFooter", removeBandNode)

	xmlquery.FindEach(m, "//band", removeNamedSubNodes)
	xmlquery.FindEach(m, "//c:table", func(i int, cTable *xmlquery.Node) {
		cTable.Prefix = ""
		cTable.Data = "component"
		cTable.Attr = []xmlquery.Attr{}
		cTable.SetAttr("kind", "table")
	})
	xmlquery.FindEach(m, "//c:column", func(i int, node *xmlquery.Node) {
		node.Prefix = ""
		node.SetAttr("kind", "single")
		xmlquery.FindEach(node, "//c:columnFooter", func(i int, node *xmlquery.Node) {
			node.Prefix = ""
		})
		xmlquery.FindEach(node, "//c:detailCell", func(i int, node *xmlquery.Node) {
			node.Prefix = ""
		})
	})
	xmlquery.FindEach(m, "//property", func(i int, property *xmlquery.Node) {
		changeBooleanPrefix(property)
		for _, attr := range property.Attr {
			if attr.Name.Local == "name" {
				switch attr.Value {
				case "com.jaspersoft.studio.data.defaultdataadapter",
					"com.jaspersoft.studio.unit.height", "com.jaspersoft.studio.unit.width":
					xmlquery.RemoveFromTree(property)
					return
				case "com.jaspersoft.studio.unit.x", "com.jaspersoft.studio.unit.y":
					xmlquery.RemoveFromTree(property)
					return
				case "com.jaspersoft.studio.layout":
					xmlquery.RemoveFromTree(property)
					return
				}
			}
		}
	})
	cleanEmptyNodes(m)
	fb := []byte(m.OutputXMLWithOptions(xmlquery.WithIndentation("\t"), xmlquery.WithEmptyTagSupport(),
		xmlquery.WithoutPreserveSpace()))
	// fb := []byte(m.OutputXML(false))
	err = os.WriteFile(destination+string(os.PathSeparator)+fileName, fb, 0644)
	if err != nil {
		return err
	}
	return nil
}

func removeNamedSubNodes(i int, band *xmlquery.Node) {
	changeBooleanPrefix(band)
	// removeSubNode(band, "property")
	cleanAllEmptySubNodes(band)
}

func removeSubNode(mainNode *xmlquery.Node, removeNodeName string) {
	var lastNode *xmlquery.Node
	log.Log.Debugf("Remove %s from %s", removeNodeName, mainNode.Data)
	log.Log.Debugf("Nodes: %d", len(mainNode.ChildNodes()))
	for node := mainNode.FirstChild; node != nil; node = node.NextSibling {
		log.Log.Debugf("Checking %s %p", string(node.Data), node)
		if node.Data == removeNodeName || node.Data[0] == 10 {
			if node.Data == removeNodeName {
				node.Data = string([]byte{10})
			}
			log.Log.Debugf("Removing %s %p", string(node.Data), lastNode)
			if lastNode == nil {
				mainNode.FirstChild = node.NextSibling
				node.NextSibling.PrevSibling = nil
			} else {
				lastNode.NextSibling = node.NextSibling
				node.PrevSibling = lastNode
			}
			if node.NextSibling == nil {
				mainNode.LastChild = lastNode
			}
			node.PrevSibling = nil
			node.Parent = nil
		} else {
			if lastNode == nil {
				mainNode.FirstChild = node
			}
			lastNode = node
			if node.NextSibling == nil {
				mainNode.LastChild = lastNode
			}
		}
	}
	log.Log.Debugf("End Nodes: %d", len(mainNode.ChildNodes()))
	for child := mainNode.FirstChild; child != nil; child = child.NextSibling {
		log.Log.Debugf("%s/%s Final: %s", mainNode.Parent.Data, mainNode.Data, child.Data)
	}
}

func ConvertDirectoryToPath(path, destination string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error opening source directory %s: %v\n", path, err)
		log.Log.Fatal(err)
	}

	sample := os.Getenv("JASPER_SAMPLE_TEST")

	successCount := 0
	failureCount := 0
	for _, e := range entries {
		if sample != "" && e.Name() != sample {
			continue
		}
		err := convertFile(path, e.Name(), destination)
		if err != nil {
			failureCount++
		} else {
			successCount++
		}
	}
	log.Log.Infof("Success: %d Error: %d", successCount, failureCount)
	return nil
}

func removeBandNode(i int, node *xmlquery.Node) {
	changeBooleanPrefix(node)
	band := node.SelectElement("band")
	node.Attr = band.Attr
	xmlquery.MoveChildNodes(band, node)
	// for _, n := range band.ChildNodes() {
	// 	if n.Data != "property" && n.Data[0] != 10 {
	// 		fmt.Println("PP", node.Data, n.Data)
	// 		// xmlquery.AddChild(node, n)
	// 		// xmlquery.AddImmediateSibling(node, n)
	// 	}
	// }
	log.Log.Debugf("Remove band from %s", node.Data)
	xmlquery.RemoveFromTree(band)
	cleanAllEmptySubNodes(node)
	if node.Data == "pageHeader" {
		logNode(docNode)
		log.Log.Debugf("XML:" + node.OutputXML(true))
	}

}

func removeCrosstabRowHeaderNode(i int, node *xmlquery.Node) {
	changeBooleanPrefix(node)
	node.Data = "header"
	cellContents := node.SelectElement("cellContents")
	if cellContents != nil {
		node.Attr = cellContents.Attr
		for _, n := range cellContents.ChildNodes() {
			xmlquery.AddSibling(cellContents, n)
		}
		xmlquery.RemoveFromTree(cellContents)
		cleanAllEmptySubNodes(node)
	}

}

func removeCrosstabTotalRowHeader(i int, node *xmlquery.Node) {
	changeBooleanPrefix(node)
	node.Data = "totalHeader"
	cellContents := node.SelectElement("cellContents")
	if cellContents != nil {
		node.Attr = cellContents.Attr
		for _, n := range cellContents.ChildNodes() {
			xmlquery.AddSibling(cellContents, n)
		}
		xmlquery.RemoveFromTree(cellContents)
		cleanAllEmptySubNodes(node)
	}
}

func workConvertElements(i int, n *xmlquery.Node) {
	saveName := n.Data
	if saveName == "componentElement" {
		saveName = "component"
	}
	n.Data = "element"
	saveAttr := n.Attr
	n.Attr = []xmlquery.Attr{}
	n.SetAttr("kind", saveName)
	reportElement := n.SelectElement("reportElement")
	if reportElement != nil {
		uuid := reportElement.SelectAttr("uuid")
		n.SetAttr("uuid", uuid)

		for _, a := range reportElement.Attr {
			if a.Name.Local == "uuid" || a.Name.Local == "evaluationTime" {
				continue
			}
			n.SetAttr(a.Name.Local, a.Value)
		}
	}
	for _, a := range saveAttr {
		n.SetAttr(a.Name.Local, a.Value)
	}

	for _, a := range saveAttr {
		n.SetAttr(a.Name.Local, a.Value)
	}
	if reportElement != nil {
		xmlquery.MoveChildNodes(reportElement, n)
		// for childNodes := reportElement.FirstChild; childNodes != nil; childNodes = childNodes.NextSibling {
		// 	if childNodes.Data[0] != 10 {
		// 		xmlquery.AddSibling(n, childNodes)
		// 	}
		// }
	}
	previous := reportElement
	textFieldExpression := n.SelectElement("textFieldExpression")
	if textFieldExpression != nil {
		textFieldExpression.Data = "expression"
		previous = textFieldExpression
	}

	subreportExpression := n.SelectElement("subreportExpression")
	if subreportExpression != nil {
		subreportExpression.Data = "expression"
		previous = subreportExpression
	}
	textElement := n.SelectElement("textElement")
	if textElement != nil {
		for _, a := range textElement.Attr {
			if a.Name.Local == "textAlignment" {
				// n.SetAttr("textAdjust", "StretchHeight")
				a.Name.Local = "hTextAlign"
			}
			n.Attr = append(n.Attr, a)
		}
		font := textElement.SelectElement("font")
		if font != nil {
			for _, a := range font.Attr {
				a.Name.Local = "fontSize"
				n.Attr = append(n.Attr, a)
			}
		}
		paragraph := textElement.SelectElement("paragraph")
		if paragraph != nil {
			xmlquery.RemoveFromTree(paragraph)
			xmlquery.AddImmediateSibling(previous, paragraph)
			previous = paragraph
		}
		cleanAllEmptySubNodes(textElement)
	}
	if t := n.SelectElement("text"); t != nil {
		previous = t
	}
	if reportElement != nil {
		property := reportElement.SelectElement("property")
		if property != nil {
			for _, attr := range property.Attr {
				if attr.Name.Local == "name" {
					fmt.Println(attr.Name.Local, attr.Value)
					switch attr.Value {
					case "com.jaspersoft.studio.unit.x", "com.jaspersoft.studio.unit.y",
						"com.jaspersoft.studio.layout", "com.jaspersoft.studio.data.defaultdataadapter":
						return
					default:
						fmt.Println("Property in", reportElement.Data, previous.Data)
						xmlquery.AddSibling(n, property)
						previous = property
					}
				}
			}
		}
	}

	box := n.SelectElement("box")
	if textElement != nil {
		xmlquery.RemoveFromTree(textElement)
	}
	if reportElement != nil {
		xmlquery.RemoveFromTree(reportElement)
	}
	if box != nil {
		xmlquery.RemoveFromTree(box)
		xmlquery.AddImmediateSibling(previous, box)
		cleanAllEmptySubNodes(box)

	}
	changeBooleanPrefix(n)
	cleanAllEmptySubNodes(n)
}

func changeBooleanPrefix(element *xmlquery.Node) {
	for i, a := range element.Attr {
		if t, _ := regexp.MatchString("^is[A-Z]", a.Name.Local); t {
			s := re.ReplaceAllString(a.Name.Local, `$1`)
			s = strings.ToLower(s[0:1]) + s[1:]
			element.Attr[i].Name.Local = s
		}
	}
}
