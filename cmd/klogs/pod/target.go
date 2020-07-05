package pod

import (
	"github.com/pkg/errors"
	"io"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"sync"
)

// Target holds the necessary information for a container inside a pod. Each container
// is converted into a Target instance in order to fetch logs.
type Target struct {
	Context   string
	Namespace string
	Pod       string
	Container string
}

// NewTarget return a target instance.
func NewTarget(context, namespace, pod, container string) *Target {
	return &Target{
		Context:   context,
		Namespace: namespace,
		Pod:       pod,
		Container: container,
	}
}

// StartThread is meant to run as a goroutine and fetch logs until the target goes down
// or process is stopped by the user.
func (l *Target) StartThread(core v1.CoreV1Interface, wg *sync.WaitGroup, w io.Writer) error {
	request := core.Pods(l.Namespace).GetLogs(l.Pod, &corev1.PodLogOptions{
		Container: l.Container,
		Follow:    false,
	})
	rc, err := request.Stream()
	if err != nil {
		return errors.Wrap(err, "failed to stream logs")
	}
	_, err = io.Copy(w, rc)
	if err != nil {
		return errors.Wrap(err, "failed to copy logs to destination writer")
	}

	wg.Done()
	return nil
}
