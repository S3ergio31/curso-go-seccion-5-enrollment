package enrollment

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
	courseSdk "github.com/S3ergio31/curso-go-seccion-5-sdk/course"
	userSdk "github.com/S3ergio31/curso-go-seccion-5-sdk/user"
)

type Service interface {
	Create(userID, courseID string) (*domain.Enrollment, error)
	GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Update(id string, status *string) error
	Count(filters Filters) (int, error)
}

type service struct {
	logger      *log.Logger
	repository  Repository
	userTrans   userSdk.Transport
	courseTrans courseSdk.Transport
}

type Filters struct {
	UserID   string
	CourseID string
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   domain.Pending,
	}

	if _, err := s.userTrans.Get(userID); err != nil {
		return nil, err
	}

	if _, err := s.courseTrans.Get(courseID); err != nil {
		return nil, err
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
	if status != nil {
		switch domain.EnrollmentStatus(*status) {
		case domain.Pending, domain.Active, domain.Studying, domain.Inactive:
		default:
			return ErrorInvalidStatus{*status}
		}
	}

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
	userTrans userSdk.Transport,
	courseTrans courseSdk.Transport,
) Service {
	return &service{
		logger:      logger,
		repository:  repository,
		userTrans:   userTrans,
		courseTrans: courseTrans,
	}
}
