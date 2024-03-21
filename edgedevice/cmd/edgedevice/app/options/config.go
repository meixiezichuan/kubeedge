package options

import (
	"net"
	"net/url"
	"os"
	"path"

	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
	"github.com/kubeedge/kubeedge/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

// DataBase indicates the database info
type DataBase struct {
	// DriverName indicates database driver name
	// default "sqlite3"
	DriverName string `json:"driverName,omitempty"`
	// AliasName indicates alias name
	// default "default"
	AliasName string `json:"aliasName,omitempty"`
	// DataSource indicates the data source path
	// default "/var/lib/kubeedge/edgedevice.db"
	DataSource string `json:"dataSource,omitempty"`
}

// EdgeDeviceConfig indicates the EdgeCore config which read from EdgeCore config file
type EdgeDeviceConfig struct {
	metav1.TypeMeta
	// DataBase indicates database info
	// +Required
	DataBase *DataBase `json:"database,omitempty"`
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

// NewDefaultEdgeCoreConfig returns a full EdgeCoreConfig object
func NewDefaultEdgeDeviceConfig() (config *EdgeDeviceConfig) {
	hostnameOverride := util.GetHostname()
	localIP, _ := util.GetLocalIP(hostnameOverride)

	config = &EdgeDeviceConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind,
			APIVersion: path.Join(GroupName, APIVersion),
		},
		DataBase: &DataBase{
			DriverName: DataBaseDriverName,
			AliasName:  DataBaseAliasName,
			DataSource: DataBaseDataSource,
		},
		Modules: &v1alpha2.Modules{
			EdgeHub: &v1alpha2.EdgeHub{
				Enable:            true,
				Heartbeat:         15,
				MessageQPS:        constants.DefaultQPS,
				MessageBurst:      constants.DefaultBurst,
				ProjectID:         "e632aba927ea4ac2b575ec1603d56f10",
				TLSCAFile:         constants.DefaultCAFile,
				TLSCertFile:       constants.DefaultCertFile,
				TLSPrivateKeyFile: constants.DefaultKeyFile,
				Quic: &v1alpha2.EdgeHubQUIC{
					Enable:           false,
					HandshakeTimeout: 30,
					ReadDeadline:     15,
					Server:           net.JoinHostPort(localIP, "10001"),
					WriteDeadline:    15,
				},
				WebSocket: &v1alpha2.EdgeHubWebSocket{
					Enable:           true,
					HandshakeTimeout: 30,
					ReadDeadline:     15,
					Server:           net.JoinHostPort(localIP, "10000"),
					WriteDeadline:    15,
				},
				HTTPServer: (&url.URL{
					Scheme: "https",
					Host:   net.JoinHostPort(localIP, "10002"),
				}).String(),
				Token:              "",
				RotateCertificates: true,
			},
			EventBus: &v1alpha2.EventBus{
				Enable:               true,
				MqttQOS:              0,
				MqttRetain:           false,
				MqttSessionQueueSize: 100,
				MqttServerExternal:   "tcp://127.0.0.1:1883",
				MqttServerInternal:   "tcp://127.0.0.1:1884",
				MqttSubClientID:      "",
				MqttPubClientID:      "",
				MqttUsername:         "",
				MqttPassword:         "",
				MqttMode:             v1alpha2.MqttModeExternal,
				TLS: &v1alpha2.EventBusTLS{
					Enable:                false,
					TLSMqttCAFile:         constants.DefaultMqttCAFile,
					TLSMqttCertFile:       constants.DefaultMqttCertFile,
					TLSMqttPrivateKeyFile: constants.DefaultMqttKeyFile,
				},
			},
			DeviceTwin: &v1alpha2.DeviceTwin{
				Enable:      true,
				DMISockPath: constants.DefaultDMISockPath,
			},
		},
	}
	return
}

// NewMinEdgeCoreConfig returns a common EdgeCoreConfig object
func NewMinEdgeDeviceConfig() (config *EdgeDeviceConfig) {
	hostnameOverride := util.GetHostname()
	localIP, _ := util.GetLocalIP(hostnameOverride)

	config = &EdgeDeviceConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind,
			APIVersion: path.Join(GroupName, APIVersion),
		},
		DataBase: &DataBase{
			DataSource: DataBaseDataSource,
		},
		Modules: &v1alpha2.Modules{
			DeviceTwin: &v1alpha2.DeviceTwin{
				DMISockPath: constants.DefaultDMISockPath,
			},
			EdgeHub: &v1alpha2.EdgeHub{
				Heartbeat:         15,
				TLSCAFile:         constants.DefaultCAFile,
				TLSCertFile:       constants.DefaultCertFile,
				TLSPrivateKeyFile: constants.DefaultKeyFile,
				WebSocket: &v1alpha2.EdgeHubWebSocket{
					Enable:           true,
					HandshakeTimeout: 30,
					ReadDeadline:     15,
					Server:           net.JoinHostPort(localIP, "10000"),
					WriteDeadline:    15,
				},
				HTTPServer: (&url.URL{
					Scheme: "https",
					Host:   net.JoinHostPort(localIP, "10002"),
				}).String(),
				Token: "",
			},
			EventBus: &v1alpha2.EventBus{
				MqttQOS:            0,
				MqttRetain:         false,
				MqttServerExternal: "tcp://127.0.0.1:1883",
				MqttServerInternal: "tcp://127.0.0.1:1884",
				MqttSubClientID:    "",
				MqttPubClientID:    "",
				MqttUsername:       "",
				MqttPassword:       "",
				MqttMode:           v1alpha2.MqttModeExternal,
			},
		},
	}
	return
}
