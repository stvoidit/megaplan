package megaplan

// v1.2

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

// APIresponse - ответ API, json
type APIresponse struct {
	Status map[string]string
	Data   map[string]interface{}
}

// GET - get запрос к API
func GET(domain string, acessid string, secretkey string, uri string, pyload map[string]string, class interface{}) APIresponse {
	const rMethod = "GET"
	urlQuery, queryHeader := queryHasher(domain, acessid, secretkey, rMethod, uri, pyload)
	responseAPI := requestQuery(rMethod, urlQuery, queryHeader, class)
	return responseAPI
}

// POST - post запрос на API
func POST(domain string, acessid string, secretkey string, uri string, pyload map[string]string, class interface{}) APIresponse {
	const rMethod = "POST"
	urlQuery, queryHeader := queryHasher(domain, acessid, secretkey, rMethod, uri, pyload)
	responseAPI := requestQuery(rMethod, urlQuery, queryHeader, class)
	return responseAPI

}

// queryHasher - задаем сигнатуру, отдает URL и Header для запросов к API
func queryHasher(domain string, acessid string, secretkey string, method string, uri string, payload map[string]string) (url.URL, http.Header) {
	const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
	URL, err := url.Parse(domain)
	if err != nil {
		panic(err.Error())
	}
	URL.Path += uri
	today := time.Now().Format(rfc2822)
	urlParams := url.Values{}
	for k, v := range payload {
		urlParams.Add(k, v)
	}
	URL.RawQuery = urlParams.Encode()
	Signature := method + "\n\n" + "application/x-www-form-urlencoded" + "\n" + today + "\n" + URL.Host + URL.Path + "?" + URL.RawQuery
	h := hmac.New(sha1.New, []byte(secretkey))
	h.Write([]byte(Signature))
	hexSha1 := hex.EncodeToString(h.Sum(nil))
	sha1Query := base64.StdEncoding.EncodeToString([]byte(hexSha1))
	queryHeader := http.Header{
		"Date":            []string{today},
		"X-Authorization": []string{acessid + ":" + sha1Query},
		"Accept":          []string{"application/json"},
		"Content-Type":    []string{"application/x-www-form-urlencoded"},
		"accept-encoding": []string{"gzip, deflate, br"},
	}
	return *URL, queryHeader
}

// requestQuery - непосредственно сам запрос к API
// В аргумент class передается экземпляр структуры в которую будут записаны данные Unmarshal
func requestQuery(method string, url url.URL, headers http.Header, class interface{}) APIresponse {
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header = headers
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	myResponse := APIresponse{}
	json.Unmarshal(body, &myResponse)
	if class != nil {
		json.Unmarshal(body, &class)
	}
	return myResponse
}
