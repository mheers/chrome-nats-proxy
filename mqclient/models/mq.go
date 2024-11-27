package models

import (
	"log"
	"time"

	"github.com/mheers/chrome-nats-proxy/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/sirupsen/logrus"
)

// MQClient descibes the message queue with connection and subscriptions
type MQClient struct {
	InstanceJWT   string
	Config        *config.Config
	Connection    *nats.Conn
	Subscriptions map[string]*nats.Subscription
}

// NewMQClient creates a new MQ
func NewMQClient(config *config.Config) (*MQClient, error) {
	conn, err := getConnection(config)
	if err != nil {
		return nil, err
	}

	mq := &MQClient{
		Config:        config,
		Connection:    conn,
		Subscriptions: make(map[string]*nats.Subscription),
	}

	return mq, nil
}

func getConnection(config *config.Config) (*nats.Conn, error) {
	url := config.MQURI // e.g. "nats://localhost:4222"
	logrus.Debugf("Connecting to nats mq url: %s", url)

	totalWait := 1 * time.Minute
	reconnectDelay := time.Second

	// Connect Options.
	opts := []nats.Option{
		nats.Name("cursed chrome Server"),
	}

	if config.MQCredsPath != "" {
		opts = append(opts, nats.UserCredentials(config.MQCredsPath))
	} else {
		instanceJWT := config.MQJWT
		uSeed := []byte(config.MQUSeed)
		opts = append(opts, nats.Token(instanceJWT))

		jwtCB := func() (string, error) {
			return instanceJWT, nil
		}
		sigCB := func(nonce []byte) ([]byte, error) {
			kp, _ := nkeys.FromSeed(uSeed)
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}
		opts = append(opts, nats.UserJWT(jwtCB, sigCB))
	}

	opts = append(opts, nats.Timeout(reconnectDelay*3))
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))

	// Connect to NATS
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	return nc, err
}
