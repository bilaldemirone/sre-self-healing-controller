package queue

import (
	"fmt"

	"github.com/bilaldemir/sre-case/controller-runtime/internal/analyzer"
)

type RemediationAction struct {
	Action    analyzer.Action
	Severity  analyzer.Severity
	Namespace string
	PodName   string
	Reason    string
	Message   string
}

type RemediationQueue struct {
	items chan RemediationAction
}

func NewRemediationQueue(size int) *RemediationQueue {
	return &RemediationQueue{
		items: make(chan RemediationAction, size),
	}
}

func (q *RemediationQueue) Enqueue(action RemediationAction) {
	q.items <- action

	fmt.Printf("Action queued: %s %s/%s reason=%s severity=%s\n",
		action.Action,
		action.Namespace,
		action.PodName,
		action.Reason,
		action.Severity,
	)
}

func (q *RemediationQueue) Items() <-chan RemediationAction {
	return q.items
}
