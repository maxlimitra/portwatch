// Package process resolves the owning process for a bound socket by
// inspecting the Linux /proc filesystem.
//
// # Overview
//
// When portwatch detects a new port binding it can optionally enrich the
// alert with the name and PID of the process responsible.  This package
// provides a [Resolver] that walks /proc/<pid>/fd, follows each file-
// descriptor symlink, and matches the "socket:[inode]" target against the
// inode reported by /proc/net/tcp (or udp).
//
// # Usage
//
//	r := process.NewResolver("/proc")
//	info, err := r.Resolve(inode)
//	if err == nil {
//		fmt.Printf("port owned by pid=%d name=%s\n", info.PID, info.Name)
//	}
//
// # Limitations
//
// Inode resolution requires that portwatch runs with sufficient privileges
// to read /proc/<pid>/fd for all processes (typically root or CAP_SYS_PTRACE).
package process
