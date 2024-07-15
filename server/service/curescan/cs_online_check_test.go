package curescan

import (
	"fmt"
	"os"
	"testing"
)

func TestOnlineCheckService_ParseFileTo(t *testing.T) {
	CSOnlineCheckService := OnlineCheckService{}
	file, err := os.Open("./onlineCheck.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	data, err := CSOnlineCheckService.ParseFileTo(file)
	if err != nil {
		panic(err)
	}
	for i, datum := range data {
		fmt.Println(i, *datum)
	}
}

func TestPortScanService_ParseFileTo(t *testing.T) {
	PortScanService := PortScanService{}
	file, err := os.Open("./portScan.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	data, err := PortScanService.ParseFileTo(file)
	if err != nil {
		panic(err)
	}
	for i, datum := range data {
		fmt.Println(i, *datum)
	}
}
