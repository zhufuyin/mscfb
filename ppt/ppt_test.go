package ppt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	simplePresPath = "../test/demo1.ppt"
)

func TestExtractTextFromPpt(t *testing.T) {
	filePaths := []string{
		//"../test/demo1.ppt",
		//"E:\\code\\github\\java\\poi\\test-data\\slideshow\\37625.ppt",
		//"../test/38256.ppt",
		//"E:\\code\\github\\java\\poi\\test-data\\slideshow\\41071.ppt",
		//"E:\\code\\github\\java\\poi\\test-data\\slideshow\\41246-1.ppt",
		//"E:\\code\\github\\java\\poi\\test-data\\slideshow\\PPT95.ppt",
		//"E:\\code\\github\\java\\poi\\test-data\\slideshow\\bug58516.ppt",
		"E:\\code\\github\\java\\poi\\test-data\\slideshow\\bug58718_008495.ppt",
	}
	for _, filePath := range filePaths {
		fmt.Println("===============================================")
		f, err := os.Open(filePath)
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
		text, err := ppt.ExtractText(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		text = strings.ReplaceAll(text, "\r", "\n")
		//strings.ReplaceAll(text, "\r\r", "\n")
		fmt.Println(text)
	}

	//text, err := ExtractText(f)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(text)
}
