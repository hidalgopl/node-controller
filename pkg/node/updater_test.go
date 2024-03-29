package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	core "k8s.io/client-go/testing"
)

func TestOptionsFromEnv(t *testing.T) {
	tt := []struct {
		testName           string
		expectedLabelKey   string
		expectedLabelValue string
		expectedTargetOS   string
	}{
		{
			testName:           "defaults",
			expectedLabelKey:   "kubermatic.io/uses-container-linux",
			expectedLabelValue: "true",
			expectedTargetOS:   "Container Linux",
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			options, _ := OptionsFromEnv()
			assert.Equal(t, tc.expectedLabelKey, options.LabelKey)
			assert.Equal(t, tc.expectedLabelValue, options.LabelValue)
			assert.Equal(t, tc.expectedTargetOS, options.TargetOS)
		})
	}

}

func TestNodeUpdater_addLabel(t *testing.T) {
	tt := []struct {
		testName string
		labelKey string
		labelVal string
	}{
		{
			testName: "happy path",
			labelKey: "kubermatic.io/uses-container-linux",
			labelVal: "true",
		},
	}
	fakeClient := &fakecorev1.FakeCoreV1{Fake: &core.Fake{}}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			node := &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{},
				},
			}
			nodeUpdater := &NodeUpdater{
				Client: fakeClient,
				Options: Options{
					TargetOS:   "",
					LabelKey:   tc.labelKey,
					LabelValue: tc.labelVal,
				},
			}
			labeledNode := nodeUpdater.addLabel(node)
			nodeLabels := labeledNode.GetLabels()
			assert.Equal(t, nodeLabels[tc.labelKey], tc.labelVal)

		})
	}
}

func TestNodeUpdater_Update(t *testing.T) {
	tt := []struct {
		testName string
		labelKey string
		labelVal string
		expected map[string]string
	}{
		{
			testName: "happy path",
			labelKey: "kubermatic.io/uses-container-linux",
			labelVal: "true",
			expected: map[string]string{
				"kubermatic.io/uses-container-linux": "true",
			},
		},
	}
	fakeClient := &fakecorev1.FakeCoreV1{Fake: &core.Fake{}}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			nodeUpdater := &NodeUpdater{
				Client: fakeClient,
				Options: Options{
					LabelKey:   tc.labelKey,
					TargetOS:   "",
					LabelValue: tc.labelVal,
				},
			}
			node := &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "my-node",
					Labels: map[string]string{},
				},
			}
			node, _ = nodeUpdater.Client.Nodes().Create(node)
			nodeUpdater.Update(context.TODO(), node)
			assert.Equal(t, tc.expected, node.GetLabels())
		})
	}
}

func TestNodeUpdater_isNodeWithOS(t *testing.T) {
	tt := []struct {
		testName          string
		nodeOSImage       string
		wantedNodeOsImage string
		expected          bool
	}{
		{
			testName:          "happy path - ContainerLinux node",
			nodeOSImage:       "Container Linux by CoreOS 2135.6.0 (Rhyolite)",
			wantedNodeOsImage: "Container Linux",
			expected:          true,
		},
		{
			testName:          "unhappy path - Ubuntu node",
			nodeOSImage:       "Ubuntu 18.04.2 LTS",
			wantedNodeOsImage: "Container Linux",
			expected:          false,
		},
		{
			testName:          "unhappy path - CentOS node",
			nodeOSImage:       "CentOS Linux 7 (Core)",
			wantedNodeOsImage: "Container Linux",
			expected:          false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			node := &v1.Node{
				Status: v1.NodeStatus{
					NodeInfo: v1.NodeSystemInfo{
						OSImage: tc.nodeOSImage,
					},
				},
			}
			fakeClient := &fakecorev1.FakeCoreV1{Fake: &core.Fake{}}
			nU := &NodeUpdater{
				Client: fakeClient,
				Options: Options{
					TargetOS:   tc.wantedNodeOsImage,
					LabelValue: "",
					LabelKey:   "",
				},
			}
			result := nU.isNodeWithOS(node)
			assert.Equal(t, tc.expected, result)

		})
	}
}

func TestNodeUpdater_alreadyHasLabelSet(t *testing.T) {
	tt := []struct {
		testName       string
		labels         map[string]string
		expectedResult bool
		expectedLabels map[string]string
	}{
		{
			testName:       "node already has correct label and value",
			labels:         map[string]string{"label-key": "label-value"},
			expectedResult: true,
		},
		{
			testName:       "node doesnt have label",
			labels:         map[string]string{},
			expectedResult: false,
		},
		{
			testName:       "node has correct label but wrong value",
			labels:         map[string]string{"label-key": "wrong-label-value"},
			expectedResult: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			node := &v1.Node{
				Status: v1.NodeStatus{
					NodeInfo: v1.NodeSystemInfo{
						OSImage: "target-os",
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "my-node",
					Labels: tc.labels,
				},
			}
			fakeClient := &fakecorev1.FakeCoreV1{Fake: &core.Fake{}}
			nU := &NodeUpdater{
				Client: fakeClient,
				Options: Options{
					TargetOS:   "target-os",
					LabelValue: "label-value",
					LabelKey:   "label-key",
				},
			}
			result := nU.alreadyHasLabelSet(node)
			assert.Equal(t, tc.expectedResult, result)

		})
	}
}
