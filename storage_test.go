package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/storage"
	assertTest "github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	assert := assertTest.New(t)

	t.Run("should use non-existing storage class", func(t *testing.T) {
		_, err := NewStorage("test")

		assert.EqualError(err, "storage class test does not exist")
	})

	t.Run("should initialize Filesystem storage", func(t *testing.T) {
		_, err := NewStorage(storage.Filesystem)
		assert.NoError(err)
	})
}
