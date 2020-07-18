package klogs

import (
	"fmt"
	internal "github.com/cinarmert/klogs/cmd/klogs/kubernetes"
	"github.com/pkg/errors"
	"os"
	"regexp"
)

func parseArgs(arg string) (_ *LogOp, err error) {
	if arg == "" {
		return nil, errors.New("pod query cannot be empty")
	}

	if config.Kubeconfig == "" {
		kubeconfigPath, err := getKubeconfigDir()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get kubeconfig path")
		}
		config.Kubeconfig = kubeconfigPath
	}

	config.PodRegex, err = regexp.Compile(arg)
	if err != nil {
		return nil, errors.Wrap(err, "could not compile pod filter regex")
	}

	config.ContainerRegex, err = regexp.Compile(config.ContainerFilter)
	if err != nil {
		return nil, errors.Wrap(err, "could not compile container filter regex")
	}

	clientConfig := internal.NewClientForContext(config.Kubeconfig, config.Context)
	clientSet, err := internal.NewClientSet(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create client set with config: %v", config.Kubeconfig))
	}

	return &LogOp{
		Config:    config,
		targets:   nil,
		clientSet: clientSet,
		tail:      config.Tail,
	}, nil
}

func getHomeDir() (string, error) {
	if os.Getenv("_FORCE_HOME_ERR") != "" { // testing
		return "", errors.New("error occurred when getting home directory")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}

	return home, nil
}

func getKubeconfigDir() (string, error) {
	if os.Getenv("_FORCE_KUBECONFIG_ERR") != "" { // testing
		return "", errors.New("error occurred when getting kubeconfig directory")
	}

	if config.Kubeconfig != "" {
		return config.Kubeconfig, nil
	}

	homePath, err := getHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}

	res := homePath + "/.kube/config"

	if _, err := os.Stat(res); os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("no kubeconfig found at %v", res))
	}

	return res, nil
}
