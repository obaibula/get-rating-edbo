package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	edboUrl          = "https://vstup.edbo.gov.ua/offer-requests/"
	facultyId        = "1338098"
	edboErrorMessage = "Помилка. Зверніться до системного адміністратора"
)

func getStudents() ([]student, error) {
	students := make([]student, 0, 1000)
	for last := 0; ; last += len(students) {
		data, err := fetchStudentsBatch(last)
		if err != nil {
			return nil, err
		}
		if len(data.Students) <= 0 {
			break
		}
		students = append(students, data.Students...)
	}
	return students, nil
}

// fetchStudentsBatch returns a batch of students starting from the given index `last`.
// The `last` parameter is a "page-like" variable used by the EDBO to paginate results.
// If `last` is 0, EDBO returns the first 200 students.
// If `last` is 200, EDBO returns the next 200 students starting from the 200th student.
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
		slog.Info("failed to parse json from response from edbo", "msg", err)
		slog.Info("retrying with workaround of `помилка` error from edbo...")
		body, repairErr := repairCorruptedBody(res.Body)
		if repairErr != nil {
			return data, errors.New(fmt.Sprintf("failed with workaround: %s", err.Error()))
		}
		err = json.NewDecoder(body).Decode(&data)
		if err != nil {
			return data, errors.New(fmt.Sprintf("failed to parse json from response from edbo: %s", err.Error()))
		}
	}
	return data, err
}

func repairCorruptedBody(corruptedBody io.ReadCloser) (io.Reader, error) {
	bodyBytes, err := io.ReadAll(corruptedBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse json from response from edbo: %s", err.Error()))
	}
	repairedBody := strings.Replace(string(bodyBytes), edboErrorMessage, "", -1)
	return strings.NewReader(repairedBody), nil
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
