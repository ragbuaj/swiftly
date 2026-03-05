package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type ActivityRepository struct {
	db DBTX
}

func NewActivityRepository(db DBTX) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Log(ctx context.Context, userID, activityType, ip, userAgent string, metadata map[string]interface{}) error {
	id := uuid.New().String()
	metaJSON, _ := json.Marshal(metadata)

	query := `INSERT INTO user_activities (id, user_id, activity_type, metadata, ip_address, user_agent, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	
	_, err := r.db.Exec(ctx, query, id, userID, activityType, metaJSON, ip, userAgent)
	return err
}
