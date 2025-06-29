package middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// 有些 response（例如 SSE）會需要即時 flush，這邊透過型別轉換（type assertion）呼叫原始的 Flush()。
func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

// 支援 WebSocket 或 raw TCP 時需要 Hijack()，這裡也是轉交給內部的 ResponseWriter。
func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
