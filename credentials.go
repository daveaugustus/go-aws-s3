package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

var (
	cred = map[string]string{}
)

func GetAWSCrtedentials() (string, string) {
	csvFile, err := os.Open(`C:\Users\davet\Etc\awskey.csv`)
	if err != nil {
		log.Printf("Error while opening the keyfile: %s", err.Error())
	}
	reader := csv.NewReader(csvFile)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		credSlice := strings.Split(line[0], "=")
		cred[credSlice[0]] = credSlice[1]
	}

	return cred["AWSAccessKeyId"], cred["AWSSecretKey"]
}
