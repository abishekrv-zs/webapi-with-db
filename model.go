package main

type department struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type employee struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	PhoneNumber string     `json:"phoneNumber"`
	Dept        department `json:"dept"`
}
