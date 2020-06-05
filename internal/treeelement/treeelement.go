package treeelement

import (
	"github.com/caos/documentation/internal/markdown"
	"path/filepath"
	"strings"
)

const (
	fileEnding = "md"
	titleAttr  = "Attribute"
	titleDesc  = "Description"
	titleDef   = "Default"
	titleCol   = "Collection"
	linkPrefix = "[here]("
	linkSuffix = ")"
)

type TreeElement struct {
	AttributeName    string
	FieldDescription string
	TypeDescription  string
	DefaultValue     string
	GoType           string
	GoName           string
	GoPackage        string
	Collection       bool
	SubElements      []*TreeElement
}

type TreeElementLine struct {
	AttributeName    string
	FieldDescription string
	DefaultValue     string
	Collection       string
}

func (t *TreeElement) GetMDFile(basePath string) ([]byte, string) {
	md := markdown.New()

	md.AddHeader1(t.GoType)

	anLength := len(titleAttr)
	fdLength := len(titleDesc)
	dvLength := len(titleDef)
	coLength := len(titleCol)
	for _, subelement := range t.SubElements {
		if subelement == nil {
			continue
		}

		treeline := subelement.GetLine()

		if len(treeline.AttributeName) > anLength {
			anLength = len(treeline.AttributeName)
		}
		if len(treeline.FieldDescription) > fdLength {
			fdLength = len(treeline.FieldDescription)
		}
		if len(treeline.DefaultValue) > dvLength {
			dvLength = len(treeline.DefaultValue)
		}
		if len(treeline.Collection) > coLength {
			coLength = len(treeline.Collection)
		}
	}

	if t.TypeDescription != "" {
		md.AddBlock(t.TypeDescription)
	}

	md.AddHeader2("Structure")

	headerEntries := []*markdown.TableEntry{
		{titleAttr, anLength},
		{titleDesc, fdLength},
		{titleDef, dvLength},
		{titleCol, coLength},
	}
	md.AddTableHeader(headerEntries)

	for _, subelement := range t.SubElements {
		if subelement == nil {
			continue
		}
		treeline := subelement.GetLine()

		entries := []*markdown.TableEntry{
			{treeline.AttributeName, anLength},
			{treeline.FieldDescription, fdLength},
			{treeline.DefaultValue, dvLength},
			{treeline.Collection, coLength},
		}
		md.AddTableLine(entries)
	}

	return md.Build(), filepath.Join(basePath, strings.Join([]string{t.GoType, fileEnding}, "."))
}

func (t *TreeElement) GetLine() *TreeElementLine {
	fieldDesc := t.FieldDescription
	if t.SubElements != nil && len(t.SubElements) > 0 {
		linkPath := filepath.Join(t.GoPackage, strings.Join([]string{t.GoType, fileEnding}, "."))

		if fieldDesc != "" {
			fieldDesc = strings.Join([]string{fieldDesc, ", ", linkPrefix, linkPath, linkSuffix}, "")
		} else {
			fieldDesc = strings.Join([]string{linkPrefix, linkPath, linkSuffix}, "")
		}
	}

	col := ""
	if t.Collection {
		col = "X"
	}

	return &TreeElementLine{
		AttributeName:    t.AttributeName,
		FieldDescription: fieldDesc,
		DefaultValue:     t.DefaultValue,
		Collection:       col,
	}
}
