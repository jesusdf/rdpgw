package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
	"github.com/jesusdf/rdpgw/cmd/rdpgw/protocol"
)

// connectionsStatusResponse is the JSON body for GET /connections.
type connectionsStatusResponse struct {
	RequestUser *connectionIdentity           `json:"requestUser,omitempty"`
	Connections []protocol.ActiveConnection `json:"connections"`
}

type connectionIdentity struct {
	Username      string `json:"username,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Email         string `json:"email,omitempty"`
	ClientIP      string `json:"clientIp,omitempty"`
	Authenticated bool   `json:"authenticated"`
	SessionID     string `json:"sessionId,omitempty"`
}

func identityFromRequest(r *http.Request) *connectionIdentity {
	id := identity.FromRequestCtx(r)
	if id == nil {
		return nil
	}
	info := &connectionIdentity{
		Username:      id.UserName(),
		DisplayName:   id.DisplayName(),
		Domain:        id.Domain(),
		Email:         id.Email(),
		Authenticated: id.Authenticated(),
		SessionID:     id.SessionId(),
		ClientIP:      identityStringAttr(id, identity.AttrClientIp),
	}
	if info.Domain == "" && info.Username != "" {
		if creds := strings.SplitN(info.Username, "@", 2); len(creds) == 2 {
			info.Domain = creds[1]
		}
	}
	return info
}

func identityStringAttr(id identity.Identity, key string) string {
	if id == nil {
		return ""
	}
	v := id.GetAttribute(key)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// ConnectionsStatus responds with the HTTP request identity (session from
// EnrichContext) and active RDP gateway tunnels. Same process scope as /metrics;
// restrict at the network layer if the endpoint must not be public.
func ConnectionsStatus(w http.ResponseWriter, r *http.Request) {
	list := protocol.SnapshotActiveConnections()
	if list == nil {
		list = []protocol.ActiveConnection{}
	}
	body := connectionsStatusResponse{
		RequestUser: identityFromRequest(r),
		Connections: list,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(body)
}
