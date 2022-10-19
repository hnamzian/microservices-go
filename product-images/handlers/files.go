package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-images/files"
)

type Files struct {
	l     *log.Logger
	store *files.Local
}

func NewFiles(l *log.Logger, store *files.Local) *Files {
	return &Files{l, store}
}

func (f *Files) SaveFile(rw http.ResponseWriter, r *http.Request) {
	f.l.Println("Handle POST files")

	params := mux.Vars(r)
	id := params["id"]
	filename := params["filename"]

	path := filepath.Join(id, filename)

	f.store.Save(path, r.Body)
}
