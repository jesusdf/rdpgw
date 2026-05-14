package protocol

import (
	"net"
	"sync"
	"time"

	"github.com/jesusdf/rdpgw/cmd/rdpgw/identity"
	"github.com/jesusdf/rdpgw/cmd/rdpgw/transport"
)

const (
	CtxTunnel = "github.com/jesusdf/rdpgw/tunnel"
)

type Tunnel struct {
	// Id identifies the connection in the server
	Id string
	// The connection-id (RDG-ConnID) as reported by the client
	RDGId string
	// The underlying incoming transport being either websocket or legacy http
	// in case of websocket transportOut will equal transportIn
	transportIn transport.Transport
	// The underlying outgoing transport being either websocket or legacy http
	// in case of websocket transportOut will equal transportOut
	transportOut transport.Transport
	// The remote desktop server (rdp, vnc etc) the clients intends to connect to
	TargetServer string
	// The obtained client ip address
	RemoteAddr string
	// User
	User identity.Identity

	// rwc is the underlying connection to the remote desktop server.
	// It is of the type *net.TCPConn
	rwc net.Conn

	// BytesSent is the total amount of bytes sent by the server to the client minus tunnel overhead
	BytesSent int64

	// BytesReceived is the total amount of bytes received by the server from the client minus tunnel overhad
	BytesReceived int64

	// ConnectedOn is when the client connected to the server
	ConnectedOn time.Time

	// LastSeen is when the server received the last packet from the client
	LastSeen time.Time

	metaMu sync.RWMutex
}

// SetTargetServer records the TCP destination after the RDP channel is established.
func (t *Tunnel) SetTargetServer(host string) {
	t.metaMu.Lock()
	t.TargetServer = host
	t.metaMu.Unlock()
}

// GetTargetServer returns the current backend address (empty until the channel is open).
func (t *Tunnel) GetTargetServer() string {
	t.metaMu.RLock()
	defer t.metaMu.RUnlock()
	return t.TargetServer
}

// ApplyPAAIdentity sets tunnel fields from a validated PAA (gateway) token.
func (t *Tunnel) ApplyPAAIdentity(remoteServer, clientIP, subject string) {
	t.metaMu.Lock()
	t.TargetServer = remoteServer
	t.RemoteAddr = clientIP
	if t.User != nil {
		t.User.SetUserName(subject)
	}
	t.metaMu.Unlock()
}

// ActiveSession returns the display name and target host for status listings.
func (t *Tunnel) ActiveSession() (displayName, target string) {
	t.metaMu.RLock()
	defer t.metaMu.RUnlock()
	target = t.TargetServer
	if t.User != nil {
		displayName = t.User.DisplayName()
	}
	return
}

// Write puts the packet on the transport and updates the statistics for bytes sent
func (t *Tunnel) Write(pkt []byte) {
	n, _ := t.transportOut.WritePacket(pkt)
	t.BytesSent += int64(n)
}

// Read picks up a packet from the transport and returns the packet type
// packet, with the header removed, and the packet size. It updates the
// statistics for bytes received
func (t *Tunnel) Read() (pt int, size int, pkt []byte, err error) {
	pt, size, pkt, err = readMessage(t.transportIn)
	t.BytesReceived += int64(size)
	t.LastSeen = time.Now()

	return pt, size, pkt, err
}
