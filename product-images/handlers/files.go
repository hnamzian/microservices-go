package handlers

import (
	hclog "github.com/hashicorp/go-hclog"
)

type Files struct {
	l hclog.Logger
}

func NewFiles(l hclog.Logger) *Files {
	return &Files{l}
}
