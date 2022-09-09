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
package main

import (
	"os"

	"github.com/echoboomer/ticker/docs"
	"github.com/echoboomer/ticker/pkg/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// @title ticker API
// @version 1.0
// @description An API that takes a stock SYMBOL and days as NDAYS and returns information about the stock.

// @contact.name Scott Hawkins
// @contact.email scott@echoboomer.net

// @host localhost:8080
// @BasePath /api/v1

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file, using environment variables instead")
	}

	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "ticker API"
	docs.SwaggerInfo.Description = "An API that takes a stock SYMBOL and days as NDAYS and returns information about the stock."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Run
	requiredEnvironmentVariables := []string{"APIKEY", "SYMBOL", "NDAYS"}
	checkEnvironmentVariables(requiredEnvironmentVariables)
	r := api.SetupRouter()
	// Listen 0.0.0.0:8080
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting API: %s", err)
	}
}

// checkEnvironmentVariables verifies whether or not required variables are set
func checkEnvironmentVariables(requiredEnvironmentVariables []string) {
	for _, v := range requiredEnvironmentVariables {
		_, check := os.LookupEnv(v)
		if !check {
			log.Fatalf("%s variable must be set.", v)
		}
	}
}
