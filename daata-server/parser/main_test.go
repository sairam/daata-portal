package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("------------")
	// TODO setup config
	os.Exit(m.Run())
}
