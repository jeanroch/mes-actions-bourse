package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// function GetStockFromXLS extract data from an xlsx
func GetStockFromXLS(fileXls string, sheetName string, headerList []string) (string, [][]string) {

	var symbolString string    // resturn a list of symbols separated by comma ',' ready for the request URL
	var resuTable [][]string   // result table to return
	var columnList []int       // column id/number to use to fecth data from the xlsx
	var rowCustom = []string{} // temporary array (like in one row) with the data to insert in the final result

	xls, err := excelize.OpenFile(fileXls)
	checkErr(err)
	defer xls.Close()

	// get all the rows from the targeted sheet (by default Sheet1)
	tableAll, err := xls.GetRows(sheetName)
	checkErr(err)

	for nb, row := range tableAll {

		// if it is the first row, then it is for the header
		if nb == 0 {
			// for each header in our list, find the correspondant column number in the xlsx file
			for _, header := range headerList {

				headerFound := false
				for position, cell := range row {

					// if the header name is the same as the cell content
					if strings.EqualFold(strings.ToLower(header), strings.ToLower(cell)) {
						headerFound = true
						//log.Println("Header present in the xlsx file:", header)
						columnList = append(columnList, position)
					}
				}
				// if header not found in the column name, print a warning, except if "symbol" as this one is mandatory
				if !headerFound {
					if strings.EqualFold(strings.ToLower(header), "symbol") {
						fmt.Println("[ERROR] The header \"Symbol\" is missing in the file:", fileXls)
						os.Exit(1)
					} else {
						fmt.Println("[INFO] The header \"", header, "\" is missing in the file:", fileXls)
					}
				}
			}

		} else {
			switch {
			// if row size is at least same as the header list
			case len(row) >= len(headerList):
				for _, position := range columnList {
					rowCustom = append(rowCustom, row[position])
				}
			// if row size is lower than header list but there is at least one column, at least the symbol
			case (len(row) > 0) && (len(row) < len(headerList)):
				rowCustom = []string{row[0]}
			}
		}
		// if our new temporary row is not empty, we can add it to the result
		if len(rowCustom) > 0 {
			// add the full row to the result table
			resuTable = append(resuTable, rowCustom)
			// and the reset the temporary info used to create a result row
			rowCustom = []string{}
		}
	}

	// prepare the symbols list
	for nb, symb := range resuTable {
		if nb == 0 {
			symbolString = symb[0]
		} else {
			symbolString = symbolString + "," + symb[0]
		}
	}

	return symbolString, resuTable
}
