package repository

import (
	model "dailyworkerroster/model"
	"database/sql"
	"strings"
)

type ShiftRepoItf interface {
	CreateShift(shift *model.Shift) (int64, error)
	GetShiftByID(id int64) (*model.Shift, error)
	GetShiftsByIDs(ids []int64) ([]*model.Shift, error)
	UpdateShiftByID(shift *model.Shift) error
	DeleteShiftByID(id int64) error
	GetListShifts(queryParam model.ShiftListQuery) ([]*model.Shift, error)
}

type ShiftRepository struct {
	DB *sql.DB
}

func NewShiftRepository(db *sql.DB) ShiftRepoItf {
	return &ShiftRepository{
		DB: db,
	}
}

// CreateShift inserts a new shift into the database
func (r *ShiftRepository) CreateShift(shift *model.Shift) (int64, error) {
	query := `
        INSERT INTO shift (date, start_time, end_time, role_assignment, location, isAvailable, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
    `
	result, err := r.DB.Exec(query, shift.Date, shift.StartTime, shift.EndTime, shift.RoleAssignment, shift.Location, shift.IsAvailable)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetShiftByID retrieves a shift by its ID
func (r *ShiftRepository) GetShiftByID(id int64) (*model.Shift, error) {
	query := `
        SELECT id, date, start_time, end_time, role_assignment, location, isAvailable, created_at, updated_at
        FROM shift WHERE id = ?
    `
	var shift model.Shift
	err := r.DB.QueryRow(query, id).Scan(
		&shift.ID, &shift.Date, &shift.StartTime, &shift.EndTime,
		&shift.RoleAssignment, &shift.Location, &shift.IsAvailable,
		&shift.CreatedAt, &shift.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// GetShiftsByIDs retrieves multiple shifts by a list of IDs
func (r *ShiftRepository) GetShiftsByIDs(ids []int64) ([]*model.Shift, error) {
	if len(ids) == 0 {
		return []*model.Shift{}, nil
	}

	// Build placeholders (?, ?, ...) for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
        SELECT id, date, start_time, end_time, role_assignment, location, isAvailable, created_at, updated_at
        FROM shift
        WHERE id IN (` + strings.Join(placeholders, ",") + `)
    `

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []*model.Shift
	for rows.Next() {
		var shift model.Shift
		err := rows.Scan(
			&shift.ID, &shift.Date, &shift.StartTime, &shift.EndTime,
			&shift.RoleAssignment, &shift.Location, &shift.IsAvailable,
			&shift.CreatedAt, &shift.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, &shift)
	}
	return shifts, nil
}

// UpdateShift updates an existing shift
func (r *ShiftRepository) UpdateShiftByID(shift *model.Shift) error {
	query := `
        UPDATE shift SET date=?, start_time=?, end_time=?, role_assignment=?, location=?, isAvailable=?, updated_at=NOW()
        WHERE id=?
    `
	_, err := r.DB.Exec(query, shift.Date, shift.StartTime, shift.EndTime, shift.RoleAssignment, shift.Location, shift.IsAvailable, shift.ID)
	return err
}

// DeleteShift deletes a shift by its ID
func (r *ShiftRepository) DeleteShiftByID(id int64) error {
	query := `DELETE FROM shift WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *ShiftRepository) GetListShifts(
	queryParam model.ShiftListQuery,
) ([]*model.Shift, error) {
	query := `
        SELECT id, date, start_time, end_time, role_assignment, location, is_available, created_at, updated_at
        FROM shift
        WHERE 1=1
    `
	args := []interface{}{}

	if queryParam.RoleAssignment != "" {
		query += " AND role_assignment = ?"
		args = append(args, queryParam.RoleAssignment)
	}
	if queryParam.Location != "" {
		query += " AND location = ?"
		args = append(args, queryParam.Location)
	}
	if queryParam.IsAvailable != nil {
		query += " AND is_available = ?"
		args = append(args, *queryParam.IsAvailable)
	}
	if queryParam.Date != "" {
		query += " AND date = ?"
		args = append(args, queryParam.Date)
	}

	query += " ORDER BY date, start_time LIMIT ? OFFSET ?"
	args = append(args, queryParam.Limit, queryParam.Offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []*model.Shift
	for rows.Next() {
		var shift model.Shift
		err := rows.Scan(
			&shift.ID, &shift.Date, &shift.StartTime, &shift.EndTime,
			&shift.RoleAssignment, &shift.Location, &shift.IsAvailable,
			&shift.CreatedAt, &shift.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, &shift)
	}
	return shifts, nil
}
