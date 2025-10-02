package enrollment

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
)

type Service interface {
	Create(userID, courseID string) (*domain.Enrollment, error)
}

type service struct {
	logger     *log.Logger
	repository Repository
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

func NewService(
	repository Repository,
	logger *log.Logger,
) Service {
	return &service{
		logger:     logger,
		repository: repository,
	}
}
