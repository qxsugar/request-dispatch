package request_dispatch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	config := CreateConfig()
	config.LogLevel = "debug"
	config.MarkHeader = "X-Test-Header"
	config.MarkHosts = map[string][]string{
		"test": {"http://localhost:8080"},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("default handler"))
	})

	disp, err := New(context.Background(), nextHandler, config, "test")
	assert.NoError(t, err)

	t.Run("test dispatch with mark header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:9000", nil)
		req.Header.Set("X-Test-Header", "test")
		rr := httptest.NewRecorder()

		disp.ServeHTTP(rr, req)

		// Check if the request was dispatched to the correct host
		assert.Equal(t, "localhost:8080", req.Host)
	})

	t.Run("test dispatch without mark header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:9000", nil)
		rr := httptest.NewRecorder()

		disp.ServeHTTP(rr, req)

		// Check if the request was routed to the default handler
		assert.Equal(t, "default handler", rr.Body.String())
	})
}
