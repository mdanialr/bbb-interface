package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	t.Run("Random string's lenght should equal as in param", func(t *testing.T) {
		out := RandString(5)

		assert.Equal(t, 5, len(out))
	})
}
