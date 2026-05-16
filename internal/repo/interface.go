package repo

import (
	"context"
	"test_effective_mobile_task/internal/models"
	"time"
	"github.com/google/uuid"
)

type SubscriptionRepoInterface interface {
	CreateSubscriptionRepo(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	GetSubscriptionByUserIDAndServiceNameRepo(ctx context.Context, userID uuid.UUID, serviceName string) (*models.Subscription, error)
	DeleteSubscriptionByIDRepo(ctx context.Context, subscriptionID uuid.UUID, serviceName string) error
	UpdateSubscriptionPlanRepo(ctx context.Context, userID uuid.UUID, serviceName string, newStartDate time.Time, newEndDate time.Time, price int) error
	GetListSubscriptionsByUserIDRepo(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
}