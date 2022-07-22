package main

import (
	"log"
	"math/rand"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	sqlcandlestick "github.com/ranjanrak/go-sqlcandlestick"
)

func main() {

	client, err := sqlcandlestick.New(sqlcandlestick.ClientParam{
		DriverName: sqlcandlestick.Clickhouse,
		DSN:        "tcp://127.0.0.1:9000?debug=true"})
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	// SQL query statement
	queryStatement := `SELECT date,
                        open,
                        close
                        max(price) AS high,
                        min(price) AS low
                        FROM candle_data
                        GROUP BY date
                        ORDER BY date ASC`

	// create a new kline instance
	kline := charts.NewKLine()

	// Fetch X and Y axis values
	kd, err := client.FetchAxisValue(queryStatement)
	if err != nil {
		log.Fatalf("Error fetching Axis values %v", err)
	}

	kline.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Candle stick chart", Width: "1400px",
			Height: "700px", Theme: "white"}),
		charts.WithTitleOpts(opts.Title{
			Title: "Candle stick with volume movement",
		}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Time", SplitNumber: 20, Scale: true,
			SplitLine: &opts.SplitLine{Show: true, LineStyle: &opts.LineStyle{Color: "Black", Type: "dotted"}},
			AxisLabel: &opts.AxisLabel{Color: "black"}}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Price", Scale: true, AxisLabel: &opts.AxisLabel{Color: "black"}}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "inside", XAxisIndex: []int{0}, Start: 0, End: 100}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "slider", XAxisIndex: []int{0}, Start: 0, End: 100}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "item", TriggerOn: "mousemove"}),
	)

	kline.SetXAxis(kd.XAxis).AddSeries("OCLH", kd.YAxis,
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "green", Color0: "red", BorderColor: "green",
			BorderColor0: "red"}))

	// Create volume bar chart
	volumeBarChart := charts.NewBar()
	volumeBarChart.SetXAxis(kd.XAxis).AddSeries("Volume", generateBarItems(len(kd.XAxis)),
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "Green"}))

	kline.Overlap(volumeBarChart)

	// Serve the modified candlestick with volume chart
	client.ServeChart(queryStatement, "", kline)
}

// Generate random bar items
func generateBarItems(kd int) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < kd; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(2100-2000) + 2000})
	}
	return items
}
