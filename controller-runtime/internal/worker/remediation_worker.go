package worker

import (
	"context"
	"fmt"

	"github.com/bilaldemir/sre-case/controller-runtime/internal/analyzer"
	"github.com/bilaldemir/sre-case/controller-runtime/internal/queue"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RemediationWorker struct {
	Client client.Client
	Queue  *queue.RemediationQueue
}

func (w *RemediationWorker) Start(ctx context.Context) {
	fmt.Println("Remediation worker started")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Remediation worker stopped")
			return
		case action := <-w.Queue.Items():
			w.handleAction(ctx, action)
		}
	}
}

func (w *RemediationWorker) handleAction(ctx context.Context, action queue.RemediationAction) {
	switch action.Action {
	case analyzer.ActionDeletePod:
		w.deletePod(ctx, action)
	default:
		fmt.Printf("Worker: unsupported action=%s reason=%s\n", action.Action, action.Reason)
	}
}

func (w *RemediationWorker) deletePod(ctx context.Context, action queue.RemediationAction) {
	pod := &corev1.Pod{}

	err := w.Client.Get(ctx, client.ObjectKey{
		Namespace: action.Namespace,
		Name:      action.PodName,
	}, pod)

	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Printf("Worker: pod already deleted: %s/%s\n", action.Namespace, action.PodName)
			return
		}
		fmt.Printf("Worker: failed to get pod: %v\n", err)
		return
	}

	if err := w.Client.Delete(ctx, pod); err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Printf("Worker: pod already deleted: %s/%s\n", action.Namespace, action.PodName)
			return
		}
		fmt.Printf("Worker: failed to delete pod: %v\n", err)
		return
	}

	fmt.Printf("Worker: deleted pod %s/%s reason=%s severity=%s\n",
		action.Namespace,
		action.PodName,
		action.Reason,
		action.Severity,
	)
}
