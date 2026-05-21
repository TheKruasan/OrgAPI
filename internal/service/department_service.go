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

type DepartmentService interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	GetDepartmentTree(ctx context.Context, id uint, depth int, includeEmployees bool) (*dto.DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id uint, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, id uint, req dto.DeleteDepartmentRequest) error
}

type departmentService struct {
	repo repository.DepartmentRepository
	db   *gorm.DB
}

func NewDepartmentService(repo repository.DepartmentRepository, db *gorm.DB) DepartmentService {

	return &departmentService{
		repo: repo,
		db:   db,
	}
}

func (s *departmentService) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {

	name, err := utils.ValidateName(req.Name)

	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		_, err := s.repo.GetByID(ctx, *req.ParentID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent department not found")
			}

			return nil, err
		}
	}

	department := &models.Department{
		Name:     name,
		ParentID: req.ParentID,
	}

	err = s.repo.Create(ctx, department)

	if err != nil {
		return nil, err
	}

	return &dto.DepartmentResponse{
		ID:        department.ID,
		Name:      department.Name,
		ParentID:  department.ParentID,
		CreatedAt: department.CreatedAt,
	}, nil
}

func (s *departmentService) GetDepartmentTree(ctx context.Context, id uint, depth int, includeEmployees bool) (*dto.DepartmentResponse, error) {

	if depth < 1 {
		depth = 1
	}

	if depth > 5 {
		depth = 5
	}

	department, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	response, err := s.buildDepartmentTree(
		ctx,
		*department,
		depth,
		includeEmployees,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *departmentService) buildDepartmentTree(ctx context.Context, department models.Department, depth int, includeEmployees bool) (*dto.DepartmentResponse, error) {

	response := &dto.DepartmentResponse{
		ID:        department.ID,
		Name:      department.Name,
		ParentID:  department.ParentID,
		CreatedAt: department.CreatedAt,
	}
	if includeEmployees {
		employees, err := s.repo.GetEmployees(
			ctx,
			department.ID,
		)

		if err != nil {
			return nil, err
		}

		for _, employee := range employees {
			response.Employees = append(
				response.Employees,
				dto.EmployeeResponse{
					ID:           employee.ID,
					DepartmentID: employee.DepartmentID,
					FullName:     employee.FullName,
					Position:     employee.Position,
					HiredAt:      employee.HiredAt,
					CreatedAt:    employee.CreatedAt,
				},
			)
		}
	}
	if depth > 1 {
		children, err := s.repo.GetChildren(
			ctx,
			department.ID,
		)

		if err != nil {
			return nil, err
		}

		for _, child := range children {

			childResponse, err := s.buildDepartmentTree(
				ctx,
				child,
				depth-1,
				includeEmployees,
			)

			if err != nil {
				return nil, err
			}

			response.Children = append(
				response.Children,
				*childResponse,
			)
		}
	}

	return response, nil
}

func (s *departmentService) UpdateDepartment(ctx context.Context, id uint, req dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {

	department, err := s.repo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}
	if req.Name != nil {

		name, err := utils.ValidateName(
			*req.Name,
		)

		if err != nil {
			return nil, err
		}

		department.Name = name
	}
	if req.ParentID != nil {

		if *req.ParentID == id {
			return nil, errors.New(
				"department cannot be parent of itself",
			)
		}

		_, err := s.repo.GetByID(
			ctx,
			*req.ParentID,
		)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(
					"parent department not found",
				)
			}

			return nil, err
		}

		hasCycle, err := s.wouldCreateCycle(
			ctx,
			id,
			*req.ParentID,
		)

		if err != nil {
			return nil, err
		}

		if hasCycle {
			return nil, errors.New(
				"cycle detected",
			)
		}

		department.ParentID = req.ParentID
	}
	err = s.repo.Update(
		ctx,
		department,
	)

	if err != nil {
		return nil, err
	}

	return &dto.DepartmentResponse{
		ID:        department.ID,
		Name:      department.Name,
		ParentID:  department.ParentID,
		CreatedAt: department.CreatedAt,
	}, nil
}

func (s *departmentService) wouldCreateCycle(ctx context.Context, departmentID uint, targetParentID uint) (bool, error) {

	currentID := targetParentID

	for {

		if currentID == departmentID {
			return true, nil
		}

		department, err := s.repo.GetByID(
			ctx,
			currentID,
		)

		if err != nil {
			return false, err
		}

		if department.ParentID == nil {
			break
		}

		currentID = *department.ParentID
	}

	return false, nil
}

func (s *departmentService) DeleteDepartment(ctx context.Context, id uint, req dto.DeleteDepartmentRequest) error {

	_, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if req.Mode == "cascade" {
		return s.repo.Delete(ctx, id)
	}
	if req.Mode != "reassign" {
		return errors.New("invalid delete mode")
	}
	if req.ReassignToDepartmentID == nil {
		return errors.New(
			"reassign_to_department_id is required",
		)
	}
	if *req.ReassignToDepartmentID == id {
		return errors.New(
			"cannot reassign to same department",
		)
	}
	_, err = s.repo.GetByID(
		ctx,
		*req.ReassignToDepartmentID,
	)

	if err != nil {
		return errors.New(
			"target department not found",
		)
	}
	hasCycle, err := s.wouldCreateCycle(
		ctx,
		id,
		*req.ReassignToDepartmentID,
	)

	if err != nil {
		return err
	}

	if hasCycle {
		return errors.New(
			"cannot reassign into subtree",
		)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		txRepo := s.repo.WithTx(tx)

		err := txRepo.ReassignEmployees(
			ctx,
			id,
			*req.ReassignToDepartmentID,
		)

		if err != nil {
			return err
		}

		err = txRepo.Delete(ctx, id)

		if err != nil {
			return err
		}

		return nil
	})
}
