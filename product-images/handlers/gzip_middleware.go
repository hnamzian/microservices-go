package handlers

import (
	gzip "compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func NewGzipHandler() *GzipHandler {
	return &GzipHandler{}
}

func (g *GzipHandler) GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		he := r.Header.Get("Accept-Encoding")
		fmt.Println("accept encoding", he)

		if strings.Contains(he, "gzip") {
			fmt.Println("Gzip")
			gw := NewGzipResponseWriterWrapper(rw)

			h.ServeHTTP(gw, r)
			
			defer gw.Flush()
			
			return
		}
		fmt.Println("without compression")

		h.ServeHTTP(rw, r)
	})
}

type GzipResponseWriteWrapper struct {
	rw  http.ResponseWriter
	gzw *gzip.Writer
}

func NewGzipResponseWriterWrapper(rw http.ResponseWriter) *GzipResponseWriteWrapper {
	gzw := gzip.NewWriter(rw)
	return &GzipResponseWriteWrapper{rw, gzw}
}

func (gw *GzipResponseWriteWrapper) Header() http.Header {
	return gw.rw.Header()
}

func (gw *GzipResponseWriteWrapper) Write(d []byte) (int, error) {
	return gw.gzw.Write(d)
}

func (gw *GzipResponseWriteWrapper) WriteHeader(statusCode int) {
	gw.rw.WriteHeader(statusCode)
}

func (gw *GzipResponseWriteWrapper) Flush() {
	gw.gzw.Flush()
	gw.gzw.Close()
}
