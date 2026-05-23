package repository

import (
	"OrgAPI/internal/models"
	"context"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Create(
	ctx context.Context,
	employee *models.Employee,
) error {

	return r.db.WithContext(ctx).
		Create(employee).
		Error
}
