package main

import (
	"io"
	"megaplang/v3"
	"os"
)

func main() {
	api := megaplang.NewClien(`mymegaplan.ru`)
	if err := api.CheckCredential("megaplan-token.json"); err != nil {
		if err := api.GetNewToken(`username@email.ru`, `password`, `megaplan-token.json`); err != nil {
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
