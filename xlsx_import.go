package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// function GetStockFromXLS extract data from an xlsx
func GetStockFromXLS(fileXls string, sheetName string, headerList []string) (string, [][]string) {

	var symbolString string    // return a list of symbols separated by comma ',' ready for the request URL
	var resuTable [][]string   // result table to return
	var columnList []int       // column id/number to use to fetch data from the xlsx
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
						log.Printf("[INFO] xlsx import : Header present in the xlsx \t name: %s \t col: %d", header, position)
						columnList = append(columnList, position)
					}
				}
				// if header not found in the column name, print a warning, except if "symbol" as this one is mandatory
				if !headerFound {
					if strings.EqualFold(strings.ToLower(header), "symbol") {
						fmt.Println("[ERROR] xlsx import : The header \"Symbol\" is missing in the file:", fileXls)
						os.Exit(1)
					} else {
						log.Println("[INFO] xlsx import : The header \"", header, "\" is missing in the file:", fileXls)
					}
				}
			}

		} else {
			// insert the selected data in a row, to include in the result table
			for _, position := range columnList {
				if len(row) > position {
					log.Printf("[INFO] xlsx import : Import data at row[%v]= %v", position, row[position])
					rowCustom = append(rowCustom, row[position])
				}
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
		if nb == 0 || symbolString == "" {
			symbolString = symb[0]
		} else {
			if symb[0] != "" {
				symbolString = symbolString + "," + symb[0]
			}
		}
	}

	log.Println("[INFO] xlsx import : result table :", resuTable)

	return symbolString, resuTable
}
