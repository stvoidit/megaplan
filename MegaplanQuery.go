package megaplan

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"

// Payload - параметры запроса. url.Values в удобной обертке
type Payload map[string]interface{}

// Encode = url.Values.Encode()
func (p Payload) Encode() string {
	var urlParams = make(url.Values)
	for k, val := range p {
		switch t := val.(type) {
		case int:
			urlParams.Add(k, strconv.FormatInt(int64(t), 10))
		case int8:
			urlParams.Add(k, strconv.FormatInt(int64(t), 10))
		case int16:
			urlParams.Add(k, strconv.FormatInt(int64(t), 10))
		case int32:
			urlParams.Add(k, strconv.FormatInt(int64(t), 10))
		case int64:
			urlParams.Add(k, strconv.FormatInt(int64(t), 10))
		case uint:
			urlParams.Add(k, strconv.FormatUint(uint64(t), 10))
		case uint8:
			urlParams.Add(k, strconv.FormatUint(uint64(t), 10))
		case uint16:
			urlParams.Add(k, strconv.FormatUint(uint64(t), 10))
		case uint32:
			urlParams.Add(k, strconv.FormatUint(uint64(t), 10))
		case uint64:
			urlParams.Add(k, strconv.FormatUint(uint64(t), 10))
		case float64:
			urlParams.Add(k, strconv.FormatFloat(t, 'f', 2, 64))
		case float32:
			urlParams.Add(k, strconv.FormatFloat(float64(t), 'f', 2, 64))
		case bool:
			urlParams.Add(k, strconv.FormatBool(t))
		case string:
			urlParams.Add(k, t)
		case nil:
			continue
		}
	}
	return urlParams.Encode()
}

// GET - get запрос к API
func (api API) GET(uri string, payload Payload) (*http.Response, error) {
	request, err := api.queryHashing(http.MethodGet, uri, payload)
	if err != nil {
		return nil, err
	}
	return api.Do(request)
}

// POST - post запрос на API
func (api API) POST(uri string, payload Payload) (*http.Response, error) {
	request, err := api.queryHashing(http.MethodPost, uri, payload)
	if err != nil {
		return nil, err
	}
	return api.Do(request)
}

// CheckUser - проверка пользователя для встроенного приложения
func (api API) CheckUser(userSign string) (*http.Response, error) {
	var payload = Payload{
		"uuid":     api.accessID,
		"userSign": userSign,
	}
	return api.POST("/BumsSettingsApiV01/Application/checkUserSign.json", payload)
}

func (api API) createSignatureSign(r *http.Request, today string) (string, error) {
	// ! специально кодируем запрос для v1 API для заголовка X-Authorization
	// ! та самая безумная структура "защиты" из-за которой возникают сложности
	// ! может быть реализована на любом ЯП
	// ! https://dev.megaplan.ru/api/API_authorization.html#id14

	var Signature = bytes.NewBuffer(make([]byte, 0, 256))
	Signature.WriteString(r.Method)
	Signature.WriteString("\n\napplication/x-www-form-urlencoded\n")
	Signature.WriteString(today)
	Signature.WriteString("\n")
	Signature.WriteString(r.URL.Host)
	Signature.WriteString(r.URL.RequestURI())
	defer Signature.Reset()
	h := hmac.New(sha1.New, api.secretKey)
	if _, err := Signature.WriteTo(h); err != nil {
		return "", err
	}
	var sig = h.Sum(nil)
	var hexbuff = make([]byte, hex.EncodedLen(len(sig)))
	hex.Encode(hexbuff, sig)
	return base64.StdEncoding.EncodeToString(hexbuff), nil
}

// queryHasher - задаем сигнатуру, отдает URL и Header для запросов к API, создаем http.Request
func (api API) queryHashing(method string, uri string, payload Payload) (request *http.Request, err error) {
	switch method {
	case http.MethodPost:
		if request, err = http.NewRequest(method, api.domain, strings.NewReader(payload.Encode())); err != nil {
			return nil, err
		}
	case http.MethodGet:
		if request, err = http.NewRequest(method, api.domain, nil); err != nil {
			return nil, err
		}
		if len(payload) > 0 {
			request.URL.RawQuery = payload.Encode()
		}
	default:
		return nil, errors.New("unavailable http method")
	}
	request.URL.Path = uri
	today := time.Now().Format(rfc2822)
	signature, err := api.createSignatureSign(request, today)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Date", today)
	request.Header.Set("X-Authorization", api.accessID+":"+signature)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return
}

func unzipResponse(response *http.Response) (err error) {
	if response.Uncompressed {
		return nil
	}
	if response.Header.Get("Content-Encoding") == "gzip" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := response.Body.Close(); err != nil {
			return err
		}
		response.Body, err = gzip.NewReader(bytes.NewReader(body))
	}
	return
}

// Do - обертка над стандартным Do(*http.Request)
func (api API) Do(request *http.Request) (response *http.Response, err error) {
	if api.enablegzip {
		request.Header.Set("Accept-Encoding", "gzip")
	}
	response, err = api.client.Do(request)
	if err != nil {
		return response, err
	}
	err = unzipResponse(response)
	return
}
