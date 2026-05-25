package unit_of_work

import (
	"errors"
	"compania-api/domain/interfaces"
	"compania-api/infrastructure/repositories"
	"gorm.io/gorm"
)

type UnitOfWorkImpl struct {
	db           *gorm.DB
	tx           *gorm.DB
	companiaRepo interfaces.ICompaniaRepository
	empleadoRepo interfaces.IEmpleadoRepository
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		db: db,
	}
}

func (u *UnitOfWorkImpl) Companias() interfaces.ICompaniaRepository {
	if u.tx != nil {
		return repositories.NewCompaniaRepository(u.tx)
	}
	if u.companiaRepo == nil {
		u.companiaRepo = repositories.NewCompaniaRepository(u.db)
	}
	return u.companiaRepo
}

func (u *UnitOfWorkImpl) Empleados() interfaces.IEmpleadoRepository {
	if u.tx != nil {
		return repositories.NewEmpleadoRepository(u.tx)
	}
	if u.empleadoRepo == nil {
		u.empleadoRepo = repositories.NewEmpleadoRepository(u.db)
	}
	return u.empleadoRepo
}

func (u *UnitOfWorkImpl) BeginTransaction() error {
	if u.tx != nil {
		return errors.New("ya existe una transacción activa")
	}
	u.tx = u.db.Begin()
	if u.tx.Error != nil {
		err := u.tx.Error
		u.tx = nil
		return err
	}
	return nil
}

func (u *UnitOfWorkImpl) Commit() error {
	if u.tx == nil {
		return errors.New("no hay una transacción activa para confirmar")
	}
	err := u.tx.Commit().Error
	u.tx = nil
	return err
}

func (u *UnitOfWorkImpl) Rollback() error {
	if u.tx == nil {
		return errors.New("no hay una transacción activa para revertir")
	}
	err := u.tx.Rollback().Error
	u.tx = nil
	return err
}

func (u *UnitOfWorkImpl) Save() error {
	return u.Commit()
}
