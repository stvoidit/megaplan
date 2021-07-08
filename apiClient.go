package megaplan

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"runtime"
	"time"
)

// NewClient - обертка над http.Client для удобной работы с API v3
func NewClient(domain, token string, opts ...ClientOption) (c *ClientV3) {
	// обмен трафиком идёт очень активный, поэтому целесообразно использовать http2 + KeepAlive
	// бэкэнд мегаплана корректно умеет работать с http и KeepAlive, что экономит время и ресурсы на соединение
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			Rand:       rand.Reader,
			Time:       time.Now,
		},
		DialContext: (&net.Dialer{
			Timeout:   time.Minute * 10,
			KeepAlive: time.Minute,
		}).DialContext,
		TLSHandshakeTimeout: 30 * time.Second,
		MaxIdleConns:        0,
		IdleConnTimeout:     time.Minute,
		ForceAttemptHTTP2:   true,
		ReadBufferSize:      256 << 10,
		WriteBufferSize:     256 << 10,
	}
	var jar, _ = cookiejar.New(nil)
	c = &ClientV3{
		client: &http.Client{
			Transport: tr,
			Jar:       jar,
			Timeout:   time.Minute * 10,
		},
		domain:         domain,
		defaultHeaders: http.Header{"User-Agent": {runtime.Version()}},
	}
	c.SetToken(token)
	c.SetOptions(opts...)
	return
}

// ClientV3 - клиент
type ClientV3 struct {
	client         *http.Client
	domain         string
	defaultHeaders http.Header
}

// Do - http.Do + установка обязательных заголовков
func (c *ClientV3) Do(req *http.Request) (*http.Response, error) {
	req.Header = c.defaultHeaders
	return c.client.Do(req)
}

// DoRequestAPI - т.к. в v3 параметры запроса для GET (json маршализируется и будет иметь вид: "*?{params}=")
func (c ClientV3) DoRequestAPI(method string, endpoint string, params QueryParams, body io.Reader) (rc io.ReadCloser, err error) {
	var args string // параметры строки запроса
	if params != nil {
		args = params.QueryEscape()
	}
	request, err := http.NewRequest(method, c.domain, body)
	if err != nil {
		return nil, err
	}
	request.URL.Path = endpoint
	request.URL.RawQuery = args
	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	return unzipResponse(response)
}

// ErrUnknownCompressionMethod - неизвестное значение в заголовке "Content-Encoding"
// не является фатальной ошибкой, должна возвращаться вместе с http.Response.Body,
// чтобы пользователь мог реализовать свой метод обработки сжатого сообщения
var ErrUnknownCompressionMethod = errors.New("unknown compression method")

// unzipResponse - распаковка сжатого ответа
func unzipResponse(response *http.Response) (rc io.ReadCloser, err error) {
	if response.Uncompressed {
		return response.Body, nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	var r = bytes.NewReader(body)
	ce := response.Header.Get("Content-Encoding")
	switch ce {
	case "":
		// кейс, когда запрашивалось сжатие, но сервер не поддерживает запрашиваемый вид сжатия
		// response.Uncompressed будет в значении false, но в тело ответа будет не сжато и заголовок "Content-Encoding" отсутствует
		rc = io.NopCloser(r)
	case "gzip":
		rc, err = gzip.NewReader(r)
	default:
		rc = io.NopCloser(r)
		err = ErrUnknownCompressionMethod
	}
	return
}

// SetOptions - применить опции
func (c *ClientV3) SetOptions(opts ...ClientOption) {
	for i := range opts {
		opts[i](c)
	}
}

// SetToken - установить или изменить токен доступа
func (c *ClientV3) SetToken(token string) { c.defaultHeaders.Set("AUTHORIZATION", "Bearer "+token) }

// ClientOption - функция применения настроект
type ClientOption func(*ClientV3)

// OptionInsecureSkipVerify - переключение флага bool в http.Client.Transport.TLSClientConfig.InsecureSkipVerify - отключать или нет проверку сертификтов
// Если домен использует самоподписанные сертифика, то удобно включать на время отладки и разработки
func OptionInsecureSkipVerify(b bool) ClientOption {
	return func(c *ClientV3) {
		if c.client.Transport != nil {
			(c.client.Transport.(*http.Transport)).TLSClientConfig.InsecureSkipVerify = b
		}
	}
}

// OptionsSetHTTPTransport - установить своб настройку http.Transport
func OptionsSetHTTPTransport(tr http.RoundTripper) ClientOption {
	return func(c *ClientV3) { c.client.Transport = tr }
}

// OptionEnableAcceptEncodingGzip - доабвить заголов Accept-Encoding=gzip к запросу
// т.е. объекм трафика на хуках может быть большим, то удобно запрашивать сжатый ответ
func OptionEnableAcceptEncodingGzip(b bool) ClientOption {
	const header = "Accept-Encoding"
	return func(c *ClientV3) {
		if b {
			c.defaultHeaders.Set(header, "gzip")
		} else {
			c.defaultHeaders.Del(header)
		}
	}
}

// OptionSetClientHTTP - установить свой экземпляр httpClient
func OptionSetClientHTTP(client *http.Client) ClientOption {
	return func(c *ClientV3) {
		if client != nil {
			c.client = client
		}
	}
}
