package megaplan

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"

// GET - get запрос к API
func (api API) GET(uri string, payload map[string]interface{}) (*http.Response, error) {
	request, err := api.queryHashing(http.MethodGet, uri, payload)
	if err != nil {
		return nil, err
	}
	return api.Do(request)
}

// POST - post запрос на API
func (api API) POST(uri string, payload map[string]interface{}) (*http.Response, error) {
	request, err := api.queryHashing(http.MethodPost, uri, payload)
	if err != nil {
		return nil, err
	}
	return api.Do(request)
}

// CheckUser - проверка пользователя для встроенного приложения
func (api API) CheckUser(userSign string) (*http.Response, error) {
	var appAPI API
	appAPI = api
	appAPI.accessID = api.appUUID
	appAPI.secretKey = api.appSecret
	var payload = map[string]interface{}{
		"uuid":     api.appUUID,
		"userSign": userSign,
	}
	return appAPI.POST("/BumsSettingsApiV01/Application/checkUserSign.json", payload)
}

// queryHasher - задаем сигнатуру, отдает URL и Header для запросов к API, создаем http.Request
func (api API) queryHashing(method string, uri string, payload map[string]interface{}) (request *http.Request, err error) {
	var urlParams = make(url.Values)
	for k, val := range payload {
		switch t := val.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
			urlParams.Add(k, fmt.Sprintf("%d", t))
		case float32, float64:
			urlParams.Add(k, fmt.Sprintf("%.2f", t))
		case bool:
			urlParams.Add(k, strconv.FormatBool(t))
		case string:
			urlParams.Add(k, t)
		case nil:
			continue
		default:
			return nil, fmt.Errorf("unrecognized type: %v", t)
		}
	}
	switch method {
	case http.MethodPost:
		if request, err = http.NewRequest(method, api.domain+uri, strings.NewReader(urlParams.Encode())); err != nil {
			return nil, err
		}
	case http.MethodGet:
		if request, err = http.NewRequest(method, api.domain+uri, nil); err != nil {
			return nil, err
		}
		if len(urlParams) > 0 {
			request.URL.RawQuery = urlParams.Encode()
		}
	default:
		return nil, errors.New("unavailable http method")
	}
	// ! специально кодируем запрос для v1 API для заголовка X-Authorization
	// ! та самая безумная структура "защиты" из-за которой возникают сложности
	// ! может быть реализована на любом ЯП
	// ! https://dev.megaplan.ru/api/API_authorization.html#id14
	today := time.Now().Format(rfc2822)
	sigURL := strings.Replace(request.URL.String(), fmt.Sprintf("%s://", request.URL.Scheme), "", 1)
	Signature := fmt.Sprintf("%s\n\napplication/x-www-form-urlencoded\n%s\n%s", method, today, sigURL)
	h := hmac.New(sha1.New, api.secretKey)
	if _, err := h.Write([]byte(Signature)); err != nil {
		return nil, err
	}
	request.Header.Set("Date", today)
	request.Header.Set("X-Authorization", api.accessID+":"+base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil)))))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return
}

// Do - обертка над стандартным Do(*httpRequest)
func (api API) Do(request *http.Request) (response *http.Response, err error) {
	return api.client.Do(request)
}
