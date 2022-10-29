package util_test

import (
	"beget/util"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestServiceModes(t *testing.T) {
	assert.Equal(t, util.ServiceMode("release"), util.ReleaseMode)
	assert.Equal(t, util.ServiceMode("debug"), util.DebugMode)
}

func TestInitConfigDefaults(t *testing.T) {
	// Add `testdata/` to path to make test configuration file discoverable
	viper.AddConfigPath("./testdata")

	err := util.InitConfig()

	assert.Nil(t, err)

	assert.Equal(t, util.DebugMode, util.Config.App.Mode)
	assert.Equal(t, 8080, util.Config.Server.Port)
}
