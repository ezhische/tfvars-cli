package kubeclient

import (
	"context"
	"fmt"
	"os"

	"github.com/ezhische/tfvar-cli/internal/config"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	cmd "k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	c coreV1Types.SecretInterface
}

func New(cfg *config.Config) (*Client, error) {
	var err error
	var config *rest.Config
	if *cfg.ClusterMode {
		config, err = buildConfigFromFlags(*cfg.Context, os.Getenv("KUBECONFIG"))
	} else {
		config, err = cmd.BuildConfigFromFlags("", *cfg.ConfigFile)
	}
	if err != nil {
		return nil, fmt.Errorf("client config creation failed: %w", err)
	}
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("client creation failed: %w", err)
	}
	secretsClient := clientset.CoreV1().Secrets(*cfg.Namespace)

	return &Client{c: secretsClient}, nil
}

func (c *Client) ReadSecret(name string) (*coreV1.Secret, error) {
	secret, err := c.c.Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s with error %w", name, err)
	}
	return secret, nil
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return cmd.NewNonInteractiveDeferredLoadingClientConfig(
		&cmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&cmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
