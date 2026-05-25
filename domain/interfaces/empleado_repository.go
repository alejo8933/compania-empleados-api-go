package interfaces

import "compania-api/domain/entities"

type IEmpleadoRepository interface {
	GetAll() ([]entities.Empleado, error)
	GetById(id uint) (*entities.Empleado, error)
	Create(e *entities.Empleado) error
	Update(e *entities.Empleado) error
	Delete(id uint) error
	FindByCondition(query interface{}, args ...interface{}) ([]entities.Empleado, error)
}
