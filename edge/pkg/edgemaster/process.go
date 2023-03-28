package edgemaster

import (
	"context"
	"encoding/json"
	"fmt"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	"github.com/kubeedge/kubeedge/edge/pkg/common/util"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (em *EdgeMaster) process(msg *model.Message) {
	_, resourceType, _, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Warningf("EdgeMaster parse message: %s resource type with error, message resource: %s, err: %v", msg.GetID(), msg.GetResource(), err)
		return
	}

	switch resourceType {
	case model.ResourceTypePod:
		// if resource type is pod, remove node
		err = em.processPodMsg(msg)
	case model.ResourceTypeConfigmap:
		err = em.processConfigMapMsg(msg)
	case model.ResourceTypeSecret:
		err = em.processSecretMsg(msg)
	case model.ResourceTypeNode:
		//TODO
	default:
		beehiveContext.SendToGroup(modules.MetaGroup, *msg)
	}
	if err != nil {
		klog.Warningf("process message: %s resource type with error, message resource: %s, err: %v", msg.GetID(), msg.GetResource(), err)
	}
}

func (em *EdgeMaster) processPodMsg(msg *model.Message) error {
	namespace, _, name, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Warningf("EdgeMaster message: %s process failure, get resource name failed with error: %v", msg.GetID(), err)
		return err
	}
	switch msg.GetOperation() {
	case model.InsertOperation:
		// TODO dispatch-strategy 最低的优先级
		newBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var pod corev1.Pod
		err = json.Unmarshal(newBytes, &pod)
		if err != nil {
			klog.Errorf("message %s content unmarshal to pod with error : %v", msg.GetID(), err)
			return err
		}
		// 清空NodeName可以理解，为什么要清空ResourceVersion？
		// remove NodeName of pod
		pod.Spec.NodeName = ""
		// remove resourceVersion of pod
		pod.ResourceVersion = ""
		_, err = em.clusterClient.CoreV1().Pods(namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("message %s send to edge cluster with error : %v", msg.GetID(), err)
			klog.Errorf("EdgeMaster create pod %v with error : %v", pod, msg.GetID())
			return err
		}
	case model.DeleteOperation:
		err = em.clusterClient.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			klog.Errorf("delete pod message %s send to edge cluster with error : %v", msg.GetID(), err)
			return err
		}
	case model.UpdateOperation:
		updateBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var pod corev1.Pod
		err = json.Unmarshal(updateBytes, &pod)
		// remove NodeName of pod
		pod.Spec.NodeName = ""
		if err != nil {
			klog.Errorf("message %s content unmarshal to pod with error : %v", msg.GetID(), err)
			return err
		}
		_, err = em.clusterClient.CoreV1().Pods(namespace).Update(context.TODO(), &pod, metav1.UpdateOptions{})
		//TODO patch action
	}
	return err
}

func (em *EdgeMaster) processConfigMapMsg(msg *model.Message) error {
	namespace, _, name, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Warningf("EdgeMaster message: %s process failure, get resource name failed with error: %v", msg.GetID(), err)
		return err
	}
	switch msg.GetOperation() {
	case model.InsertOperation:
		newBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var configMap corev1.ConfigMap
		err = json.Unmarshal(newBytes, &configMap)
		if err != nil {
			klog.Errorf("message %s content unmarshal to configMap with error : %v", msg.GetID(), err)
			return err
		}
		// 清空NodeName可以理解，为什么要清空ResourceVersion？
		// remove resourceVersion of configMap
		configMap.ResourceVersion = ""
		_, err = em.clusterClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("message %s send to edge cluster with error : %v", msg.GetID(), err)
			klog.Errorf("EdgeMaster create configMap %v with error : %v", configMap, msg.GetID())
			return err
		}
	case model.DeleteOperation:
		err = em.clusterClient.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			klog.Errorf("delete configMap message %s send to edge cluster with error : %v", msg.GetID(), err)
			return err
		}
	case model.UpdateOperation:
		updateBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var configMap corev1.ConfigMap
		err = json.Unmarshal(updateBytes, &configMap)
		if err != nil {
			klog.Errorf("message %s content unmarshal to configMap with error : %v", msg.GetID(), err)
			return err
		}
		_, err = em.clusterClient.CoreV1().ConfigMaps(namespace).Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	}
	return err
}

func (em *EdgeMaster) processSecretMsg(msg *model.Message) error {
	namespace, _, name, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Warningf("EdgeMaster message: %s process failure, get resource name failed with error: %v", msg.GetID(), err)
		return err
	}
	switch msg.GetOperation() {
	case model.InsertOperation:
		newBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var secret corev1.Secret
		err = json.Unmarshal(newBytes, &secret)
		if err != nil {
			klog.Errorf("message %s content unmarshal to secret with error : %v", msg.GetID(), err)
			return err
		}
		// remove resourceVersion of secret
		secret.ResourceVersion = ""
		_, err = em.clusterClient.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("message %s send to edge cluster with error : %v", msg.GetID(), err)
			klog.Errorf("EdgeMaster create secret %v with error : %v", secret, msg.GetID())
			return err
		}
	case model.DeleteOperation:
		err = em.clusterClient.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			klog.Errorf("delete secret message %s send to edge cluster with error : %v", msg.GetID(), err)
			return err
		}
	case model.UpdateOperation:
		updateBytes, err := msg.GetContentData()
		if err != nil {
			klog.Warningf("message: %s process failure, get data failed with error: %v", msg.GetID(), err)
			return err
		}
		var secret corev1.Secret
		err = json.Unmarshal(updateBytes, &secret)
		if err != nil {
			klog.Errorf("message %s content unmarshal to secret with error : %v", msg.GetID(), err)
			return err
		}
		_, err = em.clusterClient.CoreV1().Secrets(namespace).Update(context.TODO(), &secret, metav1.UpdateOptions{})
	}
	return err
}

func (em *EdgeMaster) podMonitor() {

	// 设置 Pod Informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	sharedInformers := informers.NewSharedInformerFactory(em.clusterClient, time.Minute*10)

	podInformer := sharedInformers.Core().V1().Pods().Informer()
	//TODO 封装message
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod added: %s\n", pod.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypePod, model.InsertOperation)
			info.FillBody(&pod)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)
			fmt.Printf("Pod updated: %s -> %s\n", oldPod.Name, newPod.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypePod, model.UpdateOperation)
			info.FillBody(&newPod)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod deleted: %s\n", pod.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypePod, model.DeleteOperation)
			info.FillBody(&pod)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
	})

	// 启动 Informer
	sharedInformers.Start(stopCh)

	// 运行直到中断
	runtime.HandleCrash()
	<-stopCh
}

func (em *EdgeMaster) configMapMonitor() {
	stopCh := make(chan struct{})
	defer close(stopCh)
	sharedInformers := informers.NewSharedInformerFactory(em.clusterClient, time.Minute*10)

	cmInformer := sharedInformers.Core().V1().ConfigMaps().Informer()
	cmInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap added: %s\n", cm.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeConfigmap, model.InsertOperation)
			info.FillBody(&cm)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldCM := oldObj.(*corev1.ConfigMap)
			newCM := newObj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap updated: %s -> %s\n", oldCM.Name, newCM.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeConfigmap, model.UpdateOperation)
			info.FillBody(&newCM)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("ConfigMap deleted: %s\n", cm.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeConfigmap, model.DeleteOperation)
			info.FillBody(&cm)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
	})

	// 启动 Informer
	sharedInformers.Start(stopCh)

	// 运行直到中断
	runtime.HandleCrash()
	<-stopCh
}

func (em *EdgeMaster) secretMonitor() {
	stopCh := make(chan struct{})
	defer close(stopCh)
	sharedInformers := informers.NewSharedInformerFactory(em.clusterClient, time.Minute*10)

	cmInformer := sharedInformers.Core().V1().Secrets().Informer()
	cmInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			secret := obj.(*corev1.Secret)
			fmt.Printf("Secret added: %s\n", secret.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeSecret, model.InsertOperation)
			info.FillBody(&secret)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldSecret := oldObj.(*corev1.Secret)
			newSecret := newObj.(*corev1.Secret)
			fmt.Printf("Secret updated: %s -> %s\n", oldSecret.Name, newSecret.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeSecret, model.UpdateOperation)
			info.FillBody(&newSecret)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
		DeleteFunc: func(obj interface{}) {
			secret := obj.(*corev1.ConfigMap)
			fmt.Printf("Secret deleted: %s\n", secret.Name)
			info := model.NewMessage("").BuildRouter(em.Name(), em.Group(), "default/"+model.ResourceTypeSecret, model.DeleteOperation)
			info.FillBody(&secret)
			beehiveContext.SendToGroup(modules.HubGroup, *info)
		},
	})

	// 启动 Informer
	sharedInformers.Start(stopCh)

	// 运行直到中断
	runtime.HandleCrash()
	<-stopCh
}
