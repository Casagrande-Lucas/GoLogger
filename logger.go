package GoLogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	file   *os.File
	writer io.Writer
}

func NewLogger(config LoggerConfig) (*Logger, error) {
	var writer io.Writer = os.Stdout

	if config.FilePath != "" {
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("could not open log file: %v", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
		return &Logger{level: config.Level, file: file, writer: writer}, nil
	}

	return &Logger{level: config.Level, writer: writer}, nil
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level >= l.level {
		log.SetOutput(l.writer)
		log.Printf("[%s] %s", level.String(), fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.log(TraceLevel, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DebugLevel, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(InfoLevel, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WarnLevel, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ErrorLevel, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FatalLevel, format, v...)
	os.Exit(1)
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
