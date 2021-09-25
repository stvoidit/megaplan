# megaplan

Смотри примеры в examples

## master

https://dev.megaplan.ru/r1905/api/index.html

В данный момент поддерживает v1 API.

Представляет простую обертку над http методами GET и POST
Алгоритм шифрования запроса см. в методе __queryHashing__, может быть реализован на любом ЯП.

    go get github.com/stvoidit/megaplan

Пример использования:

    var api = megaplan.NewAPI(accessID, secretKey, myhost)
    response, err := api.GET("/BumsCommonApiV01/UserInfo/id.api", nil)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println(response.Status)
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


## v3

https://dev.megaplan.ru/r1905/apiv3/index.html

Уже может использоваться, но не имеет кастомизации возможности сохранения токена.
Частично в процессе доработки удобства использования.

Представляет собой обертку над [oauth2](https://godoc.org/golang.org/x/oauth2).

    go get github.com/stvoidit/megaplan/v3

## Примечание

v1 и v3 координально отличаются по схемам данных, обработке и содержанию.
Многие сущности не описаны в [v3 документации](https://demo.megaplan.ru/api/v3/docs).
Например имеется endpoint на __/api/v3/department__, но при этом сущность сотрудников вообще не имеет данных об отделах.
Имеется endpoint __/api/v3/position__, но не описан в докуменатации, хотя представляет структурированные данные о должностях сотрудников.
