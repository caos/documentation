package pack

import (
	"github.com/caos/documentation/internal/object"
	"github.com/caos/documentation/internal/treeelement"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	BasePath string
	Files    []string
}

func New(path string) *Package {
	basePath := os.ExpandEnv(path)

	return &Package{
		BasePath: basePath,
	}
}

func (p *Package) GetGoFileList() []string {
	files, err := getFilesInDirectory(p.BasePath)
	if err != nil {
		return nil
	}

	goFiles := make([]string, 0)
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		goFiles = append(goFiles, file)
	}
	return goFiles
}

func parseFile(path string) ([]byte, *ast.File, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo", src, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	return src, file, nil
}

func getImports(file *ast.File) map[string]string {
	imports := make(map[string]string, 0)
	for _, imp := range file.Imports {
		name := filepath.Base(strings.TrimPrefix(strings.TrimSuffix(imp.Path.Value, "\""), "\""))
		if imp.Name != nil && imp.Name.Name != "" {
			name = imp.Name.Name
		}
		imports[name] = strings.TrimPrefix(strings.TrimSuffix(imp.Path.Value, "\""), "\"")
	}
	return imports
}

func getVariableFromField(src []byte, field *ast.Field) (varName string, impName string, typeName string, pointer bool) {
	fieldInSource := string(src[field.Pos()-1 : field.End()-1])

	parts := strings.Split(fieldInSource, " ")
	copyParts := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			copyParts = append(copyParts, part)
		}
	}
	parts = copyParts

	varName = parts[0]
	inlineType := parts[1]
	if strings.HasPrefix(parts[1], "*") {
		inlineType = strings.TrimPrefix(inlineType, "*")
		pointer = true
	}

	typeParts := strings.Split(inlineType, ".")
	if len(typeParts) > 1 {
		impName = typeParts[0]
		typeName = typeParts[1]
	} else {
		typeName = typeParts[0]
	}
	return
}

func (p *Package) GetElementForStruct(structName string) (*treeelement.TreeElement, error) {
	return p.recursiveGetElementForStruct(structName, nil)
}

func (p *Package) recursiveGetElementForStruct(structName string, obj *object.Object) (*treeelement.TreeElement, error) {
	goFiles := p.GetGoFileList()

	for _, path := range goFiles {
		treeElement, err := getElementForStructInFile(path, structName, obj)
		if err != nil {
			return nil, err
		}

		if treeElement != nil {
			return treeElement, nil
		}
	}
	return nil, nil
}

func getElementForStructInFile(path string, structName string, obj *object.Object) (*treeelement.TreeElement, error) {
	src, file, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	imports := getImports(file)

	var element *treeelement.TreeElement

	typeComments := make([]string, 0)
	for _, node := range file.Decls {
		gd, ok := node.(*ast.GenDecl)
		if ok && gd.Doc != nil {
			commentText := make([]string, 0)
			for _, comment := range gd.Doc.List {
				trimed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(comment.Text, "/*"), "*/"), "//"))
				commentText = append(commentText, trimed)
			}
			typeComments = append(typeComments, strings.Join(commentText, " "))
		}
	}

	ast.Inspect(file, func(x ast.Node) bool {
		// when type definition
		t, okt := x.(*ast.TypeSpec)
		if !okt || structName != t.Name.Name {
			return true
		}
		// when struct
		s, oks := t.Type.(*ast.StructType)
		if !oks {
			return true
		}

		element = objectToElement(obj, t.Name.Name)
		if t.Doc != nil && t.Doc.Text() != "" {
			element.TypeDescription = t.Doc.Text()
		}
		prefix := strings.Join([]string{t.Name.Name, ":"}, "")
		for _, comment := range typeComments {
			if strings.HasPrefix(comment, prefix) {
				element.TypeDescription = strings.TrimPrefix(comment, prefix)
			}
		}

		for _, field := range s.Fields.List {
			fieldObj := &object.Object{}
			if field.Doc != nil && field.Doc.Text() != "" {
				fieldObj.Comments = field.Doc.Text()
			}
			if field.Tag != nil && field.Tag.Value != "" {
				fieldObj.Tag = field.Tag.Value
			}

			v, i, t, _ := getVariableFromField(src, field)
			fieldObj.Fieldname = v

			if i != "" {
				importPath := filepath.Join(os.ExpandEnv("$GOPATH"), "src", imports[i])
				fieldObj.PackageName = filepath.Base(importPath)
				subElement, err := New(importPath).recursiveGetElementForStruct(t, fieldObj)
				if err != nil {
					return false
				}
				element.SubElements = append(element.SubElements, subElement)
			} else {
				subElement, err := New(filepath.Dir(path)).recursiveGetElementForStruct(t, fieldObj)
				if err != nil {
					return false
				}
				if subElement != nil {
					element.SubElements = append(element.SubElements, subElement)
				} else {
					element.SubElements = append(element.SubElements, objectToElement(fieldObj, t))
				}
			}
		}
		return false
	})
	return element, nil
}

func objectToElement(obj *object.Object, ty string) *treeelement.TreeElement {
	element := &treeelement.TreeElement{
		GoType:      ty,
		SubElements: make([]*treeelement.TreeElement, 0),
	}

	if obj != nil {
		element.FieldDescription = obj.GetDescription()
		attrName := obj.GetAttributeName("yaml")
		if attrName != "" {
			element.AttributeName = attrName
		} else {
			element.AttributeName = obj.GetFieldName()
		}
		element.DefaultValue = obj.GetDefaultValue()
		element.GoName = obj.GetFieldName()
		element.GoPackage = obj.GetPackageName()
	}
	return element
}

func getFilesInDirectory(dirPath string) ([]string, error) {
	files := make([]string, 0)

	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, filepath.Join(dirPath, info.Name()))
		}
	}

	return files, err
}
