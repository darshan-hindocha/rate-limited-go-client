package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type (
	APIResource func(r http.Request) (http.Response, error)

	TPSMeasure struct {
		initialTimeStamp time.Time
		tpsCount         int
		tpsTarget        int
	}
)

func RandomStringGenerator(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TPSController(tpsMeasure *TPSMeasure) error {
	if tpsMeasure.tpsTarget <= 0 {
		if tpsMeasure.tpsCount%50 == 0 {
			log.Println("no tps limit set")
		}
		return nil
	}
	duration := time.Now().Sub(tpsMeasure.initialTimeStamp).Seconds()
	tpsLevel := float64(tpsMeasure.tpsCount) / duration
	if int(tpsLevel) > tpsMeasure.tpsTarget {
		err := fmt.Errorf("current tps level of %d tps exceeds tps target of %d tps", int(tpsLevel), tpsMeasure.tpsTarget)
		return err
	}
	log.Printf("current tps level of %d tps is safely within target of %d tps", int(tpsLevel), tpsMeasure.tpsTarget)
	tpsMeasure.initialTimeStamp = time.Now()
	tpsMeasure.tpsCount = 0
	return nil
}

func handler(tps int, resource APIResource) (outputFileName string) {

	inputFile, err := os.Open("./data/MOCK_DATA.csv")
	if err != nil {
		log.Fatal("error reading input file")
	}
	r := csv.NewReader(inputFile)
	inputData, err := r.ReadAll()
	if err != nil {
		log.Fatal("error reading input data into memory")
	}

	outputFileName = "./data/outputFile" + RandomStringGenerator(5) + ".csv"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal("error creating output file")
	}
	w := csv.NewWriter(outputFile)
	defer w.Flush()

	tpsMeasure := TPSMeasure{
		initialTimeStamp: time.Now(),
		tpsCount:         0,
		tpsTarget:        tps,
	}

	for _, record := range inputData {

		duration := time.Duration(int(time.Second) / tps)
		time.Sleep(duration)
		body := ioutil.NopCloser(strings.NewReader(record[0]))
		req := http.Request{
			Method: "POST",
			URL:    nil,
			Body:   body,
		}

		resp, err := resource(req)
		if err != nil {
			log.Fatal("error requesting resource: ", err)
		}

		responseBody := new(strings.Builder)
		_, err = io.Copy(responseBody, resp.Body)
		if err != nil {
			log.Fatal("error reading response body to string")
		}
		err = w.Write([]string{record[0], responseBody.String()})
		if err != nil {
			log.Fatal("error writing response body to output file")
		}

		tpsMeasure.tpsCount += 1

		if tpsMeasure.tpsCount > 50 {
			err = TPSController(&tpsMeasure)
			if err != nil {
				log.Fatal("tps level exceeded!: ", err)
			}
		}
	}
	return
}

func main() {
	//tps:= 200
	//_ = handler(tps, apiResource)
	return
}
