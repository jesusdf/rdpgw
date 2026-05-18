package web

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/jcmturner/goidentity/v6"
	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
)

func EnrichContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := GetSessionIdentity(r)
		if err != nil {
			Debugf("session load error path=%s: %v", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createdSession := false
		if id == nil {
			createdSession = true
			id = identity.NewUser()
			if err := SaveSessionIdentity(r, w, id); err != nil {
				Debugf("new session save error path=%s sessionId=%s: %v", r.URL.Path, id.SessionId(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			Debugf("new anonymous session path=%s sessionId=%s cookieSent=%t", r.URL.Path, id.SessionId(), hasSessionCookie(r))
		}

		// Healthcheck doesn't need any log
		if r.URL.Path != "/teapot" {
			log.Printf("Identity SessionId: %s, UserName: %s: Authenticated: %t",
				id.SessionId(), id.UserName(), id.Authenticated())
			Debugf("identity path=%s method=%s host=%q user=%q auth=%t accessToken=%t cookieSession=%t newSession=%t remote=%q xfwd=%q",
				r.URL.Path, r.Method, r.Host, id.UserName(), id.Authenticated(), hasAccessToken(id),
				hasSessionCookie(r), createdSession, r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
		}

		h := r.Header.Get("X-Forwarded-For")
		if h != "" {
			var proxies []string
			ips := strings.Split(h, ",")
			for i := range ips {
				ips[i] = strings.TrimSpace(ips[i])
			}
			clientIp := ips[0]
			if len(ips) > 1 {
				proxies = ips[1:]
			}
			id.SetAttribute(identity.AttrClientIp, clientIp)
			id.SetAttribute(identity.AttrProxies, proxies)
		}

		id.SetAttribute(identity.AttrRemoteAddr, r.RemoteAddr)
		if h == "" {
			clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)
			id.SetAttribute(identity.AttrClientIp, clientIp)
		}
		next.ServeHTTP(w, identity.AddToRequestCtx(id, r))
	})
}

func TransposeSPNEGOContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gid := goidentity.FromHTTPRequestContext(r)
		if gid != nil {
			id := identity.FromRequestCtx(r)
			id.SetUserName(gid.UserName())
			id.SetAuthenticated(gid.Authenticated())
			id.SetDomain(gid.Domain())
			id.SetAuthTime(gid.AuthTime())
			r = identity.AddToRequestCtx(id, r)
		}
		next.ServeHTTP(w, r)
	})
}
