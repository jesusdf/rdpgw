package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
	"github.com/jesusdf/rdpgw/cmd/rdpgw/protocol"
)

// TestConnectionsOrderedByConnectedAt verifies the handler returns connections
// sorted by connectedAt, oldest first, regardless of registration order.
func TestConnectionsOrderedByConnectedAt(t *testing.T) {
	base := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)

	mk := func(id string, offset time.Duration) *protocol.Tunnel {
		u := identity.NewUser()
		u.SetUserName(id)
		return &protocol.Tunnel{Id: id, User: u, ConnectedOn: base.Add(offset)}
	}

	// register out of order: newest, oldest, middle
	protocol.RegisterTunnel(mk("order-new", 2*time.Hour), nil)
	protocol.RegisterTunnel(mk("order-old", 0), nil)
	protocol.RegisterTunnel(mk("order-mid", 1*time.Hour), nil)
	defer func() {
		for _, id := range []string{"order-new", "order-old", "order-mid"} {
			protocol.RemoveTunnel(&protocol.Tunnel{Id: id})
		}
	}()

	req := httptest.NewRequest(http.MethodGet, "/connections", nil)
	rec := httptest.NewRecorder()
	Connections(rec, req)

	var resp struct {
		Connections []struct {
			Username    string    `json:"username"`
			ConnectedAt time.Time `json:"connectedAt"`
		} `json:"connections"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// extract only the ones registered by this test, preserving response order
	var got []string
	for _, c := range resp.Connections {
		switch c.Username {
		case "order-old", "order-mid", "order-new":
			got = append(got, c.Username)
		}
	}

	want := []string{"order-old", "order-mid", "order-new"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
}
