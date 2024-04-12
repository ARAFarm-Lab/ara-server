package mq

import (
	"ara-server/internal/infrastructure/metric"
	"encoding/json"
	"strings"
)

func (repo *Repository) PublishJSON(topic string, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if token := repo.mq.Publish(topic, 1, false, string(bytes)); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	topicSegment := strings.Split(topic, "/")
	repo.infra.PushCounter(metric.MQOutgoingMessage, map[string]string{
		"topic": topicSegment[0],
	})

	return nil
}
