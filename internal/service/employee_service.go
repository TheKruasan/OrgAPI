package service

import (
	"OrgAPI/internal/dto"
	"OrgAPI/internal/models"
	"OrgAPI/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(
		ctx context.Context,
		employee *models.Employee,
	) error
}

type EmployeeService struct {
	employeeRepo   EmployeeRepository
	departmentRepo DepartmentRepository
}

func NewEmployeeService(employeeRepo EmployeeRepository, departmentRepo DepartmentRepository) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, departmentID int64, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {

	_, err := s.departmentRepo.GetByID(
		ctx,
		departmentID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(
				"department not found",
			)
		}

		return nil, err
	}

	fullName, err := utils.ValidateRequiredString(
		req.FullName,
		"full_name",
	)

	if err != nil {
		return nil, err
	}

	position, err := utils.ValidateRequiredString(
		req.Position,
		"position",
	)

	if err != nil {
		return nil, err
	}

	employee := &models.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      req.HiredAt,
	}

	err = s.employeeRepo.Create(
		ctx,
		employee,
	)

	if err != nil {
		return nil, err
	}

	return &dto.EmployeeResponse{
		ID:           employee.ID,
		DepartmentID: employee.DepartmentID,
		FullName:     employee.FullName,
		Position:     employee.Position,
		HiredAt:      employee.HiredAt,
		CreatedAt:    employee.CreatedAt,
	}, nil
}
