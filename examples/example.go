package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/stvoidit/megaplan"
)

const (
	myhost    = "https://example.ru"
	login     = "login@example.ru"
	password  = "password"
	accessID  = "accessId"
	secretKey = "secretKey"
)

// Получение токена - accessID и secretKey
// если уже известны и есть в конфиге - не требуется выполнять
func getToken() (api *megaplan.API) {
	if err := api.GetToken(myhost, login, password); err != nil {
		panic(err)
	}
	return
}

// пример раоты с GET запросом
func exampleGET(api *megaplan.API) {
	payload := megaplan.Payload{"Id": 1000005}
	response, err := api.GET("/BumsStaffApiV01/Employee/card.api", payload)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println(response.Status)
	w, _ := os.Create("ResponseEmployeeCard.json")
	defer w.Close()
	io.Copy(w, response.Body)
}

// пример работы с POST запросом
func examplePOST(api *megaplan.API) {
	payload := megaplan.Payload{
		"SubjectType": "task",
		"SubjectId":   1017226,
		"Model[Text]": "test api v1 message"}
	response, err := api.POST("/BumsCommonApiV01/Comment/create.api", payload)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println(response.Status)
	w, _ := os.Create("ResponseCommentCreate.json")
	defer w.Close()
	io.Copy(w, response.Body)
}

// пример удобной обертки для анмаршалинга нужных структур данных
func exampleParse(api *megaplan.API) {
	response, err := api.GET("/BumsCommonApiV01/UserInfo/id.api", nil)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println(response.Status)
	// параллельная запись ответа в файл
	w, _ := os.Create("ResponseUserInfoe.json")
	defer w.Close()
	tee := io.TeeReader(response.Body, w)

	// UserInfo - неполная модель "UserInfo"
	type UserInfo struct {
		UserID       int64  `json:"UserId"`
		EmployeeID   int64  `json:"EmployeeId"`
		ContractorID string `json:"ContractorId"`
	}
	var user = new(UserInfo)
	if err := json.NewDecoder(tee).Decode(megaplan.ExpectedResponse(user)); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)
}

func main() {
	{
		// получение токена, возвращается экземпляр API
		var api = getToken()
		fmt.Printf("%+v\n", api)
		api.SaveToken("token.json")
	}
	{
		// примеры использования
		var api = megaplan.NewAPI(accessID, secretKey, myhost)
		exampleGET(api)
		examplePOST(api)
		exampleParse(api)
	}
}
