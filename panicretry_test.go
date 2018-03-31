package panicretry_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acoshift/panicretry"
)

func TestPanicRetry(t *testing.T) {
	i := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i++
		if i < 5 {
			panic("error")
		}
	})
	createRunner := func(p http.Handler) func() {
		return func() {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			p.ServeHTTP(nil, r)
		}
	}

	t.Run("Default", func(t *testing.T) {
		i = 0
		p := createRunner(panicretry.New(panicretry.Config{})(h))
		assert.Panics(t, p)
	})

	t.Run("Fail", func(t *testing.T) {
		i = 0
		p := createRunner(panicretry.New(panicretry.Config{MaxAttempts: 1})(h))
		assert.Panics(t, p)
		assert.Equal(t, 1, i)
	})

	t.Run("Success", func(t *testing.T) {
		i = 0
		p := createRunner(panicretry.New(panicretry.Config{MaxAttempts: 5})(h))
		assert.NotPanics(t, p)
		assert.Equal(t, 5, i)
	})

	t.Run("Success-2", func(t *testing.T) {
		i = 0
		p := createRunner(panicretry.New(panicretry.Config{MaxAttempts: 10})(h))
		assert.NotPanics(t, p)
		assert.Equal(t, 5, i)
	})

	t.Run("Skip", func(t *testing.T) {
		i = 0
		p := createRunner(panicretry.New(panicretry.Config{MaxAttempts: 10, Skipper: func(r *http.Request) bool { return true }})(h))
		assert.Panics(t, p)
		assert.Equal(t, 1, i)
	})
}
