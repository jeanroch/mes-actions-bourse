package main

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
