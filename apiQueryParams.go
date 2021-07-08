package megaplan

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
)

// QueryParams - параметры запроса
type QueryParams map[string]interface{}

// QueryEscape - urlencode для запроса
func (qp QueryParams) QueryEscape() string {
	b, _ := qp.ToJSON()
	return url.QueryEscape(string(b))
}

// ToJSON - маршализация параметров в JSON
func (qp QueryParams) ToJSON() ([]byte, error) { return json.Marshal(&qp) }

// ToReader - преобразование с JSON и io.Reader для удобства записи в http.Request
func (qp QueryParams) ToReader() (io.Reader, error) {
	b, err := qp.ToJSON()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
