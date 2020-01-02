package megaplan

import (
	"encoding/json"
)

const (
	shortTimeFormat = "15:04"
	timeFormat      = "15:04:05"
	datetimeFormat  = "2006-01-02 15:04:05"
)

// EmployeeCard - Карточка сотрудника
type EmployeeCard struct {
	ID         int64  `json:"Id"`
	FirstName  string `json:"FirstName"`
	MiddleName string `json:"MiddleName"`
	LastName   string `json:"LastName"`
	Department struct {
		ID   int64  `json:"Id"`
		Name string `json:"Name"`
	} `json:"Department"`
	Position struct {
		ID   int64  `json:"Id"`
		Name string `json:"Name"`
	} `json:"Position"`
	Login     string `json:"Login"`
	Email     string `json:"Email"`
	FireDay   string `json:"FireDay"`
	Behaviour string `json:"Behaviour"`
}

// EmployeeList - Список сотрудников
type EmployeeList []EmployeeCard

// Scan - Парсинг json
func (el *EmployeeList) Scan(b []byte) error {
	var r responseEmployeeList
	json.Unmarshal(b, &r)
	if err := r.IFerror(); err != nil {
		return err
	}
	if data, ok := r.Data["employees"]; ok {
		*el = data
	}
	return nil
}

// Scan - Парсинг json
func (u *EmployeeCard) Scan(b []byte) error {
	var r userResponse
	json.Unmarshal(b, &r)
	if err := r.IFerror(); err != nil {
		return err
	}
	if data, ok := r.Data["employee"]; ok {
		*u = data
	}
	return nil
}

// TaskCard - Карточка задачи
type TaskCard struct {
	ID       int64  `json:"Id"`
	Name     string `json:"Name"`
	Status   string `json:"Status"`
	Deadline string `json:"Deadline"`
	Owner    struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Owner"`
	Responsible struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Responsible"`
	TimeCreated string `json:"TimeCreated"`
	TimeUpdated string `json:"TimeUpdated"`
	Finish      string `json:"Finish"`
	Tags        []struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Tags"`
	Statement string `json:"Statement"`
	Auditors  []struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Auditors"`
	PlannedWork            int64  `json:"PlannedWork"`
	PlannedFinish          string `json:"PlannedFinish"`
	ActualWork             int64  `json:"ActualWork"`
	ActualWorkWithSubTasks int64  `json:"ActualWorkWithSubTasks"`
	IsOverdue              bool   `json:"IsOverdue"`
}

// TaskList - Список задач
type TaskList []TaskCard

// Scan - парсинг json
func (tc *TaskList) Scan(b []byte) error {
	var r taskListResponse
	json.Unmarshal(b, &r)
	if err := r.IFerror(); err != nil {
		return err
	}
	if data, ok := r.Data["tasks"]; ok {
		*tc = data
	}
	return nil
}

// Tag - метка
type Tag struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

// TagsList - Список меток
type TagsList []Tag

// Scan - парсинг json
func (tl *TagsList) Scan(b []byte) error {
	var r tagListResponse
	json.Unmarshal(b, &r)
	if err := r.IFerror(); err != nil {
		return err
	}
	if data, ok := r.Data["tags"]; ok {
		*tl = data
	}
	return nil
}

// Comment - комментарий
type Comment struct {
	ID          int64        `json:"ID"`
	Text        string       `json:"Text"`
	Work        int64        `json:"Work"`
	WorkDate    string       `json:"WorkDate"`
	TimeCreated string       `json:"TimeCreated"`
	Author      EmployeeCard `json:"Author"`
}

// CommentsList - список комментариев
type CommentsList []Comment

// Scan - парсинг json
func (ct *CommentsList) Scan(b []byte) error {
	var r commentListResponse
	json.Unmarshal(b, &r)
	if err := r.IFerror(); err != nil {
		return err
	}
	if data, ok := r.Data["comments"]; ok {
		*ct = data
	}
	return nil
}

// UserAppVerification - структура ответа при валидации пользователя во встроенного приложения
type UserAppVerification struct {
	UserID   string `json:"id"`
	FullName string `json:"name"`
	Position string `json:"position"`
}
