package util_test

import (
	"fmt"
	"gitlab.intelligentb.com/devops/bplus/util"
	"reflect"
)

func ExampleConvertFromString() {
	fmt.Printf("%d\n", util.ConvertFromString("23", reflect.Int))
	fmt.Printf("%s\n", util.ConvertFromString("xxx", reflect.String))
	fmt.Printf("%v\n", util.ConvertFromString("true", reflect.Bool))
	fmt.Printf("%5.2f\n", util.ConvertFromString("22.88", reflect.Float32))
	fmt.Printf("%4.1f\n", util.ConvertFromString("59.2", reflect.Float64))
	// Output:
	// 23
	// xxx
	// true
	// 22.88
	// 59.2
}

func ExampleConvertToString() {
	fmt.Printf("%s\n", util.ConvertToString(23, reflect.Int))
	fmt.Printf("%s\n", util.ConvertToString("xxx", reflect.String))
	fmt.Printf("%s\n", util.ConvertToString(true, reflect.Bool))
	fmt.Printf("%s\n", util.ConvertToString(22.88, reflect.Float32))
	fmt.Printf("%s\n", util.ConvertToString(59.2, reflect.Float64))
	// Output:
	// 23
	// xxx
	// true
	// 22.880000
	// 59.200000
}
