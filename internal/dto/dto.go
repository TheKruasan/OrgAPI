package dto

import "time"

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type DepartmentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`

	Employees []EmployeeResponse   `json:"employees,omitempty"`
	Children  []DepartmentResponse `json:"children,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *uint   `json:"parent_id"`
}

type DeleteDepartmentRequest struct {
	Mode                   string
	ReassignToDepartmentID *uint
}

type EmployeeResponse struct {
	ID           uint       `json:"id"`
	DepartmentID uint       `json:"department_id"`
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
