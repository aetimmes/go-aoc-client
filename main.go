package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aetimmes/go-aoc-client/aocclient"
)

func main() {
	file, _ := ioutil.ReadFile("cookies.txt")
	sessionID := strings.TrimSpace(string(file))
	data, err := aocclient.GetInput(2020, 2, sessionID)
	if err != nil {
		panic(err)
	}
	if data != "" {
		fmt.Println("data retrieved successfully")
	}

	response_type, err := aocclient.SubmitAnswer(2021, 2, 1, 1609, sessionID)
	fmt.Println(aocclient.ResponseTypeMap[response_type], err)
}
