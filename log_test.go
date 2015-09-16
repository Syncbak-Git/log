package log_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Syncbak-Git/log"
)

func TestLog(t *testing.T) {
	var buff bytes.Buffer
	log.SetOutput(&buff)
	log.SetLogLevel(log.LevelAll ^ log.LevelDebug ^ log.LevelFatal)
	log.Debug("Hello %s", "world") // this won't get written
	log.Info("%s %d", "Hello world", 1234)
	log.Error("%v", struct{ s string }{"Hello world"})
	log.Fatal("%v", struct{ s string }{"Hello world"}) // won't trigger or write
	log.Custom("TEST", "Hello world", "extra arg")
	// we do Panic later
	b := buff.String()
	if c := strings.Count(b, "Hello"); c != 3 {
		t.Errorf("Bad \"Hello\" count: %d\n%s", c, b)
	}
	if strings.Count(b, "DEBUG") != 0 {
		t.Error("DEBUG shouldn't have been written")
	}
	if strings.Count(b, "INFO") != 1 {
		t.Error("INFO wasn't written")
	}
	if strings.Count(b, "ERROR") != 1 {
		t.Error("ERROR wasn't written")
	}
	if strings.Count(b, "FATAL") != 0 {
		t.Error("FATAL shouldn't have been written")
	}
	if strings.Count(b, "TEST") != 1 {
		t.Error("TEST wasn't written")
	}
	t.Log(b)
	buff.Truncate(0)
	log.SetLogLevel(log.LevelAll)
	defer func() {
		recover()
		if !strings.Contains(buff.String(), "Hello world 1234") {
			t.Errorf("Bad Panic() string: %s", buff.String())
		}
	}()
	log.Panic("%s %d", "Hello world", 1234)
	t.Error("Panic didn't panic")
}

func BenchmarkLog_basic(b *testing.B) {
	err := log.SetOutputFile(os.DevNull)
	if err != nil {
		b.Fatalf("Could not set output file: %s", err)
	}
	log.SetLogLevel(log.LevelAll)
	for n := 0; n < b.N; n++ {
		log.Error("Hello world")
	}
}

func BenchmarkLog_struct(b *testing.B) {
	err := log.SetOutputFile(os.DevNull)
	if err != nil {
		b.Fatalf("Could not set output file: %s", err)
	}
	log.SetLogLevel(log.LevelAll)
	for n := 0; n < b.N; n++ {
		log.Error("%s %d %v", "Hello", 1234, struct{ s string }{"World"})
	}
}

func BenchmarkLog_formatting(b *testing.B) {
	err := log.SetOutputFile(os.DevNull)
	if err != nil {
		b.Fatalf("Could not set output file: %s", err)
	}
	log.SetLogLevel(log.LevelAll)
	for n := 0; n < b.N; n++ {
		log.Error("%s %d", "Hello", 1234)
	}
}

func Example() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelAll ^ log.LevelDebug)
	log.SetTimestamp(func() string { return "2006-01-02T15:04:05.999999999Z" })
	err := log.Debug("This will not be written")
	if err != nil {
		log.Error("Could not write DEBUG entry: %s", err)
	}
	log.Info("Hello %s", "world")
	// Output: 2006-01-02T15:04:05.999999999Z	INFO	Hello world
}
