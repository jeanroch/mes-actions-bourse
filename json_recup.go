package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// the root of the json is jsonQuoteResponse
type jsonQuoteResponse struct {
	QuoteResponse jsonResult `json:"quoteResponse"`
}

// inside quoteResponse there is one field result
// made of an array of information for each symbol/stock requested
type jsonResult struct {
	Result []StockInfo `json:"result"`
}

// all the information from the symbol/stock
type StockInfo struct {
	Symbol                     string  `json:"symbol"`
	LongName                   string  `json:"longName"`
	ShortName                  string  `json:"shortName"`
	FiftyDayAverage            float64 `json:"fiftyDayAverage"`
	FiftyTwoWeekHigh           float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow            float64 `json:"fiftyTwoWeekLow"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketTime          int     `json:"regularMarketTime"`
	TwoHundredDayAverage       float64 `json:"twoHundredDayAverage"`
}

func GetDataFromURL(symbols string, urlTest bool) []byte {

	var urlTarget string
	urlBase := "https://query2.finance.yahoo.com/v7/finance/quote?lang=en-US&region=FR&corsDomain=finance.yahoo.com"
	urlFields := "symbol,longName,shortName,fiftyDayAverage,fiftyTwoWeekRange,regularMarketChange,regularMarketChangePercent,regularMarketDayRange,regularMarketPreviousClose,regularMarketPrice,regularMarketTime,twoHundredDayAverage"
	urlExternal := urlBase + "&fields=" + urlFields + "&symbols=" + symbols

	if urlTest {
		urlTarget = "http://nadev/bourse_resu.json"
		fmt.Printf("\n curl --silent \"%s\"\n", urlExternal)
		fmt.Println("")
		fmt.Println("====> TEST target URL requested:", urlTarget, " <====")
		fmt.Println("")
	} else {
		urlTarget = urlExternal
	}

	myClient := http.Client{
		Timeout: time.Second * 4, // timeout 3 sec
	}

	request, err := http.NewRequest(http.MethodGet, urlTarget, nil)
	checkErr(err)
	request.Header.Set("User-Agent", "HTTP Go Client/1.1")

	resu, err := myClient.Do(request)
	checkErr(err)
	defer resu.Body.Close()
	httpBody, err := ioutil.ReadAll(resu.Body)
	checkErr(err)

	return httpBody
}

func GetDataFromJSON(dataJSON []byte) []StockInfo {

	var quoteRes jsonQuoteResponse
	err := json.Unmarshal(dataJSON, &quoteRes)
	checkErr(err)

	var stockArr []StockInfo
	for _, val := range quoteRes.QuoteResponse.Result {
		stock := val
		stockArr = append(stockArr, stock)
	}
	return stockArr
}
