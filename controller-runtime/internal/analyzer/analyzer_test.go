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

func TestAnalyzePodCPUEqualToThresholdIsHealthy(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 0,
		CPUUsage:     0.05,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if !decision.Healthy {
		t.Fatalf("expected pod to be healthy when CPU usage equals threshold")
	}

	if decision.Action != ActionNone {
		t.Fatalf("expected action %s, got %s", ActionNone, decision.Action)
	}
}

func TestAnalyzePodRestartBelowThresholdIsHealthy(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 2,
		CPUUsage:     0.01,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if !decision.Healthy {
		t.Fatalf("expected pod to be healthy below restart threshold")
	}

	if decision.Reason != "Healthy" {
		t.Fatalf("expected reason Healthy, got %s", decision.Reason)
	}
}

func TestAnalyzePodFailedTakesPriorityOverOtherConditions(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodFailed,
			},
		},
		RestartCount: 10,
		CPUUsage:     0.90,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Reason != "PodFailed" {
		t.Fatalf("expected PodFailed to have priority, got %s", decision.Reason)
	}

	if decision.Action != ActionDeletePod {
		t.Fatalf("expected action %s, got %s", ActionDeletePod, decision.Action)
	}

	if decision.Severity != SeverityCritical {
		t.Fatalf("expected severity %s, got %s", SeverityCritical, decision.Severity)
	}
}

func TestAnalyzePodRestartTakesPriorityOverCPU(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 5,
		CPUUsage:     0.90,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Reason != "HighRestartCount" {
		t.Fatalf(
			"expected HighRestartCount to have priority over CPU, got %s",
			decision.Reason,
		)
	}

	if decision.Action != ActionCreateEvent {
		t.Fatalf("expected action %s, got %s", ActionCreateEvent, decision.Action)
	}
}

func TestAnalyzePodHighCPUSeverity(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 0,
		CPUUsage:     0.06,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Healthy {
		t.Fatalf("expected pod with high CPU usage to be unhealthy")
	}

	if decision.Severity != SeverityWarning {
		t.Fatalf("expected severity %s, got %s", SeverityWarning, decision.Severity)
	}

	if decision.ActionRequired {
		t.Fatalf("expected high CPU condition not to require remediation")
	}
}

func TestAnalyzePodHighRestartSeverity(t *testing.T) {
	input := PodHealthInput{
		Pod: corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		RestartCount: 4,
		CPUUsage:     0.01,
		Thresholds: cfg.ThresholdConfig{
			CPUThreshold:     0.05,
			RestartThreshold: 3,
		},
	}

	decision := AnalyzePod(input)

	if decision.Healthy {
		t.Fatalf("expected pod with high restart count to be unhealthy")
	}

	if decision.Severity != SeverityWarning {
		t.Fatalf("expected severity %s, got %s", SeverityWarning, decision.Severity)
	}

	if decision.ActionRequired {
		t.Fatalf("expected restart warning not to require remediation")
	}
}
