package dexpro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateProjectID(t *testing.T) {
	assert.Equal(t, "659bc69a-2b8f-5746-99f5-4ba0f0a09238", GenerateProjectID("squeeze.docker.localhost").String())
}
