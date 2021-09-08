package etl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtl(t *testing.T) {
	e, _ := NewEtl2(factoryMock)
	err := e.Run()
	assert.NoError(t, err)
}
