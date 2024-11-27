package helpers

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func HandleMQError(m *nats.Msg, err error) {
	logrus.Errorf("HandleMQError: %s", err)
	msg := &nats.Msg{
		Data: []byte(err.Error()),
	}
	m.RespondMsg(msg)
}

func HandleMQOK(m *nats.Msg) {
	type Response struct {
		Status string `json:"status"`
		Text   string `json:"text"`
	}
	response := &Response{
		Status: "OK",
		Text:   "",
	}
	responseJSON, _ := json.Marshal(response)
	msg := &nats.Msg{
		Data: responseJSON,
	}
	m.RespondMsg(msg)
}
