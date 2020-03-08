package megaplan

import (
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

// APIWithConfig - инициализация экземпляра из файла конфигурации
func APIWithConfig(file io.ReadSeeker) *API {
	var api = new(API)
	api.ReadConfig(file)
	return api
}

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

// ReadConfig - парсинг файла по указанному пути, создание конфига
func (c *Config) ReadConfig(file io.ReadSeeker) {
	file.Seek(0, 0)
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		panic(err)
	}
}

// API - Структура объекта API v1
type API struct {
	config    *Config
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

// ReadConfig - Сразу инициализирует API с указанием пути к файлу-конфигу
func (api *API) ReadConfig(file io.ReadSeeker) {
	var cnf = new(Config)
	cnf.ReadConfig(file)
	api.config = cnf
	api.client = new(http.Client)
	api.accessID = cnf.Megaplan.AccessID
	api.secretKey = []byte(cnf.Megaplan.SecretKey)
	api.domain = cnf.Megaplan.Domain
	api.appUUID = cnf.Megaplan.AppUUID
	api.appSecret = []byte(cnf.Megaplan.AppSecret)
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
