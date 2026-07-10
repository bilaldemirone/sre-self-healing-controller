package analyzer

import (
	cfg "github.com/bilaldemir/sre-case/controller-runtime/internal/config"
	corev1 "k8s.io/api/core/v1"
)

type Severity string
type Action string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"

	ActionNone        Action = "NONE"
	ActionCreateEvent Action = "CREATE_EVENT"
	ActionDeletePod   Action = "DELETE_POD"
)

type Decision struct {
	Healthy        bool
	Severity       Severity
	ActionRequired bool
	Action         Action
	Reason         string
	Message        string
}

type PodHealthInput struct {
	Pod          corev1.Pod
	RestartCount int32
	CPUUsage     float64
	Thresholds   cfg.ThresholdConfig
}

func AnalyzePod(input PodHealthInput) Decision {
	if input.Pod.Status.Phase == corev1.PodFailed {
		return Decision{
			Healthy:        false,
			Severity:       SeverityCritical,
			ActionRequired: true,
			Action:         ActionDeletePod,
			Reason:         "PodFailed",
			Message:        "Pod is in Failed phase. Remediation is required.",
		}
	}

	if input.RestartCount >= input.Thresholds.RestartThreshold {
		return Decision{
			Healthy:        false,
			Severity:       SeverityWarning,
			ActionRequired: false,
			Action:         ActionCreateEvent,
			Reason:         "HighRestartCount",
			Message:        "Pod restart count exceeded configured threshold.",
		}
	}

	if input.CPUUsage > input.Thresholds.CPUThreshold {
		return Decision{
			Healthy:        false,
			Severity:       SeverityWarning,
			ActionRequired: false,
			Action:         ActionCreateEvent,
			Reason:         "HighCPUUsage",
			Message:        "Pod CPU usage exceeded configured threshold.",
		}
	}

	return Decision{
		Healthy:        true,
		Severity:       SeverityInfo,
		ActionRequired: false,
		Action:         ActionNone,
		Reason:         "Healthy",
		Message:        "Pod is healthy.",
	}
}
