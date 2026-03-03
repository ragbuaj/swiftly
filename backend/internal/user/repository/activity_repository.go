package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	pool *pgxpool.Pool
}

func NewActivityRepository(pool *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{pool: pool}
}

func (r *ActivityRepository) Log(userID, activityType, ip, userAgent string, metadata map[string]interface{}) error {
	id := uuid.New().String()
	metaJSON, _ := json.Marshal(metadata)

	query := `INSERT INTO user_activities (id, user_id, activity_type, metadata, ip_address, user_agent, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	
	_, err := r.pool.Exec(context.Background(), query, id, userID, activityType, metaJSON, ip, userAgent)
	return err
}
