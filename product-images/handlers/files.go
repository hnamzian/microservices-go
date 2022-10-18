package handlers

import "log"

type Files struct {
	l *log.Logger
}

func NewFiles(l *log.Logger) *Files {
	return &Files{l}
}
