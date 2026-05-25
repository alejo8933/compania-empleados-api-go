package entities

type Empleado struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre     string    `gorm:"not null" json:"nombre"`
	Apellido   string    `gorm:"not null" json:"apellido"`
	Correo     string    `gorm:"uniqueIndex;not null" json:"correo"`
	Cargo      string    `gorm:"not null" json:"cargo"`
	Salario    float64   `gorm:"not null" json:"salario"`
	CompaniaID uint      `gorm:"not null" json:"compania_id"`
	Compania   *Compania `gorm:"foreignKey:CompaniaID" json:"compania,omitempty"`
}
