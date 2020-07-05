package config

import "regexp"

// Config hold the flags provided to the command line interface.
type Config struct {
	Kubeconfig      string
	Context         string
	Namespace       string
	PodFilter       string
	ContainerFilter string
	PodRegex        *regexp.Regexp
	ContainerRegex  *regexp.Regexp
}
