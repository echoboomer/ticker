package stocks

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var upstreamStockAPI string = "https://www.alphavantage.co/query"

// GetStockData takes a symbol and returns data about the requested
// stock
func GetStockData() (StockDataResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file, using environment variables instead")
	}

	apiKey := os.Getenv("APIKEY")
	symbol := os.Getenv("SYMBOL")

	req, err := http.NewRequest("GET", upstreamStockAPI, nil)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("function", "TIME_SERIES_DAILY")
	q.Add("symbol", symbol)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	resultContainer := StockDataResponse{}
	err = json.NewDecoder(resp.Body).Decode(&resultContainer)
	if err != nil {
		log.Errorf("Error during marshal of response data: %s", err)
	}
	return resultContainer, nil
}

// GetStockDataPeriod takes a symbol and ndays and returns data about the requested
// stock limited to a specific amount of days back
func GetStockDataPeriod() (StockDataResponseScoped, error) {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file, using environment variables instead")
	}

	apiKey := os.Getenv("APIKEY")
	symbol := os.Getenv("SYMBOL")
	ndays := os.Getenv("NDAYS")

	req, err := http.NewRequest("GET", upstreamStockAPI, nil)
	if err != nil {
		log.Errorf("Error: %s", err)
	}
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("function", "TIME_SERIES_DAILY")
	q.Add("symbol", symbol)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	resultContainer := StockDataResponse{}
	err = json.NewDecoder(resp.Body).Decode(&resultContainer)
	if err != nil {
		log.Errorf("Error during marshal of response data: %s", err)
	}

	// Time formatting
	const timeFormat = "2006-01-02"
	currentTime := time.Now()
	filteredResultContainer := make(map[string]StockDataResponseItem, 0)
	var averageClosePrice []float64 = []float64{}
	for key, value := range resultContainer.TimeSeriesDaily {
		dateStringAsTime, _ := time.Parse(timeFormat, key)
		daysBetween := currentTime.Sub(dateStringAsTime).Hours() / 24
		nDaysParsed, _ := strconv.ParseFloat(ndays, 64)
		if math.Round(daysBetween) <= nDaysParsed {
			filteredResultContainer[key] = value
			averageClosePriceParsed, _ := strconv.ParseFloat(value.Close, 64)
			averageClosePrice = append(averageClosePrice, math.Round(averageClosePriceParsed))
		}
	}

	// Calculate average of close values
	averageClosePriceFinal := 0.0
	for i := 0; i < len(averageClosePrice); i++ {
		averageClosePriceFinal += (averageClosePrice[i])
	}
	avg := math.Round((float64(averageClosePriceFinal)) / (float64(len(averageClosePrice))))

	return StockDataResponseScoped{Items: filteredResultContainer, AverageClosePrice: avg}, nil
}
