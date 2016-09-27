package upload

import "testing"

func TestStripPath(t *testing.T) {
	type tester struct {
		input   string
		result  string
		success bool
	}

	tests := []tester{
		tester{"ab/cd////ef/", "ab/cd/ef", true},
		tester{"ab/cd////ef", "ab/cd/ef", true},
		tester{"ab", "ab", false},
		tester{"", "", false},
		tester{"abcd/ef/gh", "abcd/ef/gh", true},
	}
	for _, test := range tests {
		str, success := getFromPath(test.input)
		if str != test.result || success != test.success {
			t.Errorf("Problem! %s is %s|%s, %v|%v", test.input, test.result, str, test.success, success)
		}
	}
}

func TestExtBasedOnContentType(t *testing.T) {
	type tester struct {
		input    []string
		expected string
	}

	tests := []tester{
		tester{[]string{"application/json"}, "json"},
		tester{[]string{"application/zip"}, "zip"},
		tester{[]string{"text/html"}, "html"},
		tester{[]string{"text/plain"}, "plain"},
	}

	for _, test := range tests {
		output := extBasedOnContentType(test.input)
		if output != test.expected {
			t.Errorf("Expected is %v, Output is %v", test.expected, output)
		}
	}

}

func TestGetAction(t *testing.T) {
	type tester struct {
		input    string
		expected fileFormat
	}

	tests := []tester{
		tester{"json", json},
		tester{"plain", text},
		tester{"html", html},
		tester{"zip", zip},
	}

	for _, test := range tests {
		output := getAction(test.input)
		if output != test.expected {
			t.Errorf("Expected is %v, Output is %v", test.expected, output)
		}
	}

}
