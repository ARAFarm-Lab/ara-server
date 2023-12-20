package mq

import (
	"ara-server/internal/infrastructure"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Repository struct {
	infra *infrastructure.Infrastructure
	mq    mqtt.Client
}

func NewRepository(infra *infrastructure.Infrastructure, mq mqtt.Client) *Repository {
	return &Repository{infra: infra, mq: mq}
}
