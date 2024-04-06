package infrastructure

import (
	"ara-server/internal/infrastructure/configuration"
	"ara-server/internal/infrastructure/metric"
	"ara-server/internal/infrastructure/tokenizer"
)

type Infrastructure struct {
	*configuration.Config
	*metric.Metric
	*tokenizer.Tokenizer
}

func NewInfrastructure(config *configuration.Config) *Infrastructure {
	appConfig := config.GetConfig()
	metric := metric.NewMetric()
	tokenizer := tokenizer.NewTokenizer(&appConfig)

	return &Infrastructure{
		config,
		metric,
		tokenizer,
	}
}
