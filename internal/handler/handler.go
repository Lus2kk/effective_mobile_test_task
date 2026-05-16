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
// @Summary Создать подписку
// @Description Создает новую подписку для пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body service.CreateSubscriptionRequest true "Данные подписки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create [post]
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

// @Summary Получить подписку
// @Description Получает подписку по ID пользователя и названию сервиса
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Param service_name query string true "Название сервиса"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get [get]
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
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H {
			  "error" : "subscription not found", 
			})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param subscription_id query string true "UUID подписки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete [delete]
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

// @Summary Обновить план подписки
// @Description Обновляет план и цену подписки по ID пользователя и названию сервиса
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Param service_name query string true "Название сервиса"
// @Param new_plan query string true "Новый план (monthly/half_yearly/yearly)"
// @Param price query int true "Новая цена"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /updateplan [put]
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

// @Summary Получить список подписок
// @Description Получает все подписки пользователя
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /list [get]
func (h *SubscriptionHandler) GetListOfSubscriptionsHandler(ctx *gin.Context) {
	userIDstr := ctx.Query("user_id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	subs, err := h.Service.GetListOfSubscriptionsByUserID(ctx.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error" : "nothing was found",
			})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get list of subscriptions", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscriptions": subs,
	})
}

// @Summary Получить общую стоимость подписок
// @Description Считает суммарную стоимость всех подписок пользователя
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /total [get]
func (h *SubscriptionHandler) GetTotalPriceHandler(ctx *gin.Context) {
    userIDstr := ctx.Query("user_id")
    userID, err := uuid.Parse(userIDstr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    total, err := h.Service.GetTotalPriceOfSubscriptionsByUserID(ctx.Request.Context(), userID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total price", "details": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "total_price": total,
    })
}