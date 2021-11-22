package tests

import (
	"fmt"
	"testing"

	"github.com/alexmolinanasaev/exterr"
)

func errMessage(expected interface{}, got interface{}) string {
	return fmt.Sprintf("\nEXPECTED: %v\nGOT:      %v", expected, got)
}

func TestNew(t *testing.T) {
	for i, tt := range []struct {
		in        string
		out       string
		isCorrect bool
	}{
		{"Error", "Error", true},
		{"", "", true},
		{"Error", "", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := exterr.New(tt.in).Error()
			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestNewWithAlt(t *testing.T) {
	for i, tt := range []struct {
		in        string
		out       string
		isCorrect bool
	}{
		{"AltError", "AltError", true},
		{"", "", true},
		{"AltError", "", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := exterr.NewWithAlt("Error", tt.in).GetAltMsg()
			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestNewWithType(t *testing.T) {
	for i, tt := range []struct {
		in        int
		out       int
		isCorrect bool
	}{
		{0, 0, true},
		{1, 1, true},
		{1, 0, false},
		{0, 1, false},
		{-100, -100, true},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := exterr.NewWithType("Error", "AltError", tt.in)
			if (result.GetErrCode() == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestNewWithExtErr(t *testing.T) {
	for i, tt := range []struct {
		in1       string
		in2       string
		out       string
		isCorrect bool
	}{
		{"Error1", "Error2", "Error1: Error2", true},
		{"Error1", "Error2", "Error1:Error2", false},
		{"Error2", "Error1", "Error1: Error2", false},
		{"", "", ": ", true},
		{"Error", "", "Error: ", true},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.New(tt.in2)
			result := exterr.NewWithExtErr(tt.in1, err2).Error()
			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestWrapMsg(t *testing.T) {
	for i, tt := range []struct {
		in1       string
		in2       string
		out       string
		isCorrect bool
	}{
		{"Error1", "Error2", "Error1: Error2", true},
		{"Error1", "Error2", "Error1:Error2", false},
		{"Error2", "Error1", "Error1: Error2", false},
		{"", "", ": ", true},
		{"Error", "", "Error: ", true},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.New(tt.in2)
			err := exterr.New(tt.in1)
			err.Wrap(err2)
			result := err.Error()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestWrapAltMsg(t *testing.T) {
	for i, tt := range []struct {
		in1       string
		in2       string
		out       string
		isCorrect bool
	}{
		{"AltMsg1", "AltMsg2", "AltMsg1: AltMsg2", true},
		{"AltMsg1", "AltMsg2", "AltMsg1:AltMsg2", false},
		{"AltMsg1", "AltMsg2", "AltMsg2: AltMsg1", false},
		{"", "", ": ", true},
		{"AltMsg1", "", "AltMsg1: ", true},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.NewWithAlt("Error", tt.in2)
			err := exterr.NewWithAlt("Error", tt.in1)
			err.Wrap(err2)
			result := err.GetAltMsg()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestTraceRawString(t *testing.T) {
	for i, tt := range []struct {
		out       string
		isCorrect bool
	}{
		{"tests:exterr_test.go:TestTraceRawString:158/tests:exterr_test.go:TestTraceRawString:157", true},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.NewWithAlt("", "")
			err := exterr.NewWithAlt("", "")
			err.Wrap(err2)
			result := err.TraceRawString()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}
func TestTraceTagged(t *testing.T) {
	for i, tt := range []struct {
		out       string
		isCorrect bool
	}{
		{"{pkg}tests:{file}exterr_test.go:{function}TestTraceTagged:{line}178/{pkg}tests:{file}exterr_test.go:{function}TestTraceTagged:{line}177", true},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.NewWithAlt("", "")
			err := exterr.NewWithAlt("", "")
			err.Wrap(err2)
			result := err.TraceTagged()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}
func TestTraceJSON(t *testing.T) {
	for i, tt := range []struct {
		out       string
		isCorrect bool
	}{
		{"[{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"TestTraceJSON\",\"line\":198},{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"TestTraceJSON\",\"line\":197}]", true},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err2 := exterr.NewWithAlt("", "")
			err := exterr.NewWithAlt("", "")
			err.Wrap(err2)
			result := err.TraceJSON()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func TestTrace(t *testing.T) {
	for i, tt := range []struct {
		out       string
		isCorrect bool
	}{
		{"tests:exterr_test.go:f3:232/tests:exterr_test.go:f2:229/tests:exterr_test.go:f1:228", true},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err := f1()
			result := err.TraceRawString()

			if (result == tt.out) != tt.isCorrect {
				t.Error(errMessage(tt.out, result))
			}
		})
	}
}

func f1() exterr.ErrExtender { return f2().AddTraceRow() }
func f2() exterr.ErrExtender { return f3().AddTraceRow().AddTraceRow() }
func f3() exterr.ErrExtender {
	return func() exterr.ErrExtender {
		return exterr.New("").AddTraceRow()
	}()
}
