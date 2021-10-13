package exterr

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
)

type ErrExtender interface {
	SetMsg(msg string) ErrExtender
	SetAltMsg(altMsg string) ErrExtender
	SetErrCode(code int) ErrExtender
	Error() string
	GetAltMsg() string
	GetErrCode() int
	GetTraceRows() []traceRow
	AddMsg(msg string) ErrExtender
	AddAltMsg(msg string) ErrExtender
	AddTraceRow() ErrExtender
	TraceTagged() string
	TraceJSON() string
	TraceRawString() string
	Wrap(w ErrExtender) ErrExtender
}

type extendedErr struct {
	msg     string
	altMsg  string
	errCode int
	trace   []traceRow
}

type traceRow struct {
	Package  string `json:"package"`
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
}

func (e *extendedErr) SetMsg(msg string) ErrExtender {
	e.msg = msg
	return e
}

func (e *extendedErr) SetAltMsg(altMsg string) ErrExtender {
	e.altMsg = altMsg
	return e
}

func (e *extendedErr) SetErrCode(code int) ErrExtender {
	e.errCode = code
	return e
}

func (e *extendedErr) Error() string {
	return e.msg
}

func (e *extendedErr) GetAltMsg() string {
	return e.altMsg
}

func (e *extendedErr) GetErrCode() int {
	return e.errCode
}

func (e *extendedErr) GetTraceRows() []traceRow {
	return e.trace
}

// AddMsg() add text to the beginning of the message
func (e *extendedErr) AddMsg(msg string) ErrExtender {
	e.msg = fmt.Sprintf("%s: %s", msg, e.msg)
	return e
}

// AddAltMsg() add text to the beginning of the alternative message
func (e *extendedErr) AddAltMsg(altMsg string) ErrExtender {
	e.altMsg = fmt.Sprintf("%s: %s", altMsg, e.altMsg)
	return e
}

// AddTraceRow() add new trace line in ErrExtender trace array
func (e *extendedErr) AddTraceRow() ErrExtender {
	w := where()
	r := e.trace[len(e.trace)-1]
	if w.Package == r.Package && w.Function == r.Function {
		return e
	}
	e.trace = append(e.trace, w)
	return e
}

// TraceRawString() return string from trace array.
// Every trace line separated by slash.
func (e *extendedErr) TraceRawString() string {
	result := ""
	for _, row := range e.trace {
		result = path.Join(result, fmt.Sprintf("%s:%s:%s:%d",
			row.Package, row.File, row.Function, row.Line))
	}
	return result
}

// TraceRawString() return tagged string from trace array.
// Every trace line separated by slash.
// Format: {pkg}:{file}:{function}:{line}
func (e *extendedErr) TraceTagged() string {
	result := ""
	for _, row := range e.trace {
		result = path.Join(result, fmt.Sprintf("{pkg}%s:{file}%s:{function}%s:{line}%d",
			row.Package, row.File, row.Function, row.Line))
	}
	return result
}

// TraceRawString() return JSON-string from trace array.
func (e *extendedErr) TraceJSON() string {
	res, _ := json.Marshal(e.trace)
	return string(res)
}

// Wrap() unite two ErrExtender objects
// Example: err.Wrap(err2)
func (e *extendedErr) Wrap(w ErrExtender) ErrExtender {
	e.msg = fmt.Sprintf("%s: %s", e.msg, w.Error())
	e.altMsg = fmt.Sprintf("%s: %s", e.altMsg, w.GetAltMsg())
	e.errCode = w.GetErrCode()
	e.trace = append(e.trace, w.GetTraceRows()...)
	return e
}

// New() create ErrExtender with {msg} message and 1 trace line.
// Example: err := New("Error message")
func New(msg string) ErrExtender {
	return &extendedErr{
		msg:   msg,
		trace: []traceRow{where()},
	}
}

// Newf() create ErrExtender with {msg} message and 1 trace line.
// Newf() allows you to format numbers, variables and strings into the first {format} parameter you give it
// Example: err := Newf("Error: %s with code %d", "SQL error", 1005)
func Newf(format string, a ...interface{}) ErrExtender {
	return &extendedErr{
		msg:   fmt.Errorf(format, a...).Error(),
		trace: []traceRow{where()},
	}
}

// With NewWithErr() You can join {msg} to {builtin.error} message
// Example: err := NewWithErr("SQL Error: ", err)
func NewWithErr(msg string, err error) ErrExtender {
	return &extendedErr{
		msg:   fmt.Sprintf("%s: %s", msg, err),
		trace: []traceRow{where()},
	}
}

// NewWithAlt() create ErrExtender with main and alternative messages.
// Example: err := NewWithAlt("SQL connection error", "<SQL_CONNECTION_ERROR>")
func NewWithAlt(msg, altMsg string) ErrExtender {
	return &extendedErr{
		msg:    msg,
		altMsg: altMsg,
		trace:  []traceRow{where()},
	}
}

// NewWithType() create ErrExtender with ErrType (error code)
// Example: err := NewWithType("SQL connection error", "<SQL_CONNECTION_ERROR>", ErrType(1005))
func NewWithType(msg, altMsg string, t int) ErrExtender {
	return &extendedErr{
		msg:     msg,
		trace:   []traceRow{where()},
		altMsg:  altMsg,
		errCode: t,
	}
}

// NewWithExtErr() create new ErrExtender from current {err} ErrExtender object
// with joining {msg} {altMsg} {errType} and {trace} stack
// Example: err := NewWithExtErr("SQL auth error", err)
func NewWithExtErr(msg string, err ErrExtender) ErrExtender {
	return &extendedErr{
		msg:     fmt.Sprintf("%s: %s", msg, err),
		altMsg:  err.GetAltMsg(),
		errCode: err.GetErrCode(),
		trace:   append(err.GetTraceRows(), where()),
	}
}

func where() traceRow {
	pc, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()

	slashIndex := strings.LastIndex(function, "/")
	function = function[slashIndex+1:]

	s := strings.Split(function, ".")
	packageName, function := s[0], s[1]

	slashIndex = strings.LastIndex(file, "/")
	fileName := file[slashIndex+1:]

	return traceRow{
		Package:  packageName,
		File:     fileName,
		Function: function,
		Line:     line,
	}
}
