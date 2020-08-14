package klogs

import (
	config2 "github.com/cinarmert/klogs/cmd/klogs/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

type Op interface {
	Run() error
}

var rootcmd = &cobra.Command{
	Use:   "klogs POD-QUERY",
	Short: "Multiple Pod/Container Target Visualizer",
	Long:  "Klogs provides a user interface for the logs from multiple Kubernetes pods/containers.",
	Run:   run,
}

var config = &config2.Config{}

func init() {
	rootcmd.Flags().StringVarP(&config.Kubeconfig, "kubeconfig", "k", "", "Kubeconfig file path")
	rootcmd.Flags().StringVarP(&config.Context, "context", "x", "", "Kubernetes context to use, default is $(kubectl config current-context)")
	rootcmd.Flags().StringVarP(&config.Namespace, "namespace", "n", "", "Kubernetes namespace to use, uses current namespace if not specified ")
	rootcmd.Flags().StringVarP(&config.ContainerFilter, "container", "c", ".*", "Regex to specify the containers to fetch logs from")
	rootcmd.Flags().Int64VarP(&config.Tail, "tail", "t", 100, "Number of lines to start tailing from")
	rootcmd.Flags().BoolP("verbose", "v", false, "print debug logs")
	log.SetLevel(log.WarnLevel)
}

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	op, err := parseArgs(args[0])
	if err != nil {
		log.Errorf("failed to parse given flags: %v", err)
		os.Exit(1)
	}

	if err := op.Run(); err != nil {
		if v, _ := cmd.Flags().GetBool("verbose"); v {
			log.Errorf("failed to fetch logs from Kubernetes pods/containers: %#v", err)
		} else {
			log.Errorf("failed to fetch logs from Kubernetes pods/containers, use -v flag for verbose logs")
		}
		os.Exit(1)
	}
}

// Execute is the entrypoint for the CLI.
func Execute() {
	if err := rootcmd.Execute(); err != nil {
		log.Fatalf("error running the init command")
	}
}
