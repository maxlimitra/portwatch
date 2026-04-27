package portscanner

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

// readProcLines opens the file at path, skips the header line, and returns
// the remaining trimmed non-empty lines.
func readProcLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	header := true
	for scanner.Scan() {
		if header {
			header = false
			continue
		}
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// parseHexAddr decodes an IPv4 "AABBCCDD:PPPP" hex address into an IP and port.
func parseHexAddr(addrHex string) (net.IP, uint16, error) {
	parts := strings.SplitN(addrHex, ":", 2)
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid addr: %s", addrHex)
	}
	ip := hexToIP(parts[0])
	if ip == nil {
		return nil, 0, fmt.Errorf("invalid ip hex: %s", parts[0])
	}
	portBytes, err := hex.DecodeString(parts[1])
	if err != nil || len(portBytes) < 2 {
		return nil, 0, fmt.Errorf("invalid port hex: %s", parts[1])
	}
	port := binary.BigEndian.Uint16(portBytes)
	return ip, port, nil
}

// parseHexAddrIPv6 decodes an IPv6 "...32hexchars...:PPPP" address into an IP and port.
func parseHexAddrIPv6(addrHex string) (net.IP, uint16, error) {
	parts := strings.SplitN(addrHex, ":", 2)
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid addr6: %s", addrHex)
	}
	ip := hexToIPv6(parts[0])
	if ip == nil {
		return nil, 0, fmt.Errorf("invalid ipv6 hex: %s", parts[0])
	}
	portBytes, err := hex.DecodeString(parts[1])
	if err != nil || len(portBytes) < 2 {
		return nil, 0, fmt.Errorf("invalid port hex: %s", parts[1])
	}
	port := binary.BigEndian.Uint16(portBytes)
	return ip, port, nil
}
