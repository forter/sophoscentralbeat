package lrlog

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	infoTagPrefix    = "[INFO] "
	warningTagPrefix = "[WARNING] "
	errorTagPrefix   = "[ERROR] "
	fatalTagPrefix   = "[FATAL] "
)

var verbosity int

func init() {
	SetLogWriter()
}

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprintf(os.Stderr, "%s %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), string(bytes))
}

// SetLogWriter sets the logging output format
func SetLogWriter() {
	log.SetOutput(new(logWriter))
}

// SetVerbosity sets the global verbosity level
func SetVerbosity(level int) {
	verbosity = level
}

// Verbose is a boolean alias that allows for convenient inline logging
type Verbose bool

// V returns true if the logging is set to `level` or higher, false otherwise
func V(level int) Verbose {
	return Verbose(verbosity >= level)
}

// Info is just a call to Info wrapped in a boolean
func (v Verbose) Info(msg string) {
	if v {
		Info(msg)
	}
}

// Infof is just a call to Infof wrapped in a boolean
func (v Verbose) Infof(msg string, objs ...interface{}) {
	if v {
		Infof(msg, objs...)
	}
}

// Warning is just a call to Warning wrapped in a boolean
func (v Verbose) Warning(msg string) {
	if v {
		Warning(msg)
	}
}

// Warningf is just a call to Warningf wrapped in a boolean
func (v Verbose) Warningf(msg string, objs ...interface{}) {
	if v {
		Warningf(msg, objs...)
	}
}

// Error is just a call to Error wrapped in a boolean
func (v Verbose) Error(msg string) {
	if v {
		Error(msg)
	}
}

// Errorf is just a call to Errorf wrapped in a boolean
func (v Verbose) Errorf(msg string, objs ...interface{}) {
	if v {
		Errorf(msg, objs...)
	}
}

// Fatal is just a call to Fatal wrapped in a boolean
func (v Verbose) Fatal(msg string) {
	if v {
		Fatal(msg)
	}
}

// Fatalf is just a call to Fatalf wrapped in a boolean
func (v Verbose) Fatalf(msg string, objs ...interface{}) {
	if v {
		Fatalf(msg, objs...)
	}
}

// Info logs args at info level
func Info(args ...interface{}) {
	log.Print(infoTagPrefix, args)
}

// Infof logs a formatted info message with args
func Infof(format string, args ...interface{}) {
	log.Printf(infoTagPrefix+format, args...)
}

// Warning logs args at warning level
func Warning(args ...interface{}) {
	log.Print(warningTagPrefix, args)
}

// Warningf logs a formatted warning message with args
func Warningf(format string, args ...interface{}) {
	log.Printf(warningTagPrefix+format, args...)
}

// Error logs args at error level
func Error(args ...interface{}) {
	log.Print(errorTagPrefix, args)
}

// Errorf logs a formatted error message with args
func Errorf(format string, args ...interface{}) {
	log.Printf(errorTagPrefix+format, args...)
}

// Fatal logs args at fatal level
func Fatal(objs ...interface{}) {
	log.Print(fatalTagPrefix, objs)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message with args
func Fatalf(format string, args ...interface{}) {
	log.Printf(fatalTagPrefix+format, args...)
	os.Exit(1)
}
