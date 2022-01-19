package request_dispatch

import (
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelDebug = LogLevel(3)
	LogLevelInfo  = LogLevel(2)
	LogLevelError = LogLevel(1)
)

// Logger
// traefik disables unsafe and syscall, which makes many logger libraries unusable
// Ref https://github.com/tomMoulard/fail2ban/blob/main/fail2ban.go#L35-L38 a simple logger is implemented
type Logger struct {
	logLevel    LogLevel
	loggerInfo  *log.Logger
	loggerDebug *log.Logger
	loggerError *log.Logger
}

func (l *Logger) Debug(args ...interface{}) {
	if l.logLevel >= LogLevelDebug {
		l.loggerDebug.Println(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.logLevel >= LogLevelInfo {
		l.loggerInfo.Println(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.logLevel >= LogLevelError {
		l.loggerError.Println(args...)
	}
}

func getLogLevel(level string) LogLevel {
	switch level {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "ERROR":
		return LogLevelError
	default:
		return LogLevelError
	}
}

func NewLogger(level string) *Logger {
	return &Logger{
		logLevel:    getLogLevel(level),
		loggerInfo:  log.New(os.Stdout, "INFO: [AB_TEST] ", log.Ldate|log.Ltime|log.Lshortfile),
		loggerDebug: log.New(os.Stdout, "DEBU: [AB_TEST] ", log.Ldate|log.Ltime|log.Lshortfile),
		loggerError: log.New(os.Stdout, "ERRO: [AB_TEST] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
