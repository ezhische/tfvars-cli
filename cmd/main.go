package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ezhische/tfvar-cli/internal/config"
	"github.com/ezhische/tfvar-cli/internal/kubeclient"
	"github.com/ezhische/tfvar-cli/internal/rcfile"
	envparse "github.com/hashicorp/go-envparse"
	coreV1 "k8s.io/api/core/v1"
)

var (
	version    string
	mirrorList = []string{"terraform-mirror.yandexcloud.net", "registry.comcloud.xyz"}
)

func main() {
	cfg := config.NewConfig()
	cfg.Parse()

	if *cfg.ShowVersion {
		printVersion()
		os.Exit(0)
	}
	if *cfg.Terraformrc {
		if err := rcfile.CreateFile(mirrorList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	client, err := kubeclient.New(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	secret, err := client.ReadSecret(*cfg.Secret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	printSecret(cfg, secret)
}

func printSecret(cfg *config.Config, secret *coreV1.Secret) {
	for key, value := range secret.Data {
		if key == *cfg.Project {
			envs, _ := envparse.Parse(strings.NewReader(string(value)))
			for key, value := range envs {
				fmt.Printf("export %s=%s\n", key, value)
			}
		}
	}
}

func printVersion() {
	fmt.Println("Version: ", version)
}
