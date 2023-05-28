package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// the root of the json is jsonQuoteResponse
type jsonQuoteResponse struct {
	QuoteResponse jsonResult `json:"quoteResponse"`
}

// in 2023 the json root was renamed to finance
// and could only get this :
// {"finance":{"result":null,"error":{"code":"Unauthorized","description":"Invalid Cookie"}}}
type jsonRootFinance struct {
	RootFinance jsonResult `json:"finance"`
}

// inside quoteResponse there is one field result (and also error)
// made of an array of information for each symbol/stock requested
type jsonResult struct {
	Result      []StockInfo `json:"result"`
	ResultError ErrorBlock  `json:"error"`
}

// the error returned by the response
type ErrorBlock struct {
	Code        string `json:"code"`
	Description string `json:"description"`
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

func GetDataFromURL(symbols string, urlTest bool, apiVer int) []byte {

	yahooVer := "v" + strconv.Itoa(apiVer)

	var urlTarget string
	urlBase := "https://query2.finance.yahoo.com/"
	urlFields1 := "/finance/quote?lang=en-US&region=FR&corsDomain=finance.yahoo.com"
	urlFields2 := "symbol,longName,shortName,fiftyDayAverage,fiftyTwoWeekRange,regularMarketChange,regularMarketChangePercent,regularMarketDayRange,regularMarketPreviousClose,regularMarketPrice,regularMarketTime,twoHundredDayAverage"
	urlExternal := urlBase + yahooVer + urlFields1 + "&fields=" + urlFields2 + "&symbols=" + symbols

	log.Printf("[INFO] URL requested : curl --silent \"%s\"\n", urlExternal)

	if urlTest {
		urlTarget = "http://ludev/bourse_resu_error.json"
		//urlTarget = "http://ludev/bourse_resu.json"
		fmt.Printf("\n curl --silent \"%s\"\n", urlExternal)
		fmt.Println("")
		fmt.Println("====> TEST target URL requested:", urlTarget)
		fmt.Println("")
	} else {
		urlTarget = urlExternal
	}

	myClient := http.Client{
		Timeout: time.Second * 4, // timeout 4 sec
	}

	request, err := http.NewRequest(http.MethodGet, urlTarget, nil)
	checkErr(err)
	request.Header.Set("User-Agent", "HTTP Go Client/1.1")

	resu, err := myClient.Do(request)
	checkErr(err)
	defer resu.Body.Close()
	httpBody, err := io.ReadAll(resu.Body)
	checkErr(err)

	log.Println("httpBody returned from web server :", string(httpBody))

	return httpBody
}

func GetDataFromJSON(dataJSON []byte) []StockInfo {

	// check the new json returned, since May 2023
	var dataFinance jsonRootFinance
	err := json.Unmarshal(dataJSON, &dataFinance)
	checkErr(err)

	log.Println("Error Code :", dataFinance.RootFinance.ResultError.Code)
	log.Println("Error Description :", dataFinance.RootFinance.ResultError.Description)

	if dataFinance.RootFinance.ResultError.Code != "" {
		fmt.Println("Returned error from Yahoo : ",
			dataFinance.RootFinance.ResultError.Code,
			"/",
			dataFinance.RootFinance.ResultError.Description)
		fmt.Println("Exit Program")
		os.Exit(2)
	}

	// This is the old response that we used to get before May 2023
	var quoteRes jsonQuoteResponse
	err = json.Unmarshal(dataJSON, &quoteRes)
	checkErr(err)

	log.Println("Struct unmarshal from web server :", quoteRes)

	var stockArr []StockInfo
	for _, val := range quoteRes.QuoteResponse.Result {
		stock := val
		stockArr = append(stockArr, stock)
	}

	return stockArr
}
