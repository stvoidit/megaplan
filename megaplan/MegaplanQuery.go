package megaplan

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// APIresponse - Ответ API, json
type APIresponse struct {
	Status map[string]string
	Data   map[string]interface{}
}

// GET - GET Запрос к API
func GET(domain string, acessid string, secretkey string, uri string, pyload map[string]string) APIresponse {
	urlQuery, queryHeader := queryHasher(domain, acessid, secretkey, "GET", uri, pyload)
	responseAPI := requestQuery(urlQuery, queryHeader)
	return responseAPI
}

// queryHasher - Задает сигнатуру, отдает Header для запросов к API
func queryHasher(d string, a string, s string, r string, uri string, payload map[string]string) (url.URL, http.Header) {
	const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
	URL, _ := url.Parse(d)
	URL.Path += uri
	today := time.Now().Format(rfc2822)
	urlParams := url.Values{}
	for k, v := range payload {
		urlParams.Add(k, v)
	}
	URL.RawQuery = urlParams.Encode()
	Signature := r + "\n\n" + "application/x-www-form-urlencoded" + "\n" + today + "\n" + URL.Host + URL.Path + "?" + URL.RawQuery
	h := hmac.New(sha1.New, []byte(s))
	h.Write([]byte(Signature))
	hexSha1 := hex.EncodeToString(h.Sum(nil))
	sha1Query := base64.StdEncoding.EncodeToString([]byte(hexSha1))
	queryHeader := http.Header{
		"Date":            []string{today},
		"X-Authorization": []string{a + ":" + sha1Query},
		"Accept":          []string{"application/json"},
		"Content-Type":    []string{"application/x-www-form-urlencoded"},
		"accept-encoding": []string{"gzip, deflate, br"},
	}
	return *URL, queryHeader
}

// requestQuery - непосредственно сам запрос к API
func requestQuery(url url.URL, h http.Header) APIresponse {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header = h
	client := http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	myResponse := APIresponse{}
	json.Unmarshal(body, &myResponse)
	return myResponse
}
