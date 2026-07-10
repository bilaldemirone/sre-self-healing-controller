package controller

import (
	"context"
	"fmt"

	"github.com/bilaldemir/sre-case/controller-runtime/internal/analyzer"
	cfg "github.com/bilaldemir/sre-case/controller-runtime/internal/config"
	promclient "github.com/bilaldemir/sre-case/controller-runtime/internal/prometheus"
	"github.com/bilaldemir/sre-case/controller-runtime/internal/queue"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodReconciler struct {
	client.Client
	Recorder   record.EventRecorder
	Prometheus *promclient.Client
	Queue      *queue.RemediationQueue
}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pod corev1.Pod

	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if pod.Namespace != "default" {
		return ctrl.Result{}, nil
	}

	if pod.Labels["app"] != "dummy-app" {
		return ctrl.Result{}, nil
	}

	restartCount := int32(0)
	for _, cs := range pod.Status.ContainerStatuses {
		restartCount += cs.RestartCount
	}

	fmt.Printf("Pod detected: %s/%s phase=%s restarts=%d\n",
		pod.Namespace,
		pod.Name,
		pod.Status.Phase,
		restartCount,
	)

	promRestartValue, err := r.Prometheus.GetPodRestartMetric(pod.Namespace, pod.Name)
	if err != nil {
		fmt.Printf("Prometheus restart query failed: %v\n", err)
	} else {
		fmt.Printf("Prometheus restart metric for pod %s: %s\n", pod.Name, promRestartValue)
	}

	cpuUsage, err := r.Prometheus.GetPodCPUUsage(pod.Namespace, pod.Name)
	if err != nil {
		fmt.Printf("Prometheus CPU query failed: %v\n", err)
	} else {
		fmt.Printf("Prometheus CPU usage for pod %s: %.4f cores\n", pod.Name, cpuUsage)

	}
	thresholds := cfg.LoadThresholdConfig(ctx, r.Client)

	decision := analyzer.AnalyzePod(analyzer.PodHealthInput{
		Pod:          pod,
		RestartCount: restartCount,
		CPUUsage:     cpuUsage,
		Thresholds:   thresholds,
	})

	if decision.Action == analyzer.ActionCreateEvent {
		r.Recorder.Eventf(
			&pod,
			corev1.EventTypeWarning,
			decision.Reason,
			"%s",
			decision.Message,
		)
	}
	if decision.Action == analyzer.ActionDeletePod {
		r.Recorder.Eventf(
			&pod,
			corev1.EventTypeWarning,
			decision.Reason,
			"%s",
			decision.Message,
		)

		r.Queue.Enqueue(queue.RemediationAction{
			Action:    decision.Action,
			Severity:  decision.Severity,
			Namespace: pod.Namespace,
			PodName:   pod.Name,
			Reason:    decision.Reason,
			Message:   decision.Message,
		})

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
