package parser

import "strings"

type TabbedTable struct {
	input []byte
}

// Parse ..
// TODO - add header
func (t *TabbedTable) Parse() ([][]string, error) {
	cleanInput := strings.TrimSpace(string(t.input))
	var rows = strings.Split(cleanInput, "\n")
	var data = make([][]string, len(rows))
	for i, row := range rows {
		data[i] = t.SplitByTab(row)
	}
	return data, nil
}

func (t *TabbedTable) SplitByTab(row string) []string {
	return strings.Split(row, "\t")
}
