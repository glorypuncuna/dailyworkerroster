package model

import "time"

type Shift struct {
	ID             int64     `json:"id"`
	Date           string    `json:"date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	RoleAssignment string    `json:"role_assignment"`
	Location       string    `json:"location"`
	IsAvailable    bool      `json:"isAvailable"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ShiftStatus struct {
	ID             int64  `json:"id"`
	Date           string `json:"date"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	RoleAssignment string `json:"role_assignment"`
	Location       string `json:"location"`
	IsAvailable    bool   `json:"isAvailable"`
	StatusWorker   string `json:"status_worker"`
}

type ShiftListQuery struct {
	Limit          int
	Offset         int
	RoleAssignment string
	Location       string
	Date           string
	IsAvailable    *bool
}
