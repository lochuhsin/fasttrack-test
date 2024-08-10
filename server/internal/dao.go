package internal

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

const DEFAULT_QUESTION_FILE_PATH = "question.yaml"

var (
	userG     *UserGroup
	uOnce     sync.Once = sync.Once{}
	submitR   *SubmitRecord
	sOnce     sync.Once = sync.Once{}
	questionD *QuestionDatabase
	qOnce     sync.Once = sync.Once{}
)

type UserGroup struct {
	users sync.Map
}

func (u *UserGroup) Create(name string) bool {
	/**
	 * return if the user is created or exists
	 */
	_, loaded := u.users.LoadOrStore(name, true)
	return !loaded
}

func (u *UserGroup) Exists(name string) bool {
	_, ok := u.users.Load(name)
	return ok
}

func newUserGroup() *UserGroup {
	return &UserGroup{
		users: sync.Map{},
	}
}

func InitUserGroup() {
	uOnce.Do(func() {
		userG = newUserGroup()
	})
}

func GetUserGroup() *UserGroup {
	return userG
}

type SubmitRecord struct {
	record sync.Map
	mu     sync.Mutex
}

func (s *SubmitRecord) Create(name string, value int) bool {
	_, loaded := s.record.LoadOrStore(name, value)
	return !loaded
}

func (s *SubmitRecord) Update(name string, value int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.record.Load(name)
	if !ok {
		return false
	}
	s.record.Swap(name, val.(int)+value)
	return true
}

func newSubmitRecord() *SubmitRecord {
	return &SubmitRecord{
		record: sync.Map{},
		mu:     sync.Mutex{},
	}
}

func InitSubmitRecord() {
	sOnce.Do(func() {
		submitR = newSubmitRecord()
	})
}

func GetSubmitRecord() *SubmitRecord {
	return submitR
}

type entry struct {
	Q       string   `yaml:"q"`
	A       int      `yaml:"a"`
	Options []string `yaml:"options"`
}

type QuestionDatabase struct {
	Entries []entry `yaml:"questions"`
}

func (q QuestionDatabase) List(limit, offset int) ([]entry, bool) {
	if limit == -1 && offset == -1 {
		return q.Entries, true
	}
	if offset >= len(q.Entries) {
		return nil, false
	}
	if limit < 0 || offset < 0 {
		return nil, false
	}
	if offset+limit > len(q.Entries) {
		return q.Entries[offset:len(q.Entries)], true
	}
	return q.Entries[offset : offset+limit], true
}

func (q QuestionDatabase) Count() int {
	return len(q.Entries)
}
func newQuestionDatabase(questions []entry) *QuestionDatabase {
	return &QuestionDatabase{
		Entries: questions,
	}
}

func InitQuestionDatabase() {
	qOnce.Do(func() {
		/**
		 * Read file from environment variable
		 */
		path, ok := os.LookupEnv("QUESTION_FILE") // Build up type settings ......etc
		if !ok {
			path = DEFAULT_QUESTION_FILE_PATH // write as constant
		}
		q := QuestionDatabase{}
		y, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(y, &q)
		if err != nil {
			panic(err)
		}
		err = validateFile(q)
		if err != nil {
			panic(err)
		}
		questionD = &q
	})
}

func validateFile(db QuestionDatabase) error {
	if db.Entries == nil {
		return fmt.Errorf("missing key or empty questions")
	}
	for _, q := range db.Entries {
		if q.Options == nil {
			return fmt.Errorf("error parsing yml, missing q, a or options")
		}
	}
	return nil
}

func GetQuestionDatabase() *QuestionDatabase {
	return questionD
}
