package parser

import "testing"

func TestUnitASCIITable(t *testing.T) {
	a := &strMiscSig{"ascii", ptnASCIITable}
	x := a.match([]byte(TestASCIIData), 0)

	if x != ptnASCIITable {
		t.Errorf("Unit Test: Tabbed table correctly")
	}
}

func TestDetectASCIITable(t *testing.T) {
	if DetectType([]byte(TestASCIIData)) != ptnASCIITable {
		t.Errorf("Could not detect ASCII table correctly")
	}
}

func TestUnitTabbedTable(t *testing.T) {
	a := &strMiscSig{"tabbed", ptnTabbedTable}
	x := a.match([]byte(TestTabbedData), 0)
	if x != ptnTabbedTable {
		t.Errorf("Unit Test: Tabbed table correctly")
	}
}

func TestDetectTabbedTable(t *testing.T) {
	if DetectType([]byte(TestTabbedData)) != ptnTabbedTable {
		t.Errorf("Could not detect Tabbed table correctly")
	}
}

func TestDetectHTMLTable(t *testing.T) {
	if DetectType([]byte("<TABLE>")) != ptnHTMLTable {
		t.Errorf("table tag not detected")
	}
	if DetectType([]byte("<table class='abc'>")) != ptnHTMLTable {
		t.Errorf("table tag not detected")
	}
}

func TestDetectMultiLineSQL(t *testing.T) {
	// input := "*************************** 78. row ***************************"
	if DetectType([]byte(TestMultiLineData)) != ptnMultiLineSQL {
		t.Errorf("multi line sql not detected")
	}
}

func TestDetectNone(t *testing.T) {
	if DetectType([]byte("TABLE")) != ptnNoMatch {
		t.Errorf("no match")
	}
}
