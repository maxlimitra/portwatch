package portscanner

import (
	"fmt"
	"net"
)

// UDPEntry represents a single UDP socket entry parsed from /proc/net/udp or /proc/net/udp6.
type UDPEntry struct {
	LocalIP   net.IP
	LocalPort uint16
	Protocol  string // "udp" or "udp6"
}

// String returns a human-readable representation of the UDP entry.
func (e UDPEntry) String() string {
	return fmt.Sprintf("%s:%d (%s)", e.LocalIP.String(), e.LocalPort, e.Protocol)
}

// ParseProcNetUDP reads /proc/net/udp and returns all bound UDP socket entries.
func ParseProcNetUDP(path string) ([]UDPEntry, error) {
	rawEntries, err := parseProcNetFile(path)
	if err != nil {
		return nil, fmt.Errorf("udp: %w", err)
	}

	entries := make([]UDPEntry, 0, len(rawEntries))
	for _, re := range rawEntries {
		entries = append(entries, UDPEntry{
			LocalIP:   re.LocalIP,
			LocalPort: re.LocalPort,
			Protocol:  "udp",
		})
	}
	return entries, nil
}

// ParseProcNetUDP6 reads /proc/net/udp6 and returns all bound UDP6 socket entries.
func ParseProcNetUDP6(path string) ([]UDPEntry, error) {
	rawEntries, err := parseProcNet6(path)
	if err != nil {
		return nil, fmt.Errorf("udp6: %w", err)
	}

	entries := make([]UDPEntry, 0, len(rawEntries))
	for _, re := range rawEntries {
		entries = append(entries, UDPEntry{
			LocalIP:   re.LocalIP,
			LocalPort: re.LocalPort,
			Protocol:  "udp6",
		})
	}
	return entries, nil
}
