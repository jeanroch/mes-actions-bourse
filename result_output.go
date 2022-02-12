package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

func printStockTable(stockArr []StockInfo) {

	colorDefault := "\x1b[0000m"
	colorRed := "\x1b[0031m"
	colorGreen := "\x1b[0032m"
	colorYellow := "\x1b[0033m"
	//colorBlue := "\x1b[0034m"
	//colorPurple := "\x1b[0035m"
	//colorCyan := "\x1b[0036m"

	if runtime.GOOS == "windows" {
		colorDefault = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		//colorBlue = ""
		//colorPurple = ""
		//colorCyan = ""
	}

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
	header := colorYellow +
		headerSymbol +
		//" \t " + headerLongName +
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

		// set color depending of the change (decreasing or stable or increasing)
		symbol := strings.ToLower(val.Symbol)
		switch {
		case val.RegularMarketChange < 0:
			symbol = colorRed + val.Symbol
		case val.RegularMarketChange == 0:
			symbol = colorDefault + val.Symbol
		case val.RegularMarketChange > 0:
			symbol = colorGreen + val.Symbol
		}

		// reduce the name size
		var name string
		var nameTmp string
		nameLower := strings.ToLower(val.ShortName)
		nameSplitSpace := strings.Split(nameLower, " ")
		if len(nameSplitSpace) > 1 {
			nameTmp = nameSplitSpace[0] + " " + nameSplitSpace[1]
		} else {
			nameTmp = nameSplitSpace[0]
		}
		nameSplitMinus := strings.Split(nameTmp, "-")
		if len(nameSplitMinus) > 1 {
			name = nameSplitMinus[0] + "-" + nameSplitMinus[1]
		} else {
			name = nameSplitMinus[0]
		}

		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		fmt.Fprintf(tabw,
			"%s \t %s \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %v \n",
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
		)
	}
	fmt.Fprint(tabw, colorDefault)
}
