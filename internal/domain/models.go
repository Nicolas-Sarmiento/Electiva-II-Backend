package domain

import (
	"time"
)

// MedicalCondition - Enfermedad / Alergia
type MedicalCondition struct {
	ID                string `gorm:"primaryKey;type:varchar(50)" json:"medicalId"`
	Name              string `gorm:"type:varchar(50)" json:"name" validate:"required,min=3"`
	Diagnostic        string `gorm:"type:text" json:"diagnostics" validate:"required"` // Cambiado a string para mapear facil en json
	AllergenType      string `gorm:"type:varchar(50)" json:"allergenType,omitempty"`
	IsContagious      *bool  `json:"isContagious,omitempty"`
	TransmissionRoute string `gorm:"type:varchar(50)" json:"transmissionRoute,omitempty"`
}

// EmergencyContact - Contacto de emergencia
type EmergencyContact struct {
	ID           string `gorm:"primaryKey;type:varchar(50)" json:"idContact"`
	FirstName    string `gorm:"type:varchar(50)" json:"firstName" validate:"required,min=2"`
	LastName     string `gorm:"type:varchar(50)" json:"lastName" validate:"required,min=2"`
	Phone        string `gorm:"type:varchar(50)" json:"phone" validate:"required,numeric,min=7"`
	Mail         string `gorm:"type:varchar(50)" json:"mail" validate:"required,email"`
	Relationship string `gorm:"-" json:"relationship,omitempty" validate:"required"` // Viene en el JSON unido de PatientContact
}

// PatientContact - Relación Paciente-Contacto (n:m con info extra)
type PatientContact struct {
	PatientID    string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	ContactID    string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	Relationship string `gorm:"type:varchar(50)" json:"relationship"`

	Patient          Patient          `gorm:"foreignKey:PatientID;references:ID" json:"-"`
	EmergencyContact EmergencyContact `gorm:"foreignKey:ContactID;references:ID" json:"-"`
}

// PatientCondition - Relación Paciente-Enfermedad
type PatientCondition struct {
	PatientID   string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	ConditionID string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	Diagnostic  string `gorm:"type:text" json:"diagnostics"`

	Patient          Patient          `gorm:"foreignKey:PatientID;references:ID" json:"-"`
	MedicalCondition MedicalCondition `gorm:"foreignKey:ConditionID;references:ID" json:"-"`
}

// Room - Habitación
type Room struct {
	ID           string `gorm:"primaryKey;type:varchar(50)" json:"idRoom"`
	Floor        int    `json:"floor" validate:"required"`
	RoomNumber   string `gorm:"type:varchar(50)" json:"roomNumber" validate:"required"`
	RoomPavilion string `gorm:"type:varchar(50)" json:"roomPavilion" validate:"required"`
}

// Wearable - Dispositivo Wearable (Device)
type Wearable struct {
	ID           string `gorm:"primaryKey;type:varchar(50)" json:"wearableId"`
	MacAddress   string `gorm:"type:varchar(20);unique" json:"macAddress" validate:"required,mac"`
	BatteryLevel int    `json:"batteryLevel" validate:"min=0,max=100"`
	IsActive     bool   `json:"isActive"`
}

// PatientWearable - Asignación
type PatientWearable struct {
	PatientID    string    `gorm:"primaryKey;type:varchar(50)" json:"-"`
	WearableID   string    `gorm:"primaryKey;type:varchar(50)" json:"-"`
	AssignedDate time.Time `gorm:"type:timestamp" json:"-"`

	Wearable Wearable `gorm:"foreignKey:WearableID;references:ID" json:"-"`
}

// Patient - Paciente
type Patient struct {
	ID          string    `gorm:"primaryKey;type:varchar(50)" json:"patientId"`
	FirstName   string    `gorm:"type:varchar(50)" json:"firstName" validate:"required,min=2"`
	LastName    string    `gorm:"type:varchar(50)" json:"lastName" validate:"required,min=2"`
	DateOfBirth time.Time `gorm:"type:timestamp" json:"dateOfBirth" validate:"required"`
	RoomID      string    `gorm:"type:varchar(50)" json:"RoomId,omitempty" validate:"required"` // Para el payload de Creación

	Room              *Room              `gorm:"foreignKey:RoomID;references:ID" json:"Room,omitempty"`
	PatientConditions []PatientCondition `gorm:"foreignKey:PatientID" json:"-"`
	PatientContacts   []PatientContact   `gorm:"foreignKey:PatientID" json:"-"`
	PatientWearables  []PatientWearable  `gorm:"foreignKey:PatientID" json:"-"`

	// Para coincidir con el payload exacto esperado (estos campos pueden poblarse manualmente)
	Allergies        []MedicalCondition `gorm:"-" json:"Allergies,omitempty"`
	Diseases         []MedicalCondition `gorm:"-" json:"Diseases,omitempty"`
	EmergencyContact *EmergencyContact  `gorm:"-" json:"emergencyContact,omitempty"`
	WearableDevices  []Wearable         `gorm:"-" json:"wearableDevices,omitempty"`
}

// AlertType
type AlertType struct {
	ID   string `gorm:"primaryKey;type:varchar(50)" json:"alertTypeId"` // Para REST devolveremos todo u omitiremos según el dto
	Name string `gorm:"type:varchar(50)" json:"name" validate:"required"`
	Code string `gorm:"type:varchar(20)" json:"code" validate:"required"`
}

// Alert - Alerta
type Alert struct {
	ID          string     `gorm:"primaryKey;type:varchar(50)" json:"alertId"`
	PatientID   string     `gorm:"type:varchar(50);index" json:"patientId"`
	WearableID  string     `gorm:"type:varchar(50);index" json:"wearableId"`
	ResolvedAt  *time.Time `gorm:"type:timestamp" json:"resolvedAt,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"createdAt,omitempty"`
	AlertStatus string     `gorm:"type:varchar(30)" json:"alertStatus"`
	AlertLevel  string     `gorm:"type:varchar(30)" json:"alertLevel"`
	AlertTypeID string     `gorm:"type:varchar(50)" json:"alertType"` // PDF usa alertType en lugar de alertTypeId
	//NurseID     string     `gorm:"type:varchar(50);index" json:"nurseId"`
}

// Shift - Turno de trabajo recurrente
type Shift struct {
	ID        string `gorm:"primaryKey;type:varchar(50)" json:"shiftId"`
	Name      string `gorm:"type:varchar(20)" json:"name" validate:"required"`
	StartTime string `gorm:"type:varchar(8)" json:"startTime" validate:"required"` // Ej: "14:00:00"
	EndTime   string `gorm:"type:varchar(8)" json:"endTime" validate:"required"`   // Ej: "18:00:00"
}

type Nurse struct {
	ID        string `gorm:"primaryKey;type:varchar(50)" json:"nurseId"`
	FirstName string `gorm:"type:varchar(50)" json:"firstName" validate:"required,min=2"`
	LastName  string `gorm:"type:varchar(50)" json:"lastName" validate:"required,min=2"`
	Phone     string `gorm:"type:varchar(20);unique" json:"phone" validate:"required,numeric,min=7"`
	ShiftID   string `gorm:"type:varchar(50)" json:"shiftId" validate:"required"`
}
