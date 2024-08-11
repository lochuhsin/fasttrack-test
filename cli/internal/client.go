package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	DEFAULT_SERVER_HOST = "http://localhost:8000"
)

func CreateUser(username string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(struct {
		Name string `json:"name"`
	}{Name: username})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", DEFAULT_SERVER_HOST, "users")
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("unable to create user")
	}
	return nil
}

func GetQuestion() (*questionResponse, error) {
	url := fmt.Sprintf("%s/%s", DEFAULT_SERVER_HOST, "questions")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("unable to get questions")
	}
	defer resp.Body.Close()

	obj := questionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func Submit(name string, ans []entry) (string, error) {
	url := fmt.Sprintf("%s/%s", DEFAULT_SERVER_HOST, "records")
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(record{
		Name:    name,
		Answers: ans,
	})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("unable to create user")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetPercentile(name string) (string, error) {
	url := fmt.Sprintf("%s/%s?name=%s", DEFAULT_SERVER_HOST, "percentile", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("unable to get percentile, invalid username perhaps ?")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
