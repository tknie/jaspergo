package jaspergo

import (
	"encoding/xml"
	"io"
	"os"
)

type JasperReports struct {
	XMLName                          xml.Name    `xml:"jasperReport"`
	Xmlns                            string      `xml:"xmlns,attr,omitempty"`
	XmlNsXsi                         string      `xml:"xmlns:xsi,attr,omitempty"`
	XsiSchemaLocation                string      `xml:"xsi:schemaLocation,attr,omitempty"`
	WhenNoDataType                   string      `xml:"whenNoDataType,attr,omitempty"`
	WhenResourceMissingType          string      `xml:"whenResourceMissingType,attr,omitempty"`
	ColumnCount                      string      `xml:"columnCount,attr,omitempty"`
	ColumnSpacing                    int         `xml:"columnSpacing,attr,omitempty"`
	Columnwidth                      string      `xml:"columnWidth,attr,omitempty"`
	IsSummaryWithPageHeaderAndFooter *bool       `xml:"isSummaryWithPageHeaderAndFooter,attr,omitempty"`
	IsFloatColumnFooter              *bool       `xml:"isFloatColumnFooter,attr,omitempty"`
	Leftmargin                       string      `xml:"leftMargin,attr,omitempty"`
	Name                             string      `xml:"name,attr,omitempty"`
	PageHeight                       int         `xml:"pageHeight,attr,omitempty"`
	Pagewidth                        int         `xml:"pageWidth,attr,omitempty"`
	Rightmargin                      string      `xml:"rightMargin,attr,omitempty"`
	Topmargin                        string      `xml:"topMargin,attr,omitempty"`
	Bottommargin                     string      `xml:"bottomMargin,attr,omitempty"`
	PrintOrder                       string      `xml:"printOrder,attr,omitempty"`
	Uuid                             string      `xml:"uuid,attr"`
	SortFields                       []sortField `xml:"sortField,omitempty"`
	// more attributes can be added as needed
	Properties     []property    `xml:"property"`
	Styles         []style       `xml:"style"`
	SubDataset     []subDataset  `xml:"subDataset,omitempty"`
	Parameters     []parameter   `xml:"parameter"`
	Field          []field       `xml:"field"`
	Variables      []variable    `xml:"variable"`
	Group          []group       `xml:"group"`
	NoData         []groupHeader `xml:"noData"`
	PageHeader     []groupHeader `xml:"pageHeader"`
	Title          []title       `xml:"title"`
	ColumnHeader   []groupHeader `xml:"columnHeader"`
	Detail         []detail      `xml:"detail"`
	ColumnFooter   []groupHeader `xml:"columnFooter"`
	PageFooter     []groupHeader `xml:"pageFooter"`
	LastPageFooter []groupHeader `xml:"lastPageFooter"`
	Summary        []groupHeader `xml:"summary"`
}

type sortField struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr,omitempty"`
}

type subDataset struct {
	Name      string      `xml:"name,attr"`
	Uuid      string      `xml:"uuid,attr"`
	Field     []field     `xml:"field"`
	Parameter []parameter `xml:"parameter"`
	Variable  []variable  `xml:"variable"`
}

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type style struct {
	Name            string `xml:"name,attr"`
	IsDefault       *bool  `xml:"isDefault,attr,omitempty"`
	FontName        string `xml:"fontName,attr,omitempty"`
	FontSize        string `xml:"fontSize,attr,omitempty"`
	IsBold          *bool  `xml:"isBold,attr,omitempty"`
	IsItalic        *bool  `xml:"isItalic,attr,omitempty"`
	IsStrikeThrough *bool  `xml:"isStrikeThrough,attr,omitempty"`
	IsUnderline     *bool  `xml:"isUnderline,attr,omitempty"`
	PdfFontName     string `xml:"pdfFontName,attr,omitempty"`
	PdfEncoding     string `xml:"pdfEncoding,attr,omitempty"`
	IsPdfEmbedded   bool   `xml:"isPdfEmbedded,attr,omitempty"`
	Pen             []pen  `xml:"pen"`
	Box             []box  `xml:"box"`
}

type variable struct {
	Name               string               `xml:"name,attr"`
	Class              string               `xml:"class,attr"`
	ResetGroup         string               `xml:"resetGroup,attr,omitempty"`
	IncrementGroup     string               `xml:"incrementGroup,attr,omitempty"`
	IncrementType      string               `xml:"incrementType,attr,omitempty"`
	ResetType          string               `xml:"resetType,attr,omitempty"`
	Calculation        string               `xml:"calculation,attr,omitempty"`
	VariableExpression []variableExpression `xml:"variableExpression"`
}

type variableExpression struct {
	XMLName xml.Name `xml:"variableExpression"`
	Value   string   `xml:",cdata"`
}

type group struct {
	Name                        string            `xml:"name,attr"`
	FooterPosition              string            `xml:"footerPosition,attr,omitempty"`
	IsReprintHeaderOnEachColumn *bool             `xml:"isReprintHeaderOnEachColumn,attr,omitempty"`
	IsReprintHeaderOnEachPage   *bool             `xml:"isReprintHeaderOnEachPage,attr,omitempty"`
	IsStartNewColumn            *bool             `xml:"isStartNewColumn,attr,omitempty"`
	MinDetailsToStartFromTop    string            `xml:"minDetailsToStartFromTop,attr,omitempty"`
	MinHeightToStartNewPage     int               `xml:"minHeightToStartNewPage,attr,omitempty"`
	KeepTogether                string            `xml:"keepTogether,attr,omitempty"`
	PreventOrphanFooter         *bool             `xml:"preventOrphanFooter,attr,omitempty"`
	GroupExpression             []groupExpression `xml:"groupExpression"`
	GroupHeader                 []groupHeader     `xml:"groupHeader"`
	GroupFooter                 []groupHeader     `xml:"groupFooter"`
	// more elements can be added as needed
}

type groupExpression struct {
	XMLName xml.Name `xml:"groupExpression"`
	Value   string   `xml:",cdata"`
}

type groupHeader struct {
	Band []band `xml:"band"`
}

type band struct {
	Height              string                `xml:"height,attr"`
	SplitType           string                `xml:"splitType,attr,omitempty"`
	Propertys           []property            `xml:"property"`
	PrintWhenExpression []printWhenExpression `xml:"printWhenExpression,omitempty"`
	Frame               []frame               `xml:"frame,omitempty"`
	Line                []line                `xml:"line,omitempty"`
	Break               []breakTag            `xml:"break,omitempty"`
	TextField           []textField           `xml:"textField,omitempty"`
	StaticText          []staticText          `xml:"staticText"`
	ComponentElement    []componentElement    `xml:"componentElement,omitempty"`
	SubReport           []subReport           `xml:"subreport,omitempty"`
	Image               []image               `xml:"image,omitempty"`
}

type componentElement struct {
	XMLName       xml.Name      `xml:"componentElement,omitempty"`
	ReportElement reportElement `xml:"reportElement"`
	Component     component     `xml:"c:table"`
}

type component struct {
	XMLName           xml.Name `xml:"c:table"`
	XmlnsC            string   `xml:"xmlns:c,attr,omitempty"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr,omitempty"`
}

type image struct {
	XMLName         xml.Name        `xml:"image,omitempty"`
	HAlign          string          `xml:"hAlign,attr,omitempty"`
	ReportElement   reportElement   `xml:"reportElement"`
	Box             []box           `xml:"box,omitempty"`
	ImageExpression imageExpression `xml:"imageExpression"`
}

type imageExpression struct {
	XMLName xml.Name `xml:"imageExpression"`
	Value   string   `xml:",cdata"`
}

type subReport struct {
	XMLName             xml.Name              `xml:"subreport,omitempty"`
	ReportElement       reportElement         `xml:"reportElement"`
	SubreportExpression []subreportExpression `xml:"subreportExpression"`
}

type subreportExpression struct {
	XMLName xml.Name `xml:"subreportExpression"`
	Value   string   `xml:",cdata"`
}

type breakTag struct {
	XMLName       xml.Name        `xml:"break,omitempty"`
	Type          string          `xml:"type,attr,omitempty"`
	ReportElement []reportElement `xml:"reportElement"`
}

type line struct {
	XMLName       xml.Name        `xml:"line,omitempty"`
	ReportElement []reportElement `xml:"reportElement"`
	Box           []box           `xml:"box,omitempty"`
}

type frame struct {
	XMLName       xml.Name        `xml:"frame,omitempty"`
	ReportElement []reportElement `xml:"reportElement"`
	Box           []box           `xml:"box,omitempty"`
	TextField     []textField     `xml:"textField,omitempty"`
}

type printWhenExpression struct {
	XMLName xml.Name `xml:"printWhenExpression"`
	Value   string   `xml:",cdata"`
}

type staticText struct {
	XMLName             xml.Name              `xml:"textField,omitempty"`
	EvaluationTime      string                `xml:"evaluationTime,attr,omitempty"`
	EvaluationGroup     string                `xml:"evaluationGroup,attr,omitempty"`
	TextAlignment       string                `xml:"textAlignment,attr,omitempty"`
	TextAdjust          string                `xml:"textAdjust,attr,omitempty"`
	ReportElement       reportElement         `xml:"reportElement"`
	Box                 []box                 `xml:"box,omitempty"`
	TextElement         []textElement         `xml:"textElement"`
	TextFieldExpression []textFieldExpression `xml:"textFieldExpression"`
}

type textField struct {
	XMLName             xml.Name              `xml:"textField,omitempty"`
	EvaluationTime      string                `xml:"evaluationTime,attr,omitempty"`
	EvaluationGroup     string                `xml:"evaluationGroup,attr,omitempty"`
	TextAlignment       string                `xml:"textAlignment,attr,omitempty"`
	TextAdjust          string                `xml:"textAdjust,attr,omitempty"`
	ReportElement       reportElement         `xml:"reportElement"`
	Box                 []box                 `xml:"box,omitempty"`
	TextElement         []textElement         `xml:"textElement"`
	TextFieldExpression []textFieldExpression `xml:"textFieldExpression"`
}

type textFieldExpression struct {
	XMLName xml.Name `xml:"textFieldExpression"`
	Value   string   `xml:",cdata"`
}

type textElement struct {
	XMLName       xml.Name    `xml:"textElement"`
	Markup        string      `xml:"markup,attr,omitempty"`
	TextAlignment string      `xml:"textAlignment,attr,omitempty"`
	Font          []font      `xml:"font,omitempty"`
	Paragraph     []paragraph `xml:"paragraph,omitempty"`
}

type font struct {
	XMLName  xml.Name `xml:"font"`
	FontName string   `xml:"fontName,attr,omitempty"`
	Size     string   `xml:"size,attr,omitempty"`
}

type paragraph struct {
	XMLName         xml.Name `xml:"paragraph"`
	LineSpacing     string   `xml:"lineSpacing,attr,omitempty"`
	LeftIndent      string   `xml:"leftIndent,attr,omitempty"`
	RightIndent     string   `xml:"rightIndent,attr,omitempty"`
	FirstLineIndent string   `xml:"firstLineIndent,attr,omitempty"`
}

type reportElement struct {
	Mode                  string     `xml:"mode,attr,omitempty"`
	StretchType           string     `xml:"stretchType,attr,omitempty"`
	PositionType          string     `xml:"positionType,attr,omitempty"`
	X                     string     `xml:"x,attr,omitempty"`
	Y                     string     `xml:"y,attr,omitempty"`
	Width                 string     `xml:"width,attr,omitempty"`
	Height                string     `xml:"height,attr,omitempty"`
	IsPrintRepeatedValues *bool      `xml:"isPrintRepeatedValues,attr,omitempty"`
	Forecolor             string     `xml:"forecolor,attr,omitempty"`
	Backcolor             string     `xml:"backcolor,attr,omitempty"`
	Uuid                  string     `xml:"uuid,attr"`
	Propertys             []property `xml:"property"`
}

type title struct {
	Band []band `xml:"band"`
}

type detail struct {
	Band []band `xml:"band"`
}

type pen struct {
	LineColor string `xml:"lineColor,attr,omitempty"`
	LineStyle string `xml:"lineStyle,attr,omitempty"`
	LineWidth string `xml:"lineWidth,attr,omitempty"`
}

type box struct {
	LineWidth     string `xml:"lineWidth,attr,omitempty"`
	TopPadding    string `xml:"topPadding,attr,omitempty"`
	LeftPadding   string `xml:"leftPadding,attr,omitempty"`
	BottomPadding string `xml:"bottomPadding,attr,omitempty"`
	RightPadding  string `xml:"rightPadding,attr,omitempty"`
	Pen           []pen  `xml:"pen"`
	TopPen        []pen  `xml:"topPen"`
	LeftPen       []pen  `xml:"leftPen"`
	BottomPen     []pen  `xml:"bottomPen"`
	RightPen      []pen  `xml:"rightPen"`
}

type parameter struct {
	Name                   string                   `xml:"name,attr"`
	Class                  string                   `xml:"class,attr"`
	DefaultValueExpression []defaultValueExpression `xml:"defaultValueExpression,omitempty"`
}

type defaultValueExpression struct {
	XMLName xml.Name `xml:"defaultValueExpression"`
	Value   string   `xml:",cdata"`
}

type field struct {
	Name      string     `xml:"name,attr"`
	Class     string     `xml:"class,attr"`
	Propertys []property `xml:"property"`
}

func LoadJasperReportsFromFile(filename string) (*JasperReports, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	m := &JasperReports{}
	err = xml.Unmarshal(byteValue, &m)
	if err != nil {
		return nil, err
	}
	m.XmlNsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	m.XsiSchemaLocation = "http://jasperreports.sourceforge.net/jasperreports http://jasperreports.sourceforge.net/xsd/jasperreport.xsd"

	return m, nil
}
