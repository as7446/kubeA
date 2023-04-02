package service

import (
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubeA/config"
)

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset
	KubeConfMap map[string]string
}

func (k *k8s) GetClient(clusterName string) (client *kubernetes.Clientset, err error) {
	var ok bool
	if client, ok = k.ClientMap[clusterName]; ok {
		return client, err
	}
	return nil, err
}

func (k *k8s) Init() {
	mp := map[string]string{}
	k.ClientMap = map[string]*kubernetes.Clientset{}

	err := json.Unmarshal([]byte(config.Kubeconfigs), &mp)
	if err != nil {
		panic("k8s kubeconfigs 反序列化失败.")
	}
	k.KubeConfMap = mp
	for cluster, kubeconfig := range mp {
		cf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8sClient配置失败：%v\n", cluster, err))
		}
		clientset, err := kubernetes.NewForConfig(cf)
		if err != nil {
			panic(fmt.Sprintf("集群%s: 创建K8sClient失败：%v\n", cluster, err))
		}
		k.ClientMap[cluster] = clientset
		logger.Info("集群%s: 创建K8sClient成功 ", cluster)
	}
}
