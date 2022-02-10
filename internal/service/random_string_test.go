package service

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	t.Run("Random string's length should equal as in param", func(t *testing.T) {
		r := RandomString{5}
		out := r.RandString()

		assert.Equal(t, 5, len(out))
	})

	t.Run("Random string should always produce different value w consistent length", func(t *testing.T) {
		r := RandomString{8}
		outOne := r.RandString()
		require.Equal(t, 8, len(outOne))

		outTwo := r.RandString()
		require.Equal(t, 8, len(outTwo))

		assert.NotEqual(t, outOne, outTwo)
	})
}
