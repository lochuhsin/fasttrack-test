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

type question struct {
	Q       *string  `yaml:"q"`
	A       *int     `yaml:"a"`
	Options []string `yaml:"options"`
}

type QuestionDatabase struct {
	Questions []question `yaml:"questions"`
}

func (q QuestionDatabase) List(limit, offset int) ([]question, bool) {
	if offset >= len(q.Questions) {
		return nil, false
	}
	return q.Questions[offset : offset+limit], true
}

func newQuestionDatabase(questions []question) *QuestionDatabase {
	return &QuestionDatabase{
		Questions: questions,
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
	if db.Questions == nil {
		return fmt.Errorf("missing key or empty questions")
	}
	for _, q := range db.Questions {
		if q.A == nil || q.Q == nil || q.Options == nil {
			return fmt.Errorf("error parsing yml, missing q, a or options")
		}

		if *q.A < 0 || *q.A >= len(q.Options) {
			return fmt.Errorf("Invalid answer index")
		}
	}
	return nil
}
