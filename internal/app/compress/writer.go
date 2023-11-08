package compress

import (
	"io"
	"net/http"
)

var supportedContentTypes = [2]string{
	"application/json",
	"text/html",
}

type compressWriter struct {
	w      http.ResponseWriter
	zw     io.WriteCloser
	engine compressEngine
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	if !cw.isContentTypeSupported() {
		return cw.w.Write(p)
	}

	cw.zw = cw.engine.NewWriter(cw.w)

	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if cw.isContentTypeSupported() {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}

	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	if cw.zw == nil {
		return nil
	}

	return cw.zw.Close()
}

func (cw *compressWriter) isContentTypeSupported() bool {
	contentType := cw.w.Header().Get("Content-Type")

	for _, v := range supportedContentTypes {
		if contentType == v {
			return true
		}
	}

	return false
}
