package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	// "github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

	/*
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
	*/

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
	if len(nameSplitSpace) > 2 {
		log.Println(nameLower, ": number of space separated words", len(nameSplitSpace), ", keep 3 words")
		nameTmp = nameSplitSpace[0] + " " + nameSplitSpace[1] + " " + nameSplitSpace[2]
	} else if len(nameSplitSpace) == 2 {
		log.Println(nameLower, ": number of space separated words", len(nameSplitSpace), ", keep 2 words")
		nameTmp = nameSplitSpace[0] + " " + nameSplitSpace[1]
	} else {
		log.Println(nameLower, ": number of space separated words", len(nameSplitSpace), ", keep it")
		nameTmp = nameSplitSpace[0]
	}

	nameSplitMinus := strings.Split(nameTmp, "-")
	if len(nameSplitMinus) > 1 {
		log.Println(nameLower, ": number of words separated by a minus", len(nameSplitMinus), ", keep 2 words")
		newName = nameSplitMinus[0] + "-" + nameSplitMinus[1]
	} else {
		log.Println(nameLower, ": number of words separated by a minus", len(nameSplitMinus), ", keep it")
		newName = nameSplitMinus[0]
	}
	return newName
}

// print the table if custom data are given from the xlsx
func printCustomTable(stockArr []StockInfo, dataXls [][]string) {

	fr := message.NewPrinter(language.French)

	//yellow := color.New(color.FgYellow).SprintFunc()
	//yellowBold := color.New(color.BgYellow).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	//redBold := color.New(color.BgRed).SprintFunc()

	myPrice := 0.0
	mySum := 0.0
	myQuantity := 0.0
	myGains := 0.0
	myChangePercent := 0.0
	myTargetSell := 0.0
	myTargetSellDiff := 0.0
	myTargetBuy := 0.0
	myTargetBuyDiff := 0.0

	myPriceTab := "--"
	mySumTab := "--"
	myGainsTab := "--"
	myChangePercentTab := "--"
	myTargetSellDiffTab := "--"
	myTargetBuyDiffTab := "--"

	totalPrice := 0.0
	totalMySum := 0.0
	totalMyGains := 0.0
	totalMyPrice := 0.0
	totalMyChangePercent := 0.0

	totalDayChangeDiff := 0.0
	yesterdayTotal := 0.0
	dayChangeTab := ""

	var xlsInfoPresent bool
	var err error
	var symbol string

	dateRequest := time.Now()
	//fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// Init the table
	resuTable := table.NewWriter()

	// Set table parameters and style
	resuTable.Style().Options.DrawBorder = false
	resuTable.Style().Options.SeparateColumns = false
	resuTable.Style().Options.SeparateFooter = false
	resuTable.Style().Options.SeparateHeader = false
	resuTable.Style().Options.SeparateRows = false

	//resuTable.SetStyle(table.StyleLight) // style avec bordure "fine", ce parametre écrase les  Style().Options.XXX
	//resuTable.SetCaption("Table using the style 'StyleLight'.\n")

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
	headerMySum := "MySum"
	headerMyTargetSellDiff := "MyTargetSell"
	headerMyTargetBuyDiff := "MyTargetBuy"

	// set the table header
	rowHeader := table.Row{
		setColor(headerSymbol, "yellow"),
		headerShortName,
		text.AlignRight.Apply(headerRegularMarketPrice, 12),
		text.AlignRight.Apply(headerMyPrice, 12),
		text.AlignRight.Apply(headerMySum, 12),
		text.AlignRight.Apply(headerMyChangePercent, 10),
		text.AlignRight.Apply(headerMyGains, 12),
		text.AlignRight.Apply(headerMyTargetSellDiff, 10),
		text.AlignRight.Apply(headerMyTargetBuyDiff, 10),
		text.AlignRight.Apply(headerRegularMarketChangePercent, 10),
		text.AlignRight.Apply(headerTwoHundredDayAverage, 11),
		text.AlignRight.Apply(headerFiftyDayAverage, 11),
		headerRegularMarketTime + setColor("", "default"),
	}
	//resuTable.AppendHeader(rowHeader)
	resuTable.AppendRow(rowHeader)

	// for each stock option, create the row and add to the table
	for _, val := range stockArr {

		// reset all the stock values
		myPrice = 0.0
		mySum = 0.0
		myQuantity = 0.0
		myGains = 0.0
		myChangePercent = 0.0
		myTargetSell = 0.0
		myTargetSellDiff = 0.0
		myTargetBuy = 0.0
		myTargetBuyDiff = 0.0

		myPriceTab = "--"
		mySumTab = "--"
		myGainsTab = "--"
		myChangePercentTab = "--"
		myTargetSellDiffTab = "--"
		myTargetBuyDiffTab = "--"

		// search for the xls stock info
		xlsInfoPresent = false
		for _, xlsRow := range dataXls {
			// if find the same symbol in the xls and row are not empty
			if strings.EqualFold(strings.ToLower(val.Symbol), strings.ToLower(xlsRow[0])) {
				xlsInfoPresent = true
				log.Println(val.Symbol, ": symbol is present in xlsx, row size is:", len(xlsRow))
				if len(xlsRow) == 1 {
					log.Println(val.Symbol, ": symbol is alone, all myVariables are set to zero, and myVariablesTab set to --")
					myPrice = 0.0
					mySum = 0.0
					myQuantity = 0.0
					myGains = 0.0
					myChangePercent = 0.0
					myTargetSell = 0.0
					myTargetSellDiff = 0.0
					myTargetBuy = 0.0
					myTargetBuyDiff = 0.0

					myPriceTab = "--"
					mySumTab = "--"
					myGainsTab = "--"
					myChangePercentTab = "--"
					myTargetSellDiffTab = "--"
					myTargetBuyDiffTab = "--"
				}
				if len(xlsRow) > 1 {
					if xlsRow[1] != "" {
						myPrice, err = strconv.ParseFloat(xlsRow[1], 64)
						log.Printf("%s : myPrice set to: %.4f", val.Symbol, myPrice)
						//log.Println(val.Symbol, ": myPrice set to:", myPrice)
						checkErr(err)
					} else {
						myPrice = 0.0
						myPriceTab = "--"
						log.Println(val.Symbol, ": myPrice is set to zero and myPriceTab to --, the value of xlsRow[1] is:", xlsRow[1])
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
				mySum = myPrice * myQuantity
				myChangePercent = ((val.RegularMarketPrice - myPrice) / myPrice) * 100
				myGains = (val.RegularMarketPrice - myPrice) * myQuantity
				mySumToday := val.RegularMarketPrice * myQuantity
				mySumYesterday := val.RegularMarketPreviousClose * myQuantity

				totalPrice = totalPrice + (val.RegularMarketPrice * myQuantity)
				totalMyPrice = totalMyPrice + (myPrice * myQuantity)
				totalMySum = totalMyPrice
				totalMyGains = totalMyGains + myGains

				yesterdayTotal = yesterdayTotal + mySumYesterday
				dayChangeDiff := mySumToday - mySumYesterday
				totalDayChangeDiff = totalDayChangeDiff + dayChangeDiff
				log.Printf("%s : price yesterday=%.3f , price today=%.3f", val.Symbol, val.RegularMarketPreviousClose, val.RegularMarketPrice)
				log.Printf("%s : mySum yesterday=%.3f , mySum today=%.3f , day change diff=%.3f , totalDayChangeDiff=%.3f",
					val.Symbol, mySumYesterday, mySumToday, dayChangeDiff, totalDayChangeDiff)

				myPriceTab = fr.Sprintf("%.2f", myPrice)
				mySumTab = fr.Sprintf("%.2f", mySum)
				myChangePercentTab = fr.Sprintf("%.2f %%", myChangePercent)
				myGainsTab = fr.Sprintf("%.2f", myGains)

				if myTargetSell > 0.0 {
					log.Printf("%s : targetSell=%.2f is processed from the xlsx", val.Symbol, myTargetSell)
					myTargetSellDiff = ((val.RegularMarketPrice - myTargetSell) / val.RegularMarketPrice) * 100
					myTargetSellDiffTab = fr.Sprintf("%.2f %%", myTargetSellDiff)
					log.Printf("%s : myTargetSellDiff=%.2f , myTargetSellDiffTab=%s", val.Symbol, myTargetSellDiff, myTargetSellDiffTab)
				} else {
					myTargetSellDiffTab = "--"
				}

				// set the row color when myPrice and myQuantity are set
				switch {
				case myGains < 0.0:
					symbol = setColor(val.Symbol, "red")
					log.Printf("%s : myGains=%.2f are negative, line color set to red", val.Symbol, myGains)
				case myGains == 0.0:
					symbol = setColor(val.Symbol, "default")
					log.Printf("%s : myGains=%.2f are equal to 0, line color set to default", val.Symbol, myGains)
				case myGains > 0.0:
					symbol = setColor(val.Symbol, "green")
					log.Printf("%s : myGains=%.2f are positive, line color set to green", val.Symbol, myGains)
				}
				if myChangePercent <= -10.0 {
					symbol = setColor(val.Symbol, "redbold")
					log.Printf("%s : myChangePercent=%.2f is less than -10%%, line color set to red bold", val.Symbol, myChangePercent)
				}
				if myTargetSell > 0 && myTargetSellDiff >= 0.0 {
					symbol = setColor(val.Symbol, "greenbold")
					log.Printf("%s : price=%.2f is higher than myTargetSell=%.2f, line color set to green bold", val.Symbol, val.RegularMarketPrice, myTargetSell)
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

		} else {
			log.Println(val.Symbol, ": symbol is not in the xlsx, all myVariables are set to zero")
			myPrice = 0.0
			mySum = 0.0
			myChangePercent = 0.0
			myGains = 0.0
			myTargetSell = 0.0
			myTargetSellDiff = 0.0
			myTargetBuy = 0.0
			myTargetBuyDiff = 0.0
			myPriceTab = "--"
			mySumTab = "--"
			myGainsTab = "--"
			myChangePercentTab = "--"
			myTargetSellDiffTab = "--"
			myTargetBuyDiffTab = "--"

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

		// if the DayChange is lower than -4%, it override previous colors to red bold
		if val.RegularMarketChangePercent <= -4.0 {
			symbol = setColor(val.Symbol, "redbold")
			log.Println(val.Symbol, ": RegularMarketChangePercent is less than -4%, line color set to red bold")
		}

		// process the target for buying
		if myTargetBuy > 0.0 {
			log.Println(val.Symbol, ": targetBuy is processed from the xlsx")
			myTargetBuyDiff = ((myTargetBuy - val.RegularMarketPrice) / val.RegularMarketPrice) * 100
			myTargetBuyDiffTab = fr.Sprintf("%.2f %%", myTargetBuyDiff)
			log.Printf("%s : myTargetBuyDiff=%.2f , myTargetBuyDiffTab=%s", val.Symbol, myTargetBuyDiff, myTargetBuyDiffTab)
			// log.Printf("%s : ", val.Symbol, )

			// if the price is lower than the target to buy, override colors to green bold
			if myTargetBuyDiff >= 0.0 {
				symbol = setColor(val.Symbol, "greenbold")
				log.Printf("%s : price=%.2f is below myTargetBuy=%.2f, line color set to green bold", val.Symbol, val.RegularMarketPrice, myTargetBuy)
			}
		} else {
			myTargetBuyDiffTab = "--"
		}

		// change just the color of dayChange
		/*
			if val.RegularMarketChangePercent < 0.0 {
				dayChangeTab = redBold(fr.Sprintf("%.2f %%", val.RegularMarketChangePercent))
			} else {
				dayChangeTab = fr.Sprintf("%.2f %%", val.RegularMarketChangePercent)
			}
		*/
		dayChangeTab = fr.Sprintf("%.2f %%", val.RegularMarketChangePercent)

		name := reduceName(val.ShortName)
		dateInfo := time.Unix(int64(val.RegularMarketTime), 0)

		// create the row with the stock option values
		tableRow := table.Row{
			symbol,
			name,
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.RegularMarketPrice), 12),
			text.AlignRight.Apply(myPriceTab, 12),
			text.AlignRight.Apply(mySumTab, 12),
			text.AlignRight.Apply(myChangePercentTab, 10),
			text.AlignRight.Apply(myGainsTab, 12),
			text.AlignRight.Apply(myTargetSellDiffTab, 10),
			text.AlignRight.Apply(myTargetBuyDiffTab, 10),
			text.AlignRight.Apply(dayChangeTab, 10),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.TwoHundredDayAverage), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.FiftyDayAverage), 11),
			dateInfo.Format("02 Jan 15:04") + setColor("", "default"),
		}
		resuTable.AppendRow(tableRow)
	}

	// print total footer row
	totalMyChangePercent = (totalMyGains / totalMyPrice) * 100.0
	symbol = setColor("--", "yellow")
	name := "TOTAL SUM"

	todayDiff := totalPrice - yesterdayTotal
	log.Printf("TOTAL SUM : yesterdayTotal= %.2f , todayDiff= %.2f , totalDayChangeDiff= %.2f", yesterdayTotal, todayDiff, totalDayChangeDiff)
	totalDayChange := ((todayDiff / yesterdayTotal) * 100.0)

	// create the footer row with the totals values
	rowFooter := table.Row{
		symbol,
		name,
		text.AlignRight.Apply(fr.Sprintf("%.2f", totalPrice), 12),
		text.AlignRight.Apply(fr.Sprintf("%.2f", totalMyPrice), 12),
		text.AlignRight.Apply(fr.Sprintf("%.2f", totalMySum), 12),
		text.AlignRight.Apply(fr.Sprintf("%.2f %%", totalMyChangePercent), 10),
		text.AlignRight.Apply(fr.Sprintf("%.2f", totalMyGains), 12),
		text.AlignRight.Apply("--", 10),
		text.AlignRight.Apply("--", 10),
		text.AlignRight.Apply(fr.Sprintf("%.2f %%", totalDayChange), 10),
		text.AlignRight.Apply("--", 11),
		text.AlignRight.Apply("--", 11),
		dateRequest.Format("02 Jan 15:04") + setColor("", "default"),
	}
	//resuTable.AppendFooter(rowFooter)
	resuTable.AppendRow(rowFooter)

	// print the whole table
	fmt.Println(resuTable.Render())

}

// print the standard info
func printDefaultTable(stockArr []StockInfo) {

	fr := message.NewPrinter(language.French)

	var symbol string

	dateRequest := time.Now()
	fmt.Println("Date request:", dateRequest.Format("02 Jan 2006 15:04"))

	// Init the table
	resuTable := table.NewWriter()

	// Set table parameters
	resuTable.Style().Options.DrawBorder = false
	resuTable.Style().Options.SeparateColumns = false
	resuTable.Style().Options.SeparateFooter = false
	resuTable.Style().Options.SeparateHeader = false
	resuTable.Style().Options.SeparateRows = false

	//resuTable.SetStyle(table.StyleLight) // style avec bordure "fine", ce parametre écrase les  Style().Options.XXX
	//resuTable.SetCaption("Table using the style 'StyleLight'.\n")
	//resuTable.SetStyle(table.StyleColoredBright)

	// define a variable for each column name
	//headerLongName := "LongName"
	headerSymbol := setColor("Symbol", "yellow")
	headerShortName := "Name"
	headerRegularMarketPrice := "Price"
	headerRegularMarketPreviousClose := "PrevClose"
	headerRegularMarketChangePercent := "DayChange"
	headerFiftyTwoWeekLow := "52WeekLow"
	headerFiftyTwoWeekHigh := "52WeekHigh"
	headerTwoHundredDayAverage := "200DayAvg"
	headerFiftyDayAverage := "50DayAvg"
	headerRegularMarketDayLow := "DayLow"
	headerRegularMarketDayHigh := "DayHigh"
	headerRegularMarketTime := "UpdateTime"

	// set the table header
	rowheader := table.Row{
		headerSymbol,
		headerShortName,
		text.AlignRight.Apply(headerRegularMarketPrice, 11),
		text.AlignRight.Apply(headerRegularMarketPreviousClose, 11),
		text.AlignRight.Apply(headerRegularMarketChangePercent, 10),
		text.AlignRight.Apply(headerRegularMarketDayLow, 11),
		text.AlignRight.Apply(headerRegularMarketDayHigh, 11),
		text.AlignRight.Apply(headerTwoHundredDayAverage, 11),
		text.AlignRight.Apply(headerFiftyDayAverage, 11),
		text.AlignRight.Apply(headerFiftyTwoWeekLow, 11),
		text.AlignRight.Apply(headerFiftyTwoWeekHigh, 11),
		headerRegularMarketTime + setColor("", "default"),
	}
	//resuTable.AppendHeader(rowheader)
	resuTable.AppendRow(rowheader)

	// create row for each stock option
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

		// create the row with the stock option values
		tableRow := table.Row{
			symbol,
			name,
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.RegularMarketPrice), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.RegularMarketPreviousClose), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f%%", val.RegularMarketChangePercent), 10),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.RegularMarketDayLow), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.RegularMarketDayHigh), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.TwoHundredDayAverage), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.FiftyDayAverage), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.FiftyTwoWeekLow), 11),
			text.AlignRight.Apply(fr.Sprintf("%.2f", val.FiftyTwoWeekHigh), 11),
			dateInfo.Format("02 Jan 15:04") + setColor("", "default"),
		}
		resuTable.AppendRow(tableRow)
	}

	// print the whole table
	fmt.Println(resuTable.Render())

}
