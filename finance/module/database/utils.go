package database

import (
	"fmt"
	"reflect"
)

func SliceToArgs(slice any) []any {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		panic(fmt.Errorf("Input is not a slice"))
	}

	result := make([]interface{}, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		result[i] = sliceValue.Index(i).Interface()
	}

	return result
}
