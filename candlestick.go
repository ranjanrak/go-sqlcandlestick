package sqlcandlestick

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// ClientParam represents interface to connect to sqldb
type ClientParam struct {
	DriverName string
	DSN        string
}

// Client represents database driver client
type Client struct {
	dbClient *sql.DB
}

// HttpInput represents input data for rendering candle stick chart
type HttpInput struct {
	Client     Client
	SqlStmt    string
	KlineChart *charts.Kline
}

// AxisValues represents X and Y-axis values
type AxisValues struct {
	XAxis []string
	YAxis []opts.KlineData
}

// Supported database drivers
const (
	Clickhouse = "clickhouse"
	Postgres   = "postgres"
	Mysql      = "mysql"
)

// New creates new data-base connection interface
func New(userParam ClientParam) (*Client, error) {
	// Both are compulsary fields
	if userParam.DriverName == "" {
		return nil, fmt.Errorf("Database driver name missing")
	}

	if userParam.DSN == "" {
		return nil, fmt.Errorf("Data source name missing")
	}

	connect, err := sql.Open(userParam.DriverName, userParam.DSN)
	if err = connect.Ping(); err != nil {
		return nil, err
	}
	return &Client{dbClient: connect}, nil
}

// Fetch rows data using sql query statement to plot X-Y Axis candlestick
func (c *Client) FetchAxisValue(sqlStatement string) (AxisValues, error) {
	rows, err := c.dbClient.Query(sqlStatement)
	if err != nil {
		return AxisValues{}, err
	}

	// Fetch column types
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return AxisValues{}, err
	}

	// Used for allocation & dereferencing
	rowValues := make([]interface{}, len(colTypes))
	for index, column := range colTypes {
		// Create go type suitable for Rows.Scan
		rowValues[index] = reflect.New(column.ScanType()).Interface()
	}

	// User must provide OCLH and time fields to plot candlestick
	if len(colTypes) < 5 {
		return AxisValues{}, fmt.Errorf("All required values not present to plot the candlestick chart")
	}

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)

	defer rows.Close()
	for rows.Next() {
		rowResult := make([]interface{}, len(colTypes))
		for i := 0; i < len(colTypes); i++ {
			rowResult[i] = &rowValues[i]
		}
		// Scan the result
		if err = rows.Scan(rowResult...); err != nil {
			return AxisValues{}, err
		}

		// X axis values
		x = append(x, fmt.Sprint(rowValues[0]))
		// Y axis values
		y = append(y, opts.KlineData{Value: [4]interface{}{rowValues[1], rowValues[2], rowValues[3], rowValues[4]}})

	}
	return AxisValues{XAxis: x, YAxis: y}, nil
}

func (h *HttpInput) httpserver(w http.ResponseWriter, _ *http.Request) {
	// Render default kline chart if no specific kline chart is assigned by the user
	if h.KlineChart != nil {
		h.KlineChart.Render(w)
	} else {
		h.DefaultKlineChart(w)
	}
}

// Default candle stick chart
func (h *HttpInput) DefaultKlineChart(w http.ResponseWriter) {
	// create a new kline instance
	kline := charts.NewKLine()

	// Fetch X and Y axis values
	kd, err := h.Client.FetchAxisValue(h.SqlStmt)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	// Set all candlestick chart variables
	kline.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Candle stick chart", Width: "1400px",
			Height: "700px", Theme: "chalk"}),
		charts.WithTitleOpts(opts.Title{
			Title: "Candle stick chart",
		}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Time", Scale: true, AxisLabel: &opts.AxisLabel{Color: "White"}}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Price", Scale: true, AxisLabel: &opts.AxisLabel{Color: "White"}}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "inside", XAxisIndex: []int{0}, Start: 0, End: 100}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "slider", XAxisIndex: []int{0}, Start: 0, End: 100}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "item", TriggerOn: "mousemove"}),
	)

	kline.SetXAxis(kd.XAxis).AddSeries("Candle", kd.YAxis,
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "green", Color0: "red", BorderColor: "green",
			BorderColor0: "red"}))
	kline.Render(w)
}

// ServeChart serves candlestick chart on the given address
func (c *Client) ServeChart(queryStatement string, addr string, klineChart *charts.Kline) {
	// Set default address
	if addr == "" {
		addr = ":8081"
	}
	httpClient := &HttpInput{Client: *c, SqlStmt: queryStatement, KlineChart: klineChart}
	http.HandleFunc("/", httpClient.httpserver)
	http.ListenAndServe(addr, nil)
}
