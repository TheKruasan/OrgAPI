package dto

import (
	"OrgAPI/internal/dto/types"
	"time"
)

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

type DepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`

	Employees []EmployeeResponse   `json:"employees,omitempty"`
	Children  []DepartmentResponse `json:"children,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name     *string       `json:"name"`
	ParentID types.NullInt `json:"parent_id"`
}

type DeleteDepartmentRequest struct {
	Mode                   string
	ReassignToDepartmentID *int64
}

type EmployeeResponse struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}
