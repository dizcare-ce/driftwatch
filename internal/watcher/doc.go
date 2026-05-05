// Package watcher provides a lightweight file-system poller that watches a
// directory of service definition files and notifies callers when any file is
// created or modified.
//
// Usage:
//
//	fw := watcher.New("/etc/driftwatch/services", 5*time.Second, func(path string) {
//		log.Printf("definition changed: %s", path)
//		// trigger a drift check for the affected service
//	})
//	if err := fw.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatalf("watcher: %v", err)
//	}
//
// The watcher uses simple mtime polling rather than inotify so that it works
// uniformly across Linux, macOS, and containerised environments where kernel
// event APIs may be restricted.
package watcher
