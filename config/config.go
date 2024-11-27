package config

import (
	"context"
	"os"
	"time"

	"dario.cat/mergo"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/mheers/chrome-nats-proxy/helpers"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

// Config describes the config
type Config struct {
	MQJWT       string `env:"CHROME_NATS_PROXY_MQ_JWT"`
	MQURI       string `env:"CHROME_NATS_PROXY_MQ_URI"`
	MQUSeed     string `env:"CHROME_NATS_PROXY_MQ_USEED"`
	MQCredsPath string `env:"CHROME_NATS_PROXY_MQ_CREDS_PATH"`
}

func (cfg *Config) OverlayConfigWithEnv() error {
	ctx := context.Background()
	overlayCfg := &Config{}
	err := envconfig.Process(ctx, overlayCfg)
	if err != nil {
		return err
	}

	err = mergo.Merge(cfg, overlayCfg, mergo.WithOverride)
	if err != nil {
		return err
	}
	return nil
}

// GetFakeConfig creates a config for testing purposes only
func GetFakeConfig() *Config {
	gofakeit.Seed(time.Now().UTC().UnixNano())

	cfg := &Config{
		MQJWT:   "eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiI2MlRRVlU0WkQ1RlVDV0NaS09ZWkg3TU1YM1Q3SzROVk80VkJEQVE2UTc0S0JZN01IRVpRIiwiaWF0IjoxNjUxMzE1Mzc4LCJpc3MiOiJBQjVVVU1XRTdMQkVQVjNSSjVUQVRIT1Y1RjJXT0xRREtZQ0VTVFBNQ1VISVAyM0VQQVpNMlBDRSIsIm5hbWUiOiJ0ZXN0Iiwic3ViIjoiVUJJSkhONlpDT0g0M0pCVFJLVzQyRk5aMjJZQ1JBM09DSFM3SUg3QVBaNjJLMzVVN1E2VzNYSVIiLCJuYXRzIjp7InB1YiI6eyJhbGxvdyI6WyJmbnhwLioiLCJfSU5CT1guXHUwMDNlIl19LCJzdWIiOnsiYWxsb3ciOlsiZm54cC4qIiwiX0lOQk9YLlx1MDAzZSJdfSwic3VicyI6LTEsImRhdGEiOi0xLCJwYXlsb2FkIjotMSwiYmVhcmVyX3Rva2VuIjp0cnVlLCJ0eXBlIjoidXNlciIsInZlcnNpb24iOjJ9fQ.t4nO6cJuTumTTF0mIzw64iYnTBjR_2DGPcFS-hYq2dQn5KS1Tuk5cfpJsqPXreuWCuQJjnM3-QzGrLLyhqeHDA",
		MQUSeed: "SUAIX7MNFMX7G7LB2P3EB53HDCAEMNRM4HLPGMJ7U4OZU55VWDXTT655EU",
		MQURI:   "ws://localhost:9222",
	}

	err := cfg.OverlayConfigWithEnv()
	if err != nil {
		return nil
	}

	logLevel := os.Getenv("LOGLEVEL")
	if logLevel != "" {
		helpers.SetLogLevel(logLevel)
	}

	return cfg
}

var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		configInstance = &Config{}
		err := configInstance.OverlayConfigWithEnv()
		if err != nil {
			logrus.Fatalf("config: %s", err)
		}
	}
	return configInstance
}
