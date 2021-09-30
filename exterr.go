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
	TraceJSON() string
	TraceRawString() string
	TraceRows() []traceRow
	AddTraceRow() ErrExtender
	Error() string
	AltError() string
	Type() ErrType
}

type extendedErr struct {
	msg     string
	altMsg  string
	errType ErrType
	trace   []traceRow
}

type traceRow struct {
	Package  string `json:"package"`
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
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
func NewWithType(msg, altMsg string, t ErrType) ErrExtender {
	return &extendedErr{
		msg:     msg,
		trace:   []traceRow{where()},
		altMsg:  altMsg,
		errType: t,
	}
}

// NewWithExtErr() create new ErrExtender from current {err} ErrExtender object
// with joining {msg} and {trace} stack
// Example: err := NewWithExtErr("SQL auth error", err)
func NewWithExtErr(msg string, err ErrExtender) ErrExtender {
	return &extendedErr{
		msg:   fmt.Sprintf("%s: %s", msg, err),
		trace: append(err.TraceRows(), where()),
	}
}

// Wrap() unite two ErrExtender objects
// Example: err.Wrap(err2)
func (e *extendedErr) Wrap(w ErrExtender) {
	e.msg = fmt.Sprintf("%s: %s", e.msg, w.Error())
	e.altMsg = fmt.Sprintf("%s: %s", e.altMsg, w.AltError())
	e.trace = append(e.trace, w.TraceRows()...)
}

// TraceRawString() return string from trace array.
// Every trace line separated by slash.
func (e *extendedErr) TraceRawString() string {
	//trace := reverse(e.trace)
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
	//trace := reverse(e.trace)
	result := ""
	for _, row := range e.trace {
		result = path.Join(result, fmt.Sprintf("{pkg}%s:{file}%s:{function}%s:{line}%d",
			row.Package, row.File, row.Function, row.Line))
	}
	return result
}

// TraceRawString() return JSON-string from trace array.
func (e *extendedErr) TraceJSON() string {
	//trace := reverse(e.trace)
	res, _ := json.Marshal(e.trace)
	return string(res)
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

func where() traceRow {
	pc, file, line, _ := runtime.Caller(exterrTraceSkip)
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

func (e *extendedErr) Error() string {
	return e.msg
}

func (e *extendedErr) AltError() string {
	return e.altMsg
}

func (e *extendedErr) Type() ErrType {
	return e.errType
}

func (e *extendedErr) TraceRows() []traceRow {
	return e.trace
}

// [NOT USED]
// Reverse trace array
func reverse(array []traceRow) []traceRow {
	for i := 0; i < len(array)/2; i++ {
		j := len(array) - i - 1
		array[i], array[j] = array[j], array[i]
	}
	return array
}
