package megaplan

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

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
	if err := OTCdata.response.IFerror(); err != nil {
		return "", err
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
	if err := AccessToken.response.IFerror(); err != nil {
		return "", "", err
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
