package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-images/files"
)

type Files struct {
	l     hclog.Logger
	store *files.Local
}

func NewFiles(l hclog.Logger, store *files.Local) *Files {
	return &Files{l, store}
}

func (f *Files) SaveFile(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	filename := params["filename"]

	f.l.Info("Handle POST files", "id", id, "filename", filename)

	path := filepath.Join(id, filename)

	f.store.Save(path, r.Body)
}
