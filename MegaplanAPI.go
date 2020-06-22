package megaplang

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// NewClien - инициализация новго экземпляра MegaplanAPI
func NewClien(domain string) *MegaplanAPI {
	var cnf = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("https://%s/api/v3/auth/access_token", domain),
		},
	}
	return &MegaplanAPI{cnf: &cnf}
}

func (mp *MegaplanAPI) setClient() {
	mp.Client = oauth2.NewClient(context.Background(), mp.ts)
}

// MegaplanAPI - клиент для работы с мегаплан v3, обертка над oauth2
type MegaplanAPI struct {
	cnf *oauth2.Config
	ts  oauth2.TokenSource
	*http.Client
}

// CheckCredential - проверить сохраненный файл токена
func (mp *MegaplanAPI) CheckCredential(tokenfile string) error {
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
func (mp *MegaplanAPI) GetNewToken(username, password, tokenfile string) error {
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
