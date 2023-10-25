package gzip

import (
	"compress/gzip"
	"net/http"
)

var supportedContentTypes = [2]string{
	"application/json",
	"text/html",
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w: w,
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.isContentTypeSupported() {
		return c.w.Write(p)
	}

	c.zw = gzip.NewWriter(c.w)

	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && c.isContentTypeSupported() {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if c.zw == nil {
		return nil
	}

	return c.zw.Close()
}

func (c *compressWriter) isContentTypeSupported() bool {
	contentType := c.w.Header().Get("Content-Type")

	for _, v := range supportedContentTypes {
		if contentType == v {
			return true
		}
	}

	return false
}
