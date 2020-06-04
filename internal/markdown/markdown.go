package markdown

import "strings"

const (
	columnPrefix    = "| "
	columnSuffix    = " "
	tableLineSuffix = " |"
	header1         = "#"
	header2         = "##"
	space           = " "
	newLine         = "\n"
)

type Markdown struct {
	lines []string
}

type TableEntry struct {
	Value string
	Width int
}

func New() *Markdown {
	return &Markdown{make([]string, 0)}
}

func (m *Markdown) Build() []byte {
	return []byte(strings.Join(m.lines, newLine))
}

func (m *Markdown) AddLine(line string) {
	m.lines = append(m.lines, line)
}

func (m *Markdown) AddBlock(line string) {
	m.lines = append(m.lines, line, newLine)
}

func (m *Markdown) AddTableHeader(entries []*TableEntry) {
	headers := make([]string, 0)
	placeholders := make([]string, 0)

	for _, entry := range entries {
		header, placeholder := getHeader(entry.Value, entry.Width)

		headers = append(headers, header)
		placeholders = append(placeholders, placeholder)
	}

	tableHeader := make([]string, 0)
	tableHeader = append(tableHeader, headers...)
	tableHeader = append(tableHeader, tableLineSuffix)
	m.AddLine(strings.Join(tableHeader, ""))

	tablePlaceholder := make([]string, 0)
	tablePlaceholder = append(tablePlaceholder, placeholders...)
	tablePlaceholder = append(tablePlaceholder, tableLineSuffix)
	m.AddLine(strings.Join(tablePlaceholder, ""))
}

func (m *Markdown) AddTableLine(entries []*TableEntry) {
	tableSlice := make([]string, 0)
	for _, entry := range entries {
		tableSlice = append(tableSlice, getColumn(entry.Value, entry.Width))
	}
	tableSlice = append(tableSlice, tableLineSuffix)

	m.AddLine(strings.Join(tableSlice, ""))
}

func (m *Markdown) AddHeader1(text string) {
	m.lines = append(m.lines, strings.Join([]string{header1, text, newLine, newLine}, space))
}
func (m *Markdown) AddHeader2(text string) {
	m.lines = append(m.lines, strings.Join([]string{header2, text, newLine, newLine}, space))
}

func getHeader(title string, length int) (string, string) {
	columnSlice := make([]string, 0)
	for i := 0; i < length; i++ {
		columnSlice = append(columnSlice, "-")
	}
	column := strings.Join(columnSlice, "")

	return getColumn(title, length), getColumn(column, length)
}

func getColumn(str string, length int) string {
	columSlice := make([]string, 0)
	columSlice = append(columSlice, columnPrefix)
	columSlice = append(columSlice, str)
	spaces := length - len(str)
	for i := 0; i < spaces; i++ {
		columSlice = append(columSlice, " ")
	}
	columSlice = append(columSlice, columnSuffix)
	return strings.Join(columSlice, "")
}
