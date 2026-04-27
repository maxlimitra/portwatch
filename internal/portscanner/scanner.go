// Package portscanner reads active network ports from the Linux /proc filesystem.
package portscanner

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single listening port observed in /proc/net.
type PortEntry struct {
	Protocol string
	IP       net.IP
	Port     uint16
}

// Scanner scans /proc/net for active ports.
type Scanner struct {
	procRoot string
}

// NewScanner returns a Scanner rooted at procRoot (usually "/proc").
func NewScanner(procRoot string) *Scanner {
	return &Scanner{procRoot: procRoot}
}

// Scan reads tcp, udp, tcp6, and udp6 tables and returns all listening entries.
func (s *Scanner) Scan() ([]PortEntry, error) {
	files := map[string]string{
		"tcp":  s.procRoot + "/net/tcp",
		"udp":  s.procRoot + "/net/udp",
		"tcp6": s.procRoot + "/net/tcp6",
		"udp6": s.procRoot + "/net/udp6",
	}

	var all []PortEntry
	for proto, path := range files {
		var entries []PortEntry
		var err error
		if proto == "tcp6" || proto == "udp6" {
			entries, err = parseProcNet6(path, proto)
		} else {
			entries, err = parseProcNet(path, proto)
		}
		if err != nil {
			// Non-fatal: file may not exist on all kernels.
			continue
		}
		all = append(all, entries...)
	}
	return all, nil
}

// ipDecoder converts a hex address string to a net.IP.
type ipDecoder func(string) (net.IP, error)

// parseProcNetFile is the shared implementation used by both IPv4 and IPv6 parsers.
func parseProcNetFile(path, protocol string, decode ipDecoder) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parseProcNetFile: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "sl") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		// Only listening sockets (state 0A)
		if fields[3] != "0A" {
			continue
		}
		local := fields[1]
		colon := strings.LastIndex(local, ":")
		if colon < 0 {
			continue
		}
		ip, err := decode(local[:colon])
		if err != nil {
			continue
		}
		portHex := local[colon+1:]
		portVal, err := strconv.ParseUint(portHex, 16, 16)
		if err != nil {
			continue
		}
		entries = append(entries, PortEntry{
			Protocol: protocol,
			IP:       ip,
			Port:     uint16(portVal),
		})
	}
	return entries, scanner.Err()
}

// parseProcNet reads an IPv4 /proc/net file.
func parseProcNet(path, protocol string) ([]PortEntry, error) {
	return parseProcNetFile(path, protocol, hexToIP)
}

// hexToIP converts an 8-character little-endian hex string to a net.IP.
func hexToIP(hexStr string) (net.IP, error) {
	if len(hexStr) != 8 {
		return nil, fmt.Errorf("hexToIP: expected 8 hex chars, got %d", len(hexStr))
	}
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("hexToIP: %w", err)
	}
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, binary.LittleEndian.Uint32(b))
	return ip, nil
}
