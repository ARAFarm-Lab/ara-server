package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	mutex      *sync.Mutex
	counterMap map[MetricKey]*prometheus.CounterVec
}

func NewMetric() *Metric {
	return &Metric{
		mutex:      new(sync.Mutex),
		counterMap: make(map[MetricKey]*prometheus.CounterVec),
	}
}

func (m *Metric) PushCounter(key MetricKey, values map[string]string) {
	var (
		counter *prometheus.CounterVec
		found   bool
	)
	if counter, found = m.counterMap[key]; !found {
		m.mutex.Lock()
		labels := make([]string, 0, len(values))
		for k := range values {
			labels = append(labels, k)
		}

		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ara_iot_" + string(key),
			},
			labels,
		)
		m.counterMap[key] = counter
		m.mutex.Unlock()
		prometheus.MustRegister(counter)
	}

	counter.With(values).Inc()
}
