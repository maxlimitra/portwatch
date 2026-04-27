package portscanner

import (
	"fmt"
	"net"
)

// ICMPEntry represents a parsed entry from /proc/net/icmp or /proc/net/icmp6.
type ICMPEntry struct {
	LocalIP   net.IP
	LocalPort uint16
	RemoteIP  net.IP
	RemotePort uint16
	State     string
	Inode     string
	Protocol  string
}

// String returns a human-readable representation of an ICMPEntry.
func (e ICMPEntry) String() string {
	return fmt.Sprintf("%s %s:%d -> %s:%d state=%s",
		e.Protocol,
		e.LocalIP.String(), e.LocalPort,
		e.RemoteIP.String(), e.RemotePort,
		e.State,
	)
}

// ParseProcNetICMP parses /proc/net/icmp and returns a slice of ICMPEntry.
func ParseProcNetICMP(path string) ([]ICMPEntry, error) {
	raw, err := parseProcNetFile(path)
	if err != nil {
		return nil, fmt.Errorf("icmp: %w", err)
	}
	entries := make([]ICMPEntry, 0, len(raw))
	for _, r := range raw {
		entries = append(entries, ICMPEntry{
			LocalIP:    r.LocalIP,
			LocalPort:  r.LocalPort,
			RemoteIP:   r.RemoteIP,
			RemotePort: r.RemotePort,
			State:      r.State,
			Inode:      r.Inode,
			Protocol:   "icmp",
		})
	}
	return entries, nil
}

// ParseProcNetICMP6 parses /proc/net/icmp6 and returns a slice of ICMPEntry.
func ParseProcNetICMP6(path string) ([]ICMPEntry, error) {
	raw, err := parseProcNet6(path)
	if err != nil {
		return nil, fmt.Errorf("icmp6: %w", err)
	}
	entries := make([]ICMPEntry, 0, len(raw))
	for _, r := range raw {
		entries = append(entries, ICMPEntry{
			LocalIP:    r.LocalIP,
			LocalPort:  r.LocalPort,
			RemoteIP:   r.RemoteIP,
			RemotePort: r.RemotePort,
			State:      r.State,
			Inode:      r.Inode,
			Protocol:   "icmp6",
		})
	}
	return entries, nil
}
