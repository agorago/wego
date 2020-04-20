package util

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConvertFromString - converts string to basic types depending on kind passed
func ConvertFromString(s string, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.Atoi(s)
		return i
	case reflect.String:
		return s
	case reflect.Bool:
		b, _ := strconv.ParseBool(s)
		return b
	case reflect.Float32:
		f, _ := strconv.ParseFloat(s, 32)
		return f
	case reflect.Float64:
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}
	return nil
}

// ConvertToString - converts from basic types to string
func ConvertToString(i interface{}, kind reflect.Kind) string {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := fmt.Sprintf("%d", i)
		return s
	case reflect.String:
		return i.(string)
	case reflect.Bool:
		b := fmt.Sprintf("%t", i)
		return b
	case reflect.Float32, reflect.Float64:
		f := fmt.Sprintf("%f", i)
		return f
	}
	return ""
}
