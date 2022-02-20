package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	fr := message.NewPrinter(language.French)

	myPrice := 0.0
	myQuantity := 0.0
	myGains := 0.0
	myChangePercent := 0.0
	myTargetSell := 0.0
	myTargetSellDiff := 0.0
	myTargetBuy := 0.0
	myTargetBuyDiff := 0.0

	totalPrice := 0.0
	totalMyGains := 0.0
	totalMyPrice := 0.0
	totalMyChangePercent := 0.0

	var xlsInfoPresent bool
	var err error
	var symbol string

	dateRequest := time.Now()
	fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// configure the tabwriter to print all in a clean table
	// output, minwidth, tabwidth, padding, padchar, flags
	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	defer tabw.Flush()

	// define a variable for each column name
	headerSymbol := "Symbol"
	headerShortName := "Name"
	headerFiftyDayAverage := "50DayAvg"
	headerRegularMarketChangePercent := "DayChange"
	headerRegularMarketPrice := "Price"
	headerRegularMarketTime := "UpdateTime"
	headerTwoHundredDayAverage := "200DayAvg"
	headerMyChangePercent := "MyChange"
	headerMyGains := "MyGains"
	headerMyPrice := "MyPrice"
	//headerMyTarget := "MyTarget"
	headerMyTargetSellDiff := "MyTargetSell"
	headerMyTargetBuyDiff := "MyTargetBuy"
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
		//" \t " + headerMyTarget +
		" \t " + headerMyTargetSellDiff +
		" \t " + headerMyTargetBuyDiff +
		" \t " + headerRegularMarketChangePercent +
		" \t " + headerTwoHundredDayAverage +
		" \t " + headerFiftyDayAverage +
		" \t " + headerRegularMarketTime +
		" \n"
	fmt.Fprint(tabw, header)

	for _, val := range stockArr {

		// reset all the stock values
		myPrice = 0.0
		myQuantity = 0.0
		myGains = 0.0
		myChangePercent = 0.0
		myTargetSell = 0.0
		myTargetSellDiff = 0.0
		myTargetBuy = 0.0
		myTargetBuyDiff = 0.0

		// search for the xls stock info
		xlsInfoPresent = false
		for _, xlsRow := range dataXls {
			// if find the same symbol in the xls and row are not empty
			if strings.EqualFold(strings.ToLower(val.Symbol), strings.ToLower(xlsRow[0])) {
				xlsInfoPresent = true
				log.Println(val.Symbol, ": symbol is present in xlsx, row size is:", len(xlsRow))
				if len(xlsRow) == 1 {
					log.Println(val.Symbol, ": symbol is alone, all myVariables are set to zero")
					myPrice = 0.0
					myQuantity = 0.0
					myGains = 0.0
					myChangePercent = 0.0
					myTargetSell = 0.0
					myTargetSellDiff = 0.0
					myTargetBuy = 0.0
					myTargetBuyDiff = 0.0
				}
				if len(xlsRow) > 1 {
					if xlsRow[1] != "" {
						myPrice, err = strconv.ParseFloat(xlsRow[1], 64)
						log.Println(val.Symbol, ": myPrice set to:", myPrice)
						checkErr(err)
					} else {
						myPrice = 0.0
						log.Println(val.Symbol, ": myPrice is set to zero, the value of xlsRow[1] is:", xlsRow[1])
					}
				}
				if len(xlsRow) > 2 {
					if xlsRow[2] != "" {
						myQuantity, err = strconv.ParseFloat(xlsRow[2], 64)
						log.Println(val.Symbol, ": myQuantity set to:", myQuantity)
						checkErr(err)
					} else {
						myQuantity = 0.0
						log.Println(val.Symbol, ": myQuantity is set to zero, the value of xlsRow[2] is:", xlsRow[2])
					}
				}
				if len(xlsRow) > 3 {
					if xlsRow[3] != "" {
						myTargetSell, err = strconv.ParseFloat(xlsRow[3], 64)
						log.Println(val.Symbol, ": myTargetSell set to:", myTargetSell)
						checkErr(err)
					} else {
						myTargetSell = 0.0
						log.Println(val.Symbol, ": myTargetSell is set to zero, the value of xlsRow[3] is:", xlsRow[3])
					}
				}
				if len(xlsRow) > 4 {
					if xlsRow[4] != "" {
						myTargetBuy, err = strconv.ParseFloat(xlsRow[4], 64)
						log.Println(val.Symbol, ": myTargetBuy set to:", myTargetBuy)
						checkErr(err)
					} else {
						myTargetBuy = 0.0
						log.Println(val.Symbol, ": myTargetBuy is set to zero, the value of xlsRow[4] is:", xlsRow[4])
					}
				}
				break
			}
		}

		if xlsInfoPresent {

			if myPrice > 0 && myQuantity > 0 {
				log.Println(val.Symbol, ": symbol, price and quantity are present in the xlsx, doing calculations")
				myChangePercent = ((val.RegularMarketPrice - myPrice) / myPrice) * 100
				myGains = (val.RegularMarketPrice - myPrice) * myQuantity

				totalPrice = totalPrice + (val.RegularMarketPrice * myQuantity)
				totalMyPrice = totalMyPrice + (myPrice * myQuantity)
				totalMyGains = totalMyGains + myGains

				if myTargetSell > 0 {
					myTargetSellDiff = ((val.RegularMarketPrice - myTargetSell) / val.RegularMarketPrice) * 100
				}

				// set the row color when myPrice and myQuantity are set
				switch {
				case myGains < 0.0:
					symbol = setColor(val.Symbol, "red")
					log.Println(val.Symbol, ": myGains are negative, line color set to red")
				case myGains == 0.0:
					symbol = setColor(val.Symbol, "default")
					log.Println(val.Symbol, ": myGains are equal to 0, line color set to default")
				case myGains > 0.0:
					symbol = setColor(val.Symbol, "green")
					log.Println(val.Symbol, ": myGains are positive, line color set to green")
				}
				if myChangePercent <= -10.0 {
					symbol = setColor(val.Symbol, "redbold")
					log.Println(val.Symbol, ": myChangePercent is less than -10%, line color set to red bold")
				}
				if myTargetSell > 0 && myTargetSellDiff >= 0.0 {
					symbol = setColor(val.Symbol, "greenbold")
					log.Println(val.Symbol, ": price is higher than myTargetSell, line color set to green bold")
				}

			} else if myTargetBuy > 0 {
				log.Println(val.Symbol, ": only targetBuy is processed from the xlsx")
				myTargetBuyDiff = ((myTargetBuy - val.RegularMarketPrice) / val.RegularMarketPrice) * 100

				if myTargetBuy > 0 && myTargetBuyDiff >= 0.0 {
					symbol = setColor(val.Symbol, "greenbold")
					log.Println(val.Symbol, ": price is below myTargetBuy, line color set to green bold")
				} else {
					symbol = setColor(val.Symbol, "default")
					log.Println(val.Symbol, ": price is still too high for myTargetBuy, line color set to default")
				}
			} else {
				// if there is only symbol in the xlsx (without myPrice, etc...), color is set from default -> on the dayChange
				switch {
				case val.RegularMarketChange < 0:
					symbol = setColor(val.Symbol, "red")
					log.Println(val.Symbol, ": DayChange is negative, line color set to red")

				case val.RegularMarketChange == 0:
					symbol = setColor(val.Symbol, "default")
					log.Println(val.Symbol, ": DayChange is equal to 0, line color set to default")

				case val.RegularMarketChange > 0:
					symbol = setColor(val.Symbol, "green")
					log.Println(val.Symbol, ": DayChange is positive, line color set to green")
				}
			}

			// if the DayChange is lower than -4%, it override previous colors
			if val.RegularMarketChangePercent <= -4.0 {
				symbol = setColor(val.Symbol, "redbold")
				log.Println(val.Symbol, ": RegularMarketChangePercent is less than -4%, line color set to red bold")
			}

		} else {
			log.Println(val.Symbol, ": symbol is not in the xlsx, all myVariables are set to zero")
			myPrice = 0.0
			myChangePercent = 0.0
			myGains = 0.0
			myTargetSell = 0.0
			myTargetSellDiff = 0.0
			myTargetBuy = 0.0
			myTargetBuyDiff = 0.0

			//symbol = strings.ToLower(val.Symbol)
			switch {
			case val.RegularMarketChange < 0:
				symbol = setColor(val.Symbol, "red")
				log.Println(val.Symbol, ": DayChange is negative, line color set to red")

			case val.RegularMarketChange == 0:
				symbol = setColor(val.Symbol, "default")
				log.Println(val.Symbol, ": DayChange is equal to 0, line color set to default")

			case val.RegularMarketChange > 0:
				symbol = setColor(val.Symbol, "green")
				log.Println(val.Symbol, ": DayChange is positive, line color set to green")
			}
		}

		name := reduceName(val.ShortName)
		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		if myPrice > 0 && myQuantity > 0 {
			// if there are data from already bought stock from the xls
			log.Println(val.Symbol, ": print row when myPrice & myQuantity are set")
			fr.Fprintf(tabw,
				"%s \t %s \t %.2f \t %.2f \t %.2f%s \t %.2f \t %.2f%s \t %s \t %.2f%s \t %.2f \t %.2f \t %v %s\n",
				symbol,
				name,
				val.RegularMarketPrice,
				myPrice,
				myChangePercent, "%",
				myGains,
				myTargetSellDiff, "%",
				" --",
				val.RegularMarketChangePercent, "%",
				val.TwoHundredDayAverage,
				val.FiftyDayAverage,
				dateInfo.Format("02 Jan 15:04"),
				setColor("", "default"),
			)
		} else if myTargetBuy > 0.0 {
			// if there is a wish to buy with targetBuy not empty
			log.Println(val.Symbol, ": print row when myPrice & myQuantity are not set, but myTargetBuy is set")
			fr.Fprintf(tabw,
				"%s \t %s \t %.2f \t %s \t %s \t %s \t %s \t %.2f%s \t %.2f%s \t %.2f \t %.2f \t %v %s\n",
				symbol,
				name,
				val.RegularMarketPrice,
				" --",
				" --",
				" --",
				" --",
				myTargetBuyDiff, "%",
				val.RegularMarketChangePercent, "%",
				val.TwoHundredDayAverage,
				val.FiftyDayAverage,
				dateInfo.Format("02 Jan 15:04"),
				setColor("", "default"),
			)
		} else {
			// if there is no data from the xls, we replace the output by "-"
			log.Println(val.Symbol, ": print row when no data are found in the xlsx")
			fr.Fprintf(tabw,
				"%s \t %s \t %.2f \t %s \t %s \t %s \t %s \t %s \t %.2f%s \t %.2f \t %.2f \t %v %s\n",
				symbol,
				name,
				val.RegularMarketPrice,
				" --",
				" --",
				" --",
				" --",
				" --",
				val.RegularMarketChangePercent, "%",
				val.TwoHundredDayAverage,
				val.FiftyDayAverage,
				dateInfo.Format("02 Jan 15:04"),
				setColor("", "default"),
			)
		}
	}

	// print total row
	totalMyChangePercent = (totalMyGains / totalMyPrice) * 100
	symbol = setColor(" --", "yellow")
	name := "TOTAL SUM"
	fr.Fprintf(tabw,
		"%s \t %s \t %.2f \t %.2f \t %.2f%s \t %.2f \t %s \t %s \t %s \t %s \t %s \t %s %s\n",
		symbol,
		name,
		totalPrice,
		totalMyPrice,
		totalMyChangePercent, "%",
		totalMyGains,
		" --",
		" --",
		" --",
		" --",
		" --",
		dateRequest.Format("02 Jan 15:04"),
		setColor("", "default"),
	)

	fmt.Fprint(tabw, setColor("", "default"))
}

// print the standard info
func printDefaultTable(stockArr []StockInfo) {

	fr := message.NewPrinter(language.French)

	var symbol string

	dateRequest := time.Now()
	fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// configure the tabwriter to print all in a clean table
	// output, minwidth, tabwidth, padding, padchar, flags
	tabw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	defer tabw.Flush()

	// define a variable for each column name
	//headerLongName := "LongName"
	//headerRegularMarketChange := "Change"
	headerSymbol := "Symbol"
	headerShortName := "Name"
	headerFiftyDayAverage := "50DayAvg"
	headerFiftyTwoWeekHigh := "52WeekHigh"
	headerFiftyTwoWeekLow := "52WeekLow"
	headerRegularMarketChangePercent := "Change"
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
		//symbol := strings.ToLower(val.Symbol)
		switch {
		case val.RegularMarketChange < 0:
			symbol = setColor(val.Symbol, "red")
			log.Println(val.Symbol, ": DayChange is negative, line color set to red")

		case val.RegularMarketChange == 0:
			symbol = setColor(val.Symbol, "default")
			log.Println(val.Symbol, ": DayChange is equal to 0, line color set to default")

		case val.RegularMarketChange > 0:
			symbol = setColor(val.Symbol, "green")
			log.Println(val.Symbol, ": DayChange is positive, line color set to green")
		}
		if val.RegularMarketChangePercent < -4.0 {
			symbol = setColor(val.Symbol, "redbold")
			log.Println(val.Symbol, ": DayChange is less than -4%, line color set to red bold")
		}

		name := reduceName(val.ShortName)
		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		fr.Fprintf(tabw,
			"%s \t %s \t %.2f \t %.2f \t %.2f%s \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %.2f \t %v %s\n",
			symbol,
			name,
			val.RegularMarketPrice,
			val.RegularMarketPreviousClose,
			//val.RegularMarketChange,
			val.RegularMarketChangePercent, "%",
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
