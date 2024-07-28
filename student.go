package main

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

type priority int

const (
	Contract = iota
	First
	Second
	Third
	Fourth
	Fifth
)

type student struct {
	Id         int      `json:"prid"`
	No         int      `json:"n"`
	Status     status   `json:"prsid"`
	Ptid       int      `json:"ptid"`
	FullName   string   `json:"fio"`
	Pa         int      `json:"pa"`
	D          int      `json:"d"`
	Cp         int      `json:"cp"`
	Acceptance string   `json:"cpt"`
	Artid      int      `json:"artid"`
	Rating     float64  `json:"kv"`
	Priority   priority `json:"p"`
	Quota      []rss    `json:"rss"`
}

type rss struct {
	Type  string `json:"t"`
	Quota string `json:"sn"`
}
