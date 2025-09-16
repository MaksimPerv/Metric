package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func ComputeHmacSha256(message []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

type CapturingResponseWriter struct {
	http.ResponseWriter
	Body []byte
}

func (w *CapturingResponseWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return w.ResponseWriter.Write(b)
}
