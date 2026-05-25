package services

import (
	"compania-api/application/dtos"
	"compania-api/domain/entities"
	"compania-api/domain/interfaces"
	"errors"

	"go.uber.org/zap"
)

type CompaniaService struct {
	uow    interfaces.IUnitOfWork
	logger *zap.Logger
}

func NewCompaniaService(uow interfaces.IUnitOfWork, logger *zap.Logger) *CompaniaService {
	return &CompaniaService{
		uow:    uow,
		logger: logger,
	}
}

func (s *CompaniaService) GetAll() ([]dtos.CompaniaResponse, error) {
	companias, err := s.uow.Companias().GetAll()
	if err != nil {
		s.logger.Error("Error al obtener todas las compañías", zap.Error(err))
		return nil, err
	}

	var response []dtos.CompaniaResponse
	for _, c := range companias {
		var empsResponse []dtos.EmpleadoResponse
		for _, e := range c.Empleados {
			empsResponse = append(empsResponse, dtos.EmpleadoResponse{
				ID:         e.ID,
				Nombre:     e.Nombre,
				Apellido:   e.Apellido,
				Correo:     e.Correo,
				Cargo:      e.Cargo,
				Salario:    e.Salario,
				CompaniaID: e.CompaniaID,
			})
		}
		response = append(response, dtos.CompaniaResponse{
			ID:            c.ID,
			Nombre:        c.Nombre,
			Direccion:     c.Direccion,
			Telefono:      c.Telefono,
			FechaCreacion: c.FechaCreacion,
			Empleados:     empsResponse,
		})
	}
	return response, nil
}

func (s *CompaniaService) GetById(id uint) (*dtos.CompaniaResponse, error) {
	c, err := s.uow.Companias().GetById(id)
	if err != nil {
		s.logger.Error("Error al obtener compañía por ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	if c == nil {
		return nil, errors.New("compañía no encontrada")
	}

	var empsResponse []dtos.EmpleadoResponse
	for _, e := range c.Empleados {
		empsResponse = append(empsResponse, dtos.EmpleadoResponse{
			ID:         e.ID,
			Nombre:     e.Nombre,
			Apellido:   e.Apellido,
			Correo:     e.Correo,
			Cargo:      e.Cargo,
			Salario:    e.Salario,
			CompaniaID: e.CompaniaID,
		})
	}

	return &dtos.CompaniaResponse{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Direccion:     c.Direccion,
		Telefono:      c.Telefono,
		FechaCreacion: c.FechaCreacion,
		Empleados:     empsResponse,
	}, nil
}

func (s *CompaniaService) Create(req *dtos.CreateCompaniaRequest) (*dtos.CompaniaResponse, error) {
	c := &entities.Compania{
		Nombre:    req.Nombre,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
	}

	err := s.uow.Companias().Create(c)
	if err != nil {
		s.logger.Error("Error al crear compañía", zap.String("nombre", req.Nombre), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Compañía creada", zap.String("nombre", c.Nombre), zap.Uint("id", c.ID))

	return &dtos.CompaniaResponse{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Direccion:     c.Direccion,
		Telefono:      c.Telefono,
		FechaCreacion: c.FechaCreacion,
	}, nil
}

func (s *CompaniaService) Update(id uint, req *dtos.UpdateCompaniaRequest) (*dtos.CompaniaResponse, error) {
	c, err := s.uow.Companias().GetById(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("compañía no encontrada")
	}

	c.Nombre = req.Nombre
	c.Direccion = req.Direccion
	c.Telefono = req.Telefono

	err = s.uow.Companias().Update(c)
	if err != nil {
		s.logger.Error("Error al actualizar compañía", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Compañía actualizada", zap.Uint("id", c.ID))

	return &dtos.CompaniaResponse{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Direccion:     c.Direccion,
		Telefono:      c.Telefono,
		FechaCreacion: c.FechaCreacion,
	}, nil
}

func (s *CompaniaService) Delete(id uint) error {
	c, err := s.uow.Companias().GetById(id)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("compañía no encontrada")
	}

	err = s.uow.Companias().Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar compañía", zap.Uint("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Compañía eliminada", zap.Uint("id", id))
	return nil
}

func (s *CompaniaService) CreateWithEmpleados(req *dtos.CreateCompaniaWithEmpleadosRequest) (*dtos.CompaniaResponse, error) {
	s.logger.Info("Iniciando transacción para crear compañía con empleados")
	err := s.uow.BeginTransaction()
	if err != nil {
		s.logger.Error("Error al iniciar transacción", zap.Error(err))
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Pánico detectado, ejecutando rollback", zap.Any("panic", r))
			s.uow.Rollback()
		}
	}()

	// 1. Crear compañía
	c := &entities.Compania{
		Nombre:    req.Nombre,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
	}

	err = s.uow.Companias().Create(c)
	if err != nil {
		s.logger.Error("Error al crear compañía en transacción, ejecutando rollback", zap.Error(err))
		s.uow.Rollback()
		return nil, err
	}
	s.logger.Info("Compañía creada en transacción", zap.String("nombre", c.Nombre), zap.Uint("id", c.ID))

	// 2. Crear empleados
	var empsCreated []entities.Empleado
	for _, empReq := range req.Empleados {
		emp := &entities.Empleado{
			Nombre:     empReq.Nombre,
			Apellido:   empReq.Apellido,
			Correo:     empReq.Correo,
			Cargo:      empReq.Cargo,
			Salario:    empReq.Salario,
			CompaniaID: c.ID,
		}

		err = s.uow.Empleados().Create(emp)
		if err != nil {
			s.logger.Error("Error al crear empleado en transacción, ejecutando rollback",
				zap.String("correo", emp.Correo),
				zap.Error(err),
			)
			s.uow.Rollback()
			return nil, err
		}
		s.logger.Info("Empleado creado en transacción",
			zap.String("nombre", emp.Nombre),
			zap.String("compañía", c.Nombre),
		)
		empsCreated = append(empsCreated, *emp)
	}

	// 3. Confirmar transacción
	err = s.uow.Commit()
	if err != nil {
		s.logger.Error("Error al confirmar (commit) transacción", zap.Error(err))
		s.uow.Rollback()
		return nil, err
	}

	s.logger.Info("Transacción completada y confirmada con éxito")

	// Mapear respuesta
	var empsResponse []dtos.EmpleadoResponse
	for _, e := range empsCreated {
		empsResponse = append(empsResponse, dtos.EmpleadoResponse{
			ID:         e.ID,
			Nombre:     e.Nombre,
			Apellido:   e.Apellido,
			Correo:     e.Correo,
			Cargo:      e.Cargo,
			Salario:    e.Salario,
			CompaniaID: e.CompaniaID,
		})
	}

	return &dtos.CompaniaResponse{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Direccion:     c.Direccion,
		Telefono:      c.Telefono,
		FechaCreacion: c.FechaCreacion,
		Empleados:     empsResponse,
	}, nil
}
