package metrics

import "go.uber.org/zap"

const GroupKeySeparator = "||"

type AppInfo struct {
	Name        string
	Environment string
	Version     string
}

type MetricContainer struct {
	Namespace    string
	Name         string
	Description  string
	Type         string
	Labels       []string
	LabelsValues []string
	Value        float64
}

func InitOnce(cfg *Config, logger *zap.Logger, app AppInfo) {
	// Создание реестра метрик
	initRegistry(logger.Sugar().Errorf)
}
