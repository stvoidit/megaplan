package main

import (
	"fmt"
	"net/http"
	"time"

	mp "github.com/stvoidit/MegaplanGO"
)

var api mp.API

// Инициаолизация экземпляра API
func init() {
	tr := http.Transport{IdleConnTimeout: 1 * time.Minute}
	c := &http.Client{Timeout: 1 * time.Minute, Transport: &tr}
	api.ParseConfig("config.yaml", c)
}

// Получение токена - accessID и secretKey
// если уже известны и есть в конфиге - не требуется выполнять
func getToken() {
	api.GetToken()
	fmt.Printf("%+v\n", api)

}

// Если модель данных не заложена в коде, то можно просто вернуть bytes
// и применить свой алгоритм матчинга ответа
func getRaw() {
	b := api.GET("/BumsCommonApiV01/UserInfo/id.api", nil)
	fmt.Println(string(b))
}

// у некоторых моделей есть метод Scan, который меняет ответ API
// в более читабельную структуру
func getOnModel1() {
	payload := map[string]int{"Id": 1000005}
	var emp mp.EmployeeCard
	b := api.GET("/BumsStaffApiV01/Employee/card.api", payload)
	emp.Scan(b)
	fmt.Printf("%+v\n", emp)
}

// второй пример - список задач. Так же имеет метод Scan,
// чтобы разложить ответ по структурам типа Task
func getOnModel2() {
	payload := map[string]int{"EmployeeId": 1000005}
	var empы mp.TaskList
	b := api.GET("/BumsTaskApiV01/Task/list.api", payload)
	empы.Scan(b)
	fmt.Printf("%+v\n", empы)
}

func main() {
	getToken()
	getRaw()
	getOnModel1()
	getOnModel2()
}
