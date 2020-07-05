package klogs

import (
	config2 "github.com/cinarmert/klogs/cmd/klogs/config"
	"github.com/cinarmert/klogs/cmd/klogs/pod"
	"github.com/cinarmert/klogs/cmd/klogs/ui"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"regexp"
)

// LogOp satisfies the Op interface and hold the filters and configs
// for the process.
type LogOp struct {
	Config          *config2.Config
	PodFilter       *regexp.Regexp
	ContainerFilter *regexp.Regexp
	targets         []*pod.Target
	clientSet       *kubernetes.Clientset
}

// Run is the entrypoint for fetching logs from all targets. It starts
// the ui thread as well.
func (l *LogOp) Run() error {
	logManager, err := pod.NewManagerFromClientSet(l.clientSet, config)
	if err != nil {
		return errors.Wrap(err, "failed to create log instances")
	}

	targets := logManager.GetTargets()
	if len(targets) == 0 {
		return errors.New("failed to connect to pods/containers (maybe none are running)")
	}

	uiManager, err := ui.NewManager(targets, l.clientSet.CoreV1())
	if err != nil {
		return errors.Wrap(err, "could not create ui manager")
	}

	uiManager.Run()
	return nil
}
