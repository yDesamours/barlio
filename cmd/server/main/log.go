package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

const (
	INFOS logType = "INFOS"
	ERROR logType = "ERROR"
	PANIC logType = "PANIC"
)

type logType string

type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

type Log struct {
	Type    logType                `json:"type"`
	Time    time.Time              `json:"time"`
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params"`
}

func newLogger(out io.Writer) *Logger {
	return &Logger{
		out: out,
	}
}

func newLog(logType logType, message string, params map[string]interface{}) *Log {
	return &Log{
		Time:    time.Now(),
		Type:    logType,
		Message: message,
		Params:  params,
	}
}

func (logger *Logger) writelog(logType logType, message string, params map[string]interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	log := newLog(logType, message, params)
	b, _ := json.Marshal(log)
	b = append(b, '\n')

	logger.out.Write(b)
}

func (logger *Logger) Fatal(err error) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	logger.out.Write([]byte(err.Error()))
	os.Exit(-1)
}

func (logger *Logger) Write(p []byte) (n int, err error) {
	return logger.out.Write(p)
}

func (app *App) setLog(log *Logger) {
	app.Logger = log
}

func (app *App) infos(message string, params map[string]interface{}) {
	app.Logger.writelog(INFOS, message, params)
}

func (app *App) error(err error) {
	_, file, line, _ := runtime.Caller(1)
	app.Logger.writelog(ERROR, err.Error(), map[string]interface{}{"file": file, "line": line - 2})
}

func (app *App) panic(err any) {
	stack := debug.Stack()
	app.Logger.writelog(PANIC, fmt.Sprintln(err), map[string]interface{}{"stack": string(stack)})
}
