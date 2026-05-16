package service

import (
	"context"
	"fmt"
	"test_effective_mobile_task/internal/models"
	"test_effective_mobile_task/internal/repo"
	"time"
	"github.com/google/uuid"
)


type SubscriptionService struct {
	Repo *repo.SubscriptionRepo
}

func NewSubscriptionService(repo *repo.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{
		Repo: repo,
	}
}

type CreateSubscriptionRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	Plan        models.Plan `json:"plan"`
}

func (s *SubscriptionService) FormatDurationByPlan(ctx context.Context, start_date time.Time, plan models.Plan) (time.Time, error) {
	switch plan {
	case models.Monthly:
		return start_date.AddDate(0, 1, 0), nil 
	case models.HalfYearly:
		return start_date.AddDate(0, 6, 0), nil 
	case models.Yearly:
		return start_date.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid plan: %s", plan)
	}
  }
  func (s *SubscriptionService) GetPlanByStartAndEndDates(ctx context.Context, start_date time.Time, end_date time.Time) (models.Plan, error) {
	duration := end_date.Sub(start_date)
	switch {
	case duration >= time.Hour * 24 * 365:
		return models.Yearly, nil
	case duration >= 6*30*24*time.Hour:
		return models.HalfYearly, nil
	case duration >= 30*24*time.Hour:
		return models.Monthly, nil
	default:
		return "", fmt.Errorf("invalid subscription duration: %s", duration)
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*models.Subscription, error) {
    
	isExist , err := s.Repo.GetSubscriptionByUserIDAndServiceNameRepo(ctx, req.UserID, req.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("trouble checking of existing subscription: %w", err)
	}
	if isExist != nil {
		return nil, fmt.Errorf("subscription already exists for user: %s and service: %s", req.UserID, req.ServiceName)
	}
	start_date := time.Now()
	end_date, err := s.FormatDurationByPlan(ctx, start_date, req.Plan)
	if err != nil {
		return nil, err
	}
	sub := &models.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		Plan:        req.Plan,
		StartDate:   start_date,
		EndDate:     end_date,
	}
	return s.Repo.CreateSubscriptionRepo(ctx, sub)
}

func (s *SubscriptionService) GetSubscriptionByUserIDAndServiceName(ctx context.Context, userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	sub, err := s.Repo.GetSubscriptionByUserIDAndServiceNameRepo(ctx, userID, serviceName)
	if err != nil {
		return nil, fmt.Errorf("trouble getting subscription: %w", err)
	}
	if sub == nil {
		return nil, fmt.Errorf("subscription not found for user: %s and service: %s", userID, serviceName)
	}
	plan , err := s.GetPlanByStartAndEndDates(ctx, sub.StartDate, sub.EndDate)
	if err != nil {
		return nil, fmt.Errorf("trouble getting subscription plan: %w", err)
	}

	sub.Plan = plan

	return sub, nil
}

func (s *SubscriptionService) DeleteSubscriptionByID(ctx context.Context, subscriptionID uuid.UUID) error {
    err := s.Repo.DeleteSubscriptionByIDRepo(ctx, subscriptionID)
    if err != nil {
        return fmt.Errorf("trouble deleting subscription: %w", err)
    }
    return nil
}

func (s *SubscriptionService) UpdateSubscriptionPlan(ctx context.Context, userID uuid.UUID, serviceName string, newPlan models.Plan, newPrice int) error {
	newStartDate := time.Now()
	newEndDate, err := s.FormatDurationByPlan(ctx, newStartDate, newPlan)
	if err != nil {
		return fmt.Errorf("trouble formatting new subscription duration: %w", err)
	}
	err = s.Repo.UpdateSubscriptionPlanRepo(ctx, userID, serviceName, newStartDate, newEndDate, newPrice)
	if err != nil {
		return fmt.Errorf("trouble updating subscription plan: %w", err)
	}
	return nil
}
func (s *SubscriptionService) GetListOfSubscriptionsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	subs, err := s.Repo.GetListSubscriptionsByUserIDRepo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("trouble getting list of subscriptions: %w", err)
	}
	for _, sub := range subs {
		plan, err := s.GetPlanByStartAndEndDates(ctx, sub.StartDate, sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("trouble getting subscription plan: %w", err)
		}
		sub.Plan = plan
	}
	return subs, nil
}