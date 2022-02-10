package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	t.Run("Random string's length should equal as in param", func(t *testing.T) {
		r := RandomString{5}
		out := r.RandString()

		assert.Equal(t, 5, len(out))
	})
}
