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
