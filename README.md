# megaplan

Обертка над oauth2 в стандартной бибилотеке.

Авторефреш токена, сохранение в указанный файл и т.п.
 
Возвращает обычный "http. Response" для дальнейшей обработки чем угодно. Прямой доступ к "*http. Client" для кастомизации.

    import (
        "io"
        "github.com/stvoidit/megaplan/v3"
        "os"
    )

    func main() {
        api := megaplang.NewClien( `mymegaplan.ru` )
        if err := api.CheckCredential("megaplan-token.json"); err != nil {
            if err := api.GetNewToken( `username@email.ru` , `password` , `megaplan-token.json` ); err != nil {
                panic(err)
            }
        }
        resp, err := api.Get("https://mymegaplan.ru/api/v3/deal/7520")
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
        io.Copy(os.Stdout, resp.Body)
    }

### Инициализация

#### Вариант 1

Если имеется файл токена, либо вы можете передать как env аргументы и создать новый экземпляр *oauth2. Token

    t, _ := megaplan.LoadTokenFromFile( `megaplan-token.json` )
    api, err := megaplan.NewClien( `mymegaplan.ru` , `username@email.ru` , `password` , t)
    if err != nil {
        panic(err)
    }
    defer func() {
        if t, err := api.Token(); err == nil && t != nil {
            megaplan.SaveToken(t, `megaplan-token.json` )
        }
    }()

#### Вариант 2

Если данных токена еще нет. Будет автоматически создан новый токен.

    api, err := megaplan.NewClien( `mymegaplan.ru` , `username@email.ru` , `password` , nil)
    if err != nil {
        panic(err)
    }
    defer func() {
        if t, err := api.Token(); err == nil && t != nil {
            megaplan.SaveToken(t, `megaplan-token.json` )
        }
    }()

### Загрузка файлов

Всегада возвращает сущность в виде, который нужен для отправки в дальнейшем POST запросе.

    ...

    r1, err := os.Open( `...\test\myfile1.xlsx` )
    if err != nil {
        panic(err)
    }
    r2, err := os.Open( `...\test\myfile2.jpg` )
    if err != nil {
        panic(err)
    }
    defer r1.Close()
    defer r2.Close()
    attchs, err := api.UploadFiles(
        megaplan.UploadFile{Filename: `myfile1.xlsx` , R: r1},
        megaplan.UploadFile{Filename: `myfile2.jpg` , R: r2})
    if err != nil {
        panic(err)
    }
