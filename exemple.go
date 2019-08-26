package main

import (
	"fmt"

	"./megaplan"
)

// Token - тип токена, сюда вписываем полученный ключ
type Token struct {
	AcessID  string `json:"AccessId"`
	SecreKey string `json:"SecretKey"`
}

// Task - пример как вернуть данные в структуру из json response с API
// Передаем тип в функцию через &, возвращаем туда данные
type Task struct {
	Status struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Data struct {
		Task struct {
			Name        string `json:"Name"`
			Statement   string `json:"Statement"`
			Status      string `json:"Status"`
			TimeCreated string `json:"TimeCreated"`
		} `json:"task"`
	} `json:"data"`
}

func main() {
	const Login = "LOGIN"
	const Password = "PASSWORD"
	const Domain = "https://mydomain.ru"
	const AcessID = "some_access_id"
	const SecretKey = "some_secret_key"

	testToken(Login, Password, Domain)   // Получение AccessID и SecretKey (см. тип Token)
	testGET(AcessID, SecretKey, Domain)  // Запрашиваем данные по задаче, как пример
	testPOST(AcessID, SecretKey, Domain) // Меняем описание у задачи, как пример

}

// testToken - возвращает 2 строки, которые можно вписать в свой тип, как указано в примере,
// либо переписать и возвращать весь тип ответа при удачном получении ключей
func testToken(login string, password string, domain string) {
	a, s, err := megaplan.Token(login, password, domain)
	if err != nil {
		panic(err.Error())
	}
	MyToken := Token{AcessID: a, SecreKey: s}
	fmt.Println(a, s)
	fmt.Println()
	fmt.Println(MyToken)
}

// testGET - в последний аргумент передается ссылка на свой тип,
// в который должны быть записаны данные из JSON
func testGET(a string, s string, h string) {
	pyload := map[string]string{
		"Id":             "1000001",
		"ExtraFields[0]": "Category130CustomFieldRezultat",
	}
	uri := "/BumsTaskApiV01/Task/card.api"
	taskData := Task{}
	res := megaplan.GET(h, a, s, uri, pyload, &taskData)
	fmt.Println(res)
	fmt.Println()
	fmt.Println(taskData)
}

// testPOST - можно так же создать свой тип и передать в функцию POST последним аргументом
func testPOST(a string, s string, h string) {
	pyload := map[string]string{
		"Id":               "1000001",
		"Model[Statement]": "testGO1",
	}
	uri := "/BumsTaskApiV01/Task/edit.api"
	res := megaplan.POST(h, a, s, uri, pyload, nil)
	fmt.Println(res)
}
