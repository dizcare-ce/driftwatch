// Package ownermap provides a thread-safe registry that maps service names to
// their declared owners.
//
// Each service may have an associated team name and a contact address (e-mail
// or Slack handle). The registry is populated at start-up from the service
// definition files and can be queried by the reporter and notifier to include
// ownership information alongside drift results.
//
// Usage:
//
//	om := ownermap.New()
//	om.Set("payment-api", "payments-squad", "payments@example.com")
//
//	if owner, ok := om.Get("payment-api"); ok {
//		fmt.Println(owner.Team, owner.Contact)
//	}
//
// Service names are normalised to lowercase on both Set and Get so that
// lookups are always case-insensitive.
package ownermap
