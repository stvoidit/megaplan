package main

import (
	"fmt"
	"io"
	"os"

	"github.com/stvoidit/megaplan/v1"
	"gopkg.in/yaml.v2"
)

var api *megaplan.API

// имплементация своего кофига, если имеется иерархия в файле
type costumeConfig struct {
	*megaplan.Config `yaml:"megaplan"`
}

// свой метод для парсинга с учетом иерархии в структуре файла
func (cc *costumeConfig) ReadConfig(r io.Reader) {
	yaml.NewDecoder(r).Decode(cc)
}

// Инициаолизация экземпляра API
func init() {
	r, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	// api.ReadConfig(file)
	var cnf costumeConfig
	cnf.ReadConfig(r)
	fmt.Printf("your costume config:\n%#v\n", cnf.Config)
	api = megaplan.NewWithConfigFile(cnf.Config)
}

// Получение токена - accessID и secretKey
// если уже известны и есть в конфиге - не требуется выполнять
func getToken() {
	if err := api.GetToken("https://exemple.com", "login", "password"); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", api)

}

// Если модель данных не заложена в коде, то можно просто вернуть bytes
// и применить свой алгоритм матчинга ответа к *bytes.Buffer
func getRaw() {
	b := api.GET("/BumsCommonApiV01/UserInfo/id.api", nil)
	fmt.Println(b.String())
}

// у некоторых моделей есть метод Scan, который меняет ответ API
// в более читабельную структуру
func getOnModel1() {
	type customeStructure struct {
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
	var emp customeStructure
	payload := map[string]interface{}{"Id": 1000005}
	if err := api.GET("/BumsStaffApiV01/Employee/card.api", payload).Scan(&emp); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", emp)
}

// второй пример - список задач. Так же имеет метод Scan,
// чтобы разложить ответ по структурам типа Task
func getOnModel2() {
	payload := map[string]interface{}{"EmployeeId": 1000005}
	var emps megaplan.TaskList
	if err := api.GET("/BumsTaskApiV01/Task/list.api", payload).Scan(&emps); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", emps)
}

func main() {
	// getToken()
	getRaw()      // return *bytes.Buffer
	getOnModel1() // Decode in your custome structure
	getOnModel2() // Decode in structure
}
