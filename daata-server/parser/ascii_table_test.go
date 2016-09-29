package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestASCIITableRow(t *testing.T) {
	ASCIIData1 := `+--+---+----+`
	var d []asciiCol
	var o string
	a := &ASCIITable{}
	d = a.findColumnOffsets(ASCIIData1)
	o = fmt.Sprintf("%v", d)
	if o != `[{1 2} {4 3} {8 4}]` {
		t.Errorf("Column Offsets are incorrect")
	}
}

func TestASCIITableHeader(t *testing.T) {
	var rows = strings.Split(string(TestASCIIData), "\n")

	a := &ASCIITable{}
	var cols = a.findColumnOffsets(rows[0])
	headers := a.findData(cols, rows[1])

	o := fmt.Sprintf("%v", headers)
	if o != "[option_name ab option_type option_width option_height]" {
		t.Errorf("headers are \n%v", headers)
	}
}

func TestASCIITableRows(t *testing.T) {
	a := &ASCIITable{[]byte(TestASCIIData)}
	output, _ := a.Parse()
	// t.Errorf("%v", len(output))
	expected := strings.Join([]string{"siteurl", "Y", "1", "20", "8"}, ",")
	if strings.Join(output[0], ",") != expected {
		t.Errorf("output[0] is %s", output[0])
	}
}

const TestASCIIData = `+----------------------+----+-------------+--------------+---------------+
| option_name          | ab | option_type | option_width | option_height |
+----------------------+----+-------------+--------------+---------------+
| siteurl              | Y  |           1 |           20 |             8 |
| blogname             | Y  |           1 |           20 |             8 |
| blogdescription      | Y  |           1 |           20 |             8 |
| wordpress_user_roles | Y  |           1 |           20 |             8 |
+----------------------+----+-------------+--------------+---------------+
`
