package parser

import (
	"strings"
	"testing"
)

func TestMultiLineSQL(t *testing.T) {
	a := &MultiLineSQL{[]byte(TestMultiLineData)}
	output, _ := a.Parse()
	// t.Errorf("%v", len(output))
	expected1 := strings.Join([]string{"start_of_week", "N", "1", "12", "18"}, ",")
	expected2 := strings.Join([]string{"wordpress_api_key", "Y", "1", "20", "8"}, ",")
	if strings.Join(output[0], ",") != expected1 || strings.Join(output[1], ",") != expected2 {
		t.Errorf("output[0:1] is %s|%s", output[0], output[1])
	}
}

func TestMultiLineSQLRowOffset(t *testing.T) {
	a := &MultiLineSQL{}
	var cleanInput = strings.TrimSpace(TestMultiLineData)
	var rows = strings.Split(cleanInput, "\n")
	var rowOffset = a.identifyRowOffset(rows[1])
	if rowOffset != 21 {
		t.Errorf("%d is rowOffset | Expected: 21", rowOffset)
	}
}

func TestMultiLineSQLRowsCount(t *testing.T) {
	a := &MultiLineSQL{}
	var rows = strings.Split(TestMultiLineData, "\n")
	var count = a.identifyCols(rows)

	if count != 6 {
		t.Errorf("%d is no of Cols | Expected: 6", a.identifyCols(rows))
	}
	var rowsCount = (len(rows) / count)
	if rowsCount*count != len(rows) {
		t.Errorf("Got %d|%d Expected", rowsCount*count, len(rows))
	}
}

func TestMultiLineSQLHeaders(t *testing.T) {
	a := &MultiLineSQL{}
	var rows = strings.Split(TestMultiLineData, "\n")
	var count = a.identifyCols(rows)
	var rowOffset = a.identifyRowOffset(rows[1])
	var headers = a.getHeaders(rows[:count], rowOffset)

	if headers[1] != "option_can_override" || headers[4] != "option_height" {
		t.Errorf("Headers are incorrect %s", headers[1])
	}
}

const TestMultiLineData = `*************************** 71. row ***************************
        option_name: start_of_week
option_can_override: N
        option_type: 1
       option_width: 12
      option_height: 18
*************************** 8. row ***************************
        option_name: wordpress_api_key
option_can_override: Y
        option_type: 1
       option_width: 20
      option_height: 8
*************************** 72. row ***************************
        option_name: default_role
option_can_override: N
        option_type: 1
       option_width: 20
      option_height: 8`
