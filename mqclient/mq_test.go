package mqclient

import (
	"testing"

	"github.com/mheers/chrome-nats-proxy/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg := config.GetFakeConfig()

	mqClient, err := Init(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, mqClient)
}
