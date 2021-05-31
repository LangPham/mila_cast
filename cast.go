package mila_cast

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/gofiber/fiber/v2"
)

func newReflectValue(oldValue reflect.Value, req string) (copy reflect.Value) {
	oldModel := oldValue.Elem()
	newModel := reflect.New(oldValue.Elem().Type()).Elem()
	if req == "update" {
		for i := 0; i < newModel.NumField(); i++ {
			newModel.Field(i).Set(oldModel.Field(i))
		}
	}
	return newModel
}

func Cast(modelIn interface{}, c *fiber.Ctx) (exchange Exchange) {

	value := c.FormValue("_METHOD")
	userReq := ""
	switch strings.ToLower(value) {
	case "post":
		userReq = "insert"
	case "put":
		userReq = "update"
	default:
		userReq = "new"
	}

	change := hashmap.New()
	mError := hashmap.New()
	dataType := reflect.TypeOf(modelIn).Name()
	modelInType := reflect.ValueOf(&modelIn).Elem()

	newModel := newReflectValue(modelInType, userReq)

	for i := 0; i < newModel.NumField(); i++ {
		//aon.Dump(newModel.Type().Field(i).Name, "NAME")
		//aon.Dump(newModel.Type().Field(i).Tag, "Tag")
		//aon.Dump(newModel.Type().Field(i).Tag.Get("cast"), "CASTTag")
		//aon.Dump(newModel.Type().Field(i).Type, "TYPE")
		//aon.Dump(newModel.Field(i).Interface(), "INT")
		//value := c.FormValue()
		//strings.Split(newModel.Type().Field(i).Tag.(string), " ")
		fieldName := newModel.Type().Field(i).Tag.Get("cast")
		if fieldName != "" {
			fieldValue := c.FormValue(fieldName)
			change.Put(fieldName, fieldValue)
			//aon.Dump(newModel.Type().Field(i).Type.String(), "TYPE")
			switch newModel.Type().Field(i).Type.String() {
			case "string":
				newModel.Field(i).SetString(fieldValue)
			case "int":
				//aon.Dump(fieldValue, "VAL INT")
				valueInt, err := strconv.Atoi(fieldValue)
				if err != nil {
					mError.Put(fieldName, "Please input number!")
				} else {
					newModel.Field(i).SetInt(int64(valueInt))
				}
			}
		}
	}

	modelInType.Set(newModel)
	exchange.Data = modelIn

	exchange.Request = userReq
	exchange.DataType = dataType
	exchange.Change = change
	exchange.Error = mError

	return
}
