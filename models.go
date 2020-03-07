package megaplan

const (
	shortTimeFormat = "15:04"
	timeFormat      = "15:04:05"
	datetimeFormat  = "2006-01-02 15:04:05"
)

// EmployeeCard - Карточка сотрудника
type EmployeeCard struct {
	ID         uint   `json:"Id"`
	FirstName  string `json:"FirstName"`
	MiddleName string `json:"MiddleName"`
	LastName   string `json:"LastName"`
	Department struct {
		ID   uint   `json:"Id"`
		Name string `json:"Name"`
	} `json:"Department"`
	Position struct {
		ID   uint   `json:"Id"`
		Name string `json:"Name"`
	} `json:"Position"`
	Login     string `json:"Login"`
	Email     string `json:"Email"`
	FireDay   string `json:"FireDay"`
	Behaviour string `json:"Behaviour"`
}

// EmployeeList - Список сотрудников
type EmployeeList []EmployeeCard

// TaskCard - Карточка задачи
type TaskCard struct {
	ID       uint   `json:"Id"`
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
	Tags        []Tag  `json:"Tags"`
	Statement   string `json:"Statement"`
	Auditors    []struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Auditors"`
	PlannedWork            string `json:"PlannedWork"`
	PlannedFinish          string `json:"PlannedFinish"`
	ActualWork             string `json:"ActualWork"`
	ActualWorkWithSubTasks string `json:"ActualWorkWithSubTasks"`
	IsOverdue              bool   `json:"IsOverdue"`
}

// TaskList - Список задач
type TaskList []TaskCard

// Tag - метка
type Tag struct {
	ID   uint   `json:"Id"`
	Name string `json:"Name"`
}

// TagsList - Список меток
type TagsList []Tag

// Comment - комментарий
type Comment struct {
	ID          uint         `json:"ID"`
	Text        string       `json:"Text"`
	Work        uint         `json:"Work"`
	WorkDate    string       `json:"WorkDate"`
	TimeCreated string       `json:"TimeCreated"`
	Author      EmployeeCard `json:"Author"`
}

// CommentsList - список комментариев
type CommentsList []Comment

// UserAppVerification - структура ответа при валидации пользователя во встроенного приложения
type UserAppVerification struct {
	UserID   string `json:"id"`
	FullName string `json:"name"`
	Position string `json:"position"`
}
