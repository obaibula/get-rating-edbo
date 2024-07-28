package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	edboUrl   = "https://vstup.edbo.gov.ua/offer-requests/"
	facultyId = "1338098"
)

func getStudents() ([]student, error) {
	students := make([]student, 0, 1000)
	for last := 0; ; last += 200 {
		req, err := createRequest(last)
		if err != nil {
			return nil, err
		}

		res, err := sendRequest(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var data studentsResponse
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		students = append(students, data.Students...)
		if len(data.Students) < 200 {
			break
		}
	}
	return students, nil
}

func createRequest(last int) (*http.Request, error) {
	payload := url.Values{
		"id":   {facultyId},
		"last": {strconv.Itoa(last)},
	}
	req, err := http.NewRequest("POST", edboUrl, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://vstup.edbo.gov.ua")
	req.Header.Set("Referer", "https://vstup.edbo.gov.ua/offer/1338098/")
	return req, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
