package handler

import (
	"net/http"
	"strconv"
	"strings"
	"test_effective_mobile_task/internal/models"
	"test_effective_mobile_task/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	Service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		Service: service,
	}
}

func (h *SubscriptionHandler) CreateSubscriptionHandler(ctx *gin.Context) {
	var req service.CreateSubscriptionRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload	"})
		return
	}
	subscription, err := h.Service.CreateSubscription(ctx.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Subscription already exists for this user and service"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Subscription created successfully",
		 "subscription": subscription,
		})
}

func (h *SubscriptionHandler) GetSubscriptionHandler(ctx *gin.Context) {
	userIDstr := ctx.Query("user_id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	serviceName := ctx.Query("service_name")
	subscription, err := h.Service.GetSubscriptionByUserIDAndServiceName(ctx.Request.Context(), userID, serviceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

func (h *SubscriptionHandler) DeleteSubscriptionHandler(ctx *gin.Context) {
    subscriptionIDstr := ctx.Query("subscription_id")
    if subscriptionIDstr == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
        return
    }
    subscriptionID, err := uuid.Parse(subscriptionIDstr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
        return
    }
    err = h.Service.DeleteSubscriptionByID(ctx.Request.Context(), subscriptionID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription", "details": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

func (h *SubscriptionHandler) UpdateSubscriptionPlanHandler(ctx *gin.Context) {
	userIDstr := ctx.Query("user_id")
	if userIDstr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	serviceName := ctx.Query("service_name")
	if serviceName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Service name is required"})
		return
	}
	newPlanStr := ctx.Query("new_plan")
	if newPlanStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New plan is required"})
		return
	}
	newPlan := models.Plan(newPlanStr)
	if newPlan != models.Monthly && newPlan != models.HalfYearly && newPlan != models.Yearly {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan format"})
    return
	}	
	pricestr := ctx.Query("price")
	if pricestr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price of plan"})
		return
	}
	price, err := strconv.Atoi(pricestr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot convert the price"})
		return
	}
	err = h.Service.UpdateSubscriptionPlan(ctx.Request.Context(), userID, serviceName, newPlan, price)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription plan", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Subscription plan updated successfully",
	})
}


func (h *SubscriptionHandler) GetListOfSubscriptionsHandler(ctx *gin.Context) {
	userIDstr := ctx.Query("user_id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	subs, err := h.Service.GetListOfSubscriptionsByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get list of subscriptions", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscriptions": subs,
	})
}