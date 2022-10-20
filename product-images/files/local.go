package files

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"golang.org/x/xerrors"
)

type Local struct {
	l        hclog.Logger
	basePath string
}

func NewLocal(l hclog.Logger, basePath string) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{l: l, basePath: p}, nil
}

func (l *Local) Save(path string, content io.Reader) error {
	fullPath := filepath.Join(l.basePath, path)

	// create directory if does not exist
	d := filepath.Dir(fullPath)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directories %w", err)
	}

	// Remove file if exists
	_, err = os.Stat(fullPath)
	if err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return xerrors.Errorf("Unable to remove existing file %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to get file info %w", err)
	}

	// // create new file
	// f, err := os.Create(fullPath)
	// if err != nil {
	// 	return xerrors.Errorf("Unable to Create new file %w", err)
	// }
	// defer f.Close()

	// // copy content to file
	// wl, err := io.Copy(f, content)
	// if err != nil {
	// 	return xerrors.Errorf("Unable to write content to file %w", err)
	// }

	// l.l.Info("Bytes written", "length", wl)
	// return nil
	// return copyToFileFromReader(fullPath, content)

	return copyToFile(fullPath, content)
}

func copyToFileFromReader(path string, content io.Reader) error {
	// create new file
	f, err := os.Create(path)
	if err != nil {
		return xerrors.Errorf("Unable to Create new file %w", err)
	}
	defer f.Close()

	// copy content to file
	_, err = io.Copy(f, content)
	if err != nil {
		return xerrors.Errorf("Unable to write content to file %w", err)
	}

	return nil
}

func copyToFile(path string, content io.Reader) error {
	data, err := io.ReadAll(content)
	if err != nil {
		return xerrors.Errorf("Unable to read data from body %w", err)
	}

	// create new file
	f, err := os.Create(path)
	if err != nil {
		return xerrors.Errorf("Unable to Create new file %w", err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return xerrors.Errorf("Unable to write content to file %w", err)
	}

	return nil
}