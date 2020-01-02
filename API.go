package megaplan

import (
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
)

// Config - формат конфига для API мегаплан
type Config struct {
	Megaplan struct {
		AccessID  string `yaml:"access_id"`
		SecretKey string `yaml:"secret_key"`
		Login     string `yaml:"login"`
		Password  string `yaml:"password"`
		Domain    string `yaml:"domain"`
		AppUUID   string `yaml:"appUUID"`
		AppSecret string `yaml:"appSecret"`
	} `yaml:"megaplan"`
}

// ParseConfig - парсинг файла по указанному пути, создание конфига
func (c *Config) ParseConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		return err
	}
	return nil
}

// API - Структура объекта API v1
type API struct {
	AccessID  string
	SecretKey []byte
	Domain    string
	login     string
	password  string
	AppUUID   string
	AppSecret []byte
	client    *http.Client
}

// ParseConfig - Сразу инициализирует API с указанием пути к файлу-конфигу
func (api *API) ParseConfig(path string) {
	var c Config
	c.ParseConfig(path)
	api.AccessID = c.Megaplan.AccessID
	api.SecretKey = []byte(c.Megaplan.SecretKey)
	api.Domain = c.Megaplan.Domain
	api.login = c.Megaplan.Login
	api.password = c.Megaplan.Password
	api.AppUUID = c.Megaplan.AppUUID
	api.client = &http.Client{Timeout: 1 * time.Minute}
	api.AppSecret = []byte(c.Megaplan.AppSecret)
}

// GetToken - Получение токена API
func (api *API) GetToken() error {
	md5p := md5Passord(api.password)
	OTCkey, err := getOTC(api.Domain, api.login, md5p)
	if err != nil {
		panic(err.Error())
	}
	AID, Skey, err := getToken(api.Domain, api.login, md5p, OTCkey)
	if err != nil {
		return err
	}
	api.AccessID = AID
	api.SecretKey = []byte(Skey)
	return nil
}
