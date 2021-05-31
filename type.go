package mila_cast

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/go-playground/validator/v10"
)

type Exchange struct {
	Error    *hashmap.Map
	Change   *hashmap.Map
	Data     interface{}
	Request  string
	DataType string
	ResultID string
	Validate *validator.Validate
}