package aocclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

var localCacheBaseDirectory string = ".aoc"

type cacheResponse int

const (
	cacheValid cacheResponse = iota
	cacheTooHigh
	cacheTooLow
)

type outputCache struct {
	upperBound    int
	correctAnswer []int
	lowerBound    int
}

func newOutputCache() outputCache {
	return outputCache{
		upperBound:    math.MaxInt,
		lowerBound:    math.MinInt,
		correctAnswer: make([]int, 0),
	}
}

func getCacheDir(year, day int, sessionID string) string {
	return fmt.Sprintf("%s/%s/%d/%d/", localCacheBaseDirectory, sessionID, year, day)
}
func getInputCachePath(year, day int, sessionID string) (string, string) {
	return getCacheDir(year, day, sessionID), "input.txt"
}
func getOutputCachePath(year, day, level int, sessionID string) (string, string) {
	return getCacheDir(year, day, sessionID), fmt.Sprintf("%d.json", level)
}

func writeToInputCache(year, day int, sessionID, data string) error {
	dir, filename := getInputCachePath(year, day, sessionID)
	err := os.MkdirAll(dir, 0700)
	f, err := os.Create(dir + filename)
	err = os.Chmod(dir+filename, 0600)
	_, err = f.WriteString(data)
	return err
}

func getInputFromCache(year, day int, sessionID string) (string, error) {
	dir, filename := getInputCachePath(year, day, sessionID)
	f, err := ioutil.ReadFile(dir + filename)
	return string(f), err
}

func getOutputCache(year, day, level int, sessionID string) (outputCache, error) {
	dir, filename := getOutputCachePath(year, day, level, sessionID)
	f, err := ioutil.ReadFile(dir + filename)
	var result outputCache
	err = json.Unmarshal(f, &result)
	if err != nil {
		result = newOutputCache()
	}
	return result, err
}

func writeOutputCache(year, day, level int, cache outputCache, sessionID string) error {
	dir, filename := getOutputCachePath(year, day, level, sessionID)
	err := os.MkdirAll(dir, 0700)
	f, err := os.Create(dir + filename)
	err = os.Chmod(dir+filename, 0600)
	data, err := json.Marshal(cache)
	_, err = f.Write(data)
	return err
}

func updateOutputCache(year, day, level, answer int, resp AOCResponseType, sessionID string) error {
	if resp == TooHigh || resp == TooLow || resp == Correct {
		data, err := getOutputCache(year, day, level, sessionID)
		if err != nil {
			log.Printf("Failed to get cache, creating new cache: %s\n", err)
		}
		if resp == TooHigh && answer < data.upperBound {
			data.upperBound = answer
		} else if resp == TooLow && answer > data.lowerBound {
			data.lowerBound = answer
		} else if resp == Correct && len(data.correctAnswer) == 0 {
			data.correctAnswer = append(data.correctAnswer, answer)
		}
		err = writeOutputCache(year, day, level, data, sessionID)
		return err
	}
	return nil
}

func validateAnswer(answer int, cache outputCache) cacheResponse {
	if answer <= cache.lowerBound {
		return cacheTooLow
	}
	if answer >= cache.upperBound {
		return cacheTooHigh
	}
	return cacheValid
}
