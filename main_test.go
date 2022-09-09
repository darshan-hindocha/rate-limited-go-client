package main

import (
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func mockApiResource(_ http.Request) (res http.Response, err error) {

	body := RandomStringGenerator(10)
	res = http.Response{
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}
	return
}

func Test(t *testing.T) {
	tps := 40
	outputFileName := handler(tps, mockApiResource)

	inputFile, err := os.Open("./data/MOCK_DATA.csv")
	if err != nil {
		t.Fatal("error reading input file")
	}
	readInput := csv.NewReader(inputFile)
	inputData, err := readInput.ReadAll()
	if err != nil {
		t.Fatal("error reading input data into memory")
	}

	outputFile, err := os.Open(outputFileName)
	if err != nil {
		t.Fatal("error reading input file")
	}
	readOutput := csv.NewReader(outputFile)
	outputData, err := readOutput.ReadAll()
	if err != nil {
		t.Fatal("error reading input data into memory")
	}

	if len(outputData) != len(inputData) {
		t.Fatalf("output data has %d rows which is not equal to %d rows in input data", len(outputData), len(inputData))
	}
}
