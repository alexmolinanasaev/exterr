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
	err := exterr.New("Simple Error")
	expectedTrace := "tests:exterr_test.go:TestTrace:21"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
	expectedJSONTrace := "{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"TestTrace\",\"line\":\"21\",\"child\":{\"package\":\"\",\"file\":\"\",\"function\":\"\",\"line\":\"\",\"child\":null}}"
	if expectedJSONTrace != err.TraceJSON() {
		t.Logf("Expected JSON: %s\n", expectedJSONTrace)
		t.Logf("Got JSON: %s\n", err.TraceJSON())
		t.Errorf(err.Error())
	}
}

func TestMultiTrace(t *testing.T) {
	err := MultiTraceFunc1()
	expectedTrace := "tests:exterr_test.go:MultiTraceFunc1:53/tests:exterr_test.go:MultiTraceFunc3:61"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
	expectedJSONTrace := "{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiTraceFunc1\",\"line\":\"53\",\"child\":{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiTraceFunc3\",\"line\":\"61\",\"child\":null}}"
	if expectedJSONTrace != err.TraceJSON() {
		t.Logf("Expected JSON: %s\n", expectedJSONTrace)
		t.Logf("Got JSON: %s\n", err.TraceJSON())
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
	expectedTrace := "tests:exterr_test.go:MultiCheckTraceFunc1:96/tests:exterr_test.go:MultiCheckTraceFunc3:104"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
	expectedJSONTrace := "{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiCheckTraceFunc1\",\"line\":\"96\",\"child\":{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiCheckTraceFunc3\",\"line\":\"104\",\"child\":null}}"
	if expectedJSONTrace != err.TraceJSON() {
		t.Logf("Expected JSON: %s\n", expectedJSONTrace)
		t.Logf("Got JSON: %s\n", err.TraceJSON())
		t.Errorf(err.Error())
	}

	err.AddTrace().AddTrace().AddTrace() // add only one trace
	err.AddTrace().AddTrace()            // skip all
	expectedTrace = "tests:exterr_test.go:TestMultiCheckTrace:79/tests:exterr_test.go:MultiCheckTraceFunc1:96/tests:exterr_test.go:MultiCheckTraceFunc3:104"
	if expectedTrace != err.TraceRaw() {
		t.Logf("Expected: %s\n", expectedTrace)
		t.Logf("Got: %s\n", err.TraceRaw())
		t.Errorf(err.Error())
	}
	expectedJSONTrace = "{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"TestMultiCheckTrace\",\"line\":\"79\",\"child\":{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiCheckTraceFunc1\",\"line\":\"96\",\"child\":{\"package\":\"tests\",\"file\":\"exterr_test.go\",\"function\":\"MultiCheckTraceFunc3\",\"line\":\"104\",\"child\":null}}}"
	if expectedJSONTrace != err.TraceJSON() {
		t.Logf("Expected JSON: %s\n", expectedJSONTrace)
		t.Logf("Got JSON: %s\n", err.TraceJSON())
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
