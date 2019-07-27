package commands

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	defer os.RemoveAll("../../internal/test/out")
	err := os.MkdirAll("../../internal/test/out", 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = Generate("../../internal/test/fixtures/openapi.json", "../../internal/test/out")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadFile("../../internal/test/out/components.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) == 0 {
		t.Fatalf("components.go should not have been empty")
	}
}
