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
