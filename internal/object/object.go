package object

import (
	"github.com/fatih/structtag"
	"strings"
)

type Object struct {
	Comments    string
	Tag         string
	Fieldname   string
	PackageName string
}

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
		return strings.Join([]string{firstLetter, rest}, "")
	}

	return yml.Name
}

func (o *Object) GetDescription() string {
	desc := ""
	lines := strings.Split(o.Comments, "\n")
	trimedLines := make([]string, 0)
	for _, line := range lines {
		trimedLines = append(trimedLines, strings.Trim(line, " "))
	}

	for _, line := range trimedLines {
		if strings.HasPrefix(line, "@default") {
			continue
		} else {
			if desc == "" {
				desc = line
			} else {
				desc = strings.Join([]string{desc, line}, " ")
			}
		}
	}

	return desc
}

func (o *Object) GetDefaultValue() string {
	def := ""

	lines := strings.Split(o.Comments, "\n")
	trimedLines := make([]string, 0)
	for _, line := range lines {
		trimedLines = append(trimedLines, strings.Trim(line, " "))
	}

	for _, line := range trimedLines {
		if strings.HasPrefix(line, "@default") {
			return def
		} else {
			continue
		}
	}

	return def
}
