package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubeA/service"
	"net/http"
)

var Pod pod

type pod struct {
}

func (p *pod) GetPod(r *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
		Cluster    string `form:"cluster"`
	})
	if err := r.Bind(params); err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		r.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	podResq, err := service.Pod.GetPods(client, params.Namespace, params.FilterName, params.Limit, params.Page)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	r.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": podResq,
	})
}

func (p *pod) GetPodLog(cxt *gin.Context) {
	params := new(struct {
		Namespace     string `form:"namespace"`
		PodName       string `form:"pod_name"`
		Cluster       string `form:"cluster"`
		ContainerName string `form:"container_name"`
	})
	if err := cxt.Bind(params); err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("pod：%s 解析参数失败. \n", params.PodName),
			"data": nil,
		})
		return
	}
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  fmt.Sprintf("pod：%s 获取k8s client失败. \n", params.PodName),
			"data": nil,
		})
		return
	}
	podLog, err := service.Pod.GetPodLog(client, params.Namespace, params.PodName, params.ContainerName)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"msg":  fmt.Sprintf("pod：%s 获取PodLog失败. \n", params.PodName),
			"data": nil,
		})
	}
	cxt.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("pod：%s 获取PodLog成功. \n", params.PodName),
		"data": podLog,
	})
}

func (p *pod) GetPodDetail(ctx *gin.Context) {

}
func (p *pod) DeletePod(ctx *gin.Context) {

}
func (p *pod) UpdatePod(ctx *gin.Context) {

}
func (p *pod) GetPodContainer(ctx *gin.Context) {

}