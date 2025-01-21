package doc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExtractTextFromDoc(t *testing.T) {
	filePaths := []string{
		//"E:\\code\\lab\\java\\demo\\src\\main\\resources\\47304.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\47950_lower.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\47950_normal.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\47950_upper.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\51921-Word-Crash067.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\52117.doc", // office 95
		//"E:\\code\\github\\java\\poi\\test-data\\document\\52420.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\53379.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\53446.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\56880.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\57603-seven_columns.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\58804.doc",
		//"E:\\code\\github\\java\\poi\\test-data\\document\\59322.doc",
		"E:\\code\\github\\java\\poi\\test-data\\document\\60279.doc",
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
		doc, err := NewDocFile(f)
		if err != nil {
			t.Fatal(err)
		}
		texts, err := doc.ExtractText(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for _, text := range texts {
			text = strings.ReplaceAll(text, "\r", "\n")
			fmt.Println(text)
		}
	}

	//text, err := ExtractText(f)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(text)
}
