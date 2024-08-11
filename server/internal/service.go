package internal

import (
	"errors"
	"sync"
)

var rMu sync.Mutex // guard create record with name
func CreateRecordWithName(name string, answers []answer) (int, error) {
	questions, _ := GetQuestionDatabase().List(-1, -1)

	rMu.Lock()
	defer rMu.Unlock()
	ok := GetUserGroup().Exists(name)
	if !ok {
		return 0, errors.New("User not exists, something went wrong ?")
	}
	count := 0
	for _, answer := range answers {
		qIndex, ans := answer.QuestionId, answer.Answer
		if questions[*qIndex].A == *ans {
			count++
		}
	}
	ok = GetRecords().Create(name, count)
	if !ok {
		return 0, errors.New("Record from this user was created")
	}
	return count, nil
}
