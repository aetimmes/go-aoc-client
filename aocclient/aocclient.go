package aocclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var localCacheBaseDirectory string = ".aoc"
var fileName string = "input.txt"

func GetInput(year, day int, sessionID string) (string, error) {
	if data, err := getInputFromFile(year, day, sessionID); err == nil {
		return data, err
	} else {
		err = redactError(err, sessionID)
		log.Printf("failed to get cached input for %d day %d: %s", year, day, err)
	}
	if data, err := getInputFromSite(year, day, sessionID); err == nil {
		if err := writeToCache(year, day, sessionID, data); err != nil {
			err = redactError(err, sessionID)
			log.Printf("failed to cache input for %d day %d: %s", year, day, err)
		}
		return data, err
	} else {
		log.Fatalf("failed to get problem input for %d day %d: %s", year, day, err)
		return data, redactError(err, sessionID)
	}
}

func getPath(year, day int, sessionID string) (string, string) {
	return fmt.Sprintf(
		"%s/%s/%d/%d/",
		localCacheBaseDirectory,
		sessionID,
		year,
		day), fileName
}

func writeToCache(year, day int, sessionID, data string) error {
	dir, filename := getPath(year, day, sessionID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil
	}
	if f, err := os.Create(dir + filename); err == nil {
		err := os.Chmod(dir+filename, 0600)
		if err != nil {
			return err
		}
		if _, err := f.WriteString(data); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func getInputFromFile(year, day int, sessionID string) (string, error) {
	dir, filename := getPath(year, day, sessionID)
	f, err := ioutil.ReadFile(dir + filename)
	return string(f), err
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
