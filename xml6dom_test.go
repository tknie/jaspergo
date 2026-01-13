package jaspergo

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/go-openapi/testify/v2/assert"
)

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

	// entries := []string{"TextReport.1.jrxml"}

	successCount := 0
	failureCount := 0
	for _, e := range entries {
		// if e.Name() != "TextReport.1.jrxml" {
		// 	continue
		// }
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
			})
			xmlquery.FindEach(m, "//style", func(i int, style *xmlquery.Node) {
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

			xmlquery.FindEach(m, "//textField", workConvertElements)
			xmlquery.FindEach(m, "//staticText", workConvertElements)
			xmlquery.FindEach(m, "//subreport", workConvertElements)

			xmlquery.FindEach(m, "//property", func(i int, property *xmlquery.Node) {
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
			xmlquery.FindEach(m, "//groupExpression", func(i int, groupExpression *xmlquery.Node) {
				groupExpression.Data = "expression"
			})
			xmlquery.FindEach(m, "//pageHeader", removeBandNode)
			xmlquery.FindEach(m, "//title", removeBandNode)
			xmlquery.FindEach(m, "//summary", removeBandNode)
			xmlquery.FindEach(m, "//lastPageFooter", removeBandNode)
			xmlquery.FindEach(m, "//pageFooter", removeBandNode)
			fb := []byte(m.OutputXMLWithOptions(xmlquery.WithIndentation("\t"), xmlquery.WithEmptyTagSupport()))
			// fb := []byte(m.OutputXML(false))
			os.WriteFile("output/"+e.Name(), fb, 0644)
			successCount++
		}
	}
	fmt.Printf("Success: %d Error: %d\n", successCount, failureCount)
}

func removeBandNode(i int, node *xmlquery.Node) {
	band := node.SelectElement("band")
	node.Attr = band.Attr
	for _, n := range band.ChildNodes() {
		xmlquery.AddImmediateSibling(band, n)
	}
	xmlquery.RemoveFromTree(band)

}

func workConvertElements(i int, n *xmlquery.Node) {
	saveName := n.Data
	n.Data = "element"
	n.Attr = []xmlquery.Attr{}
	n.SetAttr("kind", saveName)
	reportElement := n.SelectElement("reportElement")
	uuid := reportElement.SelectAttr("uuid")
	n.SetAttr("uuid", uuid)

	for _, a := range reportElement.Attr {
		if a.Name.Local == "uuid" {
			continue
		}
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

	// fmt.Printf("Modified Node:\n%#v\n", n)
}
