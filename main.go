package main

import (
	"errors"

	"github.com/mheers/chrome-nats-proxy/config"
	"github.com/mheers/chrome-nats-proxy/helpers"
	"github.com/mheers/chrome-nats-proxy/mqclient/models"
	"github.com/mheers/chrome-nats-proxy/proxy"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.GetConfig()

	helpers.SetLogLevel("trace")

	if err := initMQ(config); err != nil {
		logrus.Fatalf("initMQ: %s", err)
	}

}

func initMQ(config *config.Config) error {
	logrus.Debug("enabling feature 'Message Queue'")
	if config.MQURI == "" {
		return errors.New("no MQURI found")
	}

	if config.MQCredsPath == "" && (config.MQUSeed == "" || config.MQJWT == "") {
		return errors.New("no MQCredsPath or MQUSeed/MQJWT found")
	}

	mqClient, err := models.NewMQClient(config)
	if err != nil {
		return err
	}

	mqClient.Connection.Subscribe("instance.chrome-nats-proxy.keepalive", func(msg *nats.Msg) {
		logrus.Infof("Received message: %s", msg.Data)
		msg.RespondMsg(&nats.Msg{
			Subject: msg.Reply,
			Data:    []byte("pong"),
		})
	})

	proxy := proxy.NewProxy(mqClient)
	proxy.Start()

	return nil
}
