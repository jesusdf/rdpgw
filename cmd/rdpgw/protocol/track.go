package protocol

import (
	"fmt"
	"sync"
)

// Connections holds the currently active gateway connections keyed by tunnel id.
// It must only be accessed while holding connectionsMu. Use the helper functions
// in this file (RegisterTunnel, RemoveTunnel, Disconnect, Snapshot) instead of
// touching the map directly, as Go maps are not safe for concurrent use.
var (
	Connections   = make(map[string]*Monitor)
	connectionsMu sync.RWMutex
)

type Monitor struct {
	Processor *Processor
	Tunnel    *Tunnel
}

const (
	ctlDisconnect = -1
)

func RegisterTunnel(t *Tunnel, p *Processor) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	Connections[t.Id] = &Monitor{
		Processor: p,
		Tunnel:    t,
	}
}

func RemoveTunnel(t *Tunnel) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	delete(Connections, t.Id)
}

func Disconnect(id string) error {
	connectionsMu.RLock()
	m, ok := Connections[id]
	connectionsMu.RUnlock()

	if !ok {
		return fmt.Errorf("%s connection does not exist", id)
	}

	m.Processor.ctl <- ctlDisconnect
	return nil
}

// Snapshot returns a copy of the currently active monitors. The returned slice
// is a point-in-time snapshot taken under a read lock, so callers can iterate it
// safely without racing the connection map. The *Monitor pointers themselves are
// shared with the live connections.
func Snapshot() []*Monitor {
	connectionsMu.RLock()
	defer connectionsMu.RUnlock()

	out := make([]*Monitor, 0, len(Connections))
	for _, m := range Connections {
		out = append(out, m)
	}
	return out
}

// CalculateSpeedPerSecond calculate moving average.
/*
func CalculateSpeedPerSecond(connId string) (in int, out int) {
	now := time.Now().UnixMilli()

	c := Connections[connId]
	total := int64(0)
	for _, v := range c.Tunnel.BytesReceived {
		total += v
	}
	in = int(total / (now - c.TimeStamp) * 1000)

	total = int64(0)
	for _, v := range c.BytesSent {
		total += v
	}
	out = int(total / (now - c.TimeStamp))

	return in, out
}
*/
