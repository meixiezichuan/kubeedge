/*
Copyright 2019 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"fmt"
	"github.com/kubeedge/kubeedge/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	cliflag "k8s.io/component-base/cli/flag"

	"net"
	"net/url"
	"path"

	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
	"github.com/kubeedge/kubeedge/pkg/util/validation"
)

type EdgeDeviceOptions struct {
	ConfigFile string
}

type MqttMode int

const (
	GroupName  = "edgedevice.config.fdse.io"
	APIVersion = "v1alpha2"
	Kind       = "EdgeDevice"
)

const (
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	// DataSource
	DataBaseDataSource = "/var/lib/kubeedge/edgedevice.db"
)

var edgeDeviceOptions *EdgeDeviceOptions
var edgeDeviceConfig *EdgeDeviceConfig

func NewEdgeDeviceOptions() *EdgeDeviceOptions {
	edgeDeviceOptions = &EdgeDeviceOptions{
		ConfigFile: path.Join(constants.DefaultConfigDir, "edgedevice.yaml"),
	}
	return edgeDeviceOptions
}

func (o *EdgeDeviceOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("global")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file. Flags override values in this file.")
	return
}

func (o *EdgeDeviceOptions) Validate() []error {
	var errs []error
	if !validation.FileIsExist(o.ConfigFile) {
		errs = append(errs, field.Required(field.NewPath("config"),
			fmt.Sprintf("config file %v not exist. For the configuration file format, please refer to --minconfig and --defaultconfig command", o.ConfigFile)))
	}
	return errs
}

func (o *EdgeDeviceOptions) Config() (*EdgeDeviceConfig, error) {
	edgeDeviceConfig = NewDefaultEdgeDeviceConfig()
	if err := edgeDeviceConfig.Parse(o.ConfigFile); err != nil {
		return nil, err
	}

	return edgeDeviceConfig, nil
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
