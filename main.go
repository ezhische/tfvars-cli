package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	envparse "github.com/hashicorp/go-envparse"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	//go:embed template.txt
	terrafromrcTemplateStr string
	terrafromrcTemplate    = template.Must(template.New("terrafromrc").Parse(terrafromrcTemplateStr))
	version                string
	mirrorList             = [...]string{"terraform-mirror.yandexcloud.net", "registry.comcloud.xyz"}
	secretsClient          coreV1Types.SecretInterface
	project                *string = flag.String("project", "none", "project name for tfvars")
	secret                 *string = flag.String("secret", "test", "Secret name for tfvars")
	namespase              *string = flag.String("n", "default", "Namespace for secret")
	clusterMode            *bool   = flag.Bool("cluster", true, "Set -cluster=false for local test")
	configFile             *string = flag.String("config", "config", "Config file name. Defaut config")
	terraformrc            *bool   = flag.Bool("terraformrc", false, "Chose terraformrc mirror")
	showversion            *bool   = flag.Bool("version", false, "Print version")
	kubepath               string  = "/.kube/"
	ctx                    *string = flag.String("context", "pulse/agents:dev-test", "Context for kubeconfig")
)

func printVersion() {
	fmt.Println("Version: ", version)
	os.Exit(0)
}

func initClient() error {
	var err error
	var config *rest.Config
	if *clusterMode {
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = buildConfigFromFlags(*ctx, kubeconfig)
	} else {
		kubeconfig := os.Getenv("HOME") + kubepath + *configFile
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return err
	}
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return err
	}
	secretsClient = clientset.CoreV1().Secrets(*namespase)
	return nil
}

func main() {
	wg := new(sync.WaitGroup)
	defer wg.Wait()
	flag.Parse()
	if *showversion {
		printVersion()
	}
	if *terraformrc {
		wg.Add(1)
		createTerraformrc(wg)
	} else {
		if _, err := readAndPrintSecret(*secret); err != nil {
			log.Panic(err.Error())
		}
	}
}

func createTerraformrc(wg *sync.WaitGroup) error {
	rcfilepath := os.Getenv("HOME") + "/.terraformrc"
	for _, elem := range mirrorList {
		if checkMirror(elem) {
			if err := os.WriteFile(rcfilepath, terraformrcFile(elem), 0644); err != nil {
				return err
			}
			break
		}
	}
	wg.Done()
	return nil
}

func readAndPrintSecret(name string) (*coreV1.Secret, error) {
	if err := initClient(); err != nil {
		return nil, err
	}
	secret, err := secretsClient.Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for key, value := range secret.Data {
		if string(key) == *project {
			envs, _ := envparse.Parse(strings.NewReader(string(value)))
			for key, value := range envs {
				fmt.Printf("export %s=%s\n", key, value)
			}
		}
	}
	return secret, nil
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func checkMirror(mirr string) bool {
	link := fmt.Sprintf("https://%s/registry.terraform.io/hashicorp/random/index.json", mirr)
	resp, err := http.Get(link)
	if err != nil {
		return false
	} else {
		return resp.StatusCode == 200
	}
}

func terraformrcFile(site string) []byte {
	terraformrcBuffer := new(bytes.Buffer)
	err := terrafromrcTemplate.Execute(terraformrcBuffer, site)
	if err != nil {
		log.Panic(err.Error())
	}
	return terraformrcBuffer.Bytes()
}
