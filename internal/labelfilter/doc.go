// Package labelfilter provides selector-based filtering of services using
// key=value label pairs attached to source definitions.
//
// # Usage
//
// Construct a Filter from a slice of "key=value" selector strings.  A service
// is accepted only when its label map contains every selector key with a
// matching value (case-insensitive).  When no selectors are configured the
// filter is a no-op and all services pass through.
//
//	f := labelfilter.New([]string{"env=prod", "tier=api"})
//
//	entries := []labelfilter.Entry{
//		{Name: "payments", Labels: map[string]string{"env": "prod", "tier": "api"}},
//		{Name: "mailer",   Labels: map[string]string{"env": "prod", "tier": "worker"}},
//	}
//
//	names := f.MatchAll(entries) // ["payments"]
package labelfilter
