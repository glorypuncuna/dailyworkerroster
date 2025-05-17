package repository

import (
	model "dailyworkerroster/model"
	"database/sql"
	"time"
)

type WorkerShiftRepoItf interface {
	CreateWorkerShift(ws *model.WorkerShift) (int64, error)
	GetWorkerShiftByID(id int64) (*model.WorkerShift, error)
	GetWorkerShiftListByFilter(userAccountID *int64, status *string) ([]model.WorkerShift, error)
	UpdatesWorkerShiftStatus(id int64, status string, approvedBy *int64) error
	DeleteWorkerShiftByID(id int64) error
	ListWorkerShiftsByUser(userID int64) ([]*model.WorkerShift, error)
	ListWorkerShiftsByShift(shiftID int64) ([]*model.WorkerShift, error)
	CheckWorkerShiftLimits(userAccountID int64, date string) (bool, int, error)
	GetWorkerShiftDetailListByFilter(queryParam *model.WorkerShiftDetailQuery) ([]model.WorkerShiftDetail, error)
}

type WorkerShiftRepository struct {
	DB *sql.DB
}

func NewWorkerShiftRepository(db *sql.DB) WorkerShiftRepoItf {
	return &WorkerShiftRepository{DB: db}
}

func (r *WorkerShiftRepository) CreateWorkerShift(ws *model.WorkerShift) (int64, error) {
	query := `
        INSERT INTO worker_shift (shift_id, user_account_id, approved_by, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
    `
	result, err := r.DB.Exec(query, ws.ShiftID, ws.UserAccountID, ws.ApprovedBy, ws.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *WorkerShiftRepository) GetWorkerShiftByID(id int64) (*model.WorkerShift, error) {
	query := `
        SELECT id, shift_id, user_account_id, approved_by, status, created_at, updated_at
        FROM worker_shift WHERE id = ?
    `
	var ws model.WorkerShift
	err := r.DB.QueryRow(query, id).Scan(
		&ws.ID, &ws.ShiftID, &ws.UserAccountID, &ws.ApprovedBy, &ws.Status, &ws.CreatedAt, &ws.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (r *WorkerShiftRepository) GetWorkerShiftListByFilter(userAccountID *int64, status *string) ([]model.WorkerShift, error) {
	query := `
        SELECT id, shift_id, user_account_id, approved_by, status, created_at, updated_at
        FROM worker_shift WHERE 1=1
    `
	args := []interface{}{}

	if userAccountID != nil {
		query += " AND user_account_id = ?"
		args = append(args, *userAccountID)
	}
	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY updated_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.WorkerShift
	for rows.Next() {
		var ws model.WorkerShift
		err := rows.Scan(
			&ws.ID, &ws.ShiftID, &ws.UserAccountID, &ws.ApprovedBy, &ws.Status, &ws.CreatedAt, &ws.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, ws)
	}
	return list, nil
}

func (r *WorkerShiftRepository) UpdatesWorkerShiftStatus(id int64, status string, approvedBy *int64) error {
	query := `
        UPDATE worker_shift
        SET status = ?, approved_by = ?, updated_at = NOW()
        WHERE id = ?
    `
	_, err := r.DB.Exec(query, status, approvedBy, id)
	return err
}

func (r *WorkerShiftRepository) DeleteWorkerShiftByID(id int64) error {
	query := `DELETE FROM worker_shift WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *WorkerShiftRepository) ListWorkerShiftsByUser(userID int64) ([]*model.WorkerShift, error) {
	query := `
        SELECT id, shift_id, user_account_id, approved_by, status, created_at, updated_at
        FROM worker_shift WHERE user_account_id = ?
    `
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.WorkerShift
	for rows.Next() {
		var ws model.WorkerShift
		err := rows.Scan(
			&ws.ID, &ws.ShiftID, &ws.UserAccountID, &ws.ApprovedBy, &ws.Status, &ws.CreatedAt, &ws.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &ws)
	}
	return list, nil
}

func (r *WorkerShiftRepository) ListWorkerShiftsByShift(shiftID int64) ([]*model.WorkerShift, error) {
	query := `
        SELECT id, shift_id, user_account_id, approved_by, status, created_at, updated_at
        FROM worker_shift WHERE shift_id = ?
    `
	rows, err := r.DB.Query(query, shiftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.WorkerShift
	for rows.Next() {
		var ws model.WorkerShift
		err := rows.Scan(
			&ws.ID, &ws.ShiftID, &ws.UserAccountID, &ws.ApprovedBy, &ws.Status, &ws.CreatedAt, &ws.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &ws)
	}
	return list, nil
}

func (r *WorkerShiftRepository) CheckWorkerShiftLimits(userAccountID int64, date string) (hasShiftOnDay bool, shiftsThisWeek int, err error) {
	query := `
        SELECT 
            COUNT(CASE WHEN s.date = ? THEN 1 END) AS shifts_on_day,
            COUNT(*) AS shifts_in_week
        FROM worker_shift ws
        JOIN shift s ON ws.shift_id = s.id
        WHERE ws.user_account_id = ?
            AND ws.status = 'APPROVED'
            AND s.date BETWEEN ? AND ?
    `

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, 0, err
	}
	weekday := int(parsedDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := parsedDate.AddDate(0, 0, -weekday+1)
	sunday := monday.AddDate(0, 0, 6)

	var shiftsOnDay, shiftsInWeek int
	err = r.DB.QueryRow(query, date, userAccountID, monday.Format("2006-01-02"), sunday.Format("2006-01-02")).Scan(&shiftsOnDay, &shiftsInWeek)
	if err != nil {
		return false, 0, err
	}
	return shiftsOnDay > 0, shiftsInWeek, nil
}

func (r *WorkerShiftRepository) GetWorkerShiftDetailListByFilter(queryParam *model.WorkerShiftDetailQuery) ([]model.WorkerShiftDetail, error) {
	query := `
        SELECT ws.id, ws.shift_id, ws.user_account_id, ws.approved_by, ws.status,
               s.date, s.start_time, s.end_time, s.role_assignment, s.location, s.isAvailable
        FROM worker_shift ws
        JOIN shift s ON ws.shift_id = s.id
        WHERE 1=1
    `
	args := []interface{}{}

	if queryParam.UserAccountID != nil {
		query += " AND ws.user_account_id = ?"
		args = append(args, *queryParam.UserAccountID)
	}
	if queryParam.Status != nil {
		query += " AND ws.status = ?"
		args = append(args, *queryParam.Status)
	}
	if queryParam.Role != nil {
		query += " AND s.role_assignment = ?"
		args = append(args, *queryParam.Role)
	}
	if queryParam.Location != nil {
		query += " AND s.location = ?"
		args = append(args, *queryParam.Location)
	}
	query += " ORDER BY ws.updated_at DESC"
	if queryParam.Limit != nil {
		query += " LIMIT ?"
		args = append(args, *queryParam.Limit)
	}
	if queryParam.Offset != nil {
		query += " OFFSET ?"
		args = append(args, *queryParam.Offset)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.WorkerShiftDetail
	for rows.Next() {
		var ws model.WorkerShiftDetail
		err := rows.Scan(
			&ws.ID, &ws.ShiftID, &ws.UserAccountID, &ws.ApprovedBy, &ws.Status,
			&ws.Date, &ws.StartTime, &ws.EndTime, &ws.RoleAssignment, &ws.Location, &ws.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, ws)
	}
	return list, nil
}
