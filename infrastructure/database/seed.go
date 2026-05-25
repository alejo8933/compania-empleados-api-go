package database

import (
	"compania-api/domain/entities"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, logger *zap.Logger) error {
	logger.Info("Iniciando ejecución de semillas (seed) en la base de datos")

	companiasSeeds := []entities.Compania{
		{Nombre: "TechCorp S.A.S", Direccion: "Calle 100 # 15-30, Bogotá", Telefono: "6012345678"},
		{Nombre: "DataSolutions Ltda", Direccion: "Carrera 45 # 72-10, Medellín", Telefono: "6048765432"},
		{Nombre: "DevGroup Colombia", Direccion: "Avenida El Dorado # 68C-25, Bogotá", Telefono: "6019876543"},
	}

	companiasMap := make(map[string]entities.Compania)

	for _, compSeed := range companiasSeeds {
		var c entities.Compania
		err := db.Where("nombre = ?", compSeed.Nombre).First(&c).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c = compSeed
				if err := db.Create(&c).Error; err != nil {
					logger.Error("Error al sembrar compañía", zap.String("nombre", compSeed.Nombre), zap.Error(err))
					return err
				}
				logger.Info("Semilla: Compañía insertada", zap.String("nombre", c.Nombre), zap.Uint("id", c.ID))
			} else {
				logger.Error("Error al verificar existencia de compañía", zap.String("nombre", compSeed.Nombre), zap.Error(err))
				return err
			}
		} else {
			logger.Info("Semilla: Compañía ya existe, omitiendo inserción", zap.String("nombre", c.Nombre))
		}
		companiasMap[c.Nombre] = c
	}

	empleadosSeeds := []struct {
		Nombre         string
		Apellido       string
		Correo         string
		Cargo          string
		Salario        float64
		CompaniaNombre string
	}{
		{"Juan", "Pérez", "juan.perez@techcorp.com", "Software Engineer", 4500000, "TechCorp S.A.S"},
		{"María", "Gómez", "maria.gomez@techcorp.com", "QA Lead", 3800000, "TechCorp S.A.S"},
		{"Pedro", "Rodríguez", "pedro.rodriguez@techcorp.com", "Product Owner", 5500000, "TechCorp S.A.S"},
		{"Laura", "Sánchez", "laura.sanchez@techcorp.com", "UX Designer", 4200000, "TechCorp S.A.S"},
		{"Carlos", "Martínez", "carlos.martinez@datasolutions.com", "Data Analyst", 3600000, "DataSolutions Ltda"},
		{"Ana", "López", "ana.lopez@datasolutions.com", "Data Engineer", 4800000, "DataSolutions Ltda"},
		{"Diego", "Ramírez", "diego.ramirez@datasolutions.com", "Scrum Master", 4600000, "DataSolutions Ltda"},
		{"Sofía", "Castro", "sofia.castro@devgroup.com", "Backend Developer", 4100000, "DevGroup Colombia"},
		{"Luis", "Herrera", "luis.herrera@devgroup.com", "DevOps Engineer", 5000000, "DevGroup Colombia"},
		{"Elena", "Díaz", "elena.diaz@devgroup.com", "Frontend Developer", 3900000, "DevGroup Colombia"},
	}

	for _, empSeed := range empleadosSeeds {
		var e entities.Empleado
		err := db.Where("correo = ?", empSeed.Correo).First(&e).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				comp, ok := companiasMap[empSeed.CompaniaNombre]
				if !ok {
					logger.Error("Compañía de semilla no encontrada en mapa", zap.String("compania", empSeed.CompaniaNombre))
					continue
				}

				e = entities.Empleado{
					Nombre:     empSeed.Nombre,
					Apellido:   empSeed.Apellido,
					Correo:     empSeed.Correo,
					Cargo:      empSeed.Cargo,
					Salario:    empSeed.Salario,
					CompaniaID: comp.ID,
				}

				if err := db.Create(&e).Error; err != nil {
					logger.Error("Error al sembrar empleado", zap.String("correo", empSeed.Correo), zap.Error(err))
					return err
				}
				logger.Info("Semilla: Empleado insertado", zap.String("nombre", e.Nombre), zap.String("compañía", empSeed.CompaniaNombre))
			} else {
				logger.Error("Error al verificar existencia de empleado", zap.String("correo", empSeed.Correo), zap.Error(err))
				return err
			}
		} else {
			logger.Info("Semilla: Empleado ya existe, omitiendo inserción", zap.String("correo", e.Correo))
		}
	}

	logger.Info("Ejecución de semillas finalizada con éxito")
	return nil
}
