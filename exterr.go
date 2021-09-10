package exterr

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	exterrTraceSkip   int = 2
	exterrPackageSkip int = 7
)

type ErrType int

type ErrExtender interface {
	Wrap(w ErrExtender)
	Trace() string
	AddTrace() ErrExtender
	Error() string
	AltError() string
	Type() ErrType
}

type extendedErr struct {
	msg     string
	altMsg  string
	where   string
	errType ErrType
}

func New(msg string) ErrExtender {
	return &extendedErr{
		msg:   msg,
		where: where(),
	}
}

func NewWithAlt(msg, altMsg string) ErrExtender {
	return &extendedErr{
		msg:    msg,
		where:  where(),
		altMsg: altMsg,
	}
}

func NewWithType(msg, altMsg string, t ErrType) ErrExtender {
	return &extendedErr{
		msg:     msg,
		where:   where(),
		altMsg:  altMsg,
		errType: t,
	}
}

func (e *extendedErr) Wrap(w ErrExtender) {
	e.msg = fmt.Sprintf("%s:%s", e.msg, w.Error())
	e.altMsg = fmt.Sprintf("%s:%s", e.altMsg, w.AltError())
	e.where = fmt.Sprintf("%s/%s", where(), w.Trace())
}

// will return trace in formart: packageName:fileName:function:line
// if trace was wrapped or added, traces will be separated by slash /
func (e *extendedErr) Trace() string {
	return e.where
}

func (e *extendedErr) AddTrace() ErrExtender {
	e.where = fmt.Sprintf("%s/%s", where(), e.where)
	return e
}

func (e *extendedErr) Error() string {
	return e.msg
}

func (e *extendedErr) AltError() string {
	return e.altMsg
}

func (e *extendedErr) Type() ErrType {
	return e.errType
}

func where() string {
	pc, file, line, _ := runtime.Caller(exterrTraceSkip)
	function := runtime.FuncForPC(pc).Name()

	funcIndex := strings.LastIndex(function, "exterr")
	if funcIndex < 0 {
		funcIndex = 0
	} else {
		funcIndex += exterrPackageSkip
	}
	function = function[funcIndex:]

	dotIndex := strings.LastIndex(function, ".")
	packageName := function[:dotIndex]
	function = function[dotIndex+1:]

	fileIndex := strings.LastIndex(file, "/")
	fileName := file[fileIndex+1:]

	return fmt.Sprintf("%s:%s:%s:%d", packageName, fileName, function, line)
}
