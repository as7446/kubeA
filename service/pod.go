package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kubeA/config"
	"time"
)

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

var Pod pod

type pod struct {
}

func (p *pod) toCell(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for index := range std {
		cells[index] = podCell(std[index])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for index := range cells {
		pods[index] = corev1.Pod(cells[index].(podCell))
	}
	return pods
}

type PodsResq struct {
	Items []corev1.Pod
	Total int
}

// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(client *kubernetes.Clientset, namespace string, filterName string, limit int, page int) (podResq *PodsResq, err error) {
	podsList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	selectTableData := dataSelector{
		GenericDataList: p.toCell(podsList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := selectTableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()
	pods := p.fromCells(data.GenericDataList)
	podResq = &PodsResq{}
	podResq.Items = pods
	podResq.Total = total
	return podResq, nil
}

// 获取Pod详情
func (p *pod) GetPodDetail(client *kubernetes.Clientset, namespace string, podName string) (pod *corev1.Pod, err error) {
	pod, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pod: %s 详情失败: " + err.Error()))
	}
	return pod, nil
}

// 删除Pod
func (p *pod) DeletePod(client *kubernetes.Clientset, namespace string, podName string) (err error) {
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("删除Pod: %s失败,", podName))
	}
	return
}

// 更新Pod
func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace string, podName string, content string) (err error) {
	pod := &corev1.Pod{}
	err = json.Unmarshal([]byte(content), &pod)
	if err != nil {
		return errors.New(fmt.Sprintf("Pod: %s 反序列化失败: ", podName))
	}
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("Pod: %s 更新失败", podName))
	}
	return nil
}

// 获取Pod中容器名
func (p *pod) GetPodContainer(client *kubernetes.Clientset, namespace string, podName string) (containers []string, err error) {
	pod, err := p.GetPodDetail(client, namespace, podName)
	if err != nil {
		return
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return
}

// 获取容器日志
func (p *pod) GetPodLog(clinet *kubernetes.Clientset, namespace string, podName string, containerName string) (log string, err error) {
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	req := clinet.CoreV1().Pods(namespace).GetLogs(podName, option)
	//返回io.ReadCloser 类似 reponse.body
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", errors.New(fmt.Sprintf("获取PodLog失败:%s \n" + err.Error()))
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.New(fmt.Sprintf("读取PodLog失败:%s \n" + err.Error()))
	}
	return buf.String(), nil
}
