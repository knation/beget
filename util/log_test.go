package util_test

import (
	"beget/util"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogging(t *testing.T) {
	util.InitLogging()

	assert.Equal(t, "*zap.Logger", reflect.TypeOf(util.Logger).String())
	assert.Equal(t, "*zap.SugaredLogger", reflect.TypeOf(util.Sugar).String())
}
