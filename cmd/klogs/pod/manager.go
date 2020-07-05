package pod

import (
	"github.com/cinarmert/klogs/cmd/klogs/config"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Manager is responsible for fetching logs from the pods & containers that match
// the criteria provided by the user.
type Manager struct {
	Pods   *v1.PodList
	Config *config.Config
}

// NewManagerFromClientSet returns a Manager instance for a given clientset.
func NewManagerFromClientSet(clientSet *kubernetes.Clientset, config *config.Config) (pm *Manager, err error) {
	pm = &Manager{}
	pm.Pods, err = clientSet.CoreV1().Pods(config.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list pods")
	}

	pm.Config = config
	return pm, nil
}

// GetTargets returns the Target instances for the PodList of Manager.
func (pm *Manager) GetTargets() (res []*Target) {
	var targets []*Target
	for _, pod := range pm.Pods.Items {
		if pm.Config.PodRegex.MatchString(pod.Name) {
			targets = append(targets, pm.getTargetsFromPod(&pod)...)
		}
	}
	return targets
}

// getTargetsFromPod returns the Targets that satisfy the user-given criteria
// (e.g. container regex).
func (pm *Manager) getTargetsFromPod(pod *v1.Pod) (res []*Target) {
	var containers []v1.ContainerStatus
	containers = append(containers, pod.Status.InitContainerStatuses...)
	containers = append(containers, pod.Status.ContainerStatuses...)

	for _, c := range containers {
		if pm.Config.ContainerRegex.MatchString(c.Name) {
			t := NewTarget(pm.Config.Context, pod.Namespace, pod.Name, c.Name)
			res = append(res, t)
		}
	}
	return res
}
