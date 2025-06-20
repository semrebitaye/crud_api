package repository

import (
	"context"
	"crud_api/internal/domain/models"
	appErrors "crud_api/internal/errors"
	contexthelper "crud_api/internal/utility/context_helper"

	"github.com/jackc/pgx/v5"
	"github.com/joomcode/errorx"
)

type PgxComplaintRepo struct {
	db *pgx.Conn
}

func NewPgxComplaintRepo(db *pgx.Conn) *PgxComplaintRepo {
	return &PgxComplaintRepo{db: db}
}

func (r *PgxComplaintRepo) CreateComplaint(ctx context.Context, c *models.Complaints) error {
	if contexthelper.IsAdmin(ctx) {
		return appErrors.ErrInvalidPayload.New("users only have permission to create complaint")
	}

	query := `INSERT INTO complaints (user_id, subject, message, status) VALUES($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(ctx, query, c.UserID, c.Subject, c.Message, c.Status).Scan(&c.ID)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to query the user")
	}

	return nil
}

func (r *PgxComplaintRepo) GetComplaintByRole(ctx context.Context, UserID int) ([]*models.Complaints, error) {
	if contexthelper.IsAdmin(ctx) {
		return nil, appErrors.ErrInvalidPayload.New("this is for users")
	}

	query := `SELECT * FROM complaints WHERE user_id=$1`
	rows, err := r.db.Query(ctx, query, UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, appErrors.ErrUserNotFound.Wrap(err, "Complaint not found")
		}
		return nil, appErrors.ErrDbFailure.New("quey failed")
	}
	defer rows.Close()

	var complaint []*models.Complaints
	for rows.Next() {
		var c models.Complaints
		err := rows.Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrUserNotFound.New("failed to scan complaint row")
		}
		complaint = append(complaint, &c)
	}

	return complaint, nil
}

func (r *PgxComplaintRepo) GetComplaintByID(ctx context.Context, complaintID int) (*models.Complaints, error) {
	var c models.Complaints
	query := `SELECT * FROM complaints WHERE complaint_id=$1`
	err := r.db.QueryRow(ctx, query, complaintID).Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)

	if err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.ErrUserNotFound.Wrap(err, "complaint not found")
		}
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return &c, nil
}

func (r *PgxComplaintRepo) UpdateComplaints(ctx context.Context, ComplaintId int, status string) error {
	query := `UPDATE complaints SET status=$1 WHERE id=$2`

	_, err := r.db.Exec(ctx, query, ComplaintId, status)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to update complaint")
	}

	return nil
}

func (r *PgxComplaintRepo) GetAllComplaintByRole(ctx context.Context) ([]*models.Complaints, error) {
	if !contexthelper.IsAdmin(ctx) {
		return nil, appErrors.ErrInvalidPayload.New("admin only can retrieve all data")
	}

	query := `SELECT * FROM complaints`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}
	defer rows.Close()

	var complaints []*models.Complaints
	for rows.Next() {
		var c models.Complaints
		err := rows.Scan(&c.ID, &c.UserID, &c.Subject, &c.Message, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrDbFailure.New("Failed to scan row")
		}
		complaints = append(complaints, &c)
	}

	return complaints, nil
}

// using complaintMessage table
func (r *PgxComplaintRepo) AddMessage(ctx context.Context, cm *models.ComplaintMessages) error {
	query := `INSERT INTO complaint_messages (complaint_id, sender_id, parent_id, message) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, cm.ComplaintID, cm.SenderID, cm.ParentID, cm.Message, cm.FileUrl).Scan(&cm.ID)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "query failed")
	}

	return nil
}

func (r *PgxComplaintRepo) GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error) {
	query := `SELECT * FROM complaint_message WHERE complaint_id=$1`
	rows, err := r.db.Query(ctx, query, complaintID)
	if err != nil {
		return nil, appErrors.ErrDbFailure.Wrap(err, "query failed")
	}
	defer rows.Close()

	var complaint_messages []*models.ComplaintMessages
	for rows.Next() {
		var cm models.ComplaintMessages
		err := rows.Scan(&cm.ID, &cm.ComplaintID, &cm.SenderID, &cm.ParentID, &cm.Message, &cm.FileUrl, &cm.CreatedAt)
		if err != nil {
			return nil, appErrors.ErrDbFailure.Wrap(err, "Failed to scan row of messages")
		}
		complaint_messages = append(complaint_messages, &cm)
	}

	return complaint_messages, nil
}
