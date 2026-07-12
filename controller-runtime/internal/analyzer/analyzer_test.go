package analyzer

import (
	"testing"

	cfg "github.com/bilaldemir/sre-case/controller-runtime/internal/config"
	corev1 "k8s.io/api/core/v1"
)

func TestAnalyzePodHealthy(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 0,
		CPUUsage:     0.01,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if !decision.Healthy {
		t.Fatalf("expected pod to be healthy")
	}

	if decision.Action != ActionNone {
		t.Fatalf("expected action %s, got %s", ActionNone, decision.Action)
	}
}

func TestAnalyzePodFailed(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodFailed,
			},
		},
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Action != ActionDeletePod {
		t.Fatalf("expected action %s, got %s", ActionDeletePod, decision.Action)
	}

	if decision.Severity != SeverityCritical {
		t.Fatalf("expected severity %s, got %s", SeverityCritical, decision.Severity)
	}
}

func TestAnalyzePodHighRestartCount(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 3,
		CPUUsage:     0.01,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Reason != "HighRestartCount" {
		t.Fatalf("expected reason HighRestartCount, got %s", decision.Reason)
	}

	if decision.Action != ActionCreateEvent {
		t.Fatalf("expected action %s, got %s", ActionCreateEvent, decision.Action)
	}
}

func TestAnalyzePodHighCPUUsage(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 0,
		CPUUsage:     0.10,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Reason != "HighCPUUsage" {
		t.Fatalf("expected reason HighCPUUsage, got %s", decision.Reason)
	}

	if decision.Action != ActionCreateEvent {
		t.Fatalf("expected action %s, got %s", ActionCreateEvent, decision.Action)
	}
}
