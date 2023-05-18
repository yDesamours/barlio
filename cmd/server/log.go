package main

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

const (
	INFOS logType = "INFOS"
	ERROR logType = "ERROR"
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

func (logger *Logger) WriteInfos(message string, params map[string]interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	log := newLog(INFOS, message, params)
	b, _ := json.Marshal(log)
	b = append(b, '\n')

	logger.out.Write(b)
}
