package aocclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


func GetInput(year, day int, sessionID string) (string, error) {
	if data, err := getInputFromCache(year, day, sessionID); err == nil {
		return data, err
	} else {
		err = redactError(err, sessionID)
		log.Printf("failed to get cached input for %d day %d: %s", year, day, err)
	}
	if data, err := getInputFromSite(year, day, sessionID); err == nil {
		if err := writeToInputCache(year, day, sessionID, data); err != nil {
			err = redactError(err, sessionID)
			log.Printf("failed to cache input for %d day %d: %s", year, day, err)
		}
		return data, err
	} else {
		log.Fatalf("failed to get problem input for %d day %d: %s", year, day, err)
		return data, redactError(err, sessionID)
	}
}


func getInputFromSite(year, day int, sessionID string) (string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionID))
	resp, err := client.Do(req)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error = %s \n", err)
	}
	return string(data), err
}

func redactError(err error, sessionID string) error {
	if err == nil || !strings.Contains(err.Error(), sessionID) {
		return err
	}
	return errors.New(strings.ReplaceAll(err.Error(), sessionID, "<sessionID>"))
}
