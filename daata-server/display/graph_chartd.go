package display

import (
	"fmt"

	"../utils"
)

// ChartdGraph ..
func ChartdGraph(filename string) string {
	var graphData []float64
	_ = filename
	ymin, ymax := 0, 110
	xmin, xmax := 1, 1
	data := utils.EncodeGraphData(graphData[:], float64(ymin), float64(ymax))
	width, height := 400, 800

	fmt.Printf("https://chartd.co/a.svg?w=%d&h=%d&d0=%s&ymin=%d&ymax=%d&xmin=%d&xmax=%d&tz=%s&t=%s&hl=1&step=1&s0=%s&f0=%s",
		width, height, data, ymin, ymax, xmin, xmax, "Asia/Kolkata", "Title", "FF0000", "00FF0080")

	return `<img src="/a.svg?w=580&h=180&d0=SRWfaZHLHEDABKKTUYgpqqvws0138eZfaYtwxxsxyst">`

}
