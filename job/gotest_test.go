package job

import (
	"runtime"
	"testing"
)

func must(t *testing.T, ok bool) {
	if !ok {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\n unexcepted: %s:%d", fileName, line)
		t.FailNow()
	}
}

func TestForTest(t *testing.T) {
	Test()
}
