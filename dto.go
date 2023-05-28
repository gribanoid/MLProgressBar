package main

type Model struct {
	Vacancy string `json:"vacancy"`
	Resume  string `json:"resume"`
}

type Similarity struct {
	Score int `json:"score"`
}
