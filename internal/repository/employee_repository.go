package repository

import (
	"OrgAPI/internal/models"
	"context"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(
		ctx context.Context,
		employee *models.Employee,
	) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(
	db *gorm.DB,
) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(
	ctx context.Context,
	employee *models.Employee,
) error {

	return r.db.WithContext(ctx).
		Create(employee).
		Error
}
