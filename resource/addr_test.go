package resource

import (
	"testing"
	"github.com/aelsabbahy/goss/util"
	"github.com/aelsabbahy/goss/system/mock_system"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddr_NewAddr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAddress := mock_system.NewMockAddr(ctrl)
	mockAddress.EXPECT().Address().Return("http://goss.rocks")
	mockAddress.EXPECT().Reachable().Return(true, nil)

	result, resultErr := NewAddr(mockAddress, util.Config{})

	assert.Equal(t, "http://goss.rocks", result.Address)
	assert.Equal(t, true, result.Reachable)
	assert.Nil(t, resultErr)
}
