package metric

type MetricKey string

const (
	HTTPResponse      = "http_response"
	MQIncomingMessage = "mq_incoming_message"
	MQOutgoingMessage = "mq_outgoing_message"
)
