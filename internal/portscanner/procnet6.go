// Package portscanner reads active network ports from the Linux /proc filesystem.
package portscanner

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
)

// hexToIPv6 converts a 32-character little-endian hex string (as found in
// /proc/net/tcp6 and /proc/net/udp6) into a net.IP.
func hexToIPv6(hexStr string) (net.IP, error) {
	if len(hexStr) != 32 {
		return nil, fmt.Errorf("hexToIPv6: expected 32 hex chars, got %d", len(hexStr))
	}

	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("hexToIPv6: decode error: %w", err)
	}

	// /proc/net/tcp6 stores the address as four little-endian 32-bit words.
	ip := make(net.IP, 16)
	for i := 0; i < 4; i++ {
		word := binary.LittleEndian.Uint32(b[i*4 : i*4+4])
		binary.BigEndian.PutUint32(ip[i*4:i*4+4], word)
	}
	return ip, nil
}

// parseProcNet6 reads /proc/net/tcp6 or /proc/net/udp6 and returns PortEntries.
// It reuses parseProcNetWithTransform with an IPv6-aware address decoder.
func parseProcNet6(path, protocol string) ([]PortEntry, error) {
	return parseProcNetFile(path, protocol, hexToIPv6)
}
