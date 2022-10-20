package files

import (
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func Test_save(t *testing.T) {
	log := hclog.Default()
	l, err := NewLocal(log, "./filestore")
	assert.NoError(t, err)

	f, err := os.Open("/home/hossein/workspace/rnd/microservices-go/product-images/hossein.jpg")
	assert.NoError(t, err)
	defer f.Close()

	err = l.Save("/1/hossein.jpg", f)
	log.Error("unable to save file", "error", err)
}

func Test_Copyfile(t *testing.T) {
	log := hclog.Default()

	src := "/home/hossein/workspace/rnd/microservices-go/product-images/hossein.jpeg"
    dst := "/home/hossein/workspace/rnd/microservices-go/product-images/hossein_copy.jpeg"

    fin, err := os.Open(src)
    if err != nil {
        log.Error("%w", err)
    }
    defer fin.Close()

    fout, err := os.Create(dst)
    if err != nil {
        log.Error("%w", err)
    }
    defer fout.Close()

    _, err = io.Copy(fout, fin)

    if err != nil {
        log.Error("%w", err)
    }
}