package megaplan

import "time"

// ISO8601 - формат даты для api
const ISO8601 = `2006-01-02T15:04:05-07:00`

// QueryBuildingFunc - функция посттроения тела запроса (обычно json для post запроса)
type QueryBuildingFunc func(QueryParams)

// SetEntityField - добавить поле с сущностью
func SetEntityField(fieldName string, contentType string, value interface{}) (qbf QueryBuildingFunc) {
	var stubQBF = func(qp QueryParams) {} // зашлушка возврата, иначе возвращается nil, который может вызвать ошибки
	switch contentType {
	case "DateOnly":
		t, isTime := value.(time.Time)
		if !isTime {
			return stubQBF
		}
		qbf = func(qp QueryParams) {
			qp[fieldName] = QueryParams{
				"contentType": contentType,
				"year":        t.Year(),
				"month":       t.Month() - 1,
				"day":         t.Day(),
			}
		}
	case "DateTime":
		t, isTime := value.(time.Time)
		if !isTime {
			return func(qp QueryParams) {}
		}
		qbf = func(qp QueryParams) {
			qp[fieldName] = QueryParams{
				"contentType": contentType,
				"value":       t.Format(ISO8601),
			}
		}
	case "DateInterval":
		// если передается не время, то должно указываться кол-во секунд (актуальная документация мегаплана пишет что миллисекунды - это ошибка)
		switch v := value.(type) {
		case uint, uint32, uint64, int, int32, int64:
			qbf = func(qp QueryParams) {
				qp[fieldName] = QueryParams{
					"contentType": contentType,
					"value":       v,
				}
			}
		case time.Time:
			qbf = func(qp QueryParams) {
				qp[fieldName] = QueryParams{
					"contentType": contentType,
					"value":       v.Second(),
				}
			}
		default:
			qbf = stubQBF
		}
	default:
		// по умолчанию BaseEntity - это объект с указанием типа и ID
		qbf = func(qp QueryParams) {
			qp[fieldName] = QueryParams{
				"contentType": contentType,
				"id":          value,
			}
		}
	}
	return
}

// SetEntityArray - добавление массива сущностей в поле (например список аудиторов)
func SetEntityArray(field string, ents ...QueryBuildingFunc) QueryBuildingFunc {
	return func(qp QueryParams) {
		if len(ents) == 0 {
			return
		}
		var arr = make([]interface{}, len(ents))
		var tmpParams = make(QueryParams)
		for i, ent := range ents {
			ent(tmpParams)
			arr[i] = tmpParams[field]
		}
		qp[field] = arr
	}
}

// SetRawField - добавить поле с простым типом значения (string, int, etc.)
func SetRawField(field string, value interface{}) QueryBuildingFunc {
	return func(qp QueryParams) { qp[field] = value }
}
