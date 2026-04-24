// Package portfilter provides filtering logic for port entries, allowing
// callers to include or exclude ports based on configurable rules.
package portfilter

import "github.com/user/portwatch/internal/portscanner"

// Rule describes a single filter rule.
type Rule struct {
	// Port is the port number to match. 0 means match any port.
	Port uint16
	// Protocol is "tcp", "udp", or "" to match any protocol.
	Protocol string
	// Allow indicates whether matching entries are allowed (true) or denied (false).
	Allow bool
}

// Filter applies a set of Rules to a slice of PortEntry values.
// Entries that match an Allow rule are retained; entries that match a Deny
// rule are removed. Entries that match no rule are retained by default.
type Filter struct {
	rules []Rule
}

// NewFilter creates a Filter from the provided rules.
func NewFilter(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns only the entries that pass the filter rules.
// Rules are evaluated in order; the first matching rule wins.
func (f *Filter) Apply(entries []portscanner.PortEntry) []portscanner.PortEntry {
	if len(f.rules) == 0 {
		return entries
	}
	out := make([]portscanner.PortEntry, 0, len(entries))
	for _, e := range entries {
		if f.keep(e) {
			out = append(out, e)
		}
	}
	return out
}

// keep returns true if the entry should be retained.
func (f *Filter) keep(e portscanner.PortEntry) bool {
	for _, r := range f.rules {
		if r.matches(e) {
			return r.Allow
		}
	}
	return true // default: retain
}

// matches returns true when the rule applies to the given entry.
func (r Rule) matches(e portscanner.PortEntry) bool {
	portMatch := r.Port == 0 || r.Port == e.Port
	protoMatch := r.Protocol == "" || r.Protocol == e.Protocol
	return portMatch && protoMatch
}
