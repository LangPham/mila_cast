package mila_cast

import (
	"github.com/LangPham/mila_cast/aon"
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

func setField(field *reflect.Value, fieldName string, fieldValue string, change *hashmap.Map, mError *hashmap.Map) {
	change.Put(fieldName, fieldValue)
	rKind := field.Kind()
	switch rKind {
	case reflect.String:
		field.SetString(fieldValue)
	case reflect.Int:
		valueInt, err := strconv.Atoi(fieldValue)
		if err != nil && fieldValue != "" {
			mError.Put(fieldName, "Please input number!")
		} else {
			field.SetInt(int64(valueInt))
		}
	case reflect.Bool:
		valueBool, err := strconv.ParseBool(fieldValue)
		if err != nil {
			mError.Put(fieldName, "Please input boolean!")
		} else {
			field.SetBool(valueBool)
		}
	default:
		//aon.Dump(rKind, "rKind")
	}
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

		switch {
		case fieldName == "mixin":
			//aon.Dump(newModel.Field(i).Kind(), "mixin")
			rv := reflect.Indirect(newModel.Field(i))

			//aon.Dump(rv.Field(0).CanSet(), "mixin")
			for j := 0; j < rv.NumField(); j++ {
				fieldNameSub := rv.Type().Field(j).Tag.Get("cast")
				aon.Dump(fieldNameSub, "In mixin")
				if fieldNameSub != "" {
					fieldValue := c.FormValue(fieldNameSub)
					field := rv.Field(j)
					setField(&field, fieldNameSub, fieldValue, change, mError)
				}
			}
		case fieldName != "":
			fieldValue := c.FormValue(fieldName)
			field := newModel.Field(i)
			setField(&field, fieldName, fieldValue, change, mError)
		}
	}

	modelInType.Set(newModel)
	exchange.Data = modelIn

	exchange.Request = userReq
	exchange.DataType = dataType
	exchange.Change = change
	exchange.Error = mError
	exchange.Valid = mError.Empty()
	return
}
