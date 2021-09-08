package call

import (
	"fmt"

	"github.com/alexmolinanasaev/exterr/exterr"
)

func DO() exterr.ErrExtender {
	e := exterr.New("something wrong in DO 1")
	fmt.Println(e.Trace())
	fmt.Println(asd().Trace())
	ee := exterr.New("something wrong in DO 2")
	ee.Wrap(asd())
	fmt.Println(ee.Trace())
	return ee
}

func asd() exterr.ErrExtender {
	e := exterr.NewWithAlt("something wrong in asd", "alt asd err")
	fmt.Println(e.Trace())

	return e
}

func AsErr() error {
	return exterr.New("I'm an error")
}
