package mqclient

import (
	"sync"

	"github.com/mheers/chrome-nats-proxy/config"
	"github.com/mheers/chrome-nats-proxy/mqclient/models"
)

var MQClient *models.MQClient

var once sync.Once

// Init initializes a message queue client
func Init(appConfig *config.Config) (*models.MQClient, error) {
	var err error
	once.Do(func() {

		mqClient, err := models.NewMQClient(appConfig)
		if err != nil {
			return
		}

		MQClient = mqClient

	})
	return MQClient, err
}
