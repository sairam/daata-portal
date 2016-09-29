package parser

import "strings"

type MultiLineSQL struct {
	input []byte
}

// Parse ..
func (t *MultiLineSQL) Parse() ([][]string, error) {
	var cleanInput = strings.TrimSpace(string(t.input))
	var rows = strings.Split(cleanInput, "\n")
	var count = t.identifyCols(rows)

	var rowsCount = (len(rows) / count)
	var data = make([][]string, rowsCount)
	var rowOffset = t.identifyRowOffset(rows[1])
	var headers = t.getHeaders(rows[:count], rowOffset)
	_ = headers

	for i := 0; i < rowsCount; i++ {
		start := i * count
		end := (i + 1) * count
		data[i] = t.getValues(rows[start:end], rowOffset)
	}

	return data, nil
}

func (t *MultiLineSQL) identifyRowOffset(row string) int {
	return strings.Index(row, ": ") + 2
}

func (t *MultiLineSQL) getValues(rows []string, offset int) []string {
	var data = make([]string, len(rows)-1)
	for i, row := range rows[1:] {
		r := strings.TrimSpace(row[offset:])
		data[i] = r
	}
	return data
}

func (t *MultiLineSQL) getHeaders(rows []string, offset int) []string {
	// we are reducing to to remove ": "
	rightOffset := offset - 2
	var data = make([]string, len(rows)-1)
	for i, row := range rows[1:] {
		r := strings.TrimSpace(row[:rightOffset])
		data[i] = r
	}
	return data
}

func (t *MultiLineSQL) identifyCols(rows []string) int {
	rowDelimiter := "***********************"
	for i, row := range rows[1:] {
		if strings.HasPrefix(row, rowDelimiter) {
			return i + 1
		}
	}
	return 0
}
