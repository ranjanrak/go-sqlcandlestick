package sqlcandlestick

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"
)

// Setup mockclient
func setupMock(mockRow *sqlmock.Rows, query string) *Client {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery(query).
		WillReturnRows(mockRow)

	cli := &Client{
		dbClient: db,
	}
	return cli
}

func TestFetchAxisValue(t *testing.T) {
	// Timestamp in time.Time object
	date := time.Date(2022, 5, 18, 0, 0, 0, 0, time.Local)

	// Add mock row for test
	mockedRow := sqlmock.NewRows([]string{"date", "open", "close", "high", "low"}).
		AddRow(date, 156.35, 158.45, 156.75, 157.25).
		AddRow(date.AddDate(0, 0, 1), 159.15, 158.1, 157.2, 156.4)

	queryStatement := `SELECT date, 
					open, 
					close
					max(price) AS high,
					min(price) AS low
					FROM candle_data
					GROUP BY date
					ORDER BY date ASC`

	dbMock := setupMock(mockedRow, queryStatement)

	axisValue, err := dbMock.FetchAxisValue(queryStatement)
	if err != nil {
		log.Fatalf("Error fetching Axis values : %v", err)
	}

	excepectedAxisValue := AxisValues{XAxis: []string{"2022-05-18 00:00:00 +0000 UTC", "2022-05-19 00:00:00 +0000 UTC"},
		YAxis: []opts.KlineData{opts.KlineData{Value: [4]interface{}{156.35, 158.45, 156.75, 157.25}},
			opts.KlineData{Value: [4]interface{}{159.15, 158.1, 157.2, 156.4}}}}

	assert.Equal(t, excepectedAxisValue, axisValue, "Actual Axis values not matching with excepectedAxisValue")

}
