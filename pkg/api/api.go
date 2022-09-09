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
package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/echoboomer/ticker/pkg/stocks"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// getHealth godoc
// @Summary Return health status if the application is running.
// @Description Return health status if the application is running.
// @Tags health
// @Produce json
// @Success 200 {object} JSONHealthResponse
// @Router /health [get]
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, JSONHealthResponse{
		Status: "healthy",
	})
}

// getStock godoc
// @Summary Retrieve information about a given stock symbol parsed as env SYMBOL with full response data.
// @Description Retrieve information about a given stock symbol parsed as env SYMBOL with full response data.
// @Tags getStock
// @Accept json
// @Produce json
// @Success 200 {object} map[string]stocks.StockDataResponseItem
// @Failure 400 {object} JSONFailureResponse
// @Router /stock [get]
func getStock(c *gin.Context) {
	result, err := stocks.GetStockData()
	if err != nil {
		c.JSON(400, JSONFailureResponse{
			Result:  "failure",
			Message: fmt.Sprintf("%s", err),
		})
	}
	c.JSON(200, result.TimeSeriesDaily)
}

// getStockAvg godoc
// @Summary Retrieve information about a given stock symbol parsed as env SYMBOL over days NDAYS.
// @Description Retrieve information about a given stock symbol parsed as env SYMBOL over days NDAYS. Returns a list of matched items and average close price.
// @Tags getStockAvg
// @Accept json
// @Produce json
// @Success 200 {object} stocks.StockDataResponseScoped
// @Failure 400 {object} JSONFailureResponse
// @Router /stock/avg [get]
func getStockAvg(c *gin.Context) {
	result, err := stocks.GetStockDataPeriod()
	if err != nil {
		c.JSON(400, JSONFailureResponse{
			Result:  "failure",
			Message: fmt.Sprintf("%s", err),
		})
	}
	c.JSON(200, result)
}

// SetupRouter instantiates the gin handler instance
func SetupRouter() *gin.Engine {
	// Release mode in production
	// Omit when developing for debug
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.ForwardedByClientIP = true

	// Rate limiting middleware
	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		log.Fatalf("Error starting rate limit middleware: %s", err)
	}
	rateLimitStore := memory.NewStore()
	rateLimitMiddleware := mgin.NewMiddleware(limiter.New(rateLimitStore, rate))
	r.Use(rateLimitMiddleware)

	// Recovery middleware
	r.Use(gin.Recovery())

	// Define api/v1 group
	v1 := r.Group("api/v1")
	{
		v1.GET("/stock", getStock)
		v1.GET("/stock/avg", getStockAvg)
		v1.GET("/health", getHealth)
	}

	// Render HTML
	r.LoadHTMLGlob("static/*")

	// Home
	r.GET("/", func(c *gin.Context) {
		// Render template(s)
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "ticker",
			},
		)
	})

	// swagger-ui
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
