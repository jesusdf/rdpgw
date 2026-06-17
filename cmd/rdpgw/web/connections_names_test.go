package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
	"github.com/jesusdf/rdpgw/cmd/rdpgw/protocol"
)

// TestConnectionsTwoNames verifies that a user whose name has several words is
// reported correctly (spaces preserved, valid JSON) by the /conexiones handler.
func TestConnectionsTwoNames(t *testing.T) {
	u := identity.NewUser()
	u.SetUserName("jdieguez")
	u.SetDisplayName("Jesús María Diéguez Pérez")
	u.SetDomain("SIVSA")

	protocol.RegisterTunnel(&protocol.Tunnel{
		Id:           "names-test-1",
		User:         u,
		TargetServer: "10.0.0.15:3389",
		RemoteAddr:   "203.0.113.7",
	}, nil)
	defer func() {
		protocol.RemoveTunnel(&protocol.Tunnel{Id: "names-test-1"})
	}()

	req := httptest.NewRequest(http.MethodGet, "/conexiones", nil)
	rec := httptest.NewRecorder()
	Connections(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp struct {
		Count       int `json:"count"`
		Connections []struct {
			User        string `json:"user"`
			DisplayName string `json:"displayName"`
		} `json:"connections"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON produced: %v\nbody: %s", err, rec.Body.String())
	}

	var found bool
	for _, c := range resp.Connections {
		if c.User == "jdieguez" {
			found = true
			if c.DisplayName != "Jesús María Diéguez Pérez" {
				t.Errorf("displayName = %q, want %q", c.DisplayName, "Jesús María Diéguez Pérez")
			}
		}
	}
	if !found {
		t.Fatalf("connection for jdieguez not found in: %s", rec.Body.String())
	}

	t.Logf("JSON válido y nombre completo preservado:\n%s", rec.Body.String())
}
