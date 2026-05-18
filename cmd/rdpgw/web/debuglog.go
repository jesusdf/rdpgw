package web

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
)

var debugLog bool

// SetDebugLog enables verbose request/session logging (see Server.DebugLog / RDPGW_SERVER_DEBUGLOG).
func SetDebugLog(enabled bool) {
	debugLog = enabled
}

// DebugEnabled reports whether extended debug logging is active.
func DebugEnabled() bool {
	return debugLog
}

// Debugf writes a log line prefixed with [DEBUG] when debug logging is enabled.
func Debugf(format string, args ...interface{}) {
	if debugLog {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func hasSessionCookie(r *http.Request) bool {
	for _, c := range r.Cookies() {
		if c.Name == rdpGwSession || strings.HasPrefix(c.Name, rdpGwSession+".") {
			return true
		}
	}
	return false
}

func hasAccessToken(id identity.Identity) bool {
	if id == nil {
		return false
	}
	v := id.GetAttribute(identity.AttrAccessToken)
	if v == nil {
		return false
	}
	s, ok := v.(string)
	return ok && s != ""
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// DebugRequestLog logs method, path, status, duration and session cookie presence per request.
func DebugRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !debugLog {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		Debugf("HTTP %s %s host=%q remote=%q xfwd=%q cookieSession=%t -> %d in %s",
			r.Method, r.URL.RequestURI(), r.Host, r.RemoteAddr, r.Header.Get("X-Forwarded-For"),
			hasSessionCookie(r), rec.status, time.Since(start))
	})
}
