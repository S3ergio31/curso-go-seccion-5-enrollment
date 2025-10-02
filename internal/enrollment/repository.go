package enrollment

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(enrollment *domain.Enrollment) error
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

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}
