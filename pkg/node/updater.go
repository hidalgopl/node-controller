package node

import (
	"context"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/api/core/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

const envPrefix = ""

type Options struct {
	TargetOS   string `envconfig:"TARGET_OS" default:"Container Linux"`
	LabelKey   string `envconfig:"LABEL_KEY" default:"kubermatic.io/uses-container-linux"`
	LabelValue string `envconfig:"LABEL_VALUE" default:"true"`
}

func OptionsFromEnv() (Options, error) {
	var options Options
	err := envconfig.Process(envPrefix, &options)
	return options, err
}

func PrintEnvOptionsUsage() error {
	return envconfig.Usage(envPrefix, &Options{})
}

type NodeUpdater struct {
	Client  corev1client.NodesGetter
	Options Options
}

func (nu *NodeUpdater) IsNodeWithOS(node *v1.Node) bool {
	osImage := node.Status.NodeInfo.OSImage
	if !strings.Contains(osImage, nu.Options.TargetOS) {
		return false
	}
	return true
}

func (nu *NodeUpdater) AddLabel(node *v1.Node) *v1.Node {
	if node.ObjectMeta.Labels == nil {
		node.ObjectMeta.Labels = map[string]string{}
	}
	if nu.AlreadyHasLabelSet(node) {
		return node
	}
	node.ObjectMeta.Labels[nu.Options.LabelKey] = nu.Options.LabelValue
	return node
}

func (nu *NodeUpdater) AlreadyHasLabelSet(node *v1.Node) bool {
	// check if node has label set
	if val, ok := node.ObjectMeta.Labels[nu.Options.LabelKey]; ok {
		// if it has, check if value is correct
		if val != nu.Options.LabelValue {
			// if not, return false to allow AddLabel to be set correctly
			return false
		}
		return true
	}
	return false
}

func (nu *NodeUpdater) Update(ctx context.Context, node *v1.Node) *v1.Node {
	node = nu.AddLabel(node)
	result, err := nu.Client.Nodes().Update(node)
	if err != nil {
		log.Printf("error while updating node: %s", err.Error())
	}
	return result
}
