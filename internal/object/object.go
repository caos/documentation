package object

import (
	"github.com/fatih/structtag"
	"strings"
)

type Object struct {
	Collection  bool
	Comments    string
	Tag         string
	Fieldname   string
	PackageName string
}

const (
	defPrefix = "@default:"
	newLine   = "\n"
	space     = " "
	empty     = ""
)

func (o *Object) GetFieldName() string {
	return o.Fieldname
}
func (o *Object) GetPackageName() string {
	return o.PackageName
}

func (o *Object) GetAttributeName(format string) string {
	tag := strings.TrimSuffix(strings.TrimPrefix(o.Tag, "`"), "`")
	tags, err := structtag.Parse(tag)
	if err != nil {
		return ""
	}

	yml, err := tags.Get(format)
	if err != nil {
		firstLetter := strings.ToLower(string(o.Fieldname[0]))
		rest := o.Fieldname[1:]
		return strings.Join([]string{firstLetter, rest}, empty)
	}

	return yml.Name
}

func (o *Object) IsCollection() bool {
	return o.Collection
}

func (o *Object) GetDescription() string {
	desc := empty
	lines := strings.Split(o.Comments, newLine)
	trimedLines := make([]string, 0)
	for _, line := range lines {
		trimedLines = append(trimedLines, strings.Trim(line, space))
	}

	for _, line := range trimedLines {
		if strings.HasPrefix(line, defPrefix) {
			continue
		} else {
			if desc == empty {
				desc = line
			} else {
				desc = strings.Join([]string{desc, line}, space)
			}
		}
	}

	return desc
}

func (o *Object) GetDefaultValue() string {
	def := empty

	lines := strings.Split(o.Comments, newLine)
	trimedLines := make([]string, 0)
	for _, line := range lines {
		trimedLine := strings.Trim(line, space)
		if trimedLine != empty {
			trimedLines = append(trimedLines, trimedLine)
		}
	}

	for _, line := range trimedLines {
		if strings.HasPrefix(line, defPrefix) {
			def = strings.TrimPrefix(line, defPrefix)
		} else {
			continue
		}
	}
	return def
}
