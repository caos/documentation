package code

import (
	"fmt"
	"github.com/caos/documentation/pkg/modules"
	"github.com/caos/documentation/pkg/object"
	"github.com/caos/documentation/pkg/treeelement"
	"go/ast"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NodeToTreeElements(
	filePath string,
	imports map[string]string,
	elementName string,
	element *treeelement.TreeElement,
	x ast.Node,
	obj *object.Object,
	typeComments []string,
) error {
	t, okt := x.(*ast.TypeSpec)
	ident, okIdent := x.(*ast.Ident)
	i, oki := x.(*ast.InterfaceType)
	if (!okt || elementName != t.Name.Name) && !oki && (!okIdent || elementName != ident.Name) {
		return nil
	}
	if okt && t != nil {
		return TypeToTreeElements(filePath, imports, elementName, element, t, obj, typeComments)
	}
	if okIdent && ident != nil {
		if ident.Obj != nil {
			t, okt := ident.Obj.Decl.(*ast.TypeSpec)
			f, okf := ident.Obj.Decl.(*ast.Field)
			if okt && t != nil {
				return TypeToTreeElements(filePath, imports, elementName, element, t, obj, typeComments)
			}
			if okf && f != nil {
				fi, okfi := f.Type.(*ast.Ident)
				if okfi && fi != nil {
					if fi.Obj != nil {
						ft, okft := fi.Obj.Decl.(*ast.TypeSpec)
						if okft && ft != nil {
							return TypeToTreeElements(filePath, imports, elementName, element, ft, obj, typeComments)
						}
					} else {
						element = objectToElement(nil, fi.Name)
						return nil
					}
				}
			}
		}
	}
	if oki && i != nil {
		src, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		iFace := string(src[i.Pos()-1 : i.End()-1])
		fmt.Println(iFace)
		return nil
	}
	return nil
}

func TypeToTreeElements(
	filePath string,
	imports map[string]string,
	elementName string,
	element *treeelement.TreeElement,
	t *ast.TypeSpec,
	obj *object.Object,
	typeComments []string,
) error {
	a, oka := t.Type.(*ast.ArrayType)
	s, oks := t.Type.(*ast.StructType)
	f, okf := t.Type.(*ast.FuncType)
	i, oki := t.Type.(*ast.InterfaceType)
	m, okm := t.Type.(*ast.MapType)
	p, okp := t.Type.(*ast.StarExpr)
	sel, oksel := t.Type.(*ast.SelectorExpr)
	ident, okIdent := t.Type.(*ast.Ident)

	if !oka && !oks && !okf && !oki && !okm && !oksel && !okIdent && !okp {
		return nil
	}

	*element = *objectToElement(obj, t.Name.Name)
	if t.Doc != nil && t.Doc.Text() != "" {
		element.TypeDescription = t.Doc.Text()
	}
	prefix := strings.Join([]string{t.Name.Name, ":"}, "")
	for _, comment := range typeComments {
		if strings.HasPrefix(comment, prefix) {
			element.TypeDescription = strings.TrimPrefix(comment, prefix)
		}
	}
	if okp && p != nil {
		ident, okpIdent := p.X.(*ast.Ident)
		if okpIdent && ident != nil {
			return TypeOfTypeToTreeElement(filePath, imports, elementName, element, ident)
		}
	}

	if oka && a != nil {
		return SliceToTreeElement(filePath, imports, elementName, element, a)
	}
	if oks && s != nil {
		return StructToTreeElements(filePath, imports, element, s)
	}
	if okIdent && ident != nil {
		return TypeOfTypeToTreeElement(filePath, imports, elementName, element, ident)
	}
	if oksel && sel != nil {
		return SelectorToTreeElement(filePath, imports, elementName, element, sel)
	}
	if oki && i != nil {
		//TODO not implemented
	}
	if okm && m != nil {
		//TODO not implemented
	}
	if okf && f != nil {
		//TODO not implemented
	}

	return nil
}

func SelectorToTreeElement(
	filePath string,
	imports map[string]string,
	elementName string,
	element *treeelement.TreeElement,
	sel *ast.SelectorExpr,
) error {
	fieldObj := &object.Object{}

	ident := sel.X.(*ast.Ident)
	fieldObj.Fieldname = elementName
	fieldObj.Collection = true

	importPath := modules.CachedModule(filePath).GetPathForImport(imports[ident.Name])
	fieldObj.PackageName = filepath.Base(importPath)
	subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), sel.Sel.Name, fieldObj)
	if err != nil {
		return err
	}
	if subElement != nil {
		if subElement.SubElements != nil {
			element.Map = subElement.Map
			element.Collection = subElement.Collection
			element.SubElements = append(element.SubElements, subElement.SubElements...)
		}
	}
	return nil
}

func TypeOfTypeToTreeElement(
	filePath string,
	imports map[string]string,
	elementName string,
	element *treeelement.TreeElement,
	ident *ast.Ident,
) error {
	fieldObj := &object.Object{}

	fieldObj.Fieldname = elementName

	importPath := modules.CachedModule(filePath).GetPathForImport(imports[ident.Name])
	fieldObj.PackageName = filepath.Base(importPath)
	subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), ident.Name, fieldObj)
	if err != nil {
		return err
	}
	if subElement != nil {
		element.Collection = subElement.Collection
		element.Map = subElement.Map
		if subElement.SubElements != nil {
			element.SubElements = append(element.SubElements, subElement.SubElements...)
		}
	}
	return nil
}

func VariableToTreeElement() {

}

func SliceToTreeElement(
	filePath string,
	imports map[string]string,
	elementName string,
	element *treeelement.TreeElement,
	a *ast.ArrayType,
) error {
	fieldObj := &object.Object{}
	sel, oksel := a.Elt.(*ast.SelectorExpr)
	i, oki := a.Elt.(*ast.Ident)

	if oksel && sel != nil {
		identX := sel.X.(*ast.Ident)

		fieldObj.Fieldname = elementName
		fieldObj.Collection = true

		importPath := modules.CachedModule(filePath).GetPathForImport(imports[identX.Name])
		fieldObj.PackageName = filepath.Base(importPath)
		subElement, err := recursiveGetElementForStruct(modules.CachedModule(importPath).CachePackage(importPath), sel.Sel.Name, fieldObj)
		if err != nil {
			return err
		}
		if subElement != nil {
			if subElement.SubElements != nil {
				element.Map = subElement.Map
				element.Collection = subElement.Collection
				element.SubElements = append(element.SubElements, subElement.SubElements...)
			}
		}
	}
	if oki && i != nil {
		fieldObj.Fieldname = elementName
		fieldObj.Collection = true

		dir := filepath.Dir(filePath)
		fieldObj.PackageName = filepath.Base(filePath)
		subElement, err := recursiveGetElementForStruct(modules.CachedModule(dir).CachePackage(dir), i.Name, fieldObj)
		if err != nil {
			return err
		}
		if subElement != nil {
			if subElement.SubElements != nil {
				element.Map = subElement.Map
				element.Collection = subElement.Collection
				element.SubElements = append(element.SubElements, subElement.SubElements...)
			}
		} else {
			// basic types
			bt := objectToElement(fieldObj, i.Name)
			element.Collection = bt.Collection
			element.GoType = bt.GoType
		}
	}
	return nil
}

func StructToTreeElements(
	filePath string,
	imports map[string]string,
	element *treeelement.TreeElement,
	struc *ast.StructType,
) error {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	for _, field := range struc.Fields.List {
		if err := StructFieldToTreeElement(filePath, imports, element, field, src); err != nil {
			return err
		}
	}
	return nil
}

func StructFieldToTreeElement(
	filePath string,
	imports map[string]string,
	element *treeelement.TreeElement,
	field *ast.Field,
	src []byte,
) error {
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

	dir := filepath.Dir(filePath)
	if i != "" {
		importPath := modules.CachedModule(filePath).GetPathForImport(imports[i])
		fieldObj.PackageName = filepath.Base(importPath)
		dir = importPath
	}
	subElement, err := recursiveGetElementForStruct(modules.CachedModule(dir).CachePackage(dir), t, fieldObj)
	if err != nil {
		return err
	}

	if subElement != nil {
		if subElement.Inline {
			if subElement.SubElements != nil {
				element.SubElements = append(element.SubElements, subElement.SubElements...)
			}
		} else {
			element.SubElements = append(element.SubElements, subElement)
		}
	} else {
		// basic types
		element.SubElements = append(element.SubElements, objectToElement(fieldObj, t))
	}
	return err
}
