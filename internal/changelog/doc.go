// Package changelog provides an append-only, per-service record of drift state
// transitions observed by driftwatch.
//
// Each time a service moves from clean→drifted or drifted→clean the runner (or
// any other caller) should call [Log.Record] with the new state.  Entries are
// stored as newline-delimited JSON in a configurable directory, one file per
// service, so they can be tailed, grepped, or shipped to an external log
// aggregator without any additional tooling.
//
// Usage:
//
//	log, err := changelog.New("/var/lib/driftwatch/changelog")
//	if err != nil { ... }
//
//	log.Record(changelog.Entry{
//		Service:   "api-gateway",
//		State:     "drifted",
//		DiffCount: 3,
//	})
package changelog

import "bytes"

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
