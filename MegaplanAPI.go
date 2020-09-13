package megaplan

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Response - структура стандартного ответа API
type Response struct {
	Status struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Data interface{} `json:"data"`
}

// ExpectedResponse - оборачивает ожидаемый ответ в стандартную структуру.
// Ожидаемый интерфейс будет находиться в поле Response.Data.
// После обработки необходимо сделать assert вложенного интерфейса к ожидаемому (см. примеры)
func ExpectedResponse(data interface{}) *Response {
	return &Response{Data: data}
}

// API - Структура объекта API v1
type API struct {
	accessID  string
	secretKey []byte
	domain    string
	appUUID   string
	appSecret []byte
	client    *http.Client
}

// SaveToken - сохранить конфигурацию в json
func (api API) SaveToken(filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(map[string]string{
		"accessID":  api.accessID,
		"secretKey": string(api.secretKey),
		"domain":    api.domain,
		"appUUID":   api.appUUID,
		"appSecret": string(api.appSecret)})
}

// accessToken - структура ответа с токеном доступа
type accessToken struct {
	Response
	Data struct {
		UserID       int    `json:"UserId"`
		EmployeeID   int    `json:"EmployeeId"`
		ContractorID string `json:"ContractorId"`
		AccessID     string `json:"AccessId"`
		SecretKey    string `json:"SecretKey"`
	} `json:"data"`
}

type otcData struct {
	Response
	Data struct {
		OneTimeKey string `json:"OneTimeKey"`
	} `json:"data"`
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
	return hex.EncodeToString(hashPassword.Sum(nil))
}

// getOTC - получение временного ключа
func getOTC(domain string, login string, md5password string) (OneTimeKey string, err error) {
	const uriOTC = "/BumsCommonApiV01/User/createOneTimeKeyAuth.api"
	var payload = url.Values{"Login": {login}, "Password": {md5password}}
	req, err := http.NewRequest(http.MethodPost, domain+uriOTC, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var otc otcData
	if err := json.NewDecoder(resp.Body).Decode(&otc); err != nil {
		return "", err
	}
	if otc.Response.Status.Code == "error" {
		return "", errors.New(otc.Response.Status.Message)
	}
	return otc.Data.OneTimeKey, nil
}

// getToken - AccessId, SecretKey
func getToken(domain string, login string, md5password string, otc string) (AccessID string, SecretKey string, err error) {
	const uriToken = "/BumsCommonApiV01/User/authorize.api"
	var payload = url.Values{"Login": {login}, "Password": {md5password}, "OneTimeKey": {otc}}
	req, err := http.NewRequest(http.MethodPost, domain+uriToken, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var token accessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", "", err
	}
	if token.Response.Status.Code == "error" {
		return "", "", errors.New(token.Response.Status.Message)
	}
	return token.Data.AccessID, token.Data.SecretKey, nil
}

// GetToken - Получение токена API
func (api *API) GetToken(domain, login, password string) (err error) {
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
	return
}
