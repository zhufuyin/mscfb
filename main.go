package main

import (
	"context"
	"fmt"
	"github.com/zhufuyin/mscfb/ppt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("input file path need to be specified")
		return
	}
	filePath := args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	ppt, err := ppt.NewPptFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	text, err := ppt.ExtractText(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) > 2 {
		outputFilePath := args[2]
		outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer outputFile.Close()
		if _, err := outputFile.WriteString(text); err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	fmt.Println(text)
}
