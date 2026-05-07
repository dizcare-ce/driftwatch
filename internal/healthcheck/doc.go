// Package healthcheck provides a lightweight HTTP health endpoint for the
// driftwatch daemon.
//
// The endpoint GET /healthz returns a JSON body describing the outcome of
// the most recent drift-detection run:
//
//	{
//	  "ok": true,
//	  "last_run_at": "2024-05-01T12:00:00Z",
//	  "drifted_services": 0,
//	  "total_services": 5,
//	  "message": "all services in sync"
//	}
//
// HTTP 200 is returned when no drift is present; HTTP 409 (Conflict) is
// returned when one or more services have drifted, making the endpoint
// suitable for use with container-orchestration liveness/readiness probes
// or external monitoring systems.
package healthcheck
