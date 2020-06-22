# megaplan

Смотри примеры в examples

## master

https://dev.megaplan.ru/r1905/api/index.html

В данный момент поддерживает v1 API.

___не рекомендуется использоваться метод Scan. Методы запроса GET и POST могут возвращать буфер.___

    go get github.com/stvoidit/megaplan

## v3

https://dev.megaplan.ru/r1905/apiv3/index.html

Уже может использоваться, но не имеет кастомизации возможности сохранения токена.

Представляет собой обертку над [oauth2](https://godoc.org/golang.org/x/oauth2).

    go get github.com/stvoidit/megaplan/v3
