package kubeclient

import (
	"io/fs"
	"os"
	"testing"

	"github.com/ezhische/tfvar-cli/internal/config"
)

var cfg *config.Config

func TestConfig(t *testing.T) {
	cfg = config.NewConfig()
}
func TestInitClient(t *testing.T) {
	*cfg.ClusterMode = true
	*cfg.Context = "test"
	t.Setenv("KUBECONFIG", "/dev/null")
	_, got := New(cfg)
	if got == nil {
		t.Errorf("got %s; want nil", got)
	}
}

func TestInitClient2(t *testing.T) {
	*cfg.ClusterMode = false
	*cfg.ConfigFile = "/tmp/kubetestconfig"
	if err := os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777)); err != nil {
		t.Errorf("got %t; want nil", err)
	}
	_, got := New(cfg)
	os.Remove("/tmp/kubetestconfig")
	if got != nil {
		t.Errorf("got %t; want nil", got)
	}
}

func TestBuildConfigFromFlags(t *testing.T) {
	if err := os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777)); err != nil {
		t.Errorf("got %t; want nil", err)
	}
	_, got := buildConfigFromFlags("test", "/tmp/kubetestconfig")
	os.Remove("/tmp/kubetestconfig")
	if got == nil {
		t.Errorf("got %t; want Error", got)
	}
}

func TestBuildConfigFromFlags2(t *testing.T) {
	_, got := buildConfigFromFlags("test", "/tmp/fakeconfig")
	if got == nil {
		t.Errorf("got %t; want Error", got)
	}
}

func TestBuildConfigFromFlags3(t *testing.T) {
	if err := os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777)); err != nil {
		t.Errorf("got %t; want nil", err)
	}
	_, got := buildConfigFromFlags("default", "/tmp/kubetestconfig")
	os.Remove("/tmp/kubetestconfig")
	if got != nil {
		t.Errorf("got %t; want Error", got)
	}
}

const configTemplate = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data:
    server: https://10.10.10.10:6443
  name: dr-prod-cluster
contexts:
- context:
    cluster: dr-prod-cluster
    user: kubernetes-admin
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data:
    client-key-data:
`
