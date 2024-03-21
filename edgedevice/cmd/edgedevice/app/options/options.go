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
	"path"

	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	cliflag "k8s.io/component-base/cli/flag"
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
