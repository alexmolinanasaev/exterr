package exterr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

var (
	BadRequest            = errors.New("Bad request")
	WrongCredentials      = errors.New("Wrong Credentials")
	NotFound              = errors.New("Not Found")
	Unauthorized          = errors.New("Unauthorized")
	Forbidden             = errors.New("Forbidden")
	PermissionDenied      = errors.New("Permission Denied")
	ExpiredCSRFError      = errors.New("Expired CSRF token")
	WrongCSRFToken        = errors.New("Wrong CSRF token")
	CSRFNotPresented      = errors.New("CSRF not presented")
	NotRequiredFields     = errors.New("No such required fields")
	BadQueryParams        = errors.New("Invalid query params")
	InternalServerError   = errors.New("Internal Server Error")
	RequestTimeoutError   = errors.New("Request Timeout")
	ExistsEmailError      = errors.New("User with given email already exists")
	InvalidJWTToken       = errors.New("Invalid JWT token")
	InvalidJWTClaims      = errors.New("Invalid JWT claims")
	NotAllowedImageHeader = errors.New("Not allowed image header")
	NoCookie              = errors.New("not found cookie header")
)

type ErrExtender interface {
	SetMsg(msg string) ErrExtender
	SetAltMsg(altMsg string) ErrExtender
	SetErrCode(code int) ErrExtender
	Error() string
	GetAltMsg() string
	GetErrCode() int
	GetTraceRows() string
	AddMsg(msg string) ErrExtender
	AddAltMsg(msg string) ErrExtender
	// AddTraceRow() ErrExtender
	// TraceTagged() string
	// TraceJSON() string
	// TraceRawString() string
	// Wrap(w ErrExtender) ErrExtender
}

type extendedErr struct {
	ErrStatus int      `json:"status,omitempty"`
	ErrError  string   `json:"error,omitempty"`
	ErrMsg    string   `json:"msg,omitempty"`
	Causes    string   `json:"-"`
	trace     []string `json:"-"`
}

func (e *extendedErr) SetMsg(msg string) ErrExtender {
	e.ErrError = msg
	return e
}

func (e *extendedErr) SetAltMsg(altMsg string) ErrExtender {
	e.ErrMsg = altMsg
	return e
}

func (e *extendedErr) SetErrCode(code int) ErrExtender {
	e.ErrStatus = code
	return e
}

func (e *extendedErr) Error() string {
	return e.ErrError
}

func (e *extendedErr) GetAltMsg() string {
	return e.ErrMsg
}

func (e *extendedErr) GetErrCode() int {
	return e.ErrStatus
}

func (e *extendedErr) GetTraceRows() string {
	e.trace = append([]string{e.Causes}, e.trace...)
	return strings.Join(e.trace, "\n")
}

// AddMsg() add text to the beginning of the message
func (e *extendedErr) AddMsg(msg string) ErrExtender {
	e.ErrError = fmt.Sprintf("%s: %s", msg, e.ErrError)
	return e
}

// AddAltMsg() add text to the beginning of the alternative message
func (e *extendedErr) AddAltMsg(altMsg string) ErrExtender {
	e.ErrMsg = fmt.Sprintf("%s: %s", altMsg, e.ErrMsg)
	return e
}

// AddTraceRow() add new trace line in ErrExtender trace array
// func (e *extendedErr) AddTraceRow() ErrExtender {
// 	w := where()
// 	r := e.trace[len(e.trace)-1]
// 	if w.Package == r.Package && w.Function == r.Function {
// 		return e
// 	}
// 	e.trace = append(e.trace, w)
// 	return e
// }

// // TraceRawString() return string from trace array.
// // Every trace line separated by slash.
// func (e *extendedErr) TraceRawString() string {
// 	result := ""
// 	for _, row := range e.trace {
// 		result = path.Join(result, fmt.Sprintf("%s:%s:%s:%d",
// 			row.Package, row.File, row.Function, row.Line))
// 	}
// 	return result
// }

// // TraceRawString() return tagged string from trace array.
// // Every trace line separated by slash.
// // Format: {pkg}:{file}:{function}:{line}
// func (e *extendedErr) TraceTagged() string {
// 	result := ""
// 	for _, row := range e.trace {
// 		result = path.Join(result, fmt.Sprintf("{pkg}%s:{file}%s:{function}%s:{line}%d",
// 			row.Package, row.File, row.Function, row.Line))
// 	}
// 	return result
// }

// TraceRawString() return JSON-string from trace array.
func (e *extendedErr) TraceJSON() string {
	res, _ := json.Marshal(e.trace)
	return string(res)
}

// Wrap() unite two ErrExtender objects
// Example: err.Wrap(err2)
// func (e *extendedErr) Wrap(w ErrExtender) ErrExtender {
// 	e.msg = fmt.Sprintf("%s: %s", e.msg, w.Error())
// 	e.altMsg = fmt.Sprintf("%s: %s", e.altMsg, w.GetAltMsg())
// 	e.errCode = w.GetErrCode()
// 	e.trace = append(e.trace, w.GetTraceRows()...)
// 	return e
// }

// New() create ErrExtender with {msg} message and 1 trace line.
// Example: err := New("Error message")
func New(msg string) ErrExtender {
	return &extendedErr{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  msg,
		trace:     []string{where()},
	}
}

// Newf() create ErrExtender with {msg} message and 1 trace line.
// Newf() allows you to format numbers, variables and strings into the first {format} parameter you give it
// Example: err := Newf("Error: %s with code %d", "SQL error", 1005)
func Newf(format string, a ...interface{}) ErrExtender {
	return &extendedErr{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  fmt.Errorf(format, a...).Error(),
		trace:     []string{where()},
	}
}

func Wrap(msg string, err error) ErrExtender {
	if e, ok := err.(*extendedErr); ok {
		e.Causes = fmt.Sprintf("%s: %s", msg, e.Causes)
		e.trace = append(e.trace, where())
		return e
	}
	return &extendedErr{
		ErrError:  InternalServerError.Error(),
		ErrMsg:    msg,
		Causes:    err.Error(),
		ErrStatus: http.StatusInternalServerError,
		trace:     []string{where()},
	}
}

// With NewWithErr() You can join {msg} to {builtin.error} message
// Example: err := NewWithErr("SQL Error: ", err)
func NewWithErr(msg string, err error) ErrExtender {
	return &extendedErr{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  fmt.Sprintf("%s: %s", msg, err),
		trace:     []string{where()},
	}
}

// NewWithAlt() create ErrExtender with main and alternative messages.
// Example: err := NewWithAlt("SQL connection error", "<SQL_CONNECTION_ERROR>")
func NewWithAlt(err, msg string) ErrExtender {
	return &extendedErr{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  err,
		ErrMsg:    msg,
		trace:     []string{where()},
	}
}

// NewWithType() create ErrExtender with ErrType (error code)
// Example: err := NewWithType("SQL connection error", "<SQL_CONNECTION_ERROR>", ErrType(1005))
func NewWithType(err error, causes string, msg string, t int) ErrExtender {
	return &extendedErr{
		ErrStatus: t,
		ErrError:  err.Error(),
		Causes:    causes,
		ErrMsg:    msg,
		trace:     []string{where()},
	}
}

// NewWithExtErr() create new ErrExtender from current {err} ErrExtender object
// with joining {msg} {altMsg} {errType} and {trace} stack
// Example: err := NewWithExtErr("SQL auth error", err)
func NewWithExtErr(msg string, err ErrExtender) ErrExtender {
	return &extendedErr{
		ErrStatus: err.GetErrCode(),
		ErrError:  fmt.Sprintf("%s: %s", msg, err),
		ErrMsg:    err.GetAltMsg(),
		trace:     nil,
	}
}

func ParseErr(err error) ErrExtender {
	if e, ok := err.(*extendedErr); ok {
		return e
	}
	return &extendedErr{
		ErrError:  InternalServerError.Error(),
		ErrMsg:    "",
		Causes:    err.Error(),
		ErrStatus: http.StatusInternalServerError,
		trace:     []string{where()},
	}
}

func ErrorResponse(err error) (int, interface{}) {
	res := ParseErr(err)

	return res.GetErrCode(), res
}

// func where() traceRow {
// 	// pc, file, line, _ := runtime.Caller(2)
// 	// function := runtime.FuncForPC(pc).Name()

// 	// slashIndex := strings.LastIndex(function, "/")
// 	// function = function[slashIndex+1:]

// 	// s := strings.Split(function, ".")
// 	// packageName, function := s[0], s[1]

// 	// slashIndex = strings.LastIndex(file, "/")
// 	// fileName := file[slashIndex+1:]

// 	_, filename, line, _ := runtime.Caller(1)

// 	return traceRow{
// 		Where: fmt.Sprintf("%s:%d", filename, line),
// 	}
// }

func where() string {

	pc, filename, line, _ := runtime.Caller(2)
	trace := fmt.Sprintf("%s:%d", filename, line)
	function := runtime.FuncForPC(pc).Name()

	slashIndex := strings.LastIndex(function, "/")
	function = function[slashIndex+1:]

	s := strings.Split(function, ".")
	packageName, function := s[0], s[1]

	return fmt.Sprintf("%s %s %s", trace, packageName, function)
}
