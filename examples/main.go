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
	log.Println(err.Trace())
	fmt.Println()

	err = Wrap()
	log.Println(err.Error())
	log.Println(err.Trace())
	e := exterr.New("wraping err")
	e.Wrap(err)
	log.Println(e.Error())
	log.Println(e.Trace())
	fmt.Println()

	log.Println(AddTrace().Trace())
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
func AddTrace() exterr.ErrExtender {
	// func1
	err := func() exterr.ErrExtender {
		// func2
		err := func() exterr.ErrExtender {
			// func3
			err := func() exterr.ErrExtender {
				return exterr.New("trace me")
			}()
			return err.AddTrace()
		}()
		return err.AddTrace()
	}()
	return err.AddTrace()
}
