package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckHealth(t *testing.T) {
	checkers := map[string]Checker{
		"fail": fail{},
		"pass": Nop(),
	}
	components, healthy := CheckHealth(checkers)
	require.Len(t, components, 1)
	require.False(t, healthy)
	require.NotNil(t, components["fail"])
	require.Nil(t, components["pass"])
}

type fail struct{}

func (c fail) HealthCheck() error {
	return errors.New("fail")
}

func TestHealthzHandler(t *testing.T) {
	failing := Handler(map[string]Checker{
		"mock": healthcheckFunc(func() error {
			return errors.New("health check failed")
		})})

	ok := Handler(map[string]Checker{
		"mock": healthcheckFunc(func() error {
			return nil
		})})

	var httpTests = []struct {
		name       string
		wantHeader int
		handler    http.Handler
	}{
		{"ok", 200, ok},
		{"fail", 500, failing},
	}
	for _, tt := range httpTests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/healthz", nil)
			tt.handler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, tt.wantHeader)
		})
	}

}

type healthcheckFunc func() error

func (fn healthcheckFunc) HealthCheck() error {
	return fn()
}
