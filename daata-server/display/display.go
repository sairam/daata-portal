package display

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../config"
)

/*
 based on format
 1. Regular format is to display the data as-is
 2. Data Points format or JSON data format has key / value injected into a file
 3. Display appended file from multiple hosts
 4. Graphs / Dashboards
*/

// DisplayPrefix is a
const DisplayPrefix = "/d"

type ChartjsLineGraph struct {
	Data   []string
	Labels []string
	Title  string
}

// DisplayGraph ..
func DisplayGraph(w http.ResponseWriter, r *http.Request) {
	filename := config.DataDirectory + "code-coverage/mycode/coverage.datapoint"
	bytes, _ := ioutil.ReadFile(filename)

	params := map[string]string{}
	for _, dint := range strings.Split(r.URL.RawQuery, "&") {
		str := strings.Split(dint, "=")
		i, j := str[0], str[1]
		params[i] = j
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
	generateChartjsGraph(w, graphData, graphLabel)
	return
}

// _ = filename
// ymin, ymax := 0, 110
// xmin, xmax := 1, 1
// data := utils.EncodeGraphData(graphData[:], float64(ymin), float64(ymax))
// fmt.Fprintf(w,
// 	"https://chartd.co/a.svg?w=%d&h=%d&d0=%s&ymin=%d&ymax=%d&xmin=%d&xmax=%d&tz=%s&t=%s&hl=1&step=1&s0=%s&f0=%s",
// 	width, height, data, ymin, ymax, xmin, xmax, "Asia/Kolkata", "Title", "FF0000", "00FF0080")

// <img src="/a.svg?w=580&h=180&d0=SRWfaZHLHEDABKKTUYgpqqvws0138eZfaYtwxxsxyst">

func generateChartjsGraph(w http.ResponseWriter, graphData, graphLabel []string) {
	templatefile := "line_graph_chartjs.tmpl"
	t, err := template.New(templatefile).ParseFiles(config.DisplayDirectory + "tmpl/" + templatefile)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, ChartjsLineGraph{
		Data:   graphData,
		Labels: graphLabel,
		Title:  "coverage",
	})
}

func openDir(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/d/test" {
		DisplayGraph(w, r)
		return
	}
	// fmt.Println(r.URL.Path)
	if p := strings.TrimPrefix(r.URL.Path, DisplayPrefix); len(p) < len(r.URL.Path) {
		r.URL.Path = p
		// check auth here
		// TODO - fix directory here from config
		http.FileServer(http.Dir(config.DataDirectory)).ServeHTTP(w, r)
	} else {
		// fmt.Println(p)
		http.NotFound(w, r)
	}
}

// Prefix specifies the download/display location of a file
func Prefix() string {
	return DisplayPrefix + "/"
}

func init() {
	http.HandleFunc(Prefix(), openDir)
}
