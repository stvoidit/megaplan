package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/stvoidit/megaplan/v3"
)

func main() {
	t, _ := megaplan.LoadTokenFromFile(`megaplan-token.json`)
	api, err := megaplan.NewClien(`mymegaplan.ru`, `username@email.ru`, `password`, t)
	if err != nil {
		panic(err)
	}
	defer func() {
		if t, err := api.Token(); err == nil && t != nil {
			megaplan.SaveToken(t, `megaplan-token.json`)
		}
	}()
	resp, err := api.Get("https://mymegaplan.ru/api/v3/deal/1011111")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	// createTask(api)
}

// Пример создания задачи и загрузки файлов в описание.
// Если нужно закреплять определенные файлы для определенных полей,
// то нужно вызывать UploadFiles для нужных вам файлов и сохранять результат
// для дальнейшего закрепления объекта в нужное поле
func createTask(api *megaplan.APImegaplan) {
	r1, err := os.Open(`...\test\myfile1.xlsx`)
	if err != nil {
		panic(err)
	}
	r2, err := os.Open(`...\test\myfile2.jpg`)
	if err != nil {
		panic(err)
	}
	defer r1.Close()
	defer r2.Close()
	attchs, err := api.UploadFiles(
		megaplan.UploadFile{Filename: `myfile1.xlsx`, R: r1},
		megaplan.UploadFile{Filename: `myfile2.jpg`, R: r2})
	if err != nil {
		panic(err)
	}
	task := map[string]interface{}{
		"isUrgent":   false,
		"name":       "test v3",
		"isTemplate": false,
		"responsible": map[string]interface{}{
			"contentType": "Employee",
			"id":          1000001,
		},
		"owner": map[string]interface{}{
			"contentType": "Employee",
			"id":          1000001,
		},
		"subject":  "some text",
		"attaches": attchs,
	}
	b, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	resp, err := api.Post("https://mymegaplan/api/v3/task/1013973", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	w, _ := os.Create(`response.json`)
	defer w.Close()
	fmt.Println(resp.Status)
	io.Copy(w, resp.Body)
}
