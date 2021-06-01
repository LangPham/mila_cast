package mila_cast

import "reflect"

func (exchange *Exchange) PutField(field string, value interface{}) {
	rData := reflect.ValueOf(&exchange.Data).Elem()
	oldModel := rData.Elem()
	newModel := reflect.New(rData.Elem().Type()).Elem()
	for i := 0; i < newModel.NumField(); i++ {
		if newModel.Type().Field(i).Name == field {
			// Put data
			switch newModel.Type().Field(i).Type.String() {
			case "string":
				newModel.Field(i).SetString(value.(string))
			case "int":
				newModel.Field(i).SetInt(value.(int64))
			}
		} else {
			newModel.Field(i).Set(oldModel.Field(i))
		}
	}
	rData.Set(newModel)
}
