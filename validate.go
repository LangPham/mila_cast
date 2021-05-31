package mila_cast

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/huandu/xstrings"

)

// for validate
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// For validate
func (exchange *Exchange) putError(field string, err string) {
	key := xstrings.ToSnakeCase(field)
	val, has := exchange.Error.Get(key)
	if has {
		exchange.Error.Put(key, val.(string)+", "+err)
	} else {
		exchange.Error.Put(key, err)
	}
}

func (exchange *Exchange) ValidateModel(properties ...string) {

	validate := validator.New()

	err := validate.Struct(exchange.Data)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		//if _, ok := err.(*validator.InvalidValidationError); ok {
		//	fmt.Println(err)
		//	return
		//}

		for _, err := range err.(validator.ValidationErrors) {

			switch {
			case contains([]string{"email"}, err.Tag()):
				exchange.putError(err.Field(), "Please input email!")
			case contains([]string{"required"}, err.Tag()):
				exchange.putError(err.Field(), "Can't blank!")
			case contains([]string{"eq", "gt", "gte", "lt", "lte", "ne"}, err.Tag()):
				exchange.putError(err.Field(), "not "+err.Tag()+" "+err.Param())
			default:
				exchange.putError(err.Field(), "not "+err.Tag()+" "+err.Param())
				//fmt.Println(err.Namespace())
				//fmt.Println(err.Field())
				//fmt.Println(err.StructNamespace())
				//fmt.Println(err.StructField())
				//fmt.Println(err.Tag())
				//fmt.Println(err.ActualTag())
				//fmt.Println(err.Kind())
				//fmt.Println(err.Type())
				//fmt.Println(err.Value())
				//fmt.Println(err.Param())
				//fmt.Println()
				//fmt.Println(err)
			}
		}

		return
	}
}

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
