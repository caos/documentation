package docu

import (
	"github.com/caos/documentation/internal/pack"
	"github.com/caos/documentation/internal/treeelement"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Documentation struct {
	tree []*treeelement.TreeElement
}

func New() *Documentation {
	elements := make([]*treeelement.TreeElement, 0)

	return &Documentation{tree: elements}
}
func (d *Documentation) Parse(path string, structNames []string) error {
	elements := make([]*treeelement.TreeElement, 0)

	p := pack.New(path)
	for _, structName := range structNames {
		element, err := p.GetElementForStruct(structName)
		if err != nil {
			return err
		}
		elements = append(elements, element)
	}

	d.tree = elements
	return nil
}

func (d *Documentation) GenerateMarkDown(basePath string) error {
	for _, element := range d.tree {
		if err := generateMarkDownPerElement(basePath, element); err != nil {
			return err
		}
	}
	return nil
}

func generateMarkDownPerElement(basePath string, element *treeelement.TreeElement) error {
	if element == nil {
		return nil
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	data, filePath := element.GetMDFile(basePath)
	if err := ioutil.WriteFile(filePath, data, os.ModePerm); err != nil {
		return err
	}

	if element.SubElements != nil {
		for _, subelement := range element.SubElements {
			if subelement != nil && subelement.SubElements != nil && len(subelement.SubElements) > 0 {
				if err := generateMarkDownPerElement(filepath.Join(basePath, subelement.GoPackage), subelement); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
