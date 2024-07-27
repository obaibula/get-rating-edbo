package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type studentsResponse struct {
	Students []student `json:"requests"`
}

type status int

const (
	Uncategorized status = iota
	Applicated
	Unknown1
	Unknown2
	Rejected
	Registered
	Allowed
)

type student struct {
	Id         int     `json:"prid"`
	No         int     `json:"n"`
	Status     status  `json:"prsid"`
	Ptid       int     `json:"ptid"`
	FullName   string  `json:"fio"`
	Pa         int     `json:"pa"`
	D          int     `json:"d"`
	Cp         int     `json:"cp"`
	Acceptance string  `json:"cpt"`
	Artid      int     `json:"artid"`
	Rating     float64 `json:"kv"`
	Priority   int     `json:"p"`
	Quota      []rss   `json:"rss"`
}

type rss struct {
	Type  string `json:"t"`
	Quota string `json:"sn"`
}

const (
	facultyId   = "1338098"
	edboUrl     = "https://vstup.edbo.gov.ua/offer-requests/"
	alinaId     = 13670738
	alinaRating = 154.640
)

func main() {
	students, err := getStudents()
	if err != nil {
		fmt.Println("Error getting students", err)
		return
	}
	result := make([]student, 0, 1000)
	for _, s := range students {
		hasQuota := false
		for _, q := range s.Quota {
			if q.Quota != "" {
				hasQuota = true
			}
		}
		predicate := s.Status != Rejected &&
			s.Rating >= alinaRating &&
			s.Priority <= 1 &&
			s.Priority > 0 &&
			!hasQuota

		if predicate {
			result = append(result, s)
		}
	}
	currentRating := len(result) + 10
	fmt.Println("Your current rating is:", currentRating)
}

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

		/*respBody := new(bytes.Buffer)
		respBody.ReadFrom(res.Body)
		fmt.Println("Response Status:", res.Status)
		fmt.Println("Response Body:", respBody.String())*/
		bodyerr, err := io.ReadAll(res.Body)
		body := strings.Replace(string(bodyerr), "Помилка. Зверніться до системного адміністратора", "", -1)
		var data studentsResponse
		//err = json.NewDecoder(res.Body).Decode(&data)
		err = json.NewDecoder(strings.NewReader(body)).Decode(&data)
		if err != nil {
			return nil, err
		}
		if len(data.Students) == 0 {
			break
		}
		students = append(students, data.Students...)
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
