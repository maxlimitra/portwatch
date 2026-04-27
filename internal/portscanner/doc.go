// Package portscanner provides utilities for reading active network socket
// information from the Linux /proc filesystem.
//
// # TCP
//
// ParseProcNet reads /proc/net/tcp (IPv4 TCP sockets).
// ParseProcNet6 reads /proc/net/tcp6 (IPv6 TCP sockets).
//
// # UDP
//
// ParseProcNetUDP reads /proc/net/udp (IPv4 UDP sockets).
// ParseProcNetUDP6 reads /proc/net/udp6 (IPv6 UDP sockets).
//
// # Scanner
//
// NewScanner returns a Scanner that aggregates results from all four proc
// files and exposes a unified list of PortEntry values, each tagged with
// the appropriate protocol string ("tcp", "tcp6", "udp", "udp6").
//
// Entries are keyed by (protocol, localIP, localPort) and deduplicated
// before being returned to callers.
package portscanner
