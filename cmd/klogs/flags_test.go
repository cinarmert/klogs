package klogs

import (
	config2 "github.com/cinarmert/klogs/cmd/klogs/config"
	"os"
	"path/filepath"
	"testing"
)

func Test_getHomeDir(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		env     string
	}{
		{
			name:    "successful",
			wantErr: false,
			env:     "",
		},
		{
			name:    "unsuccessful",
			wantErr: true,
			env:     "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("_FORCE_HOME_ERR", tt.env)
				defer os.Unsetenv("_FORCE_HOME_ERR")
			}
			_, err := getHomeDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("getHomeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getKubeconfigDir(t *testing.T) {
	createKubeconfigIfNotExists()
	tests := []struct {
		name    string
		wantErr bool
		env     string
	}{
		{
			name:    "successful",
			wantErr: false,
			env:     "",
		},
		{
			name:    "unsuccessful",
			wantErr: true,
			env:     "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("_FORCE_KUBECONFIG_ERR", tt.env)
				defer os.Unsetenv("_FORCE_KUBECONFIG_ERR")
			}
			_, err := getKubeconfigDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("getHomeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_parseArgs(t *testing.T) {
	createKubeconfigIfNotExists()
	type env struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		arg     string
		env     *env
		config  *config2.Config
		wantErr bool
	}{
		{
			name:    "invalid regex",
			arg:     "[",
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name:    "no args given",
			arg:     "",
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name: "home directory error, no args",
			arg:  "",
			env: &env{
				key: "_FORCE_HOME_ERR",
				val: "true",
			},
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name: "kubeconfig directory error, no args",
			arg:  "",
			env: &env{
				key: "_FORCE_KUBECONFIG_ERR",
				val: "true",
			},
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name: "kubeconfig directory error",
			arg:  "test",
			env: &env{
				key: "_FORCE_KUBECONFIG_ERR",
				val: "true",
			},
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name: "home directory error",
			arg:  "test",
			env: &env{
				key: "_FORCE_HOME_ERR",
				val: "true",
			},
			config:  &config2.Config{},
			wantErr: true,
		},
		{
			name: "invalid container filter",
			arg:  "test",
			config: &config2.Config{
				ContainerFilter: "[",
			},
			wantErr: true,
		},
		{
			name: "invalid context name",
			arg:  "test",
			config: &config2.Config{
				Context: "invalid context",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config = tt.config
			if tt.env != nil {
				os.Setenv(tt.env.key, tt.env.val)
				defer os.Unsetenv(tt.env.key)
			}
			_, err := parseArgs(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createKubeconfigIfNotExists() {
	home, _ := os.UserHomeDir()
	kubeconfigPath := filepath.Join(home, "/.kube")
	if _, err := os.Stat(filepath.Join(kubeconfigPath + "/config")); err != nil {
		os.MkdirAll(kubeconfigPath, os.ModePerm)
		os.Create(filepath.Join(kubeconfigPath, "config"))
	}
}
