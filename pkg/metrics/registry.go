package metrics

import (
	"log"
	"sync"
)

type loggerFunc func(format string, args ...interface{})

type registry struct {
	mutex   sync.RWMutex
	metrics map[string]RegistryMetric
	logf    loggerFunc
}

type RegistryMetric interface {
	getMetrics() ([]*MetricContainer, error)
}

var globalRegistry registry

func init() {
	initRegistry(log.Printf)
}

func initRegistry(logf loggerFunc) {
	globalRegistry = registry{
		mutex:   sync.RWMutex{},
		metrics: map[string]RegistryMetric{},
		logf:    logf,
	}
}

func (r *registry) addMetric(name string, metric RegistryMetric) {
	r.mutex.Lock()
	r.metrics[name] = metric
	r.mutex.Unlock()
}

func (r *registry) getMetric(name string) RegistryMetric {
	r.mutex.RLock()
	metric := r.metrics[name]
	r.mutex.RUnlock()

	return metric
}

func (r *registry) getMetrics() map[string][]*MetricContainer {
	metricsMap := map[string][]*MetricContainer{}

	r.mutex.RLock()
	for _, sm := range r.metrics {
		if gsm, err := sm.getMetrics(); err == nil && len(gsm) > 0 {
			metricsMap[gsm[0].Namespace+gsm[0].Name] = gsm
		}
	}
	r.mutex.RUnlock()

	return metricsMap
}

func (r *registry) markCollision(metricName string, structName string) {
	r.logf(
		"error: %s metric is not a %s",
		metricName,
		structName,
	)
	panic("Metrics names collision")
}

func (r *registry) checkLabels(received []string, expected []string, name string) {
	if len(received) != len(expected) {
		r.logf(
			"error: %s metric has %d labels but %d labels were received",
			name,
			len(expected),
			len(received),
		)
		panic("Metrics labels collision")
	}
}
