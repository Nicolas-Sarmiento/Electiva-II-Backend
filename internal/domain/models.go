package domain

import (
	"encoding/json"
	"time"
)

// MedicalCondition - Enfermedad / Alergia
type MedicalCondition struct {
	ID                string `gorm:"primaryKey;type:varchar(50)" json:"medicalId"`
	Name              string `gorm:"type:varchar(50)" json:"name" validate:"required,min=3"`
	Diagnostic        string `gorm:"type:text" json:"diagnostics" validate:"required"` // Cambiado a string para mapear facil en json
	AllergenType      string `gorm:"type:varchar(50)" json:"allergenType"`
	IsContagious      *bool  `json:"isContagious"`
	TransmissionRoute string `gorm:"type:varchar(50)" json:"transmissionRoute"`
}

// EmergencyContact - Contacto de emergencia
type EmergencyContact struct {
	ID           string `gorm:"primaryKey;type:varchar(50)" json:"idContact"`
	FirstName    string `gorm:"type:varchar(50)" json:"firstName" validate:"required,min=2"`
	LastName     string `gorm:"type:varchar(50)" json:"lastName" validate:"required,min=2"`
	Phone        string `gorm:"type:varchar(50)" json:"phone" validate:"required"`
	Mail         string `gorm:"type:varchar(50)" json:"mail" validate:"required,email"`
	Relationship string `gorm:"-" json:"relationship,omitempty" validate:"required"` // Viene en el JSON unido de PatientContact
}

// PatientContact - Relación Paciente-Contacto (n:m con info extra)
type PatientContact struct {
	PatientID    string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	ContactID    string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	Relationship string `gorm:"type:varchar(50)" json:"relationship"`

	Patient          Patient          `gorm:"foreignKey:PatientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	EmergencyContact EmergencyContact `gorm:"foreignKey:ContactID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// PatientCondition - Relación Paciente-Enfermedad
type PatientCondition struct {
	PatientID   string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	ConditionID string `gorm:"primaryKey;type:varchar(50)" json:"-"`
	Diagnostic  string `gorm:"type:text" json:"diagnostics"`

	Patient          Patient          `gorm:"foreignKey:PatientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	MedicalCondition MedicalCondition `gorm:"foreignKey:ConditionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
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
	ID             string  `gorm:"primaryKey;type:varchar(50)" json:"wearableId"`
	MacAddress     string  `gorm:"type:varchar(20);unique" json:"macAddress" validate:"required,mac"`
	BatteryLevel   int     `json:"batteryLevel" validate:"min=0,max=100"`
	BatteryVoltage float64 `json:"batteryVoltage"`
	IsCharging     bool    `json:"isCharging"`
	IsActive       bool    `json:"isActive"`
}

func (w *Wearable) UnmarshalJSON(data []byte) error {
	// Intenta decodificar como string (ej: "uuid-device")
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		w.ID = id
		return nil
	}

	// Si no es string, decodifica como struct original usando un alias
	type Alias Wearable
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*w = Wearable(aux)
	return nil
}

// PatientWearable - Asignación
type PatientWearable struct {
	PatientID    string    `gorm:"primaryKey;type:varchar(50)" json:"-"`
	WearableID   string    `gorm:"primaryKey;type:varchar(50)" json:"-"`
	AssignedDate time.Time `gorm:"type:timestamp" json:"-"`

	Wearable Wearable `gorm:"foreignKey:WearableID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Patient - Paciente
type Patient struct {
	ID          string    `gorm:"primaryKey;type:varchar(50)" json:"patientId"`
	FirstName   string    `gorm:"type:varchar(50)" json:"firstName" validate:"required,min=2"`
	LastName    string    `gorm:"type:varchar(50)" json:"lastName" validate:"required,min=2"`
	DateOfBirth time.Time `gorm:"type:timestamp" json:"dateOfBirth" validate:"required"`
	RoomID      string    `gorm:"type:varchar(50)" json:"RoomId,omitempty" validate:"required"` // Para el payload de Creación

	Room              *Room              `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"Room,omitempty"`
	PatientConditions []PatientCondition `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PatientContacts   []PatientContact   `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PatientWearables  []PatientWearable  `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Para coincidir con el payload exacto esperado (estos campos pueden poblarse manualmente)
	Allergies        []MedicalCondition `gorm:"-" json:"Allergies"`
	Diseases         []MedicalCondition `gorm:"-" json:"Diseases"`
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
	PatientID   string     `gorm:"type:varchar(50) REFERENCES patients(id) ON DELETE CASCADE;index" json:"patientId"`
	WearableID  string     `gorm:"type:varchar(50) REFERENCES wearables(id) ON DELETE CASCADE;index" json:"wearableId"`
	ResolvedAt  *time.Time `gorm:"type:timestamp" json:"resolvedAt,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp" json:"createdAt,omitempty"`
	AlertStatus string     `gorm:"type:varchar(30)" json:"alertStatus"`
	AlertLevel  string     `gorm:"type:varchar(30)" json:"alertLevel"`
	AlertType   string     `gorm:"type:varchar(50) REFERENCES alert_types(id) ON DELETE CASCADE" json:"alertType"` // PDF usa alertType en lugar de alertTypeId
	BPM         int        `gorm:"type:integer" json:"bpm,omitempty"`
	AlertCode   string     `gorm:"type:varchar(50)" json:"alertCode,omitempty"`
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
