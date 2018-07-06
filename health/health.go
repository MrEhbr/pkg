// Package health adds methods for checking the health of service dependencies.
package health

import (
	"net/http"

	"encoding/json"
)

// Checker returns an error indicating if a service is in an unhealthy state.
// Checkers should be implemented by dependencies which can fail, like a DB or mail service.
type Checker interface {
	HealthCheck() error
}

// Handler returns an http.Handler that checks the status of all the dependencies.
// Handler responds with either:
// 200 OK if the server can successfully communicate with it's backends or
// 500 if any of the backends are reporting an issue with list of failed checkers.
func Handler(checkers map[string]Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		components, healthy := CheckHealth(checkers)
		if !healthy {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(components)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// CheckHealth checks multiple checkers returning and return checkers that fail with errors.
func CheckHealth(checkers map[string]Checker) (map[string]error, bool) {
	components := make(map[string]error, len(checkers))
	healthy := true
	for name, hc := range checkers {
		if err := hc.HealthCheck(); err != nil {
			healthy = false
			components[name] = err
		}
	}
	return components, healthy
}

// Nop creates a noop checker. Useful in tests.
func Nop() Checker {
	return nop{}
}

type nop struct{}

func (c nop) HealthCheck() error {
	return nil
}
