package repository

import (
	"context"
	"errors"

	"asset-backend/internal/user/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type EmployeeRepository interface {
	Create(ctx context.Context, employee *domain.Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetManagerOf(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error)
	List(ctx context.Context) ([]domain.Employee, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.WithContext(ctx).First(&employee, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	var employee domain.Employee
	if err := r.db.WithContext(ctx).First(&employee, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return &employee, nil
}

// GetManagerOf follows an employee's ManagerID and returns that manager's record.
func (r *employeeRepository) GetManagerOf(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error) {
	employee, err := r.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	if employee.ManagerID == nil {
		return nil, ErrEmployeeNotFound
	}
	return r.GetByID(ctx, *employee.ManagerID)
}

func (r *employeeRepository) List(ctx context.Context) ([]domain.Employee, error) {
	var employees []domain.Employee
	err := r.db.WithContext(ctx).Find(&employees).Error
	return employees, err
}
