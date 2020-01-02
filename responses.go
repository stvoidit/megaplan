package megaplan

import "errors"

// Дефолтные структуры ответов от API

type response struct {
	Status struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

func (r *response) IFerror() error {
	if r.Status.Code == "error" {
		return errors.New(r.Status.Message)
	}
	return nil
}

type userResponse struct {
	response
	Data map[string]EmployeeCard `json:"data"`
}
type responseEmployeeList struct {
	response
	Data map[string][]EmployeeCard `json:"data"`
}
type taskListResponse struct {
	response
	Data map[string][]TaskCard `json:"data"`
}
type tagListResponse struct {
	response
	Data map[string][]Tag `json:"data"`
}
type commentListResponse struct {
	response
	Data map[string]CommentsList `json:"data"`
}

// UserVerifyResponse - тип для верификации юзеров во встроенном приложении
type UserVerifyResponse struct {
	response
	Data UserAppVerification `json:"data"`
}
