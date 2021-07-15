package megaplan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Response - ответ API
type Response struct {
	Meta Meta        `json:"meta"` // metainfo ответа
	Data interface{} `json:"data"` // поле для декодирования присвоенной структуры
}

// Next - есть ли следующая страница
func (res Response) Next() bool { return res.Meta.Pagination.HasMoreNext }

// Prev - есть ли предыдущая страница
func (res Response) Prev() bool { return res.Meta.Pagination.HasMorePrev }

// Decode - парсинг ответа API
func (res *Response) Decode(r io.Reader, i interface{}) (err error) {
	res.Data = i
	if err := json.NewDecoder(r).Decode(res); err != nil {
		return err
	}
	return res.Meta.Error()
}

// Pagination - пагинация
type Pagination struct {
	Count       int64 `json:"count"`
	Limit       int64 `json:"limit"`
	CurrentPage int64 `json:"currentPage"`
	HasMoreNext bool  `json:"hasMoreNext"`
	HasMorePrev bool  `json:"hasMorePrev"`
}

// UnmarshalJSON - json.Unmarshaler
func (p *Pagination) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte{91, 93}) {
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, isdelim := t.(json.Delim); isdelim {
			continue
		}
		if field, ok := t.(string); ok {
			switch field {
			case "count":
				err = dec.Decode(&p.Count)
			case "limit":
				err = dec.Decode(&p.Limit)
			case "currentPage":
				err = dec.Decode(&p.CurrentPage)
			case "hasMoreNext":
				err = dec.Decode(&p.HasMoreNext)
			case "hasMorePrev":
				err = dec.Decode(&p.HasMorePrev)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON - json.Marshaler
// TODO: вообще обратный маршалинг на практике не нужен, поэтому нужно доделать позже
func (p Pagination) MarshalJSON() ([]byte, error) { return nil, nil }

// Meta - metainfo
type Meta struct {
	Errors []struct {
		Fields  interface{} `json:"field"`
		Message interface{} `json:"message"`
	} `json:"errors"`
	Status     int64      `json:"status"`
	Pagination Pagination `json:"pagination"`
}

// Error - если была ошибка, переданная в meta, то вернется ошибка с описание мегаплана, если нет, то вернется nil
func (m Meta) Error() (err error) {
	if len(m.Errors) > 0 {
		var errorsStr = make([]string, len(m.Errors))
		for i := range m.Errors {
			errorsStr[i] = fmt.Sprintf("FIELD: %v MESSAGE: %v", m.Errors[i].Fields, m.Errors[i].Message)
		}
		err = errors.New(strings.Join(errorsStr, "\n"))
	}
	return
}

// ParseResponse - обертка над методов Response.Decode + данные о пагинации
// utility-функция для упрощения чтения ответа API
func ParseResponse(r io.Reader, i interface{}) (next bool, prev bool, err error) {
	var res Response
	if err := res.Decode(r, i); err != nil {
		return res.Next(), res.Prev(), err
	}
	return res.Next(), res.Prev(), nil
}
