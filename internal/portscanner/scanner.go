package portscanner

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single bound port observed in /proc/net/tcp.
type PortEntry struct {
	Port    int
	Address string
	PID     int
}

// ScannerFunc is a function type that implements the Scanner interface.
type ScannerFunc func() ([]PortEntry, error)

// Scan calls the underlying function.
func (f ScannerFunc) Scan() ([]PortEntry, error) { return f() }

// Scanner reads current port bindings from the OS.
type Scanner struct {
	procNetTCP string
}

// NewScanner returns a Scanner reading from the standard /proc/net/tcp path.
func NewScanner() *Scanner {
	return &Scanner{procNetTCP: "/proc/net/tcp"}
}

// Scan returns all currently bound TCP ports.
func (s *Scanner) Scan() ([]PortEntry, error) {
	return parseProcNet(s.procNetTCP)
}

func parseProcNet(path string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	first := true
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// Skip non-LISTEN state (0A == 10 == LISTEN).
		if fields[3] != "0A" {
			continue
		}
		local := fields[1]
		parts := strings.SplitN(local, ":", 2)
		if len(parts) != 2 {
			continue
		}
		ip, err := hexToIP(parts[0])
		if err != nil {
			continue
		}
		port, err := strconv.ParseInt(parts[1], 16, 32)
		if err != nil {
			continue
		}
		entries = append(entries, PortEntry{Port: int(port), Address: ip})
	}
	return entries, scanner.Err()
}

func hexToIP(h string) (string, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return "", err
	}
	if len(b) == 4 {
		// Little-endian for IPv4.
		return net.IPv4(b[3], b[2], b[1], b[0]).String(), nil
	}
	return "", fmt.Errorf("unsupported address length %d", len(b))
}
