package portscanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single listening port binding.
type PortEntry struct {
	Protocol string
	LocalAddr string
	Port     int
	PID      int
	Process  string
}

// Scanner reads active port bindings from the system.
type Scanner struct{}

// NewScanner creates a new Scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan returns the current list of listening port entries.
func (s *Scanner) Scan() ([]PortEntry, error) {
	entries, err := parseProcNet("/proc/net/tcp")
	if err != nil {
		return nil, fmt.Errorf("scanning tcp ports: %w", err)
	}
	udpEntries, err := parseProcNet("/proc/net/udp")
	if err != nil {
		return nil, fmt.Errorf("scanning udp ports: %w", err)
	}
	return append(entries, udpEntries...), nil
}

func parseProcNet(path string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	protocol := "tcp"
	if strings.Contains(path, "udp") {
		protocol = "udp"
	}

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A = listening for TCP
		if protocol == "tcp" && fields[3] != "0A" {
			continue
		}
		addr := fields[1]
		parts := strings.Split(addr, ":")
		if len(parts) != 2 {
			continue
		}
		portHex := parts[1]
		port, err := strconv.ParseInt(portHex, 16, 32)
		if err != nil {
			continue
		}
		entries = append(entries, PortEntry{
			Protocol:  protocol,
			LocalAddr: hexToIP(parts[0]),
			Port:      int(port),
		})
	}
	return entries, scanner.Err()
}

func hexToIP(hex string) string {
	if len(hex) != 8 {
		return hex
	}
	d, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return hex
	}
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(d), byte(d>>8), byte(d>>16), byte(d>>24))
}
