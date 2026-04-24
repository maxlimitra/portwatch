// Package alerting provides alert creation, dispatch, and rate-limited
// delivery for portwatch.
//
// # Overview
//
// An [Alerter] accepts [Alert] values and fans them out to one or more
// registered handler functions.  Handlers are plain Go functions with the
// signature
//
//	func(Alert) error
//
// so they are easy to compose and test in isolation.
//
// # Built-in handlers
//
// Three ready-made handlers are provided:
//
//   - [StdoutHandler] – writes a human-readable line to os.Stdout.
//   - [JSONHandler]   – writes a JSON-encoded alert to os.Stdout.
//   - [FileHandler]   – appends JSON-encoded alerts to a file on disk.
//
// # Rate limiting
//
// [NewRateLimitedAlerter] wraps any Alerter with a per-port cooldown so that
// repeated alerts for the same port are suppressed until the cooldown period
// has elapsed.  This prevents log spam during transient port flapping.
//
// # Alert levels
//
// Alerts carry one of three severity levels:
//
//   - [LevelInfo]    – expected or informational binding detected.
//   - [LevelWarning] – unexpected binding that may warrant attention.
//   - [LevelCritical] – binding that violates a hard policy rule.
package alerting
