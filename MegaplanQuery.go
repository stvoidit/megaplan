package megaplan

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"

// GET - get запрос к API
func (api *API) GET(uri string, payload interface{}) []byte {
	const rMethod = "GET"
	urlQuery, queryHeader := api.queryHasher(rMethod, uri, payload)
	return api.requestQuery(rMethod, urlQuery, queryHeader)
}

// POST - post запрос на API
func (api *API) POST(uri string, payload interface{}) []byte {
	const rMethod = "POST"
	urlQuery, queryHeader := api.queryHasher(rMethod, uri, payload)
	return api.requestQuery(rMethod, urlQuery, queryHeader)
}

// CheckUser - проверка пользователя для встроенного приложения
func (api *API) CheckUser(userSign string) (UserAppVerification, error) {
	var appAPI API
	appAPI = *api
	appAPI.AccessID = api.AppUUID
	appAPI.SecretKey = api.AppSecret
	var payload = map[string]string{
		"uuid":     api.AppUUID,
		"userSign": userSign,
	}
	b := appAPI.POST("/BumsSettingsApiV01/Application/checkUserSign.json", payload)
	var i UserVerifyResponse
	if err := json.Unmarshal(b, &i); err != nil {
		return i.Data, err
	}
	if i.response.Status.Code != "ok" {
		return i.Data, errors.New("invalide user data")
	}
	return i.Data, nil
}

// queryHasher - задаем сигнатуру, отдает URL и Header для запросов к API
func (api *API) queryHasher(method string, uri string, payload interface{}) (url.URL, http.Header) {
	var normalizePayload = make(map[string]string)
	switch p := payload.(type) {
	case map[string]string:
		normalizePayload = p
	case map[string]int:
		for k, v := range p {
			normalizePayload[k] = strconv.Itoa(v)
		}
	case map[string]int64:
		for k, v := range p {
			normalizePayload[k] = strconv.FormatInt(v, 10)
		}
	case map[string]bool:
		for k, v := range p {
			normalizePayload[k] = strconv.FormatBool(v)
		}
	case nil:
		break
	default:
		log.Fatalln(errors.New("cant parse paylaod interface"))
	}
	URL, err := url.Parse(api.Domain)
	if err != nil {
		panic(err.Error())
	}
	URL.Path += uri
	today := time.Now().Format(rfc2822)
	urlParams := url.Values{}
	for k, v := range normalizePayload {
		urlParams.Add(k, v)
	}
	if len(urlParams) > 0 {
		URL.RawQuery = urlParams.Encode()
	}
	sigURL := strings.Replace(URL.String(), fmt.Sprintf("%s://", URL.Scheme), "", 1)
	Signature := fmt.Sprintf("%s\n\napplication/x-www-form-urlencoded\n%s\n%s", method, today, sigURL)
	h := hmac.New(sha1.New, api.SecretKey)
	h.Write([]byte(Signature))
	hexSha1 := hex.EncodeToString(h.Sum(nil))
	sha1Query := base64.StdEncoding.EncodeToString([]byte(hexSha1))
	queryHeader := http.Header{
		"Date":            []string{today},
		"X-Authorization": []string{api.AccessID + ":" + sha1Query},
		"Accept":          []string{"application/json"},
		"Content-Type":    []string{"application/x-www-form-urlencoded"},
		"accept-encoding": []string{"gzip, deflate, br"},
	}
	return *URL, queryHeader
}

// requestQuery - итоговый запрос к API предварительно сформированный Request с правильным набором headers
func (api *API) requestQuery(method string, url url.URL, headers http.Header) []byte {
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header = headers
	resp, err := api.client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	return body
}
