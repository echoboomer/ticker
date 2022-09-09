/*
Copyright © 2022 Scott Hawkins <scott@echoboomer.net>
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package stocks

// StockDataResponse returns successfully feftched data from the upstream API
type StockDataResponse struct {
	MetaData        map[string]string                `json:"Meta Data"`
	TimeSeriesDaily map[string]StockDataResponseItem `json:"Time Series (Daily)"`
}

// StockDataResponseItem describes iterable results from the API's response type
type StockDataResponseItem struct {
	Open   string `json:"1. open" example:"279.0800"`
	High   string `json:"2. high" example:"280.3400"`
	Low    string `json:"3. low" example:"267.9800"`
	Close  string `json:"4. close" example:"268.0900"`
	Volume string `json:"5. volume" example:"27549307"`
}

// StockDataResponseScoped returns items that fit within the NDAYS range along with the
// average close price over the same period
type StockDataResponseScoped struct {
	Items             map[string]StockDataResponseItem `json:"items"`
	AverageClosePrice float64                          `json:"average_close_price"`
}
