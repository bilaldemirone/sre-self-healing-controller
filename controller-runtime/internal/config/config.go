package config

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ThresholdConfig struct {
	CPUThreshold        float64
	RestartThreshold    int32
	NodeCPUThreshold    float64
	NodeMemoryThreshold float64
}

func LoadThresholdConfig(ctx context.Context, c client.Client) ThresholdConfig {
	cfg := ThresholdConfig{
		CPUThreshold:        0.05,
		RestartThreshold:    3,
		NodeCPUThreshold:    0.80,
		NodeMemoryThreshold: 0.85,
	}

	var cm corev1.ConfigMap
	err := c.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "sre-controller-config",
	}, &cm)

	if err != nil {
		return cfg
	}

	if v, ok := cm.Data["cpuThreshold"]; ok {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.CPUThreshold = parsed
		}
	}

	if v, ok := cm.Data["restartThreshold"]; ok {
		if parsed, err := strconv.Atoi(v); err == nil {
			cfg.RestartThreshold = int32(parsed)
		}
	}
	if v, ok := cm.Data["nodeCPUThreshold"]; ok {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.NodeCPUThreshold = parsed
		}
	}

	if v, ok := cm.Data["nodeMemoryThreshold"]; ok {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.NodeMemoryThreshold = parsed
		}
	}
	return cfg
}
