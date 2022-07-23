package main

import (
	"log"

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

	// Simple moving average line chart
	smaLineChart := charts.NewLine()
	smaLineChart.SetGlobalOptions(charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20, GridIndex: 0}), charts.WithYAxisOpts(opts.YAxis{Scale: true, GridIndex: 0}))
	smaLineChart.AddSeries("SMA", generateSMAItems(kd.YAxis),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "Black"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 0, YAxisIndex: 0}))
	kline.Overlap(smaLineChart)

	// Create exponential moving avg line chart
	emaLineChart := charts.NewLine()
	emaLineChart.SetGlobalOptions(charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20, GridIndex: 0}), charts.WithYAxisOpts(opts.YAxis{Scale: true, GridIndex: 0}))
	emaLineChart.AddSeries("EMA", generateEMAItems(kd.YAxis),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "Blue"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 0, YAxisIndex: 0}))
	kline.Overlap(emaLineChart)

	// Serve the modified candlestick with volume chart
	client.ServeChart(queryStatement, "", kline)
}

// Create simple moving average values
func generateSMAItems(kd []opts.KlineData) []opts.LineData {
	items := make([]opts.LineData, 0)
	var sum, avg float32
	for i := 0; i < len(kd); i++ {
		// Use only close value in OCHL to calculate SMA
		sum = sum + kd[i].Value.([4]interface{})[1].(float32)
		avg = sum / float32(i+1)
		items = append(items, opts.LineData{Value: avg})
	}
	return items
}

// Create exponential moving average values
func generateEMAItems(kd []opts.KlineData) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(kd); i++ {
		ochl := kd[i].Value.([4]interface{})
		avg := (ochl[0].(float32) + ochl[1].(float32) + ochl[2].(float32) + ochl[3].(float32)) / 4
		items = append(items, opts.LineData{Value: avg})
	}
	return items
}
