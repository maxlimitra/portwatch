// Package portscanner provides utilities for reading active TCP/UDP port
// bindings from the Linux /proc/net filesystem.
//
// The primary entry point is NewScanner, which returns a Scanner capable of
// listing all currently bound ports on the host. Ports are represented as
// [PortEntry] values containing the local address, port number, and protocol.
//
// Example usage:
//
//	s, err := portscanner.NewScanner()
//	if err != nil {
//		log.Fatal(err)
//	}
//	ports, err := s.Scan()
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range ports {
//		fmt.Println(p)
//	}
package portscanner
