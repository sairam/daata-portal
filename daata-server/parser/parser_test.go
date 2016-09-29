package parser

import "testing"

// TODO - check for 1 column data in all cases and 2+ row data

func TestParseASCII(t *testing.T) {
	data, _ := Parse([]byte(TestASCIIData))
	// should be 5 once headers are included
	if len(data) != 4 {
		t.Errorf("Should be 4")
	}
}

func TestParseTabbed(t *testing.T) {
	data, _ := Parse([]byte(TestTabbedData))
	if len(data) != 5 {
		t.Errorf("Should be 5")
	}
}

func TestParseMultiLine(t *testing.T) {
	data, _ := Parse([]byte(TestMultiLineData))
	// should be 3 once headers are included
	if len(data) != 3 {
		t.Errorf("Should be 3")
	}
}

func TestParseNone(t *testing.T) {
	data, _ := Parse([]byte("table\nfdakdfl"))
	if len(data) != 0 {
		t.Errorf("Data should not be parsable")
	}
}
