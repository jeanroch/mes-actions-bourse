package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

/*
Version is filled during compilation with

	VERSION="$(git tag -l | tail -1)_$(date +%F)"
	VERSION="$(git tag -l | tail -1)_$(git log --pretty=format:"%ad" -n 1 --date=format:%F)"
	go build -ldflags "-X main.version=$VERSION"

const version = "v0.4 2022-06-13"
*/
var version string = ""

func main() {

	var symbols string = "^FCHI,^SBF120,^GDAXI,^STOXX50E,^IXIC,^GSPC,BTC-USD,ETH-USD" // list of symbols separated by a comma , (set by from the cli)
	var dataXls [][]string                                                            // data from xls file if given from the cli

	flagVersion := flag.Bool("ver", false, "Print version and exit\n")
	flagVerbose := flag.Bool("verb", false, "Print more verbose informations from the application run\n")
	flagTest := flag.Bool("test", false, "Use my local URL, for dev/testing purpose\n")
	flagIndice := flag.Bool("cac", false, "Get the values for few selected index and crypto-money: CAC40, SBF120, DAX, EuroStoxx50, Nasdaq, S&P500, Bitcoin and Ethereum\n")
	flagSymbols := flag.String("sym", "", "Need to provide a list of symbols separated by a comma, for example BTC-USD,BNP.PA,LI.PA,etc...\n")
	flagXLS := flag.String("xls", "", "Need to provide an xlsx file with the stock options to request\nThe First row of the file must be column title from this: Symbol, Price (PRU), Quantity, TargetSell, TargetBuy\nIt need at least the column Symbol\n")
	flagApiVer := flag.Int("apiver", 6, "API version to request on Yahoo webserver")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "\nThis small app is requesting some stock option values from the site finance.yahoo.com\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetOutput(io.Discard)
	if *flagVerbose {
		log.SetOutput(os.Stdout)
	}

	if *flagApiVer < 1 || *flagApiVer > 10 {
		fmt.Println("Yahoo API version must be between 1 and 10")
		os.Exit(1)
	}

	switch {
	case *flagVersion:
		fmt.Println("Version:", version)
		os.Exit(0)

	case *flagIndice:
		symbols = "^FCHI,^SBF120,^GDAXI,^STOXX50E,^IXIC,^GSPC,BTC-USD,ETH-USD"

	case *flagSymbols != "":
		symbols = *flagSymbols

	case *flagXLS != "":
		sheet := "Sheet1"
		// headers := []string{"symbol", "pru", "nombre", "objectif"}
		headers := []string{"symbol", "price", "quantity", "targetSell", "targetBuy"}
		symbols, dataXls = GetStockFromXLS(*flagXLS, sheet, headers)
	}

	fmt.Println("")

	dataJson := GetDataFromURL(symbols, *flagTest, *flagApiVer)
	dataInternet := GetDataFromJSON(dataJson)

	if len(dataXls) > 0 {
		printCustomTable(dataInternet, dataXls)
	} else {
		printDefaultTable(dataInternet)
	}

	fmt.Println("")
}
