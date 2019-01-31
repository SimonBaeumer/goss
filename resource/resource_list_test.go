package resource

import (
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

const SuccessStatusCode = 200


func TestAddrMap_AppendSysResource(t *testing.T) {
	conf := util.Config{}
	systemMock := &system.System{
		NewHTTP: func(src string, sys *system.System, config util.Config) system.HTTP {
			return system.NewDefHTTP("http://goss.rocks", nil, conf)
		},
	}
	httpMap := HTTPMap{}

	got, err := httpMap.AppendSysResource("http://goss.rocks", systemMock, conf)

	assert.Nil(t, err)
	assert.Equal(t, "http://goss.rocks", got.HTTP)
	assert.Equal(t, SuccessStatusCode, got.Status)
	assert.Empty(t, got.Headers)
}
