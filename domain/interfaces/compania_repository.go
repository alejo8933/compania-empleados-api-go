package interfaces

import "compania-api/domain/entities"

type ICompaniaRepository interface {
	GetAll() ([]entities.Compania, error)
	GetById(id uint) (*entities.Compania, error)
	Create(c *entities.Compania) error
	Update(c *entities.Compania) error
	Delete(id uint) error
	FindByCondition(query interface{}, args ...interface{}) ([]entities.Compania, error)
}
