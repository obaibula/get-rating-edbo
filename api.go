package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		data, err := fetchStudentsBatch(last)
		if err != nil {
			return nil, err
		}
		//todo: what if data is empty?
		students = append(students, data.Students...)
		if len(data.Students) < 200 {
			break
		}
	}
	return students, nil
}

// fetchStudentsBatch returns a batch of students starting from the given index `last`.
// The `last` parameter is a "page-like" variable used by the EDBO to paginate results.
// If `last` is 0, EDBO returns the first 200 students.
// If `last` is 200, EDBO returns the next 200 students starting from the 200th student.
// The function always attempts to return 200 students, but if there are fewer than 200 remaining,
// it will return all the remaining students.
func fetchStudentsBatch(last int) (studentsResponse, error) {
	var data studentsResponse
	req, err := createRequest(last)
	if err != nil {
		return data, err
	}

	res, err := sendRequest(req)
	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, errors.New(fmt.Sprintf("failed to parse json from response from edbo: %s", err.Error()))
	}
	return data, err
}

func createRequest(last int) (*http.Request, error) {
	payload := url.Values{
		"id":   {facultyId},
		"last": {strconv.Itoa(last)},
	}
	req, err := http.NewRequest("POST", edboUrl, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed creating request to edbo: %s", err.Error()))
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
		return nil, errors.New(fmt.Sprintf("failed to send the request to edbo: %s", err.Error()))
	}
	return res, nil
}
