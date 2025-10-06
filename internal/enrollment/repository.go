package enrollment

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(enrollment *domain.Enrollment) error
	GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Update(id string, status *string) error
	Count(filters Filters) (int, error)
}

type repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func (r repository) Create(enrollment *domain.Enrollment) error {
	if err := r.db.Create(enrollment).Error; err != nil {
		r.logger.Println(err)
		return err
	}

	r.logger.Println("enrollment created with id: ", enrollment.ID)
	return nil
}

func (r repository) GetAll(filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment

	tx := r.db.Model(&enrollments)

	tx = applyFilters(tx, filters)

	tx = tx.Limit(limit).Offset(offset)

	if err := tx.Order("created_at desc").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r repository) Update(id string, status *string) error {
	values := make(map[string]any, 0)

	if status != nil {
		values["status"] = *status
	}

	result := r.db.Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorEnrollmentNotFound{id}
	}

	return nil
}

func (r repository) Count(filters Filters) (int, error) {
	var count int64

	tx := r.db.Model(domain.Enrollment{})

	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserID != "" {
		tx = tx.Where("user_id = ?", filters.UserID)
	}

	if filters.CourseID != "" {
		tx = tx.Where("course_id = ?", filters.CourseID)
	}

	return tx
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}
