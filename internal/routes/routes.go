package routes

import (
	"test_effective_mobile_task/internal/handler"
	"github.com/gin-gonic/gin"
)


func SubscriptionRoutes(router *gin.Engine, handler *handler.SubscriptionHandler) {
	subscriptionGroup := router.Group("/subscriptions")
	{
		subscriptionGroup.POST("/create", handler.CreateSubscriptionHandler)
		subscriptionGroup.GET("/get", handler.GetSubscriptionHandler)
		subscriptionGroup.DELETE("/delete", handler.DeleteSubscriptionHandler)
		subscriptionGroup.PUT("/updateplan", handler.UpdateSubscriptionPlanHandler)
		subscriptionGroup.GET("/list", handler.GetListOfSubscriptionsHandler)
	}
}