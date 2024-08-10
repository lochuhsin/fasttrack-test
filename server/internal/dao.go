package internal

import "sync"

var (
	userG   *UserGroup
	uOnce   sync.Once = sync.Once{}
	submitR *SubmitRecord
	sOnce   sync.Once = sync.Once{}
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
	return loaded
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
