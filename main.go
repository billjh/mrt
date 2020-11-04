package main

import (
	"fmt"
	"os"
)

var allStations []Station

func main() {
	// load all stations from csv file
	csvFile, err := os.Open("./data/StationMap.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	allStations, err = ReadStations(csvFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d stations from csv file '%s'\n", len(allStations), csvFile.Name())
}
