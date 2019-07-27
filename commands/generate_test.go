package commands

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	defer os.RemoveAll("../test/out")
	os.MkdirAll("../test/out", 0700)
	err := Generate("../test/fixtures/openapi.json", "../test/out")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadFile("../test/out/components.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) == 0 {
		t.Fatalf("components.go should not have been empty")
	}
}
