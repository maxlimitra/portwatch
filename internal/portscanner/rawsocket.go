package portscanner

import (
	"fmt"
	"net"
)

// RawEntry represents a raw socket entry parsed from /proc/net/raw or /proc/net/raw6.
type RawEntry struct {
	LocalIP   net.IP
	LocalPort uint16
	Protocol  uint8
	Inode     uint64
}

// String returns a human-readable representation of a RawEntry.
func (r RawEntry) String() string {
	return fmt.Sprintf("%s:%d (proto=%d inode=%d)", r.LocalIP, r.LocalPort, r.Protocol, r.Inode)
}

// ParseProcNetRaw parses /proc/net/raw for raw IPv4 socket entries.
func ParseProcNetRaw(path string) ([]RawEntry, error) {
	lines, err := readProcLines(path)
	if err != nil {
		return nil, fmt.Errorf("raw: read %s: %w", path, err)
	}
	var entries []RawEntry
	for _, line := range lines {
		e, ok := parseRawLine(line, false)
		if !ok {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ParseProcNetRaw6 parses /proc/net/raw6 for raw IPv6 socket entries.
func ParseProcNetRaw6(path string) ([]RawEntry, error) {
	lines, err := readProcLines(path)
	if err != nil {
		return nil, fmt.Errorf("raw6: read %s: %w", path, err)
	}
	var entries []RawEntry
	for _, line := range lines {
		e, ok := parseRawLine(line, true)
		if !ok {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func parseRawLine(line string, ipv6 bool) (RawEntry, bool) {
	var slot, localAddrHex, remAddrHex string
	var state, txq, rxq, timer, retrans int
	var inode uint64
	var uid, timeout int
	var proto uint8

	n, err := fmt.Sscanf(line,
		"%s %s %s %x %x:%x %d %d %d %d %d",
		&slot, &localAddrHex, &remAddrHex,
		&proto, &txq, &rxq,
		&timer, &retrans, &uid, &timeout, &inode,
	)
	_ = state
	if n < 4 || err != nil {
		// fallback: simpler parse
		var localHex string
		var inodeVal uint64
		nn, e2 := fmt.Sscanf(line, "%s %s %s %x", &slot, &localHex, &remAddrHex, &proto)
		if nn < 4 || e2 != nil {
			return RawEntry{}, false
		}
		localAddrHex = localHex
		inode = inodeVal
	}

	var ip net.IP
	var port uint16
	if ipv6 {
		ip, port, err = parseHexAddrIPv6(localAddrHex)
	} else {
		ip, port, err = parseHexAddr(localAddrHex)
	}
	if err != nil {
		return RawEntry{}, false
	}
	return RawEntry{
		LocalIP:   ip,
		LocalPort: port,
		Protocol:  proto,
		Inode:     inode,
	}, true
}
