package mq

import "encoding/json"

func (repo *Repository) PublishJSON(topic string, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if token := repo.mq.Publish(topic, 1, false, string(bytes)); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
