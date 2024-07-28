package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	alinaRating = 154.64
	maxQuota    = 10 // number of places for people with quota (will be occupied in any case)
)

func main() {
	students, err := getStudents()
	if err != nil {
		fmt.Println("Error getting students", err)
		return
	}
	currentRating := getCurrentRatingBasedOnFirstPriority(students)
	fmt.Println("Your current rating is:", currentRating)

}

type studentPredicate func(student) bool

func (sp studentPredicate) and(other studentPredicate) studentPredicate {
	return func(s student) bool {
		return sp(s) && other(s)
	}
}

func (sp studentPredicate) or(other studentPredicate) studentPredicate {
	return func(s student) bool {
		return sp(s) || other(s)
	}
}

func (sp studentPredicate) negate() studentPredicate {
	return func(s student) bool {
		return !sp(s)
	}
}

func hasQuota() studentPredicate {
	return func(s student) bool {
		for _, q := range s.Quota {
			if q.Quota != "" {
				return true
			}
		}
		return false
	}
}

func notRejected() studentPredicate {
	return func(s student) bool {
		return s.Status != Rejected
	}
}

func ratingAboveOrEqual(threshold float64) studentPredicate {
	return func(s student) bool {
		return s.Rating >= threshold
	}
}

func priorityAbove(threshold int) studentPredicate {
	return func(s student) bool {
		return s.Priority <= threshold && s.Priority > 0 // 0 priority is reserved for contract
	}
}

func getCurrentRatingBasedOnFirstPriority(students []student) int {
	predicate := hasQuota().negate().
		and(notRejected()).
		and(ratingAboveOrEqual(alinaRating)).
		and(priorityAbove(1))

	return filter(students, predicate)
}

func filter(students []student, predicate studentPredicate) int {
	result := make([]student, 0, len(students))
	for _, s := range students {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return len(result) + maxQuota
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
