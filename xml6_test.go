package jaspergo

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestJasperTest(t *testing.T) {
	entries, err := os.ReadDir("./tests6")
	if err != nil {
		log.Fatal(err)
	}

	successCount := 0
	failureCount := 0
	for _, e := range entries {
		fmt.Println(e.Name())
		m, err := LoadJasperReportsFromFile("tests6/" + e.Name())
		if !assert.NoError(t, err) {
			fmt.Printf("Error '%s': %v\n", e.Name(), err)
			failureCount++
		} else {
			b, err := xml.MarshalIndent(m, "", "\t")
			if !assert.NoError(t, err) {
				return
			}
			fb := []byte(xml.Header + "<!-- Created with Jaspersoft Studio version 6.4.2.qualifier using JasperReports Library version 6.4.2  -->")
			fb = append(fb, b...)
			os.WriteFile("output/"+e.Name(), fb, 0644)
			successCount++
		}
	}
	fmt.Printf("Success: %d Error: %d\n", successCount, failureCount)
}

func TestCheckName(t *testing.T) {
	checkList := []string{"isABC", "abcXisDfff", "isaaa", "isDefIs"}
	var re = regexp.MustCompile(`^is([A-Z])`)
	for _, result := range checkList {
		if t, _ := regexp.MatchString("^is[A-Z]", result); t {
			s := re.ReplaceAllString(result, `$1`)
			s = strings.ToLower(s[0:1]) + s[1:]
			fmt.Println("Found it ", result, "to", s)

		} else {
			fmt.Println("Found it not " + result)
		}
	}
}
