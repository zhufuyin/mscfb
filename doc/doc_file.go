package doc

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhufuyin/mscfb/cfb"
	"github.com/zhufuyin/mscfb/global"
	"io"
	"regexp"
	"strings"
)

const (
	WordDocumentStreamName = "worddocument"
	Table0StreamName       = "0table"
	Table1StreamName       = "1table"
	ObjectPoolStreamName   = "ObjectPool"
)

var (
	newlineReg = regexp.MustCompile(`[\a\r]+`)
	spaceReg   = regexp.MustCompile(`[ \t\r\f\v]+`)
)

type DocFile struct {
	WordDocumentStream io.ReaderAt
	Table0Stream       io.ReaderAt
	Table1Stream       io.ReaderAt
	fib                *Fib
	cft                *ComplexFileTable
}

func NewDocFile(file io.Reader) (*DocFile, error) {
	ra := global.NewReaderAt(file)
	cfbFile, err := cfb.New(ra)
	if err != nil {
		return nil, err
	}
	doc := &DocFile{}
	for _, stream := range cfbFile.File {
		streamName := strings.ToLower(stream.Name)
		if streamName == WordDocumentStreamName && len(stream.Path) == 0 {
			doc.WordDocumentStream = stream
		} else if streamName == Table0StreamName && len(stream.Path) == 0 {
			doc.Table0Stream = stream
		} else if streamName == Table1StreamName && len(stream.Path) == 0 {
			doc.Table1Stream = stream
		}
	}
	if doc.WordDocumentStream == nil || (doc.Table0Stream == nil && doc.Table1Stream == nil) {
		return nil, errors.New("invalid doc file")
	}
	fib, err := NewFib(doc.WordDocumentStream)
	if err != nil {
		return nil, err
	}
	doc.fib = fib
	if fib.fibBase.field_2_nFib < 106 {
		fmt.Printf("The document is too old - Word 95 or older\n")
		return nil, errors.New("the document is too old - Word 95 or older")
	}
	fWhichTblStream := fib.fibBase.isFWhichTblStm()
	var tableStream io.ReaderAt
	if fWhichTblStream {
		tableStream = doc.Table1Stream
		if doc.Table1Stream == nil {
			fmt.Printf("Table Stream %s was not found\n", Table1StreamName)
			return nil, fmt.Errorf("Table Stream %s was not found", Table1StreamName)
		}
	} else {
		tableStream = doc.Table0Stream
		if doc.Table0Stream == nil {
			fmt.Printf("Table Stream %s was not found\n", Table0StreamName)
			return nil, fmt.Errorf("Table Stream %s was not found", Table0StreamName)
		}
	}
	cft, err := NewComplexFileTable(doc.WordDocumentStream, tableStream, int64(fib.fibRgFcLcb97.fcClx))
	if err != nil {
		return nil, err
	}
	doc.cft = cft
	return doc, nil
}

func (doc *DocFile) ExtractText(ctx context.Context) ([]string, error) {
	var texts []string
	if doc == nil || doc.cft == nil || doc.cft.tpt == nil {
		return texts, nil
	}
	//reg := regexp.MustCompile(`[\a\r]+`)
	//placeReg := regexp.MustCompile(`\s+`)
	for _, tp := range doc.cft.tpt.textPieces {
		text := newlineReg.ReplaceAllString(tp.text, "\n")
		text = spaceReg.ReplaceAllString(text, " ")
		texts = append(texts, text)
	}
	return texts, nil
}
