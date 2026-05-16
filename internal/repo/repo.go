package repo

import (
	"context"
	"database/sql"
	"fmt"
	"test_effective_mobile_task/internal/models"
	"time"

	"github.com/google/uuid"
)

 type SubscriptionRepo struct {
	db *sql.DB
 }

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}

func (r *SubscriptionRepo) CreateSubscriptionRepo(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	err := r.db.QueryRowContext(ctx, "INSERT INTO subscriptions (user_id,service_name,price,start_date,end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		 sub.UserID,
	  	 sub.ServiceName,
		 sub.Price,
	     sub.StartDate,
		 sub.EndDate,
		).Scan(&sub.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	return sub, nil
}


func (r *SubscriptionRepo) GetSubscriptionByUserIDAndServiceNameRepo(ctx context.Context, userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE user_id = $1 AND service_name = $2", userID, serviceName).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
	)
	 if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
	return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return &sub, nil
}

func (r *SubscriptionRepo) DeleteSubscriptionByIDRepo(ctx context.Context, subscriptionID uuid.UUID) error {
    result, err := r.db.ExecContext(ctx,
        "DELETE FROM subscriptions WHERE id = $1",
        subscriptionID,
    )
    if err != nil {
        return fmt.Errorf("failed to delete subscription: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("subscription not found for id: %s", subscriptionID)
    }
    return nil
}

func (r *SubscriptionRepo) UpdateSubscriptionPlanRepo(ctx context.Context, userID uuid.UUID, serviceName string, startDate time.Time, endDate time.Time, price int) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE subscriptions SET start_date = $1, end_date = $2, price = $3 WHERE user_id = $4 AND service_name = $5",
		startDate, endDate, price, userID, serviceName,
	)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("subscription not found for user: %s and service: %s", userID, serviceName)
	}
	return nil
}

func (r *SubscriptionRepo) GetListSubscriptionsByUserIDRepo(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.ServiceName,
			&sub.Price,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}
