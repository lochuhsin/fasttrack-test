package internal

type CreateUser struct {
	Name *string `json:"name"`
}

type answer struct {
	QuestionId *int `json:"questionId"`
	Answer     *int `json:"answer"`
}

type CreateRecord struct {
	Name    *string  `json:"name"`
	Answers []answer `json:"answers"`
}

type Question struct {
	Q       string   `json:"q"`
	Options []string `json:"options"`
}

type ListQuestions struct {
	Questions []Question `json:"questions"`
	Next      *int       `json:"next"`
}
