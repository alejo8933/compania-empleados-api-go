package entities

import (
	"time"
)

type Compania struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre        string     `gorm:"not null" json:"nombre"`
	Direccion     string     `gorm:"not null" json:"direccion"`
	Telefono      string     `gorm:"not null" json:"telefono"`
	FechaCreacion time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"fecha_creacion"`
	Empleados     []Empleado `gorm:"foreignKey:CompaniaID;constraint:OnDelete:CASCADE" json:"empleados,omitempty"`
}
