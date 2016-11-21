package redirect

import "testing"

func TestStripPrefix(t *testing.T) {
	str := stripPrefix(RedirectPrefix + "/abc")
	if str != "abc" {
		t.Errorf("StripPrefix is not removing %s", RedirectPrefix)
	}
}

// func ExampleRedirect_second() {
// 	fmt.Println(stripPrefix(RedirectPrefix + "/abc"))
// 	// Output: abc
// }

func TestStripPrefixPlain(t *testing.T) {
	path := RedirectPrefix
	str := stripPrefix(path)
	if str != "" {
		t.Errorf("StripPrefix is not removing %s", RedirectPrefix)
	}
}

func TestStripPrefixWithoutValue(t *testing.T) {
	str := stripPrefix("/abc")
	if str != "abc" {
		t.Errorf("StripPrefix does not working out Prefix(%s)", RedirectPrefix)
	}
}

func TestStripIfRedirectBlank(t *testing.T) {
	str, isRedirect := stripIfRedirect("")
	if str != "" || isRedirect != false {
		t.Errorf("String is blank or Redirect is set for empty path")
	}
}
func TestStripIfRedirect(t *testing.T) {
	str, isRedirect := stripIfRedirect("test")
	if str != "test" || isRedirect != true {
		t.Errorf("String is incorrect or Redirect is set to false")
	}
}

func TestStripIfRedirect2(t *testing.T) {
	str, isRedirect := stripIfRedirect("test+")
	if str != "test" || isRedirect != false {
		t.Errorf("String is incorrect or Redirect is set to true")
	}
}
