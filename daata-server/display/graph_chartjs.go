package display

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ChartjsLineGraph ..
type ChartjsLineGraph struct {
	Height int
	Width  int

	Data   []string
	Labels []string
	Title  string
	Name   string
}

// Graph ..
func Graph(w http.ResponseWriter, r *http.Request, filename string) {
	// fmt.Println(filename)
	pathSplit := strings.Split(filename, string(os.PathSeparator))
	name := pathSplit[len(pathSplit)-1]
	bytes, _ := ioutil.ReadFile(filename)

	params := make(map[string]string)
	if r.URL.RawQuery != "" {
		// TODO - fails when no RawQuery is sent
		// fmt.Println(r.URL.RawQuery)
		for _, dint := range strings.Split(r.URL.RawQuery, "&") {
			str := strings.Split(dint, "=")
			i, j := str[0], str[1]
			params[i] = j
		}
	}
	height := 400
	width := 800
	if params["h"] != "" {
		height, _ = strconv.Atoi(params["h"])
		height -= 20
	}
	if params["w"] != "" {
		width, _ = strconv.Atoi(params["w"])
		width -= 20
	}

	noOfEntries := float64(30)
	if count, ok := params["count"]; ok {
		t, ok := strconv.ParseFloat(count, 64)
		if ok == nil && t > 0 {
			noOfEntries = t
		}
	}

	reader := csv.NewReader(strings.NewReader(string(bytes)))
	records, _ := reader.ReadAll()
	noOfRecords := float64(len(records))
	count := int(math.Min(noOfEntries, noOfRecords))

	var graphData = make([]string, count)
	var graphLabel = make([]string, count)

	for i, record := range records[(int(noOfRecords) - count):] {
		if len(record) >= 2 {
			graphData[i] = record[0]
			datetime, _ := strconv.ParseInt(record[1], 10, 64)
			parsedDate := strings.Split(time.Unix(datetime, 0).String(), " ")[1]
			graphLabel[i] = parsedDate // time.Unix(datetime, 0).String()
			_ = parsedDate
		}
	}
	linegraph := &ChartjsLineGraph{
		Data:   graphData,
		Labels: graphLabel,
		Title:  name,
		Name:   name,
		Height: height,
		Width:  width,
	}
	generateChartjsGraph(w, linegraph)
	return
}

func generateChartjsGraph(w http.ResponseWriter, linegraph *ChartjsLineGraph) {
	templatefile := "line_graph_chartjs.tmpl"
	t, err := template.New(templatefile).ParseFiles(config("directory") + "tmpl/" + templatefile)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, linegraph)
}
