package jaspergo

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestJasperDomConvertTest(t *testing.T) {
	err := ConvertDirectoryToPath("./tests6", "./output")
	assert.NoError(t, err)
}

func TestJasperFileDomConvertTest(t *testing.T) {
	InitLogLevelWithFile("file.log")

	fileName := "./tests6/KeepTogetherReport.20.jrxml"
	sample := os.Getenv("JASPER_SAMPLE_TEST")
	if sample != "" {
		fileName = sample
	}
	err := ConvertFileToPath(fileName, "./output")
	assert.NoError(t, err)
}
