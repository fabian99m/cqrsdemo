package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log/slog"
)

func loadAWSConfig() aws.Config {
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("eror loading aws config", "error", err)
		panic(err)
	}

	return sdkConfig
}
