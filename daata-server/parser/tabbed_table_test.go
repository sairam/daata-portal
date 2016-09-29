package parser

import (
	"strings"
	"testing"
)

func TestTabbedTableRow(t *testing.T) {
	a := &TabbedTable{[]byte(TestTabbedData)}
	output, _ := a.Parse()
	// t.Errorf("%v", len(output))
	expected := strings.Join([]string{"siteurl", "Y", "1", "20", "8"}, ",")
	if strings.Join(output[1], ",") != expected {
		t.Errorf("output[1] is %s", output[1])
	}

}

const TestTabbedData = `option_name	ov	option_type	option_width	option_height
siteurl	Y	1	20	8
blogname	Y	1	20	8
blogdescription	Y	1	20	8
wordpress_user_roles	Y	1	20	8
`
