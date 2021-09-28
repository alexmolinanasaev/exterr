package main

import (
	"fmt"
	"log"

	"github.com/alexmolinanasaev/exterr"
)

const (
	internalServerErrorType exterr.ErrType = 1
)

func main() {
	log.Println(New().Error())
	fmt.Println()

	log.Println(NewWithAlt().Error())
	log.Println(NewWithAlt().AltError())
	fmt.Println()

	err := NewWithType()
	if err.Type() == internalServerErrorType {
		log.Println(err.AltError())
	}
	fmt.Println()

	log.Println(LikeErr().Error())
	fmt.Println()

	err = TraceErr()
	log.Println(err.Error())
	log.Println(err.TraceTagged())
	fmt.Println()

	err = Wrap()
	log.Println(err.Error())
	log.Println(err.TraceTagged())
	e := exterr.New("wraping err")
	e.Wrap(err)
	log.Println(e.Error())
	log.Println(e.TraceTagged())
	fmt.Println()

	log.Println(AddTraceExample().AddTrace().TraceRaw())
	log.Println(AddTraceExample().AddTrace().TracePretty())

}

// is simple to create
func New() exterr.ErrExtender {
	return exterr.New("i am an extended error")
}

// can have alternative error message
func NewWithAlt() exterr.ErrExtender {
	return exterr.NewWithAlt("this is main message", "this is alt message")
}

// can have type identificator and can be procceed in different ways
func NewWithType() exterr.ErrExtender {
	return exterr.NewWithType("sql no rows", "user not found", internalServerErrorType)
}

// exterr can be used like a standard golang error
func LikeErr() error {
	return exterr.New("can be used like standard error")
}

// can store info about place where was
func TraceErr() exterr.ErrExtender {
	return exterr.New("there is where i was created")
}

// can wrap other errors
func Wrap() exterr.ErrExtender {
	return exterr.New("wrap me!")
}

// if error will be just passed higher you can add trace manually
func AddTraceExample() exterr.ErrExtender {
	return f1().AddTrace()
}

func f1() exterr.ErrExtender { return f2().AddTrace() }
func f2() exterr.ErrExtender { return f3().AddTrace() }
func f3() exterr.ErrExtender { return exterr.New("trace me") }
