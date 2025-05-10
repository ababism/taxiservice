package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"strings"
	"sync"
)

// Counter (Счетчик) — это метрика, представляющая одно числовое значение (не только целочисленное), которое всегда увеличивается.
type Counter interface {
	Inc()
	Add(val float64)
	getMetric() (*MetricContainer, error)
}

// counter - реализация Counter интерфейса.
type counter struct { // implement Counter
	prom prometheus.Counter
}

func (c *counter) Inc() {
	c.prom.Inc()
}

func (c *counter) Add(val float64) {
	c.prom.Add(val)
}

func (c *counter) getMetric() (*MetricContainer, error) {
	mC := &dto.Metric{}
	if err := c.prom.Write(mC); err != nil {
		return nil, err
	}

	return &MetricContainer{
		Type:  "Counter",
		Value: mC.Counter.GetValue(),
	}, nil
}

// CounterOpts опции для создания метрики
type CounterOpts struct {
	Namespace   string
	Name        string
	Description string
}

// CounterVec метрика счетчика с поддержкой labels
type CounterVec struct { // implement RegistryMetric
	opts    CounterOpts
	mutex   sync.RWMutex
	labels  []string
	metrics map[string]Counter

	constructor *prometheus.CounterVec
}

// GetOrRegisterCounterVec создает или возвращает уже созданную фабрику метрик Counter.
func GetOrRegisterCounterVec(opts CounterOpts, labels []string) *CounterVec {
	if rm := globalRegistry.getMetric(opts.Name); rm != nil {
		if vm, ok := rm.(*CounterVec); ok {
			return vm
		}
		globalRegistry.markCollision(opts.Name, "CounterVec")
	}

	promConstructor := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Name:      opts.Name,
			Help:      opts.Description,
		},
		labels,
	)

	counterVec := &CounterVec{
		opts:        opts,
		labels:      labels,
		metrics:     map[string]Counter{},
		constructor: promConstructor,
	}

	globalRegistry.addMetric(opts.Name, counterVec)

	return counterVec
}

// WithLabelValues возвращает необходимый счетчик Counter в зависимости от переданных значений
func (c *CounterVec) WithLabelValues(labelValues ...string) Counter {
	return c.getOrRegisterLabelMetric(labelValues)
}

// getMetric дает альтернативу по получению метрик как пакет promhttp
func (c *CounterVec) getMetrics() ([]*MetricContainer, error) {
	var metrics []*MetricContainer
	for group, counter := range c.metrics {
		metric, _ := counter.getMetric()
		metric.Namespace = c.opts.Namespace
		metric.Name = c.opts.Name
		metric.Description = c.opts.Description

		metric.Labels = c.labels

		// Получаем нормальный список labels
		labelsValues := strings.Split(group, GroupKeySeparator)
		metric.LabelsValues = labelsValues

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// getOrRegisterLabelMetric это специальная конструкция - фабрика,
// которая в зависимости от переданных значений создаст вам новый и полноценный счетчик
func (c *CounterVec) getOrRegisterLabelMetric(labelValues []string) Counter {
	globalRegistry.checkLabels(labelValues, c.labels, c.opts.Name)

	key := strings.Join(labelValues, GroupKeySeparator)

	c.mutex.RLock()
	cm, exist := c.metrics[key]
	c.mutex.RUnlock()

	if !exist {
		prom := c.constructor.WithLabelValues(labelValues...)

		// оболочка метрики с которой и взаимодействует приложение
		cm = &counter{
			prom: prom,
		}

		c.mutex.Lock()
		c.metrics[key] = cm
		c.mutex.Unlock()
	}

	return cm
}
