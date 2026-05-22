package service

import (
	"OrgAPI/internal/dto"
	"OrgAPI/internal/models"
	"OrgAPI/internal/repository"
	"OrgAPI/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(
		ctx context.Context,
		departmentID int64,
		req dto.CreateEmployeeRequest,
	) (*dto.EmployeeResponse, error)
}

type employeeService struct {
	employeeRepo   repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	departmentRepo repository.DepartmentRepository,
) EmployeeService {
	return &employeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *employeeService) CreateEmployee(
	ctx context.Context,
	departmentID int64,
	req dto.CreateEmployeeRequest,
) (*dto.EmployeeResponse, error) {

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
