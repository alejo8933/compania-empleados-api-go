package services

import (
	"compania-api/application/dtos"
	"compania-api/domain/entities"
	"compania-api/domain/interfaces"
	"errors"

	"go.uber.org/zap"
)

type EmpleadoService struct {
	uow    interfaces.IUnitOfWork
	logger *zap.Logger
}

func NewEmpleadoService(uow interfaces.IUnitOfWork, logger *zap.Logger) *EmpleadoService {
	return &EmpleadoService{
		uow:    uow,
		logger: logger,
	}
}

func (s *EmpleadoService) GetAll() ([]dtos.EmpleadoResponse, error) {
	empleados, err := s.uow.Empleados().GetAll()
	if err != nil {
		s.logger.Error("Error al obtener todos los empleados", zap.Error(err))
		return nil, err
	}

	var response []dtos.EmpleadoResponse
	for _, e := range empleados {
		response = append(response, dtos.EmpleadoResponse{
			ID:         e.ID,
			Nombre:     e.Nombre,
			Apellido:   e.Apellido,
			Correo:     e.Correo,
			Cargo:      e.Cargo,
			Salario:    e.Salario,
			CompaniaID: e.CompaniaID,
		})
	}
	return response, nil
}

func (s *EmpleadoService) GetById(id uint) (*dtos.EmpleadoResponse, error) {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		s.logger.Error("Error al obtener empleado por ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	if e == nil {
		return nil, errors.New("empleado no encontrado")
	}

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Create(req *dtos.CreateEmpleadoRequest) (*dtos.EmpleadoResponse, error) {
	// Verificar si la compañía existe
	c, err := s.uow.Companias().GetById(req.CompaniaID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("compañía asociada no encontrada")
	}

	e := &entities.Empleado{
		Nombre:     req.Nombre,
		Apellido:   req.Apellido,
		Correo:     req.Correo,
		Cargo:      req.Cargo,
		Salario:    req.Salario,
		CompaniaID: req.CompaniaID,
	}

	err = s.uow.Empleados().Create(e)
	if err != nil {
		s.logger.Error("Error al crear empleado", zap.String("correo", req.Correo), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Empleado creado",
		zap.String("nombre", e.Nombre),
		zap.String("compañía", c.Nombre),
	)

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Update(id uint, req *dtos.UpdateEmpleadoRequest) (*dtos.EmpleadoResponse, error) {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("empleado no encontrado")
	}

	e.Nombre = req.Nombre
	e.Apellido = req.Apellido
	e.Correo = req.Correo
	e.Cargo = req.Cargo
	e.Salario = req.Salario

	err = s.uow.Empleados().Update(e)
	if err != nil {
		s.logger.Error("Error al actualizar empleado", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Empleado actualizado", zap.Uint("id", e.ID))

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Delete(id uint) error {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return err
	}
	if e == nil {
		return errors.New("empleado no encontrado")
	}

	err = s.uow.Empleados().Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar empleado", zap.Uint("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Empleado eliminado", zap.Uint("id", id))
	return nil
}
