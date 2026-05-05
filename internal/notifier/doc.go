// Package notifier provides threshold-based alerting for drift results.
//
// A Notifier inspects a slice of drift.Result values and emits a one-line
// summary for each result that meets or exceeds the configured Level:
//
//	"all"   – notify for every service, whether drifted or clean.
//	"drift" – notify only for services where drift was detected (default).
//	"none"  – suppress all notifications.
//
// Notifications are written to any io.Writer (defaults to os.Stderr) so they
// can be redirected to a log aggregator, a file, or captured in tests.
//
// Example:
//
//	n := notifier.New(notifier.LevelDrift, os.Stderr)
//	count, err := n.Notify(results)
package notifier
