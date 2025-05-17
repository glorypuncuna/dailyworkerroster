package model

import "time"

const (
	WORKER_SHIFT_PENDING  = "PENDING"
	WORKER_SHIFT_APPROVED = "APPROVED"
	WORKER_SHIFT_REJECTED = "REJECTED"
	WORKER_SHIFT_DONE     = "DONE"
	WORKER_SHIFT_EXPIRED  = "EXPIRED"

	MAXIMUM_WORKER_SHIFT_WEEK = 5
)

type WorkerShift struct {
	ID            int64     `json:"id"`
	ShiftID       int64     `json:"shift_id"`
	UserAccountID int64     `json:"user_account_id"`
	ApprovedBy    *int64    `json:"approved_by"` // nullable
	Status        string    `json:"status"`      // PENDING, APPROVED, REJECTED
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListShiftDetail struct {
	Name          string              `json:"name"`
	UserAccountID int64               `json:"user_account_id"`
	ShiftDetails  []WorkerShiftDetail `json:"shift_details"`
}

type WorkerShiftDetail struct {
	ID             int64  `json:"id"`
	ShiftID        int64  `json:"shift_id"`
	ApprovedBy     *int64 `json:"approved_by"` // nullable
	Status         string `json:"status"`      // PENDING, APPROVED, REJECTED
	Date           string `json:"date"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	RoleAssignment string `json:"role_assignment"`
	Location       string `json:"location"`
	IsAvailable    bool   `json:"isAvailable"`
	UserAccountID  int64  `json:"user_account_id"`
}

type WorkerShiftDetailQuery struct {
	UserAccountID *int64
	Status        *string
	Role          *string
	Location      *string
	Limit         *int
	Offset        *int
}
