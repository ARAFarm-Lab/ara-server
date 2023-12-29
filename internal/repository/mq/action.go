package mq

import "encoding/json"

func (repo *Repository) PublishJSON(topic string, payload PublishJSONPayload) error {
	bytes, err := json.Marshal([]interface{}{payload.ActionType, payload.Value})
	if err != nil {
		return err
	}

	if token := repo.mq.Publish(topic, 1, false, string(bytes)); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
