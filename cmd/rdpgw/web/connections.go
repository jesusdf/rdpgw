package web

import (
	"encoding/json"
	"net/http"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/protocol"
)

// ConnectionsStatus responds with active RDP gateway tunnels (authenticated user
// display name and backend target). Same process scope as /metrics; restrict at
// the network layer if the endpoint must not be public.
func ConnectionsStatus(w http.ResponseWriter, r *http.Request) {
	list := protocol.SnapshotActiveConnections()
	if list == nil {
		list = []protocol.ActiveConnection{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(list)
}
