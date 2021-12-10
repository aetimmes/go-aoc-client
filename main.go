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
	data, err := aocclient.GetInput(2021, 2, sessionID)
	if err != nil {
		panic(err)
	}
	if data != "" {
		fmt.Println("data retrieved successfully")
	}
}
