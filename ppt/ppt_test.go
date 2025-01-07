package ppt

import (
	"context"
	"fmt"
	"os"
	"testing"
)

const (
	simplePresPath = "../test/demo1.ppt"
)

func TestExtractTextFromPpt(t *testing.T) {
	f, err := os.Open(simplePresPath)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		f.Close()
	})
	ppt, err := NewPptFile(f)
	if err != nil {
		t.Fatal(err)
	}
	text, err := ppt.extractText(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(text)
	//text, err := ExtractText(f)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(text)
}
