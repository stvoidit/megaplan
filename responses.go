package megaplan

// func (r *response) IFerror() error {
// 	if r.Status.Code == "error" {
// 		return errors.New(r.Status.Message)
// 	}
// 	return nil
// }

// // ResponseBuffer - ответ от API
// type ResponseBuffer struct {
// 	bytes.Buffer
// }

// type semiResponse struct {
// 	response
// 	Data map[string]interface{} `json:"data"`
// }

// // Scan - парсинг структуры
// func (rb *ResponseBuffer) Scan(i interface{}) error {
// 	var res = new(semiResponse)
// 	var dec = json.NewDecoder(rb)
// 	dec.UseNumber()
// 	if err := dec.Decode(&res); err != nil {
// 		return err
// 	}
// 	if err := res.IFerror(); err != nil {
// 		return err
// 	}
// 	for _, v := range res.Data {
// 		var buff = new(bytes.Buffer)
// 		var dec = json.NewDecoder(buff)
// 		dec.UseNumber()
// 		if err := json.NewEncoder(buff).Encode(&v); err != nil {
// 			return err
// 		}
// 		if err := dec.Decode(i); err == nil {
// 			break
// 		}
// 	}
// 	return nil
// }

// // UserVerifyResponse - тип для верификации юзеров во встроенном приложении
// type UserVerifyResponse struct {
// 	response
// 	Data UserAppVerification `json:"data"`
// }
