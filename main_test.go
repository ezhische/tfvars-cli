package main

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
	"testing"
)

func TestCheckMirror(t *testing.T) {
	got := checkMirror("mail.com")
	want := false
	if got != want {
		t.Errorf("got %t; want %t", got, want)
	}
}

func TestTerraformrcFile(t *testing.T) {
	got := string(terraformrcFile("mail.ru"))
	want := testTemplate
	if got != want {
		fmt.Println(len(got), len(want))
		t.Errorf("got %s; want %s", got, want)
	}
}

var testTemplate = `provider_installation {
  network_mirror {
    url = "https://mail.ru/"
    include = ["registry.terraform.io/*/*"]
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}`

func TestInitClient(t *testing.T) {
	*clusterMode = true
	os.Setenv("KUBECONFIG", "/dev/null")
	got := initClient()
	if got == nil {
		t.Errorf("got %s; want nil", got.Error())
	}
}

func TestInitClient2(t *testing.T) {
	kubepath = "/"
	os.Setenv("HOME", "/tmp")
	*configFile = "kubetestconfig"
	os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777))
	*clusterMode = false
	got := initClient()
	os.Remove("/tmp/kubetestconfig")
	if got != nil {
		t.Errorf("got %t; want nil", got)
	}
}

func TestCreateTerraformrc(t *testing.T) {
	os.Setenv("HOME", "/s/sda")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	got := createTerraformrc(wg)
	if got == nil {
		t.Errorf("got %t; want Error", got)
	}
}

func TestCreateTerraformrc2(t *testing.T) {
	os.Setenv("HOME", "/tmp")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	got := createTerraformrc(wg)
	if got != nil {
		t.Errorf("got %t; want nil", got)
	}
}

func TestBuildConfigFromFlags(t *testing.T) {
	os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777))
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
	os.WriteFile("/tmp/kubetestconfig", []byte(configTemplate), fs.FileMode(0777))
	_, got := buildConfigFromFlags("default", "/tmp/kubetestconfig")
	os.Remove("/tmp/kubetestconfig")
	if got != nil {
		t.Errorf("got %t; want Error", got)
	}
}

var configTemplate = `apiVersion: v1
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
