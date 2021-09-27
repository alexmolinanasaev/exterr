package exterr

import (
	"encoding/json"
	"fmt"
	"path"
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
	TraceTagged() string
	TracePretty() string
	TraceRaw() string
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

type PrettyTrace struct {
	Package  string       `json:"package"`
	File     string       `json:"file"`
	Function string       `json:"function"`
	Line     string       `json:"line"`
	Child    *PrettyTrace `json:"child"`
}

func New(msg string) ErrExtender {
	return &extendedErr{
		msg:   msg,
		where: where(),
	}
}

func NewWithErr(msg string, err error) ErrExtender {
	msg = fmt.Sprintf("%s: %s", msg, err)
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
	e.where = fmt.Sprintf("%s/%s", where(), w.TraceRaw())
}

// will return trace in formart: packageName:fileName:function:line
// if trace was wrapped or added, traces will be separated by slash /
func (e *extendedErr) TraceRaw() string {
	return e.where
}

func (e *extendedErr) TraceTagged() string {
	result := ""
	parsedFullTrace := strings.Split(e.where, "/")
	for _, t := range parsedFullTrace {
		parsedTrace := strings.Split(t, ":")
		result = path.Join(result, fmt.Sprintf("{pkg}%s:{file}%s:{function}%s:{line}%s",
			parsedTrace[0], parsedTrace[1], parsedTrace[2], parsedTrace[3]))

	}
	return result
}

func (e *extendedErr) TracePretty() string {
	parsedFullTrace := strings.Split(e.where, "/")
	parsedTrace := strings.Split(parsedFullTrace[0], ":")
	prettyTrace := &PrettyTrace{
		Package:  parsedTrace[0],
		File:     parsedTrace[1],
		Function: parsedTrace[2],
		Line:     parsedTrace[3],
		Child:    &PrettyTrace{},
	}
	currTrace := prettyTrace.Child
	lastTraceIndex := len(parsedFullTrace)
	for i, t := range parsedFullTrace[1:] {
		parsedTrace := strings.Split(t, ":")
		currTrace.Package = parsedTrace[0]
		currTrace.File = parsedTrace[1]
		currTrace.Function = parsedTrace[2]
		currTrace.Line = parsedTrace[3]
		if i+2 < lastTraceIndex {
			currTrace.Child = &PrettyTrace{}
			currTrace = currTrace.Child
		}
	}

	res, _ := json.Marshal(prettyTrace)

	return string(res)
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