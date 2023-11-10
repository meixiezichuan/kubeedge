package edgemaster

import (
	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	emconfig "github.com/kubeedge/kubeedge/edge/pkg/edgemaster/config"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

type EdgeMaster struct {
	clusterConfig string
	clusterClient kubernetes.Interface
	enable        bool
	nodeName      string
}

var _ core.Module = (*EdgeMaster)(nil)

func newEdgeMaster(enable bool, configFile string, nodeName string) *EdgeMaster {
	return &EdgeMaster{
		enable:        enable,
		clusterConfig: configFile,
		nodeName:      nodeName,
	}
}

// Register register EdgeMaster
func Register(em *v1alpha2.EdgeMaster, nodeName string) {
	emconfig.InitConfigure(em, nodeName)
	core.Register(newEdgeMaster(em.Enable, em.ClusterConfig, emconfig.Config.NodeName))
}

// Name returns the name of EdgeMaster module
func (em *EdgeMaster) Name() string {
	return modules.EdgeMasterModuleName
}

// Group returns EdgeMaster group
func (em *EdgeMaster) Group() string {
	return modules.MasterGroup
}

// Enable indicates whether this module is enabled
func (em *EdgeMaster) Enable() bool {
	return em.enable
}

// Start sets context and starts the controller
func (em *EdgeMaster) Start() {
	defer func() {
		if err := recover(); err != nil {
			klog.Errorf("EdgeMaster panic occurred:", err)
		}
	}()

	conf, er := clientcmd.BuildConfigFromFlags("", em.clusterConfig)
	if er != nil {
		klog.Error(er, "error on getting conf from clusterConfig")
		return
	}
	em.clusterClient = kubernetes.NewForConfigOrDie(conf)
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
			klog.V(2).Infof("EdgeMaster get a message %+v", msg)
			em.process(&msg)
		}
	}()

	klog.Infof("Start monitoring changes in edge cluster resources")

	time.Sleep(10 * time.Second)
	//监控资源状态并上报
	go em.podMonitor()
	go em.configMapMonitor()
	go em.secretMonitor()

}
