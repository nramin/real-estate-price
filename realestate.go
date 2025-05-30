package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nramin/real-estate-price/structs"
	"gopkg.in/yaml.v3"
)

func main() {
	var result structs.RealEstatePrices
	yamlConfig := readYamlFile("properties.yaml", &result)
	properties := yamlConfig.Properties

	if len(properties) <= 0 {
		printError(&result, "No homes listed in properties yaml file.")
		os.Exit(0)
	}

	for _, v := range properties {
		zillowEstimate := getZillowEstimate(v.Zillow, &result)
		if zillowEstimate == 0 {
			printError(&result, "Home value could not be retrieved.")
			os.Exit(0)
		}
		var property structs.PropertyDetails
		property.Address = v.Address
		property.Balance = v.Balance
		property.Price = zillowEstimate
		property.Equity = zillowEstimate - v.Balance
		result.Properties = append(result.Properties, property)

		result.TotalEquity += property.Equity
		result.TotalValue += property.Price
	}

	success := new(bool)
	*success = true
	result.Success = success

	marshaledResult, err := json.Marshal(result)
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func getZillowEstimate(zillow string, result *structs.RealEstatePrices) float64 {
	var estimate float64
	c := colly.NewCollector(
		colly.AllowedDomains(
			"https://www.zillow.com/",
			"zillow.com/",
			"https://zillow.com/",
			"www.zillow.com",
			"zillow.com",
			"https://zillow.com",
		),
	)

	c.OnRequest(func(r *colly.Request) {
		userAgent := "1 Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"
		r.Headers.Set("User-Agent", userAgent)
	})

	c.OnHTML("span[data-testid=price] > span", func(e *colly.HTMLElement) {
		re := regexp.MustCompile("[^0-9]+")
		priceFormatted := re.ReplaceAllString(e.Text, "")
		priceFloat, _ := strconv.ParseFloat(strings.TrimSpace(priceFormatted), 64)
		estimate = priceFloat
	})

	c.OnError(func(r *colly.Response, err error) {
		printError(result, fmt.Sprintf("Request URL: %s, Failed with response: %s, Error: %s", r.Request.URL, string(r.Body), err))
		os.Exit(0)
	})

	if err := c.Visit(zillow); err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	return estimate
}

func readYamlFile(filePath string, result *structs.RealEstatePrices) structs.YamlConfig {
	b, err := os.ReadFile(filePath)
	if err != nil {
		printError(result, "Unable to read input file "+filePath)
		os.Exit(0)
	}
	var yamlConfig structs.YamlConfig

	err = yaml.Unmarshal([]byte(b), &yamlConfig)
	if err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	return yamlConfig
}

func printError(result *structs.RealEstatePrices, error string) {
	success := new(bool)
	*success = false

	result.Success = success
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}
