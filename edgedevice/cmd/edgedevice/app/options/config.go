package app

import (
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/yaml"
)

// EdgeDeviceConfig indicates the EdgeCore config which read from EdgeCore config file
type EdgeDeviceConfig struct {
	metav1.TypeMeta
	// Modules indicates EdgeCore modules config
	// +Required
	Modules *v1alpha2.Modules `json:"modules,omitempty"`
	// FeatureGates is a map of feature names to bools that enable or disable alpha/experimental features.
	FeatureGates map[string]bool `json:"featureGates,omitempty"`
}

func (c *EdgeDeviceConfig) Parse(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		klog.Errorf("Failed to read configfile %s: %v", filename, err)
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		klog.Errorf("Failed to unmarshal configfile %s: %v", filename, err)
		return err
	}
	return nil
}
