package portscanner

import (
	"encoding/hex"
	"fmt"
	"net"
)

// hexToIPv6 converts a 32-character hex string (as found in /proc/net/tcp6)
// into a net.IP. The kernel stores the address as four little-endian 32-bit
// words, so each word must be byte-reversed before assembly.
func hexToIPv6(h string) (net.IP, error) {
	if len(h) != 32 {
		return nil, fmt.Errorf("hexToIPv6: expected 32 hex chars, got %d", len(h))
	}
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("hexToIPv6: %w", err)
	}
	// Reverse each 4-byte word.
	ip := make(net.IP, 16)
	for w := 0; w < 4; w++ {
		off := w * 4
		ip[off+0] = b[off+3]
		ip[off+1] = b[off+2]
		ip[off+2] = b[off+1]
		ip[off+3] = b[off+0]
	}
	return ip, nil
}

// parseProcNet6 reads /proc/net/tcp6 (or udp6) and appends the resulting
// PortEntries to dst, reusing the same line-parsing logic as parseProcNetFile
// but delegating IP conversion to hexToIPv6.
func parseProcNet6(path string, proto string, dst []PortEntry) ([]PortEntry, error) {
	return parseProcNetFile(path, proto, dst, hexToIPv6)
}
