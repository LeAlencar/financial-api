package repositories

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/models"
)

type ValidationLogRepository struct {
	session *gocql.Session
}

func NewValidationLogRepository(session *gocql.Session) *ValidationLogRepository {
	return &ValidationLogRepository{
		session: session,
	}
}

// SaveValidationLog saves a validation log entry to Cassandra
func (r *ValidationLogRepository) SaveValidationLog(ctx context.Context, log *models.ValidationLog) error {
	if log.ID.Version() == 0 {
		log.ID = gocql.TimeUUID()
	}

	query := `INSERT INTO validation_logs (
		id, event_type, action, event_id, user_id, status, message, 
		details, raw_payload, processed_at, source
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		log.ID,
		string(log.EventType),
		string(log.Action),
		log.EventID,
		log.UserID,
		string(log.Status),
		log.Message,
		log.Details,
		log.RawPayload,
		log.ProcessedAt,
		log.Source,
	).WithContext(ctx).Exec()
}

// GetValidationLogsByEventType retrieves validation logs by event type
func (r *ValidationLogRepository) GetValidationLogsByEventType(ctx context.Context, eventType models.EventType, limit int) ([]*models.ValidationLog, error) {
	query := `SELECT id, event_type, action, event_id, user_id, status, message, 
			  details, raw_payload, processed_at, source 
			  FROM validation_logs WHERE event_type = ? LIMIT ?`

	iter := r.session.Query(query, string(eventType), limit).WithContext(ctx).Iter()
	defer iter.Close()

	var logs []*models.ValidationLog
	for {
		var log models.ValidationLog
		var eventTypeStr, actionStr, statusStr string

		if !iter.Scan(
			&log.ID,
			&eventTypeStr,
			&actionStr,
			&log.EventID,
			&log.UserID,
			&statusStr,
			&log.Message,
			&log.Details,
			&log.RawPayload,
			&log.ProcessedAt,
			&log.Source,
		) {
			break
		}

		log.EventType = models.EventType(eventTypeStr)
		log.Action = models.ActionType(actionStr)
		log.Status = models.ValidationStatus(statusStr)

		logs = append(logs, &log)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetValidationLogsByUserID retrieves validation logs for a specific user
func (r *ValidationLogRepository) GetValidationLogsByUserID(ctx context.Context, userID string, limit int) ([]*models.ValidationLog, error) {
	query := `SELECT id, event_type, action, event_id, user_id, status, message, 
			  details, raw_payload, processed_at, source 
			  FROM validation_logs_by_user WHERE user_id = ? LIMIT ?`

	iter := r.session.Query(query, userID, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var logs []*models.ValidationLog
	for {
		var log models.ValidationLog
		var eventTypeStr, actionStr, statusStr string

		if !iter.Scan(
			&log.ID,
			&eventTypeStr,
			&actionStr,
			&log.EventID,
			&log.UserID,
			&statusStr,
			&log.Message,
			&log.Details,
			&log.RawPayload,
			&log.ProcessedAt,
			&log.Source,
		) {
			break
		}

		log.EventType = models.EventType(eventTypeStr)
		log.Action = models.ActionType(actionStr)
		log.Status = models.ValidationStatus(statusStr)

		logs = append(logs, &log)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetValidationLogsByStatus retrieves validation logs by status
func (r *ValidationLogRepository) GetValidationLogsByStatus(ctx context.Context, status models.ValidationStatus, limit int) ([]*models.ValidationLog, error) {
	query := `SELECT id, event_type, action, event_id, user_id, status, message, 
			  details, raw_payload, processed_at, source 
			  FROM validation_logs_by_status WHERE status = ? LIMIT ?`

	iter := r.session.Query(query, string(status), limit).WithContext(ctx).Iter()
	defer iter.Close()

	var logs []*models.ValidationLog
	for {
		var log models.ValidationLog
		var eventTypeStr, actionStr, statusStr string

		if !iter.Scan(
			&log.ID,
			&eventTypeStr,
			&actionStr,
			&log.EventID,
			&log.UserID,
			&statusStr,
			&log.Message,
			&log.Details,
			&log.RawPayload,
			&log.ProcessedAt,
			&log.Source,
		) {
			break
		}

		log.EventType = models.EventType(eventTypeStr)
		log.Action = models.ActionType(actionStr)
		log.Status = models.ValidationStatus(statusStr)

		logs = append(logs, &log)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetValidationLogsByDateRange retrieves validation logs within a date range
func (r *ValidationLogRepository) GetValidationLogsByDateRange(ctx context.Context, start, end time.Time, limit int) ([]*models.ValidationLog, error) {
	query := `SELECT id, event_type, action, event_id, user_id, status, message, 
			  details, raw_payload, processed_at, source 
			  FROM validation_logs_by_date WHERE processed_at >= ? AND processed_at <= ? LIMIT ?`

	iter := r.session.Query(query, start, end, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var logs []*models.ValidationLog
	for {
		var log models.ValidationLog
		var eventTypeStr, actionStr, statusStr string

		if !iter.Scan(
			&log.ID,
			&eventTypeStr,
			&actionStr,
			&log.EventID,
			&log.UserID,
			&statusStr,
			&log.Message,
			&log.Details,
			&log.RawPayload,
			&log.ProcessedAt,
			&log.Source,
		) {
			break
		}

		log.EventType = models.EventType(eventTypeStr)
		log.Action = models.ActionType(actionStr)
		log.Status = models.ValidationStatus(statusStr)

		logs = append(logs, &log)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return logs, nil
}
