package config

import (
	"sync"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
)

var Config Configure
var once sync.Once

type Configure struct {
	v1alpha2.EdgeMaster
	NodeName string
}

func InitConfigure(em *v1alpha2.EdgeMaster, nodeName string) {
	once.Do(func() {
		Config = Configure{
			EdgeMaster: *em,
			NodeName:   nodeName,
		}
	})
}
