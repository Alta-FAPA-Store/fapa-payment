package payment

import (
	"gorm.io/gorm"
)

//GormRepository The implementation of pet.Repository object
type GormRepository struct {
	DB *gorm.DB
}

//NewGormDBRepository Generate Gorm DB pet repository
func NewGormDBRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db,
	}
}
