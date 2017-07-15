package logs

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	var log = NewLogger("main")
	if log == nil {
		t.Error("Cannot create logger")
	}
}

func Test_NewLoggerJSON(t *testing.T) {

	oldosGetenv := OsGetenv
	// as we are exiting, revert sqlOpen back to oldSqlOpen at end of function
	defer func() { OsGetenv = oldosGetenv }()
	var newVal = func(key string) string {
		return "JSON"
	}
	OsGetenv = newVal

	var log = NewLogger("main")
	if log == nil {
		t.Error("Cannot create logger")
	}

}

func Test_SubLogger(t *testing.T) {
	var log = NewLogger("main")
	var log1 = log.SubLogger("subLogger")
	if log1 == nil {
		t.Error("Cannot create logger")
	}
}