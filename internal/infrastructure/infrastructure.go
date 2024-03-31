package infrastructure

import (
	"ara-server/internal/infrastructure/configuration"
	"ara-server/internal/infrastructure/tokenizer"
)

type Infrastructure struct {
	*configuration.Config
	*tokenizer.Tokenizer
}

func NewInfrastructure(config *configuration.Config) *Infrastructure {
	appConfig := config.GetConfig()

	tokenizer := tokenizer.NewTokenizer(&appConfig)

	return &Infrastructure{
		config,
		tokenizer,
	}
}
