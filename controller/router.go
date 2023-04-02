package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/healthz", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg":  "ok",
			"data": nil,
		})
	})
	router.GET("/api/k8s/pods", Pod.GetPod)
	router.GET("/api/k8s/pod/log", Pod.GetPodLog)
	router.DELETE("/api/k8s/pod", Pod.DeletePod)
	router.GET("/api/k8s/pod/containers", Pod.GetPodContainer)
	router.GET("/api/k8s/pod", Pod.GetPodDetail)
	router.PUT("/api/k8s/pod", Pod.UpdatePod)
}
