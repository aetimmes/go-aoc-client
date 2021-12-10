package aocclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getInput(year, day int, sessionID string) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", fmt.Sprintf("session %s", sessionID))
	resp, err := client.Do(req)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error = %s \n", err)
	}
	fmt.Printf("Response = %s", string(data))
}
