package jaspergo

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/go-openapi/testify/v2/assert"
)

var re = regexp.MustCompile(`^is([A-Z])`)

func validateAttr(element *xmlquery.Node) {
	for i, a := range element.Attr {
		if t, _ := regexp.MatchString("^is[A-Z]", a.Name.Local); t {
			s := re.ReplaceAllString(a.Name.Local, `$1`)
			s = strings.ToLower(s[0:1]) + s[1:]
			element.Attr[i].Name.Local = s
		}
	}
}

func TestJasperDomTest(t *testing.T) {
	entries, err := os.ReadDir("./tests6")
	if err != nil {
		log.Fatal(err)
	}

	successCount := 0
	failureCount := 0
	for _, e := range entries {
		fmt.Println(e.Name())
		m, err := LoadJasperReportsDomFromFile("tests6/" + e.Name())
		if !assert.NoError(t, err) {
			fmt.Printf("Error '%s': %v\n", e.Name(), err)
			failureCount++
		} else {
			fb := []byte(m.OutputXML(true))
			os.WriteFile("output/"+e.Name(), fb, 0644)
			successCount++
		}
	}
	fmt.Printf("Success: %d Error: %d\n", successCount, failureCount)
}

func TestJasperDomConvertTest(t *testing.T) {
	entries, err := os.ReadDir("./tests6")
	if err != nil {
		log.Fatal(err)
	}

	sample := os.Getenv("JASPER_SAMPLE_TEST")

	successCount := 0
	failureCount := 0
	for _, e := range entries {
		if sample != "" && e.Name() != sample {
			continue
		}
		fmt.Println(e.Name())
		m, err := LoadJasperReportsDomFromFile("tests6/" + e.Name())

		if !assert.NoError(t, err) {
			fmt.Printf("Error '%s': %v\n", e.Name(), err)
			failureCount++
		} else {
			if !assert.NotNil(t, m) {
				failureCount++
				continue
			}
			xmlquery.FindEach(m, "//band", func(i int, band *xmlquery.Node) {
				boxProperty := band.SelectElement("property")
				if boxProperty != nil {
					xmlquery.RemoveFromTree(boxProperty)
				}
				validateAttr(band)
			})
			xmlquery.FindEach(m, "//style", func(i int, style *xmlquery.Node) {
				validateAttr(style)
				for i, attr := range style.Attr {
					if attr.Name.Local == "fontSize" {
						f, _ := strconv.ParseFloat(attr.Value, 64)
						style.Attr[i].Value = fmt.Sprintf("%.1f", f)
					}
					if attr.Name.Local == "isDefault" {
						style.Attr[i].Name.Local = "default"
					}
					if attr.Name.Local == "isBold" {
						style.Attr[i].Name.Local = "bold"
					}
					if attr.Name.Local == "isItalic" {
						style.Attr[i].Name.Local = "italic"
					}
					if attr.Name.Local == "isUnderline" {
						style.Attr[i].Name.Local = "underline"
					}
					if attr.Name.Local == "isStrikeThrough" {
						style.Attr[i].Name.Local = "strikeThrough"
					}
				}
			})
			jasperReprt := m.SelectElement("jasperReport")
			rs := []string{"xmlns", "xsi", "schemaLocation", "name"}
			attr := jasperReprt.Attr
			name := e.Name()
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
			validateAttr(jasperReprt)

			xmlquery.FindEach(m, "//textField", workConvertElements)
			xmlquery.FindEach(m, "//staticText", workConvertElements)
			xmlquery.FindEach(m, "//subreport", workConvertElements)
			xmlquery.FindEach(m, "//line", workConvertElements)
			xmlquery.FindEach(m, "//componentElement", workConvertElements)

			xmlquery.FindEach(m, "//subDataset", func(i int, subDataset *xmlquery.Node) {
				subDataset.Data = "dataset"
			})
			xmlquery.FindEach(m, "//property", func(i int, property *xmlquery.Node) {
				validateAttr(property)
				for _, attr := range property.Attr {
					if attr.Name.Local == "name" {
						switch {
						case attr.Value == "com.jaspersoft.studio.data.defaultdataadapter":
							xmlquery.RemoveFromTree(property)
						case attr.Value == "com.jaspersoft.studio.unit.x":
							xmlquery.RemoveFromTree(property)
						}
					}
				}
			})
			expressionList := []string{"groupExpression", "bucketExpression", "variableExpression", "datasetParameterExpression"}
			for _, e := range expressionList {
				xmlquery.FindEach(m, "//"+e, func(i int, expression *xmlquery.Node) {
					expression.Data = "expression"
				})
			}
			xmlquery.FindEach(m, "//group", func(i int, group *xmlquery.Node) {
				validateAttr(group)
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
			fb := []byte(m.OutputXMLWithOptions(xmlquery.WithIndentation("\t"), xmlquery.WithEmptyTagSupport()))
			// fb := []byte(m.OutputXML(false))
			os.WriteFile("output/"+e.Name(), fb, 0644)
			successCount++
		}
	}
	fmt.Printf("Success: %d Error: %d\n", successCount, failureCount)
}

func removeBandNode(i int, node *xmlquery.Node) {
	validateAttr(node)
	band := node.SelectElement("band")
	node.Attr = band.Attr
	for _, n := range band.ChildNodes() {
		xmlquery.AddImmediateSibling(band, n)
	}
	xmlquery.RemoveFromTree(band)

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
	uuid := reportElement.SelectAttr("uuid")
	n.SetAttr("uuid", uuid)

	for _, a := range reportElement.Attr {
		if a.Name.Local == "uuid" || a.Name.Local == "evaluationTime" {
			continue
		}
		n.SetAttr(a.Name.Local, a.Value)
	}
	for _, a := range saveAttr {
		n.SetAttr(a.Name.Local, a.Value)
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
				n.SetAttr("textAdjust", "StretchHeight")
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
	}
	if t := n.SelectElement("text"); t != nil {
		previous = t
	}
	property := reportElement.SelectElement("property")
	if property != nil {
		xmlquery.AddImmediateSibling(previous, property)
		previous = property
	}
	box := n.SelectElement("box")

	if textElement != nil {
		xmlquery.RemoveFromTree(textElement)
	}
	xmlquery.RemoveFromTree(reportElement)
	if box != nil {
		xmlquery.RemoveFromTree(box)
		xmlquery.AddImmediateSibling(previous, box)
	}
	validateAttr(n)

	// fmt.Printf("Modified Node:\n%#v\n", n)
}
