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
		retTreeElement := objectToElement(obj, treeElement.GoType)
		retTreeElement.SubElements = treeElement.SubElements
		return retTreeElement, nil
	}

	goFiles := p.GetGoFileList()

	for _, path := range goFiles {
		treeElement, err := getElementForStructInFile(path, structName, obj)
		if err != nil {
			return nil, err
		}

		if treeElement != nil {
			p.CachedElements[structName] = treeElement
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
		a, oka := t.Type.(*ast.ArrayType)
		sel, oksel := t.Type.(*ast.SelectorExpr)
		if !oks && !oka && !oksel {
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

		// if type of a type
		if sel != nil {
			fieldObj := &object.Object{}
			identX := sel.X.(*ast.Ident)

			fieldObj.Fieldname = structName

			importPath := modules.CachedModule(path).GetPathForImport(imports[identX.Name])
			fieldObj.PackageName = filepath.Base(importPath)
			subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), sel.Sel.Name, fieldObj)
			if err != nil {
				return false
			}
			if subElement != nil {
				if subElement.SubElements != nil {
					element.SubElements = append(element.SubElements, subElement.SubElements...)
				}
			}
			return true
		}

		// if array type
		if a != nil {
			fieldObj := &object.Object{}
			sel, oksel := a.Elt.(*ast.SelectorExpr)
			i, oki := a.Elt.(*ast.Ident)

			if oksel {
				identX := sel.X.(*ast.Ident)

				fieldObj.Fieldname = structName
				fieldObj.Collection = true

				importPath := modules.CachedModule(path).GetPathForImport(imports[identX.Name])
				fieldObj.PackageName = filepath.Base(importPath)
				subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), sel.Sel.Name, fieldObj)
				if err != nil {
					return false
				}
				if subElement != nil {
					if subElement.SubElements != nil {
						element.SubElements = append(element.SubElements, subElement.SubElements...)
					}
				}
				return true
			}

			if oki {
				ty := i.Obj.Decl.(*ast.TypeSpec)
				strc := ty.Type.(*ast.StructType)

				for _, field := range strc.Fields.List {
					fieldObj := &object.Object{}
					if field.Doc != nil && field.Doc.Text() != "" {
						fieldObj.Comments = field.Doc.Text()
					}
					if field.Tag != nil && field.Tag.Value != "" {
						fieldObj.Tag = field.Tag.Value
					}

					v, i, t, _, c, m, mkey := getVariableFromField(src, field)
					fieldObj.Fieldname = v
					fieldObj.Collection = c
					fieldObj.Mapkey = mkey
					fieldObj.MapType = m

					if i != "" {
						importPath := modules.CachedModule(path).GetPathForImport(imports[i])
						fieldObj.PackageName = filepath.Base(importPath)
						subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), t, fieldObj)
						if err != nil {
							return false
						}
						if subElement != nil {
							if subElement.Inline {
								if subElement.SubElements != nil {
									element.SubElements = append(element.SubElements, subElement.SubElements...)
								}
							} else {
								element.SubElements = append(element.SubElements, subElement)
							}
						}
					} else {
						dir := filepath.Dir(path)
						subElement, err := recursiveGetElementForStruct(modules.CachedModule(dir).CachePackage(dir), t, fieldObj)
						if err != nil {
							return false
						}
						if subElement != nil {
							// another struct type
							element.SubElements = append(element.SubElements, subElement)
						} else {
							// basic types
							element.SubElements = append(element.SubElements, objectToElement(fieldObj, t))
						}
					}
				}
				return true
			}
		}

		//if struct type
		for _, field := range s.Fields.List {
			fieldObj := &object.Object{}
			if field.Doc != nil && field.Doc.Text() != "" {
				fieldObj.Comments = field.Doc.Text()
			}
			if field.Tag != nil && field.Tag.Value != "" {
				fieldObj.Tag = field.Tag.Value
			}

			v, i, t, _, c, m, mkey := getVariableFromField(src, field)
			fieldObj.Fieldname = v
			fieldObj.Collection = c
			fieldObj.Mapkey = mkey
			fieldObj.MapType = m

			if i != "" {
				importPath := modules.CachedModule(path).GetPathForImport(imports[i])
				fieldObj.PackageName = filepath.Base(importPath)
				subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), t, fieldObj)
				if err != nil {
					return false
				}
				if subElement != nil {
					if subElement.Inline {
						if subElement.SubElements != nil {
							element.SubElements = append(element.SubElements, subElement.SubElements...)
						}
					} else {
						element.SubElements = append(element.SubElements, subElement)
					}
				}
			} else {
				dir := filepath.Dir(path)
				subElement, err := recursiveGetElementForStruct(modules.CachedModule(dir).CachePackage(dir), t, fieldObj)
				if err != nil {
					return false
				}
				if subElement != nil {
					// another struct type
					element.SubElements = append(element.SubElements, subElement)
				} else {
					// basic types
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
	parts = copyParts

	varName = parts[0]
	inlineType := parts[1]
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
