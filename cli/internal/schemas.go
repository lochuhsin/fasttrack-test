package internal

type question struct {
	Q       string   `json:"q"`
	Options []string `json:"options"`
}

type questionResponse struct {
	Questions []question `json:"questions"`
	Next      *int       `json:"next"`
}

type entry struct {
	QuestionId int `json:"questionId"`
	Answer     int `json:"answer"`
}

type record struct {
	Name    string  `json:"name"`
	Answers []entry `json:"answers"`
}
