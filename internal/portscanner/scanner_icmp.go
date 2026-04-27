package portscanner

import (
	"fmt"
	"net"
)

// ScanICMP returns all ICMP (v4 and v6) entries currently visible in /proc/net.
// It aggregates both /proc/net/icmp and /proc/net/icmp6.
func (s *Scanner) ScanICMP() ([]ICMPEntry, error) {
	v4, err := ParseProcNetICMP(fmt.Sprintf("%s/net/icmp", s.procRoot))
	if err != nil {
		// Non-fatal: kernel may not expose this file on all systems.
		v4 = nil
	}

	v6, err := ParseProcNetICMP6(fmt.Sprintf("%s/net/icmp6", s.procRoot))
	if err != nil {
		v6 = nil
	}

	all := make([]ICMPEntry, 0, len(v4)+len(v6))
	all = append(all, v4...)
	all = append(all, v6...)
	return all, nil
}

// ICMPEntryToPortEntry converts an ICMPEntry to a PortEntry so it can be
// processed uniformly by the rest of the pipeline.
func ICMPEntryToPortEntry(e ICMPEntry) PortEntry {
	return PortEntry{
		Protocol:  e.Protocol,
		LocalIP:   e.LocalIP,
		LocalPort: e.LocalPort,
		State:     e.State,
		Inode:     e.Inode,
	}
}

// PortEntry is a generic representation of a bound local port.
type PortEntry struct {
	Protocol  string
	LocalIP   net.IP
	LocalPort uint16
	State     string
	Inode     string
}

// String returns a concise description of a PortEntry.
func (p PortEntry) String() string {
	return fmt.Sprintf("%s %s:%d state=%s", p.Protocol, p.LocalIP, p.LocalPort, p.State)
}
