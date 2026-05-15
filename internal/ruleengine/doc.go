// Package ruleengine provides a lightweight rule evaluation engine for
// driftwatch. Operators define named Rule values that express constraints
// such as "this service must always be clean" or "no more than N diffs
// are tolerated". Each Rule may target a specific subset of services via
// a regular expression matched against the service name.
//
// Usage:
//
//	engine, err := ruleengine.New([]ruleengine.Rule{
//		{Name: "api-must-be-clean", ServicePattern: "^api", RequireClean: true},
//		{Name: "global-max-diffs",  MaxDiffs: 5},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, result := range results {
//		for _, v := range engine.Evaluate(result) {
//			log.Printf("rule %q violated for %s: %s", v.Rule, v.Service, v.Reason)
//		}
//	}
package ruleengine
