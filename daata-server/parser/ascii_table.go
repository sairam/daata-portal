package parser

import (
	"regexp"
	"strings"
)

type asciiCol struct {
	start  int8
	length int8
}

type ASCIITable struct {
	input []byte
}

// Parse ..
func (t *ASCIITable) Parse() ([][]string, error) {
	cleanInput := strings.TrimSpace(string(t.input))
	var rows = strings.Split(cleanInput, "\n")
	var count = int8(len(rows) - 1)
	var cols = t.findColumnOffsets(rows[0])
	headers := t.findData(cols, rows[1])
	_ = headers
	var data = make([][]string, count-3)
	for i, row := range rows[3:count] {
		data[i] = t.findData(cols, row)
	}
	return data, nil
}

func (t *ASCIITable) findData(cols []asciiCol, row string) []string {
	data := []string{}
	for _, col := range cols {
		str := row[col.start:(col.start + col.length)]
		data = append(data, strings.TrimSpace(str))
	}
	return data
}

// finds header split
// start and end column locations in a line per column
func (t *ASCIITable) findColumnOffsets(row string) []asciiCol {
	// TODO - check if we can use strings instead of regexp
	re := regexp.MustCompile(`\+`)
	d := re.FindAllStringSubmatchIndex(row, 1000)
	cols := []asciiCol{}
	var c = len(d) - 1
	for _, s := range d {
		col := asciiCol{start: int8(s[1])}
		cols = append(cols, col)
	}
	for i := 0; i < c; i++ {
		cols[i].length = cols[i+1].start - cols[i].start - 1
	}
	return cols[:c]
}
