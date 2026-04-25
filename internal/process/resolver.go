// Package process provides utilities for resolving the owning process
// of a bound port using /proc/net inode information.
package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Info holds basic information about the process that owns a socket inode.
type Info struct {
	PID  int
	Name string
}

// Resolver maps socket inodes to process information by scanning /proc.
type Resolver struct {
	procRoot string
}

// NewResolver returns a Resolver that reads from the given procRoot
// (typically "/proc").
func NewResolver(procRoot string) *Resolver {
	return &Resolver{procRoot: procRoot}
}

// Resolve returns the process Info for the given socket inode, or an error
// if no matching process is found.
func (r *Resolver) Resolve(inode uint64) (*Info, error) {
	entries, err := os.ReadDir(r.procRoot)
	if err != nil {
		return nil, fmt.Errorf("read proc root: %w", err)
	}

	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // skip non-numeric entries
		}

		fdDir := filepath.Join(r.procRoot, e.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == fmt.Sprintf("socket:[%d]", inode) {
				name := r.readProcessName(pid)
				return &Info{PID: pid, Name: name}, nil
			}
		}
	}

	return nil, fmt.Errorf("no process found for inode %d", inode)
}

func (r *Resolver) readProcessName(pid int) string {
	commPath := filepath.Join(r.procRoot, strconv.Itoa(pid), "comm")
	data, err := os.ReadFile(commPath)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}
