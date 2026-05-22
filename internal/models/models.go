package models

import "time"

type Department struct {
	ID int64 `gorm:"primaryKey"`

	Name string `gorm:"size:200;not null"`

	ParentID *int64

	Parent *Department

	Children []Department `gorm:"foreignKey:ParentID"`

	Employees []Employee

	CreatedAt time.Time
}

type Employee struct {
	ID int64 `gorm:"primaryKey"`

	DepartmentID int64 `gorm:"not null"`

	Department Department

	FullName string `gorm:"size:200;not null"`

	Position string `gorm:"size:200;not null"`

	HiredAt *time.Time

	CreatedAt time.Time
}
