package dto

type Challenge struct {
	Value      string `json:"value"`
	Difficulty int    `json:"complexity"`
}

type Solution struct {
	Body string `json:"body"`
}

type Quote struct {
	Body string `json:"body"`
}
