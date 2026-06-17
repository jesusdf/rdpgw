package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/protocol"
)

// connectionInfo is the JSON representation of a single active connection.
type connectionInfo struct {
	Id            string    `json:"id"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"displayName"`
	Domain        string    `json:"domain"`
	Target        string    `json:"target"`
	RemoteAddr    string    `json:"remoteAddr"`
	ConnectedAt   time.Time `json:"connectedAt"`
	LastSeen      time.Time `json:"lastSeen"`
	BytesSent     int64     `json:"bytesSent"`
	BytesReceived int64     `json:"bytesReceived"`
}

type connectionsResponse struct {
	Connections []connectionInfo `json:"connections"`
}

// Connections returns the list of currently active gateway connections as JSON,
// reporting who is connected and to which target host.
func Connections(w http.ResponseWriter, r *http.Request) {
	monitors := protocol.Snapshot()
	list := make([]connectionInfo, 0, len(monitors))

	for _, m := range monitors {
		if m == nil || m.Tunnel == nil {
			continue
		}
		t := m.Tunnel

		c := connectionInfo{
			Id:            t.Id,
			Target:        t.TargetServer,
			RemoteAddr:    t.RemoteAddr,
			ConnectedAt:   t.ConnectedOn,
			LastSeen:      t.LastSeen,
			BytesSent:     t.BytesSent,
			BytesReceived: t.BytesReceived,
		}
		if t.User != nil {
			c.Username = t.User.UserName()
			c.DisplayName = t.User.DisplayName()
			c.Domain = t.User.Domain()
		}
		list = append(list, c)
	}

	resp := connectionsResponse{
		Connections: list,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
