package aocclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AOCResponseType int64

const (
	Correct AOCResponseType = iota
	TooHigh
	TooLow
	RateLimited
	WrongAccount
	BadResponseCode
	WrongLevel
	HTMLParseError
	UnknownError
)

var ResponseTypeMap = map[AOCResponseType]string{
	Correct:         "Answer Correct!",
	TooHigh:         "Answer too high",
	TooLow:          "Answer too low",
	RateLimited:     "Rate limited, please wait",
	WrongAccount:    "Right answer, but for the wrong account",
	BadResponseCode: "Non-200 HTTP response code received for submission",
	WrongLevel:      "Maybe not the right level?",
	HTMLParseError:  "Failed to parse response HTML",
	UnknownError:    "Unknown Error",
}

func postAnswer(year, day, level, answer int, sessionID string) (AOCResponseType, error) {
	uri := fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", year, day)
	values := url.Values{}
	values.Add("level", fmt.Sprintf("%d", level))
	values.Add("answer", fmt.Sprintf("%d", answer))
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionID))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error = %s \n", err)
		return UnknownError, err
	}
	result, err := ParseResponse(resp)
	resp.Body.Close()
	return result, err
}

func ParseResponse(resp *http.Response) (AOCResponseType, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	data := ""
	if err != nil {
		return HTMLParseError, err
	}
	doc.Find("html body main article").Each(func(i int, s *goquery.Selection) {
		data += s.Find("p").Text()
	})
	return GetAOCResponseType(data)
}

func GetAOCResponseType(data string) (AOCResponseType, error) {
	if strings.Contains(data, "You don't seem to be solving the right level") {
		return WrongLevel, nil
	}
	if strings.Contains(data, "you might be logged in to the wrong account") {
		return WrongAccount, nil
	}
	if strings.Contains(data, "your answer is too high") {
		return TooHigh, nil
	}
	if strings.Contains(data, "your answer is too low") {
		return TooLow, nil
	}
	if strings.Contains(data, "you have to wait") {
		return RateLimited, nil
	}
	if strings.Contains(data, "That's the right answer!") {
		return Correct, nil
	}

	return UnknownError, errors.New(fmt.Sprintf("Couldn't determine status from response: %s", data))
}

func SubmitAnswer(year, day, level, answer int, sessionID string) (AOCResponseType, error) {
	cache, err := getOutputCache(year, day, level, sessionID)
	isValid := validateAnswer(answer, cache)
	if isValid == cacheTooLow {
		return UnknownError, errors.New(fmt.Sprintf("Answer %d for %d day %d level %d smaller than previously observed lower bound", answer, year, day, level))
	}
	if isValid == cacheTooHigh {
		return UnknownError, errors.New(fmt.Sprintf("Answer %d for %d day %d level %d larger than previously observed upper bound", answer, year, day, level))
	}
	resp, err := postAnswer(year, day, level, answer, sessionID)
	updateOutputCache(year, day, level, answer, resp, sessionID)
	return resp, err
}
