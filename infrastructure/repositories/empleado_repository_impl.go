package repositories

import (
	"compania-api/domain/entities"

	"gorm.io/gorm"
)

type EmpleadoRepositoryImpl struct {
	db *gorm.DB
}

func NewEmpleadoRepository(db *gorm.DB) *EmpleadoRepositoryImpl {
	return &EmpleadoRepositoryImpl{db: db}
}

func (r *EmpleadoRepositoryImpl) GetAll() ([]entities.Empleado, error) {
	var empleados []entities.Empleado
	err := r.db.Find(&empleados).Error
	return empleados, err
}

func (r *EmpleadoRepositoryImpl) GetById(id uint) (*entities.Empleado, error) {
	var e entities.Empleado
	err := r.db.First(&e, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *EmpleadoRepositoryImpl) Create(e *entities.Empleado) error {
	return r.db.Create(e).Error
}

func (r *EmpleadoRepositoryImpl) Update(e *entities.Empleado) error {
	return r.db.Save(e).Error
}

func (r *EmpleadoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Empleado{}, id).Error
}

func (r *EmpleadoRepositoryImpl) FindByCondition(query interface{}, args ...interface{}) ([]entities.Empleado, error) {
	var empleados []entities.Empleado
	err := r.db.Where(query, args...).Find(&empleados).Error
	return empleados, err
}
