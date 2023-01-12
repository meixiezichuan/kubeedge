package edgemaster

import (
	"context"
	"encoding/json"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/cloud/pkg/common/messagelayer"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (em *EdgeMaster) process(msg model.Message) {
	resourceType, err := messagelayer.GetResourceType(msg)
	if err != nil {
		klog.Warningf("parse message: %s resource type with error, message resource: %s, err: %v", msg.GetID(), msg.GetResource(), err)
		return
	}

	switch resourceType {
	case model.ResourceTypePod:
		// if resource type is pod, remove node
		err = em.processPodMsg(msg)
	case model.ResourceTypeNode:
		// todo
	default:
		beehiveContext.SendToGroup(modules.MetaGroup, msg)
	}
	if err != nil {
		klog.Warningf("process message: %s resource type with error, message resource: %s, err: %v", msg.GetID(), msg.GetResource(), err)
	}
}

func (em *EdgeMaster) processPodMsg(msg model.Message) error {
	namespace, err := messagelayer.GetNamespace(msg)
	if err != nil {
		klog.Warningf("message: %s process failure, get namespace failed with error: %v", msg.GetID(), err)
		return err
	}
	name, err := messagelayer.GetResourceName(msg)
	if err != nil {
		klog.Warningf("message: %s process failure, get resource name failed with error: %v", msg.GetID(), err)
		return err
	}
	switch msg.GetOperation() {
	case model.InsertOperation:
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
		// remove NodeName of pod
		pod.Spec.NodeName = ""
		_, err = em.clusterClient.CoreV1().Pods(namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("message %s send to edge cluster with error : %v", msg.GetID(), err)
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
	}
	return err
}
