// Package scorecard computes a health score for each monitored service
// based on the proportion of drifted fields detected during a run.
//
// Scores range from 0 (no drift) to 100 (every tracked field has drifted).
// Each score maps to a letter grade:
//
//	 A  – score == 0   (clean)
//	 B  – score 1–25   (minor drift)
//	 C  – score 26–50  (moderate drift)
//	 D  – score 51–75  (significant drift)
//	 F  – score 76–100 (critical drift)
//
// Typical usage:
//
//	entries := scorecard.Build(results)
//	fmt.Println(scorecard.Summary(entries))
//	_ = scorecard.Format(os.Stdout, entries, "text")
package scorecard
