package megaplan

import (
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

// NewWithConfigFile - инициализация экземпляра из файла конфигурации
func NewWithConfigFile(config *Config) *API {
	return &API{
		client:    http.DefaultClient,
		accessID:  config.AccessID,
		secretKey: []byte(config.SecretKey),
		domain:    config.Domain,
		appUUID:   config.AppUUID,
		appSecret: []byte(config.AppSecret),
	}
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

// Config - формат конфига для API мегаплан
type Config struct {
	AccessID  string `yaml:"access_id"`
	SecretKey string `yaml:"secret_key"`
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
	Domain    string `yaml:"domain"`
	AppUUID   string `yaml:"appUUID"`
	AppSecret string `yaml:"appSecret"`
}

// ReadConfig - парсинг файла по указанному пути, создание конфига
func (c *Config) ReadConfig(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(&c)
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

// SetCostumeClient - установить свой http.Client для API
func (api *API) SetCostumeClient(c *http.Client) {
	api.client = c
}

// SetEmbeddedApplication - установить ключ от встроенного приложения
func (api *API) SetEmbeddedApplication(appuuid, appsecret string) {
	api.appUUID = appuuid
	api.secretKey = []byte(appsecret)
}

// ReadConfig - Сразу инициализирует API с указанием пути к файлу-конфигу
func (api *API) ReadConfig(r io.Reader) {
	var cnf Config
	cnf.ReadConfig(r)
	api.client = http.DefaultClient
	api.accessID = cnf.AccessID
	api.secretKey = []byte(cnf.SecretKey)
	api.domain = cnf.Domain
	api.appUUID = cnf.AppUUID
	api.appSecret = []byte(cnf.AppSecret)
}
