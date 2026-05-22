package repository

import (
	"OrgAPI/internal/models"
	"context"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	GetByID(ctx context.Context, id int64) (*models.Department, error)
	GetChildren(ctx context.Context, parentID int64) ([]models.Department, error)
	GetEmployees(ctx context.Context, departmentID int64) ([]models.Employee, error)
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, id int64) error
	ReassignEmployees(ctx context.Context, fromDepartmentID int64, toDepartmentID int64) error
	WithTx(tx *gorm.DB) DepartmentRepository
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) WithTx(tx *gorm.DB) DepartmentRepository {

	return &departmentRepository{
		db: tx,
	}
}

func (r *departmentRepository) Create(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).
		Create(department).
		Error
}

func (r *departmentRepository) GetByID(ctx context.Context, id int64) (*models.Department, error) {
	var department models.Department

	err := r.db.WithContext(ctx).
		First(&department, id).
		Error

	if err != nil {
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) GetChildren(ctx context.Context, parentID int64) ([]models.Department, error) {

	var children []models.Department

	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&children).
		Error

	if err != nil {
		return nil, err
	}

	return children, nil
}

func (r *departmentRepository) GetEmployees(ctx context.Context, departmentID int64) ([]models.Employee, error) {

	var employees []models.Employee

	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Order("full_name ASC").
		Find(&employees).
		Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *departmentRepository) Update(ctx context.Context, department *models.Department) error {

	return r.db.WithContext(ctx).
		Save(department).
		Error
}

func (r *departmentRepository) Delete(ctx context.Context, id int64) error {

	return r.db.WithContext(ctx).
		Delete(&models.Department{}, id).
		Error
}

func (r *departmentRepository) ReassignEmployees(ctx context.Context, fromDepartmentID int64, toDepartmentID int64) error {

	return r.db.WithContext(ctx).
		Model(&models.Employee{}).
		Where(
			"department_id = ?",
			fromDepartmentID,
		).
		Update(
			"department_id",
			toDepartmentID,
		).
		Error
}
