package log_test

import (
	"bytes"
	"fmt"
	"io"
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
	log.Close() // should have no effect, because Buffer is not a WriteCloser
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

func TestClose(t *testing.T) {
	var buff bytes.Buffer
	log.SetOutput(&buff)
	log.SetLogLevel(log.LevelAll)
	err := log.Close() // should have no effect, because Buffer is not a WriteCloser
	if err != nil {
		t.Errorf("Close shouldn't have returned error: %s", err)
	}
	log.Info("Hello")
	if !strings.Contains(buff.String(), "Hello") {
		t.Errorf("Close killed output of Buffer: %s", buff.String())
	}
	read, write := io.Pipe()
	log.SetOutput(write)
	go func() { // we need to read before writing to the pipe
		b := make([]byte, 100)
		n, err := read.Read(b)
		if !strings.Contains(string(b), "Hello") {
			t.Errorf("Pipe read failed: %d, %s, %v", n, string(b), err)
		}
		n, err = read.Read(b)
		if n != 0 || err != io.EOF {
			t.Errorf("Close didn't work: %d bytes read, %s (%v)", n, string(b), err)
		}
	}()
	log.Info("Hello")
	err = log.Close()
	if err != nil {
		t.Errorf("Close shouldn't have returned error on pipe: %s", err)
	}
	log.Info("Hello")
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

// Example of using the global log.
func Example() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelAll ^ log.LevelDebug)
	// we call SetTimestamp so that the timestamp will be deterministic
	log.SetTimestamp(func() string { return "2006-01-02T15:04:05.999999999Z" })
	err := log.Debug("This will not be written")
	if err != nil {
		log.Error("Could not write DEBUG entry: %s", err)
	}
	log.Info("Hello %s", "world")
	// Output: 2006-01-02T15:04:05.999999999Z	INFO	Hello world
}

// Example of using a private log.
func ExampleLog() {
	l := log.NewLog()
	var buf bytes.Buffer
	l.SetOutput(&buf)
	l.SetLogLevel(log.LevelAll ^ log.LevelDebug)
	// we call SetTimestamp so that the timestamp will be deterministic
	l.SetTimestamp(func() string { return "2006-01-02T15:04:05.999999999Z" })
	err := l.Debug("This will not be written")
	if err != nil {
		log.Error("Could not write DEBUG entry: %s", err)
	}
	l.Info("Hello %s", "world")
	fmt.Print(buf.String())
	l.Close()
	// Output: 2006-01-02T15:04:05.999999999Z	INFO	Hello world
}
