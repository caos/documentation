package code

import (
	"github.com/caos/documentation/pkg/modules"
	"github.com/caos/documentation/pkg/modules/pack"
	"github.com/caos/documentation/pkg/object"
	"github.com/caos/documentation/pkg/treeelement"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

var basicTypes []string = []string{
	"string",
	"bool",
	"int8",
	"uint8",
	"int16",
	"uint16",
	"int32",
	"uint32",
	"int64",
	"uint64",
	"int",
	"uint",
	"uintptr",
	"byte",
	"rune",
}

func GetElementForStruct(packagePath string, structName string) (*treeelement.TreeElement, error) {
	p := modules.CachedModule(packagePath).CachePackage(packagePath)
	return recursiveGetElementForStruct(p, structName, nil)
}

func recursiveGetElementForStruct(p *pack.Package, structName string, obj *object.Object) (*treeelement.TreeElement, error) {
	for _, basic := range basicTypes {
		if basic == structName {
			return nil, nil
		}
	}

	treeElement, cached := p.CachedElements[structName]
	if cached {
		retTreeElement := *treeElement
		return &retTreeElement, nil
		/*retTreeElement := objectToElement(obj, treeElement.GoType)
		if treeElement.Collection {
			retTreeElement.Collection = true
		}
		if treeElement.Map {
			retTreeElement.Map = true
		}
		retTreeElement.SubElements = treeElement.SubElements
		return retTreeElement, nil*/
	}

	goFiles := p.GetGoFileList()

	for _, path := range goFiles {
		treeElement, err := getElementForStructInFile(path, structName, obj)
		if err != nil {
			return nil, err
		}

		if treeElement != nil && !reflect.DeepEqual(*treeElement, treeelement.TreeElement{}) {
			p.CachedElements[structName] = treeElement
			return treeElement, nil
		}
	}
	return nil, nil
}

func getElementForStructInFile(path string, structName string, obj *object.Object) (*treeelement.TreeElement, error) {
	_, file, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	imports := getImports(file)
	element := &treeelement.TreeElement{}
	typeComments := getTypeComments(file)

	ast.Inspect(file, func(x ast.Node) bool {
		if err := NodeToTreeElements(path, imports, structName, element, x, obj, typeComments); err != nil {
			return false
		}
		return true
	})
	return element, nil
}

func getTypeComments(file *ast.File) []string {
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
	return typeComments
}

func objectToElement(obj *object.Object, ty string) *treeelement.TreeElement {
	element := &treeelement.TreeElement{
		GoType:      ty,
		SubElements: make([]*treeelement.TreeElement, 0),
	}
	format := "yaml"

	if obj != nil {
		element.FieldDescription = obj.GetDescription()
		attrName := obj.GetAttributeName(format)
		if attrName != "" {
			element.AttributeName = attrName
		} else {
			element.AttributeName = obj.GetFieldName()
		}
		element.DefaultValue = obj.GetDefaultValue()
		element.GoName = obj.GetFieldName()
		element.GoPackage = obj.GetPackageName()
		element.Collection = obj.IsCollection()
		element.Inline = obj.IsInline(format)
		element.Map = obj.MapType
	}
	return element
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

func getVariableFromField(src []byte, field *ast.Field) (varName string, impName string, typeName string, pointer bool, slice bool, mp bool, mapKeyType string) {
	fieldInSource := string(src[field.Pos()-1 : field.End()-1])

	parts := strings.Split(fieldInSource, " ")
	copyParts := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			copyParts = append(copyParts, part)
		}
	}

	startTrim := false
	parts = make([]string, 0)
	for _, part := range copyParts {
		if strings.HasPrefix(part, "`") {
			startTrim = true
		}
		if !startTrim {
			parts = append(parts, part)
		}
	}

	inlineType := ""
	if len(parts) >= 2 {
		varName = parts[0]
		inlineType = parts[1]
	} else if len(parts) < 2 {
		varName = ""
		inlineType = parts[0]
	}

	if strings.HasPrefix(inlineType, "map") {
		inlineType = strings.TrimPrefix(inlineType, "map")
		mp = true
		typeParts := strings.SplitAfter(inlineType, "]")
		mapKeyType = strings.TrimPrefix(strings.TrimSuffix(typeParts[0], "]"), "[")
		inlineType = strings.TrimPrefix(inlineType, typeParts[0])
	}

	if strings.HasPrefix(inlineType, "[]") {
		inlineType = strings.TrimPrefix(inlineType, "[]")
		slice = true
	}

	if strings.HasPrefix(inlineType, "*") {
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
