package megaplan

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// API - Структура объекта API v1
type API struct {
	accessID  string
	secretKey []byte
	domain    string
	appUUID   string
	appSecret []byte
	client    *http.Client
}

// NewAPI - новый экземпляр api
func NewAPI(accessID, secretKey, domain, appUUID, appSecret string) *API {
	return &API{
		client:    http.DefaultClient,
		accessID:  accessID,
		secretKey: []byte(secretKey),
		domain:    domain,
		appUUID:   appUUID,
		appSecret: []byte(appSecret),
	}
}

// SetCustomClient - установить свой http.Client для API
func (api *API) SetCustomClient(c *http.Client) {
	api.client = c
}

// SetEmbeddedApplication - установить ключ от встроенного приложения
func (api *API) SetEmbeddedApplication(appuuid, appsecret string) {
	api.appUUID = appuuid
	api.secretKey = []byte(appsecret)
}

// md5Passord - хэшируем пароль в md5
func md5Passord(p string) string {
	hashPassword := md5.New()
	hashPassword.Write([]byte(p))
	md5String := hex.EncodeToString(hashPassword.Sum(nil))
	return md5String
}

// getOTC - получение временного ключа
func getOTC(domain string, login string, md5password string) (string, error) {
	const uriOTC = "/BumsCommonApiV01/User/createOneTimeKeyAuth.api"
	var payload = url.Values{}
	payload.Add("Login", login)
	payload.Add("Password", md5password)
	req, _ := http.NewRequest("POST", domain+uriOTC, strings.NewReader(payload.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var OTCdata = new(struct {
		response
		Data struct {
			OneTimeKey string `json:"OneTimeKey"`
		} `json:"data"`
	})
	if err := json.NewDecoder(resp.Body).Decode(OTCdata); err != nil {
		return "", err
	}
	if OTCdata.response.Status.Code == "error" {
		return "", errors.New(OTCdata.response.Status.Message)
	}
	return OTCdata.Data.OneTimeKey, nil
}

// getToken - AccessId, SecretKey
func getToken(domain string, login string, md5password string, otc string) (string, string, error) {
	const uriToken = "/BumsCommonApiV01/User/authorize.api"
	var payload = url.Values{}
	payload.Add("Login", login)
	payload.Add("Password", md5password)
	payload.Add("OneTimeKey", otc)
	req, _ := http.NewRequest("POST", domain+uriToken, strings.NewReader(payload.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	AccessToken := new(struct {
		response
		Data struct {
			UserID       int    `json:"UserId"`
			EmployeeID   int    `json:"EmployeeId"`
			ContractorID string `json:"ContractorId"`
			AccessID     string `json:"AccessId"`
			SecretKey    string `json:"SecretKey"`
		} `json:"data"`
	})
	if err := json.NewDecoder(resp.Body).Decode(AccessToken); err != nil {
		return "", "", err
	}
	if AccessToken.response.Status.Code == "error" {
		return "", "", errors.New(AccessToken.response.Status.Message)
	}
	return AccessToken.Data.AccessID, AccessToken.Data.SecretKey, nil
}

// GetToken - Получение токена API
func (api *API) GetToken(domain, login, password string) error {
	md5p := md5Passord(password)
	OTCkey, err := getOTC(domain, login, md5p)
	if err != nil {
		return err
	}
	AID, Skey, err := getToken(domain, login, md5p, OTCkey)
	if err != nil {
		return err
	}
	api.accessID = AID
	api.secretKey = []byte(Skey)
	api.domain = domain
	return nil
}
