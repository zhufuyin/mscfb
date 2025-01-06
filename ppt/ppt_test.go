package ppt

import (
	"fmt"
	"os"
	"testing"
)

const (
	simplePresPath = "../resources/demo1.ppt"
)

func TestExtractTextFromPpt(t *testing.T) {
	f, err := os.Open(simplePresPath)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	text, err := ExtractText(f)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(text)
}
