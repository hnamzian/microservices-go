package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-images/files"
)

type Files struct {
	l    hclog.Logger
	stor *files.Local
}

func NewFiles(stor *files.Local, l hclog.Logger) *Files {
	return &Files{stor: stor, l: l}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	filename := params["filename"]

	f.l.Info("Handle POST", "id", id, "filename", filename)
	f.saveFile(id, filename, rw, r)
}

func (f *Files) saveFile(id string, filename string, rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Save file for product", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)
	err := f.stor.Save(fp, r.Body)
	if err != nil {
		f.l.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
