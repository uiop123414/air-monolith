package rww

import "net/http"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	written bool
	status  int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	if rw.written {
		// Avoid superfluous calls
		return
	}
	rw.written = true
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(data)
}

func (rw *ResponseWriterWrapper) HasWritten() bool {
	return rw.written
}
