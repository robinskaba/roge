package ux

import (
	"fmt"
	"io"
	"reflect"
)

func ListStruct(str any, out io.Writer) {
	val := reflect.ValueOf(str)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("json")
		fieldValue := val.Field(i).Interface()
		fmt.Fprintf(out, "- %s=%v\n", fieldName, fieldValue)
	}
}
