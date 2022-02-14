package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// start printing with color
func setColor(word string, color string) string {

	var result string // the string with the color info included

	colorDefault := "\x1b[0000m"
	colorRed := "\x1b[0031m"
	colorGreen := "\x1b[0032m"
	colorYellow := "\x1b[0033m"
	colorBlue := "\x1b[0034m"
	colorPurple := "\x1b[0035m"
	colorCyan := "\x1b[0036m"
	colorRedBold := "\x1b[0041m"
	colorGreenBold := "\x1b[0042m"
	colorYellowBold := "\x1b[0043m"

	if runtime.GOOS == "windows" {
		colorDefault = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		colorBlue = ""
		colorPurple = ""
		colorCyan = ""
		colorRedBold = ""
		colorGreenBold = ""
		colorYellowBold = ""
	}

	switch {
	case color == "default":
		result = colorDefault + word
	case color == "red":
		result = colorRed + word
	case color == "green":
		result = colorGreen + word
	case color == "yellow":
		result = colorYellow + word
	case color == "purple":
		result = colorPurple + word
	case color == "cyan":
		result = colorCyan + word
	case color == "blue":
		result = colorBlue + word
	case color == "redbold":
		result = colorRedBold + word
	case color == "greenbold":
		result = colorGreenBold + word
	case color == "yellowbold":
		result = colorYellowBold + word
	}
	return result
}

// reduce the newName size
func reduceName(oldName string) string {

	var newName string
	var nameTmp string

	nameLower := strings.ToLower(oldName)
	nameSplitSpace := strings.Split(nameLower, " ")
	if len(nameSplitSpace) > 1 {
		nameTmp = nameSplitSpace[0] + " " + nameSplitSpace[1]
	} else {
		nameTmp = nameSplitSpace[0]
	}
	nameSplitMinus := strings.Split(nameTmp, "-")
	if len(nameSplitMinus) > 1 {
		newName = nameSplitMinus[0] + "-" + nameSplitMinus[1]
	} else {
		newName = nameSplitMinus[0]
	}
	return newName
}

// print the table if custom data are given from the xlsx
func printCustomTable(stockArr []StockInfo, dataXls [][]string) {

	myPrice := 0.0
	myQuantity := 0.0
	myTarget := 0.0
	myGains := 0.0
	myTargetDiff := 0.0
	myChangePercent := 0.0

	totalPrice := 0.0
	totalMyGains := 0.0
	totalMyChangePercent := 0.0
	totalMyPrice := 0.0

	var xlsInfoPresent bool
	var err error
	var symbol string

	dateRequest := time.Now()
	fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// configure the tabwriter to print all in a clean table
	// output, minwidth, tabwidth, padding, padchar, flags
	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer tabw.Flush()

	// define a variable for each column name
	headerSymbol := "Symbol"
	headerShortName := "Name"
	headerFiftyDayAverage := "50DayAvg"
	headerRegularMarketChangePercent := "%DayChange"
	headerRegularMarketPrice := "Price"
	headerRegularMarketTime := "UpdateTime"
	headerTwoHundredDayAverage := "200DayAvg"
	headerMyChangePercent := "%MyChange"
	headerMyGains := "MyGains"
	headerMyPrice := "MyPrice"
	headerMyTarget := "MyTarget"
	headerMyTargetDiff := "MyTargetDiff"
	/*
		headerLongName := "LongName"
		headerRegularMarketChange := "Diff"
		headerFiftyTwoWeekHigh := "52WeekHigh"
		headerFiftyTwoWeekLow := "52WeekLow"
		headerRegularMarketDayHigh := "DayHigh"
		headerRegularMarketDayLow := "DayLow"
		headerRegularMarketPreviousClose := "PrevClose"
	*/

	// set the table header
	header := setColor(headerSymbol, "yellow") +
		" \t " + headerShortName +
		" \t " + headerRegularMarketPrice +
		" \t " + headerMyPrice +
		" \t " + headerMyChangePercent +
		" \t " + headerMyGains +
		" \t " + headerMyTarget +
		" \t " + headerMyTargetDiff +
		" \t " + headerRegularMarketChangePercent +
		" \t " + headerTwoHundredDayAverage +
		" \t " + headerFiftyDayAverage +
		" \t " + headerRegularMarketTime +
		" \n"
	fmt.Fprint(tabw, header)

	for _, val := range stockArr {

		// search for the xls stock info
		xlsInfoPresent = false
		for _, xlsRow := range dataXls {
			// if find the same symbol in the xls and row are not empty
			if strings.EqualFold(strings.ToLower(val.Symbol), strings.ToLower(xlsRow[0])) && len(xlsRow) > 3 {
				xlsInfoPresent = true
				myPrice, err = strconv.ParseFloat(xlsRow[1], 64)
				checkErr(err)
				myQuantity, err = strconv.ParseFloat(xlsRow[2], 64)
				checkErr(err)
				myTarget, err = strconv.ParseFloat(xlsRow[3], 64)
				checkErr(err)
			}
		}

		if xlsInfoPresent {
			// calculation with the xls data
			myChangePercent = ((val.RegularMarketPrice - myPrice) / myPrice) * 100
			myTargetDiff = myTarget - val.RegularMarketPrice
			myGains = (val.RegularMarketPrice - myPrice) * myQuantity

			totalPrice = totalPrice + (val.RegularMarketPrice * myQuantity)
			totalMyPrice = totalMyPrice + (myPrice * myQuantity)
			totalMyGains = totalMyGains + myGains

			symbol = strings.ToLower(val.Symbol)
			switch {
			case myGains < 0.0:
				symbol = setColor(val.Symbol, "red")
			case myGains == 0.0:
				symbol = setColor(val.Symbol, "default")
			case myGains > 0.0:
				symbol = setColor(val.Symbol, "green")
			}
			if myTargetDiff < 0.0 {
				symbol = setColor(val.Symbol, "greenbold")
			}
			if val.RegularMarketChangePercent < -4.0 {
				symbol = setColor(val.Symbol, "redbold")
			}
			if myChangePercent < -10.0 {
				symbol = setColor(val.Symbol, "redbold")
			}

		} else {
			myPrice = 0.0
			myChangePercent = 0.0
			myTargetDiff = 0.0
			myGains = 0.0
			myTarget = 0.0
			myTargetDiff = 0.0

			symbol = strings.ToLower(val.Symbol)
			switch {
			case val.RegularMarketChange < 0:
				symbol = setColor(val.Symbol, "red")
			case val.RegularMarketChange == 0:
				symbol = setColor(val.Symbol, "default")
			case val.RegularMarketChange > 0:
				symbol = setColor(val.Symbol, "green")
			}
		}

		name := reduceName(val.ShortName)
		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		fmt.Fprintf(tabw,
			"%s \t %s \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %v %s\n",
			symbol,
			name,
			val.RegularMarketPrice,
			myPrice,
			myChangePercent,
			myGains,
			myTarget,
			myTargetDiff,
			val.RegularMarketChangePercent,
			val.TwoHundredDayAverage,
			val.FiftyDayAverage,
			dateInfo.Format("02 Jan 15:04"),
			setColor("", "default"),
		)
	}

	// print total row
	totalMyChangePercent = (totalMyGains / totalMyPrice) * 100
	symbol = setColor("-", "yellowbold")
	name := "Total sum"
	fmt.Fprintf(tabw,
		"%s \t %s \t %.2f \t %.2f \t %.2f \t %.2f \t %s \t %s \t %s \t %s \t %s \t %s %s\n",
		symbol,
		name,
		totalPrice,
		totalMyPrice,
		totalMyChangePercent,
		totalMyGains,
		"-",
		"-",
		"-",
		"-",
		"-",
		dateRequest.Format("02 Jan 15:04"),
		setColor("", "default"),
	)

	fmt.Fprint(tabw, setColor("", "default"))
}

// print the standard info
func printDefaultTable(stockArr []StockInfo) {

	dateRequest := time.Now()
	fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// configure the tabwriter to print all in a clean table
	// output, minwidth, tabwidth, padding, padchar, flags
	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer tabw.Flush()

	// define a variable for each column name
	//headerLongName := "LongName"
	// headerRegularMarketChange := "Diff"
	headerSymbol := "Symbol"
	headerShortName := "Name"
	headerFiftyDayAverage := "50DayAvg"
	headerFiftyTwoWeekHigh := "52WeekHigh"
	headerFiftyTwoWeekLow := "52WeekLow"
	headerRegularMarketChangePercent := "% Change"
	headerRegularMarketDayHigh := "DayHigh"
	headerRegularMarketDayLow := "DayLow"
	headerRegularMarketPreviousClose := "PrevClose"
	headerRegularMarketPrice := "Price"
	headerRegularMarketTime := "UpdateTime"
	headerTwoHundredDayAverage := "200DayAvg"

	// set the table header
	header := setColor(headerSymbol, "yellow") +
		" \t " + headerShortName +
		" \t " + headerRegularMarketPrice +
		" \t " + headerRegularMarketPreviousClose +
		// " \t " + headerRegularMarketChange +
		" \t " + headerRegularMarketChangePercent +
		" \t " + headerRegularMarketDayLow +
		" \t " + headerRegularMarketDayHigh +
		" \t " + headerTwoHundredDayAverage +
		" \t " + headerFiftyDayAverage +
		" \t " + headerFiftyTwoWeekLow +
		" \t " + headerFiftyTwoWeekHigh +
		" \t " + headerRegularMarketTime +
		" \n"

	fmt.Fprint(tabw, header)

	for _, val := range stockArr {

		// set color depending of the change (decreasing / stable / increasing)
		symbol := strings.ToLower(val.Symbol)
		switch {
		case val.RegularMarketChange < 0:
			symbol = setColor(val.Symbol, "red")
		case val.RegularMarketChange == 0:
			symbol = setColor(val.Symbol, "default")
		case val.RegularMarketChange > 0:
			symbol = setColor(val.Symbol, "green")
		}
		if val.RegularMarketChangePercent < -4.0 {
			symbol = setColor(val.Symbol, "redbold")
		}

		name := reduceName(val.ShortName)
		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		fmt.Fprintf(tabw,
			"%s \t %s \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %v %s\n",
			symbol,
			name,
			val.RegularMarketPrice,
			val.RegularMarketPreviousClose,
			//val.RegularMarketChange,
			val.RegularMarketChangePercent,
			val.RegularMarketDayLow,
			val.RegularMarketDayHigh,
			val.TwoHundredDayAverage,
			val.FiftyDayAverage,
			val.FiftyTwoWeekLow,
			val.FiftyTwoWeekHigh,
			dateInfo.Format("02 Jan 15:04"),
			setColor("", "default"),
		)
	}
	fmt.Fprint(tabw, setColor("", "default"))
}
