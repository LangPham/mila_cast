package mila_cast

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Exchange struct {
	Valid    bool
	Error    *hashmap.Map
	Change   *hashmap.Map
	Data     interface{}
	Request  string
	DataType string
	ResultID string
	Validate *validator.Validate
}

func NewExchange(modelIn interface{}) *Exchange {
	switch v := modelIn.(type) {
	case string:
		return &Exchange{
			Error:    hashmap.New(),
			Change:   hashmap.New(),
			Request:  "new",
			DataType: v,
			Validate: validator.New(),
		}
	default:
		dataType := reflect.TypeOf(modelIn).Name()
		change := hashmap.New()
		modelInType := reflect.ValueOf(&modelIn).Elem()
		for i := 0; i < modelInType.Elem().NumField(); i++ {
			fieldName := modelInType.Elem().Type().Field(i).Tag.Get("cast")
			if fieldName != "" {
				change.Put(fieldName, modelInType.Elem().Field(i).Interface())
			}
		}
		return &Exchange{
			Error:    hashmap.New(),
			Change:   change,
			Data:     modelIn,
			Request:  "edit",
			DataType: dataType,
			Validate: validator.New(),
		}
	}
}
