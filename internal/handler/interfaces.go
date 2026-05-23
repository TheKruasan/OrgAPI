package handler

import (
	"OrgAPI/internal/dto"
	"context"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	GetDepartmentTree(ctx context.Context, id int64, depth int64, includeEmployees bool) (*dto.DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id int64, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, id int64, req dto.DeleteDepartmentRequest) error
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, departmentID int64, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
}
