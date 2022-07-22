# go-sqlcandlestick

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

client, err := sqlcandlestick.New(sqlcandlestick.ClientParam{
                DriverName: sqlcandlestick.Clickhouse,
                DSN: "tcp://127.0.0.1:9000?debug=true"})
if err != nil {
    log.Fatalf("Error connecting to db: %v", err)
}
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

```

### Default candlestick chart

![image](https://user-images.githubusercontent.com/29432131/180370745-73637dbc-a020-440e-973d-ead2bf5089ec.png)

### Run unit tests

```
go test -v
```
