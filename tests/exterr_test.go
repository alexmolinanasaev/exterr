package tests

import (
	"testing"

	"github.com/alexmolinanasaev/exterr"
)

func TestNew(t *testing.T) {
	var message string = "Test message"
	var err = exterr.NewWithAlt(message, "alternative message")
	if message != err.Error() {
		t.Error(
			"expected", message,
			"got", err,
		)
	}
}

func TestTrace(t *testing.T) {
	err := exterr.New("trace does not match to expected")
	expectedTrace := "tests:exterr_test.go:TestTrace:21"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
}

func TestMultiTrace(t *testing.T) {
	err := MultiTraceFunc1()
	expectedTrace := "tests:exterr_test.go:MultiTraceFunc1:41/tests:exterr_test.go:MultiTraceFunc3:49"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
}

func MultiTraceFunc1() exterr.ErrExtender {
	return MultiTraceFunc2().AddTrace()
}

func MultiTraceFunc2() exterr.ErrExtender {
	return MultiTraceFunc3() // <- no add trace here
}

func MultiTraceFunc3() exterr.ErrExtender {
	return exterr.New("MultiTraceError")
}

func TestMultiCheckTrace(t *testing.T) {
	err := MultiCheckTraceFunc1()
	expectedTrace := "tests:exterr_test.go:MultiCheckTraceFunc1:71/tests:exterr_test.go:MultiCheckTraceFunc3:79"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
	err.AddTrace().AddTrace().AddTrace() // add only one trace
	err.AddTrace().AddTrace()            // skip all
	expectedTrace = "tests:exterr_test.go:TestMultiCheckTrace:60/tests:exterr_test.go:MultiCheckTraceFunc1:71/tests:exterr_test.go:MultiCheckTraceFunc3:79"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
}

func MultiCheckTraceFunc1() exterr.ErrExtender {
	return MultiCheckTraceFunc2().AddTrace().AddTrace() // add only one trace
}

func MultiCheckTraceFunc2() exterr.ErrExtender {
	return MultiCheckTraceFunc3() // <- no add trace here
}

func MultiCheckTraceFunc3() exterr.ErrExtender {
	return exterr.New("MultiCheckTraceError")
}
