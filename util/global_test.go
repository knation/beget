package util_test

import (
	"beget/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceModes(t *testing.T) {
	assert.Equal(t, util.ServiceMode("release"), util.ReleaseMode)
	assert.Equal(t, util.ServiceMode("debug"), util.DebugMode)
}
