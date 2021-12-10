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

	fmt.Println(aocclient.getInput(2021, 1, sessionID))

}
