package megaplan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// NewClien - инициализация новго экземпляра MegaplanAPI
func NewClien(domain string) *APImegaplan {
	var cnf = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("https://%s/api/v3/auth/access_token", domain),
		},
	}
	return &APImegaplan{cnf: &cnf, domain: domain}
}

func (mp *APImegaplan) setClient() {
	mp.Client = oauth2.NewClient(context.Background(), mp.ts)
}

// APImegaplan - клиент для работы с мегаплан v3, обертка над oauth2
type APImegaplan struct {
	domain string
	cnf    *oauth2.Config
	ts     oauth2.TokenSource
	*http.Client
}

// CheckCredential - проверить сохраненный файл токена
func (mp *APImegaplan) CheckCredential(tokenfile string) error {
	r, err := os.Open(tokenfile)
	if err != nil {
		return err
	}
	defer r.Close()
	t := new(oauth2.Token)
	if err = json.NewDecoder(r).Decode(t); err != nil {
		return err
	}
	mp.ts = mp.cnf.TokenSource(context.Background(), t)
	if !t.Valid() {
		fmt.Println("refresh token")
		t, err := mp.ts.Token()
		if err != nil {
			return err
		}
		if err := saveToken(t, tokenfile); err != nil {
			return err
		}
	}
	mp.setClient()
	return nil
}

// GetNewToken - получить новый ключ, сохранить и применить в текущем экземпляре API
func (mp *APImegaplan) GetNewToken(username, password, tokenfile string) error {
	t, err := mp.cnf.PasswordCredentialsToken(context.Background(), username, password)
	if err != nil {
		return err
	}
	saveToken(t, tokenfile)
	mp.ts = mp.cnf.TokenSource(context.Background(), t)
	mp.setClient()
	fmt.Println("create new token file")
	return nil
}

func saveToken(t *oauth2.Token, filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(t)
}

// UploadFiles - загрузка файла, в ответ приходит список из объектов,
// которые могут быть переданы в json запросе как сущности File для содания вложения
func (mp *APImegaplan) UploadFiles(files ...UploadFile) (attachments []Attachment, err error) {
	if len(files) == 0 {
		return nil, nil
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for _, file := range files {
		fw, err := mw.CreateFormFile("files[]", file.Filename)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fw, file.R); err != nil {
			return nil, err
		}
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	resp, err := mp.Post(fmt.Sprintf("https://%s/api/file", mp.domain), mw.FormDataContentType(), &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var ats metaAttachments
	if err := json.NewDecoder(resp.Body).Decode(&ats); err != nil {
		return nil, err
	}
	return ats.Data, ats.Error()
}

// UploadFile - файл для загрузки
type UploadFile struct {
	Filename string
	R        io.Reader
}

// Attachment - вложение
type Attachment struct {
	ContentType string `json:"contentType"`
	ID          string `json:"id"`
}

type metaAttachments struct {
	MetaInfo `json:"meta"`
	Data     []Attachment `json:"data"`
}

// MetaInfo - meta из ответа API
type MetaInfo struct {
	Status uint64   `json:"status"`
	Errors []string `json:"errors"`
}

func (mi *MetaInfo) Error() error {
	if len(mi.Errors) == 0 {
		return nil
	}
	var msgs string
	for _, errMsg := range mi.Errors {
		msgs += errMsg
	}
	return errors.New(msgs)
}
