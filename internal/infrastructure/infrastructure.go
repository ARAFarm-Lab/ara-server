package infrastructure

import "ara-server/internal/infrastructure/configuration"

type Infrastructure struct {
	configuration.Config
}

func NewInfrastructure(config configuration.Config) *Infrastructure {
	return &Infrastructure{
		config,
	}
}
