package models

import "time"

type Department struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"size:200;not null"`

	ParentID *uint

	Parent *Department

	Children []Department `gorm:"foreignKey:ParentID"`

	Employees []Employee

	CreatedAt time.Time
}

type Employee struct {
	ID uint `gorm:"primaryKey"`

	DepartmentID uint `gorm:"not null"`

	Department Department

	FullName string `gorm:"size:200;not null"`

	Position string `gorm:"size:200;not null"`

	HiredAt *time.Time

	CreatedAt time.Time
}
