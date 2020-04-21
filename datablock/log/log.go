package log // import "github.com/neoul/gostudy/datablock/log"

import (
	"io"

	"github.com/op/go-logging"
)

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Password string

// Redacted - what?
func (p Password) Redacted() interface{} {
    return logging.Redact(string(p))
}

// Log - Wrapper structure for logging.Logger
type Log struct {
    *logging.Logger
}

// NewLog return *Logger to log
func NewLog(module string, out io.Writer) Log {
    // log - Logging object
    var log = logging.MustGetLogger(module)

    // Example format string. Everything except the message has a custom color
    // which is dependent on the log level. Many fields have a custom output
    // formatting too, eg. the time returns the hour down to the milli second.
    var format = logging.MustStringFormatter(
        `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
    )
    // For demo purposes, create two backend for out.
    backend1 := logging.NewLogBackend(out, "", 0)
    backend2 := logging.NewLogBackend(out, "", 0)

    // For messages written to backend2 we want to add some additional
    // information to the output, including the used log level and the name of
    // the function.
    backend2Formatter := logging.NewBackendFormatter(backend2, format)

    // Only errors and more severe messages should be sent to backend1
    backend1Leveled := logging.AddModuleLevel(backend1)
    backend1Leveled.SetLevel(logging.ERROR, "")

    // Set the backends to be used.
    logging.SetBackend(backend1Leveled, backend2Formatter)

    return Log {Logger: log}
}


// Fatal is equivalent to l.Logger.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func (l Log) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// Fatalf is equivalent to l.Logger.Critical followed by a call to os.Exit(1).
func (l Log) Fatalf(format string, args ...interface{}) {
    l.Logger.Fatalf(format, args...)
}

// Panic is equivalent to l.Logger.Critical(fmt.Sprint()) followed by a call to panic().
func (l Log) Panic(args ...interface{}) {
	l.Logger.Panic(args...)
}

// Panicf is equivalent to l.Logger.Critical followed by a call to panic().
func (l Log) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// Critical logs a message using logging.CRITICAL as log level.
func (l Log) Critical(args ...interface{}) {
	l.Logger.Critical(args...)
}

// Criticalf logs a message using logging.CRITICAL as log level.
func (l Log) Criticalf(format string, args ...interface{}) {
	l.Logger.Criticalf(format, args...)
}

// Error logs a message using logging.ERROR as log level.
func (l Log) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

// Errorf logs a message using logging.ERROR as log level.
func (l Log) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Warning logs a message using logging.WARNING as log level.
func (l Log) Warning(args ...interface{}) {
	l.Logger.Warning(args...)
}

// Warningf logs a message using logging.WARNING as log level.
func (l Log) Warningf(format string, args ...interface{}) {
	l.Logger.Warningf(format, args...)
}

// Notice logs a message using logging.NOTICE as log level.
func (l Log) Notice(args ...interface{}) {
	l.Logger.Notice(args...)
}

// Noticef logs a message using logging.NOTICE as log level.
func (l Log) Noticef(format string, args ...interface{}) {
	l.Logger.Noticef(format, args...)
}

// Info logs a message using logging.INFO as log level.
func (l Log) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

// Infof logs a message using logging.INFO as log level.
func (l Log) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Debug logs a message using logging.DEBUG as log level.
func (l Log) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

// Debugf logs a message using logging.DEBUG as log level.
func (l Log) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}