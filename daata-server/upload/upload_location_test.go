package upload

import "testing"

func TestLocationBasic(t *testing.T) {

	_ = &uploadLocation{"snapdeal/test", "", []string{}, "hello", "zip"}
}
