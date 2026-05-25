package repositories

import (
	"compania-api/domain/entities"

	"gorm.io/gorm"
)

type CompaniaRepositoryImpl struct {
	db *gorm.DB
}

func NewCompaniaRepository(db *gorm.DB) *CompaniaRepositoryImpl {
	return &CompaniaRepositoryImpl{db: db}
}

func (r *CompaniaRepositoryImpl) GetAll() ([]entities.Compania, error) {
	var companias []entities.Compania
	err := r.db.Preload("Empleados").Find(&companias).Error
	return companias, err
}

func (r *CompaniaRepositoryImpl) GetById(id uint) (*entities.Compania, error) {
	var c entities.Compania
	err := r.db.Preload("Empleados").First(&c, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CompaniaRepositoryImpl) Create(c *entities.Compania) error {
	return r.db.Create(c).Error
}

func (r *CompaniaRepositoryImpl) Update(c *entities.Compania) error {
	return r.db.Save(c).Error
}

func (r *CompaniaRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Compania{}, id).Error
}

func (r *CompaniaRepositoryImpl) FindByCondition(query interface{}, args ...interface{}) ([]entities.Compania, error) {
	var companias []entities.Compania
	err := r.db.Preload("Empleados").Where(query, args...).Find(&companias).Error
	return companias, err
}
