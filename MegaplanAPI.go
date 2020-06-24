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
func NewClien(domain, username, password string, token *oauth2.Token) (api *APImegaplan, err error) {
	var cnf = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("https://%s/api/v3/auth/access_token", domain),
		},
	}
	api = &APImegaplan{cnf: &cnf, domain: domain}
	if token != nil {
		err = api.checkCredential(token)
	} else {
		err = api.getNewToken(username, password)
	}
	api.Client = oauth2.NewClient(oauth2.NoContext, api.ts)
	return
}

// APImegaplan - клиент для работы с мегаплан v3, обертка над oauth2
type APImegaplan struct {
	domain string
	cnf    *oauth2.Config
	ts     oauth2.TokenSource
	*http.Client
}

// Token - вернуть актуальный токен
func (mp *APImegaplan) Token() (*oauth2.Token, error) {
	return mp.ts.Token()
}

// CheckCredential - проверить сохраненный файл токена
func (mp *APImegaplan) checkCredential(token *oauth2.Token) error {
	mp.ts = mp.cnf.TokenSource(oauth2.NoContext, token)
	return nil
}

// GetNewToken - получить новый ключ, сохранить и применить в текущем экземпляре API
func (mp *APImegaplan) getNewToken(username, password string) error {
	t, err := mp.cnf.PasswordCredentialsToken(oauth2.NoContext, username, password)
	if err != nil {
		return err
	}
	mp.ts = mp.cnf.TokenSource(context.Background(), t)
	return nil
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

// SaveToken - сохранить токен в файл
func SaveToken(t *oauth2.Token, filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(t)
}

// LoadTokenFromFile - загрузить токен из файла
func LoadTokenFromFile(filename string) (*oauth2.Token, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	t := new(oauth2.Token)
	if err = json.NewDecoder(r).Decode(t); err != nil {
		return nil, err
	}
	return t, nil
}
