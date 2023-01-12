package edgemaster

import (
	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	"github.com/kubeedge/kubeedge/edge/pkg/edgemaster/config"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type EdgeMaster struct {
	clusterConfig string
	clusterClient kubernetes.Interface
	enable        bool
}

var _ core.Module = (*EdgeMaster)(nil)

func newEdgeMaster(enable bool, configFile string) *EdgeMaster {
	return &EdgeMaster{
		enable:        enable,
		clusterConfig: configFile,
	}
}

// Register register EdgeMaster
func Register(em *v1alpha2.EdgeMaster, nodeName string) {
	config.InitConfigure(em, nodeName)
	core.Register(newEdgeMaster(em.Enable, em.ClusterConfig))
}

// Name returns the name of EdgeMaster module
func (eh *EdgeMaster) Name() string {
	return modules.EdgeMasterModuleName
}

// Group returns EdgeMaster group
func (eh *EdgeMaster) Group() string {
	return modules.MasterGroup
}

// Enable indicates whether this module is enabled
func (eh *EdgeMaster) Enable() bool {
	return eh.enable
}

// Start sets context and starts the controller
func (em *EdgeMaster) Start() {
	defer func() {
		if err := recover(); err != nil {
			klog.Errorf("EdgeMaster panic occurred:", err)
		}
	}()

	config, er := clientcmd.BuildConfigFromFlags("", em.clusterConfig)
	if er != nil {
		klog.Error(er, "error on getting config from clusterConfig")
		return
	}
	em.clusterClient = kubernetes.NewForConfigOrDie(config)
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				klog.Warning("EdgeMaster main loop stop")
				return
			default:
			}
			msg, err := beehiveContext.Receive(em.Name())
			if err != nil {
				klog.Errorf("get a message %+v: %v", msg, err)
				continue
			}
			klog.V(2).Infof("get a message %+v", msg)
			em.process(msg)
		}
	}()
}
