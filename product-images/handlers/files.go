package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-images/files"
)

type Files struct {
	log     hclog.Logger
	store *files.Local
}

func NewFiles(log hclog.Logger, store *files.Local) *Files {
	return &Files{log, store}
}

func (f *Files) SaveFile(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	filename := params["filename"]

	f.log.Info("Handle POST files", "id", id, "filename", filename)

	path := filepath.Join(id, filename)

	// defer r.Body.Close()

	err := f.store.Save(path, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (f *Files) SaveFileMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Unable to parse multipart form", "error", err)
		http.Error(rw, "Unable to parse form", http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		f.log.Error("Invalid id", "error", err)
		http.Error(rw, "Invalid Id", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Unable to read file from form", "error", err)
		http.Error(rw, "Unable to read file from form", http.StatusBadRequest)
		return
	}

	path := filepath.Join(id, mh.Filename)
	err = f.store.Save(path, ff)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
		return
	}
}
