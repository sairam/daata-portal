package parser

// PatternType ..
type PatternType int

const (
	ptnNoMatch PatternType = iota
	ptnHTMLTable
	ptnMultiLineSQL
	ptnASCIITable
	ptnTabbedTable
)

type Parser interface {
	Parse() [][]string
}

// Parse takes in input data and returns a CSV
func Parse(data []byte) ([][]string, error) {
	var p Parser
	_ = p
	switch DetectType(data) {
	case ptnHTMLTable:
		return HTMLTable(data)
	case ptnMultiLineSQL:
		p := &MultiLineSQL{data}
		return p.Parse()
	case ptnASCIITable:
		p := &ASCIITable{data}
		return p.Parse()
	case ptnTabbedTable:
		p := &TabbedTable{data}
		return p.Parse()
	default:
		return [][]string{}, nil
	}
}

func HTMLTable(data []byte) ([][]string, error) {
	return [][]string{}, nil
}
