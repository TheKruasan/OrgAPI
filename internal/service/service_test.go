package service

import (
	"OrgAPI/internal/dto"
	"OrgAPI/internal/models"
	mock_service "OrgAPI/internal/service/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateDepartment(t *testing.T) {

	tests := []struct {
		name      string
		req       dto.CreateDepartmentRequest
		mockSetup func(repo *mock_service.MockDepartmentRepository)
		wantErr   bool
	}{
		{
			name: "ok create root department",
			req: dto.CreateDepartmentRequest{
				Name: "IT",
			},
			mockSetup: func(repo *mock_service.MockDepartmentRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, d *models.Department) error {
						d.ID = 1
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "parent not found",
			req: dto.CreateDepartmentRequest{
				Name:     "IT",
				ParentID: ptr(int64(10)),
			},
			mockSetup: func(repo *mock_service.MockDepartmentRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), int64(10)).
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "repo create error",
			req: dto.CreateDepartmentRequest{
				Name: "IT",
			},
			mockSetup: func(repo *mock_service.MockDepartmentRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockDepartmentRepository(ctrl)
			db := &gorm.DB{}

			s := NewDepartmentService(repo, db)

			tt.mockSetup(repo)

			_, err := s.CreateDepartment(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDepartmentTree(t *testing.T) {

	tests := []struct {
		name      string
		id        int64
		depth     int64
		include   bool
		mockSetup func(repo *mock_service.MockDepartmentRepository)
		wantErr   bool
	}{
		{
			name:    "depth 1 - no children",
			id:      1,
			depth:   1,
			include: true,
			mockSetup: func(repo *mock_service.MockDepartmentRepository) {

				repo.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&models.Department{ID: 1, Name: "Root"}, nil)

				repo.EXPECT().
					GetEmployees(gomock.Any(), int64(1)).
					Return([]models.Employee{}, nil)
			},
			wantErr: false,
		},
		{
			name:    "depth 2 - with children",
			id:      1,
			depth:   2,
			include: true,
			mockSetup: func(repo *mock_service.MockDepartmentRepository) {

				// ROOT
				repo.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&models.Department{ID: 1, Name: "Root"}, nil)

				repo.EXPECT().
					GetEmployees(gomock.Any(), int64(1)).
					Return([]models.Employee{}, nil)

				repo.EXPECT().
					GetChildren(gomock.Any(), int64(1)).
					Return([]models.Department{
						{ID: 2, Name: "Child"},
					}, nil)

				// CHILD (ВАЖНО: AnyTimes убирает строгую зависимость порядка)
				repo.EXPECT().
					GetByID(gomock.Any(), int64(2)).
					Return(&models.Department{ID: 2, Name: "Child"}, nil).
					AnyTimes()

				repo.EXPECT().
					GetEmployees(gomock.Any(), int64(2)).
					Return([]models.Employee{}, nil).
					AnyTimes()

				repo.EXPECT().
					GetChildren(gomock.Any(), int64(2)).
					Return([]models.Department{}, nil).
					AnyTimes()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockDepartmentRepository(ctrl)
			s := NewDepartmentService(repo, &gorm.DB{})

			tt.mockSetup(repo)

			res, err := s.GetDepartmentTree(
				context.Background(),
				tt.id,
				tt.depth,
				tt.include,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}

func ptr(v int64) *int64 {
	return &v
}
