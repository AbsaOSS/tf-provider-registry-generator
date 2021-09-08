package etl

import (
	"testing"

	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/stretchr/testify/assert"
)

// todo:

func TestEtl(t *testing.T) {
	e, _ := NewEtl(config.Config{})
	err := e.Run()
	assert.NoError(t, err)
}
