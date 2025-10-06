package enrollment

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
)

type Service interface {
	Create(userID, courseID string) (*domain.Enrollment, error)
	GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Update(id string, status *string) error
	Count(filters Filters) (int, error)
}

type service struct {
	logger     *log.Logger
	repository Repository
}

type Filters struct {
	UserID   string
	CourseID string
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}

	if err := s.repository.Create(enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repository.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (s service) Update(id string, status *string) error {
	if err := s.repository.Update(id, status); err != nil {
		return err
	}
	return nil
}

func (s service) Count(filters Filters) (int, error) {
	return s.repository.Count(filters)
}

func NewService(
	repository Repository,
	logger *log.Logger,
) Service {
	return &service{
		logger:     logger,
		repository: repository,
	}
}
