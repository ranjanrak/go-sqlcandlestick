# go-sqlcandlestick

[![Run Tests](https://github.com/ranjanrak/go-sqlcandlestick/actions/workflows/go-test.yml/badge.svg)](https://github.com/ranjanrak/go-sqlcandlestick/actions/workflows/go-test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ranjanrak/go-sqlcandlestick.svg)](https://pkg.go.dev/github.com/ranjanrak/go-sqlcandlestick)

Tiny go package for serving candlestick chart for relational database.

## Installation

```
go get -u github.com/ranjanrak/go-sqlcandlestick
```

## Usage

```go
package main
import (
	"log"

	sqlcandlestick "github.com/ranjanrak/go-sqlcandlestick"
)

func main() {
    client, err := sqlcandlestick.New(sqlcandlestick.ClientParam{
                    DriverName: sqlcandlestick.Clickhouse,
                    DSN: "tcp://127.0.0.1:9000?debug=true"})
    if err != nil {
        log.Fatalf("Error connecting to db: %v", err)
    }
    // Prepare sql statement
    queryStatement := `SELECT date,
                        open,
                        close
                        max(price) AS high,
                        min(price) AS low
                        FROM candle_data
                        GROUP BY date
                        ORDER BY date ASC`

    // Serve the candlestick chart
    client.ServeChart(queryStatement, "", nil)
}
```

### Default candlestick chart

![image](https://user-images.githubusercontent.com/29432131/180370745-73637dbc-a020-440e-973d-ead2bf5089ec.png)

## Create your own candlestick pattern

You can create your own candlestick chart types and pass the chart config to `ServeChart(..., chart *charts.Kline)`.
Few example chart config are shown under [examples](https://github.com/ranjanrak/go-sqlcandlestick/tree/main/examples) folder.

1. Candlestick OCLH chart along with volume movement
   ![image](https://user-images.githubusercontent.com/29432131/180378371-8665436f-3bb1-48d5-9dd9-e4b748e89a3d.png)

2. Candlestick OCLH chart with EMA and SMA lines
   ![image](https://user-images.githubusercontent.com/29432131/180597744-36ec1cdc-10d8-4992-923d-f3fd3333d3cc.png)

### Run unit tests

```
go test -v
```
