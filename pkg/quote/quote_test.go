package quote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	quotes := []string{"Hello, World!", "Hello, Test!"}
	q := New(quotes)
	assert.NotNil(t, q, "Expected New() to return a non-nil Quotes instance")
}

func TestNewWithEmptyQuotes(t *testing.T) {
	quotes := []string{}
	assert.Panics(t, func() { New(quotes) }, "Expected New() to panic when no quotes are provided")
}

func TestGet(t *testing.T) {
	quotes := []string{"Hello, World!", "Hello, Test!"}
	q := New(quotes)

	for i := 0; i < 10; i++ {
		quote := q.Get()
		assert.Contains(t, quotes, string(quote), "Expected quote to be one of the original quotes")
	}
}
