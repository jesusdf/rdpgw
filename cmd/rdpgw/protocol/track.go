package protocol

import (
	"fmt"
	"sync"
)

var (
	Connections   map[string]*Monitor
	connectionsMu sync.RWMutex
)

type Monitor struct {
	Processor *Processor
	Tunnel    *Tunnel
}

// ActiveConnection describes one gateway tunnel for observability APIs.
type ActiveConnection struct {
	Username string `json:"username"`
	Target   string `json:"target"`
}

const (
	ctlDisconnect = -1
)

func RegisterTunnel(t *Tunnel, p *Processor) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()
	if Connections == nil {
		Connections = make(map[string]*Monitor)
	}

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
	connectionsMu.Lock()
	defer connectionsMu.Unlock()
	if Connections == nil {
		return fmt.Errorf("%s connection does not exist", id)
	}

	m, ok := Connections[id]
	if !ok {
		return fmt.Errorf("%s connection does not exist", id)
	}
	m.Processor.ctl <- ctlDisconnect
	return nil
}

// SnapshotActiveConnections returns a copy of current tunnels and their user/target.
func SnapshotActiveConnections() []ActiveConnection {
	connectionsMu.RLock()
	defer connectionsMu.RUnlock()
	if Connections == nil {
		return nil
	}
	out := make([]ActiveConnection, 0, len(Connections))
	for _, m := range Connections {
		if m == nil || m.Tunnel == nil {
			continue
		}
		u, tgt := m.Tunnel.ActiveSession()
		out = append(out, ActiveConnection{Username: u, Target: tgt})
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
