package routes

import (
	"test_effective_mobile_task/internal/handler"
	"github.com/gin-gonic/gin"
	_ "test_effective_mobile_task/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SubscriptionRoutes(router *gin.Engine, handler *handler.SubscriptionHandler) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	subscriptionGroup := router.Group("/subscriptions")
	{
		subscriptionGroup.POST("/create", handler.CreateSubscriptionHandler)
		subscriptionGroup.GET("/get", handler.GetSubscriptionHandler)
		subscriptionGroup.DELETE("/delete", handler.DeleteSubscriptionHandler)
		subscriptionGroup.PUT("/updateplan", handler.UpdateSubscriptionPlanHandler)
		subscriptionGroup.GET("/list", handler.GetListOfSubscriptionsHandler)
		subscriptionGroup.GET("/total", handler.GetTotalPriceHandler)
	
	}
}