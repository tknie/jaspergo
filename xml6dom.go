package jaspergo

import (
	"fmt"
	"io"
	"os"

	"github.com/antchfx/xmlquery"
)

func LoadJasperReportsDomFromFile(path string) (*xmlquery.Node, error) {
	fmt.Println("Loading ...", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := io.Reader(f)
	// Convert UTF-16 XML to UTF-8
	doc, err := xmlquery.Parse(r)
	return doc, err
}
