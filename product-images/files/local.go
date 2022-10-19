package files

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	l        *log.Logger
	basePath string
}

func NewLocal(l *log.Logger, basePath string) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{l: l, basePath: p}, nil
}

func (l *Local) Save(path string, content io.Reader) {
	fullPath := filepath.Join(l.basePath, path)

	// create directory if does not exist
	d := filepath.Dir(fullPath)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		l.l.Println("Unable to create directories", err)
		return
	}

	// Remove file if exists
	_, err = os.Stat(fullPath)
	if err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			l.l.Println("Unable to remove existing file", err)
			return
		}
	} else if !os.IsNotExist(err) {
		l.l.Println("Unable to get file info", err)
		return
	}

	// create new file
	f, err := os.Create(fullPath)
	if err != nil {
		l.l.Println("Unable to Create new file")
		return
	}
	defer f.Close()

	// copy content to file
	wl, err := io.Copy(f, content)
	if err != nil {
		l.l.Println("Unable to write content to file", err)
		return
	}

	l.l.Println("%i Bytes written", wl)
}
