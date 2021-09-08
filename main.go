package main

import (
	"fmt"
	"reflect"

	"github.com/alexmolinanasaev/exterr/call"
	"github.com/alexmolinanasaev/exterr/exterr"
)

func main() {
	err := call.DO()
	fmt.Println(err.Trace())
	err.AddTrace()
	fmt.Println(err.Trace())
	fmt.Println(err.Error())

	e := exterr.New("")
	fmt.Println(e.Trace())
	fmt.Println(e.Error())

	ee := call.AsErr()
	fmt.Println(ee.Error())

	fmt.Println(reflect.TypeOf(ee))
}
