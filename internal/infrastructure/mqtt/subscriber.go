package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"ancianato-backend/internal/domain"
	ws "ancianato-backend/internal/infrastructure/websocket"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// AlertMessage es el JSON que envía el ESP32 al tópico de alertas
type AlertMessage struct {
	AlertCode  string `json:"alertCode"`
	AlertLevel string `json:"alertLevel"`
	BPM        int    `json:"bpm,omitempty"`
}

// BatteryMessage es el JSON que envía el ESP32 al tópico de batería
type BatteryMessage struct {
	Level    int     `json:"level"`
	Voltage  float64 `json:"voltage"`
	Charging bool    `json:"charging"`
}

// StatusMessage es el JSON de estado (online/offline vía LWT)
type StatusMessage struct {
	Status string `json:"status"`
}

// MQTTSubscriber gestiona la conexión MQTT y procesa mensajes IoT
type MQTTSubscriber struct {
	client        pahomqtt.Client
	deviceRepo    domain.DeviceRepository
	alertRepo     domain.AlertRepository
	alertTypeRepo domain.AlertTypeRepository
	patientRepo   domain.PatientRepository
}

// NewMQTTSubscriber crea e inicializa el subscriber MQTT
func NewMQTTSubscriber(
	brokerURL string,
	deviceRepo domain.DeviceRepository,
	alertRepo domain.AlertRepository,
	alertTypeRepo domain.AlertTypeRepository,
	patientRepo domain.PatientRepository,
) *MQTTSubscriber {
	sub := &MQTTSubscriber{
		deviceRepo:    deviceRepo,
		alertRepo:     alertRepo,
		alertTypeRepo: alertTypeRepo,
		patientRepo:   patientRepo,
	}

	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("ancianato-backend-" + uuid.New().String()[:8])
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
		log.Printf("MQTT: Conexión perdida: %v. Reintentando...", err)
	})
	opts.SetOnConnectHandler(func(c pahomqtt.Client) {
		log.Println("MQTT: Conectado al broker. Suscribiéndose a tópicos IoT...")
		sub.subscribe()
	})

	sub.client = pahomqtt.NewClient(opts)
	return sub
}

// Connect establece la conexión con el broker MQTT
func (s *MQTTSubscriber) Connect() error {
	token := s.client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("error conectando al broker MQTT: %w", err)
	}
	log.Println("MQTT Subscriber conectado exitosamente al broker")
	return nil
}

// subscribe registra los tópicos de escucha
func (s *MQTTSubscriber) subscribe() {
	// Suscribirse a todos los dispositivos con wildcard
	topics := map[string]byte{
		"ancianato/device/+/alert":   1, // QoS 1 para alertas (at-least-once)
		"ancianato/device/+/battery": 0, // QoS 0 para batería (best-effort)
		"ancianato/device/+/status":  1, // QoS 1 para estado
	}

	token := s.client.SubscribeMultiple(topics, func(client pahomqtt.Client, msg pahomqtt.Message) {
		go s.handleMessage(msg)
	})
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("MQTT: Error al suscribirse: %v", err)
		return
	}
	log.Println("MQTT: Suscrito exitosamente a tópicos ancianato/device/+/alert, battery, status")
}

// handleMessage enruta cada mensaje al handler correspondiente
func (s *MQTTSubscriber) handleMessage(msg pahomqtt.Message) {
	topic := msg.Topic()
	parts := strings.Split(topic, "/")

	// Formato esperado: ancianato/device/{mac}/alert|battery|status
	if len(parts) != 4 {
		log.Printf("MQTT: Tópico inesperado: %s", topic)
		return
	}

	mac := parts[2]
	messageType := parts[3]

	log.Printf("MQTT: Mensaje recibido [%s] de MAC %s: %s", messageType, mac, string(msg.Payload()))

	switch messageType {
	case "alert":
		s.handleAlert(mac, msg.Payload())
	case "battery":
		s.handleBattery(mac, msg.Payload())
	case "status":
		s.handleStatus(mac, msg.Payload())
	default:
		log.Printf("MQTT: Tipo de mensaje desconocido: %s", messageType)
	}
}

// handleAlert procesa una alerta recibida del wearable
func (s *MQTTSubscriber) handleAlert(mac string, payload []byte) {
	var alertMsg AlertMessage
	if err := json.Unmarshal(payload, &alertMsg); err != nil {
		log.Printf("MQTT: Error deserializando alerta de MAC %s: %v", mac, err)
		return
	}

	ctx := context.Background()

	// 1. Buscar dispositivo por MAC
	device, err := s.deviceRepo.GetByMacAddress(ctx, mac)
	if err != nil {
		log.Printf("MQTT: Dispositivo con MAC %s no registrado en el sistema", mac)
		return
	}

	// 2. Buscar paciente asignado al dispositivo
	patientID, err := s.findPatientByDevice(ctx, device.ID)
	if err != nil {
		log.Printf("MQTT: No se encontró paciente asignado al dispositivo %s (MAC: %s)", device.ID, mac)
		return
	}

	// 3. Buscar tipo de alerta por código
	alertTypeID, err := s.findAlertTypeByCode(ctx, alertMsg.AlertCode)
	if err != nil {
		log.Printf("MQTT: Tipo de alerta '%s' no encontrado. Creando tipo genérico.", alertMsg.AlertCode)
		// Crear el tipo de alerta si no existe
		newType := &domain.AlertType{
			ID:   uuid.New().String(),
			Name: alertMsg.AlertCode,
			Code: alertMsg.AlertCode,
		}
		if err := s.alertTypeRepo.Create(ctx, newType); err != nil {
			log.Printf("MQTT: Error creando tipo de alerta: %v", err)
			return
		}
		alertTypeID = newType.ID
	}

	// 3.5 Throttle: Si ya existe una alerta Activa o En Revisión para este paciente del mismo tipo, omitir
	hasActive, err := s.alertRepo.HasActiveAlert(ctx, patientID, alertTypeID)
	if err == nil && hasActive {
		log.Printf("MQTT: ⚠️ Alerta '%s' ya está activa/revisión para el paciente %s. Omitiendo duplicado.", alertMsg.AlertCode, patientID)
		return
	}

	// 4. Crear la alerta en la base de datos
	alert := &domain.Alert{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		WearableID:  device.ID,
		AlertType:   alertTypeID,
		AlertLevel:  alertMsg.AlertLevel,
		AlertStatus: "Activa",
		CreatedAt:   time.Now(),
		BPM:         alertMsg.BPM,
		AlertCode:   alertMsg.AlertCode,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		log.Printf("MQTT: Error guardando alerta en DB: %v", err)
		return
	}

	log.Printf("MQTT: ✅ Alerta creada [%s] para paciente %s desde dispositivo %s (MAC: %s)",
		alertMsg.AlertCode, patientID, device.ID, mac)

	// 5. Notificar a todos los clientes WebSocket
	ws.Broadcast("new_alert", map[string]interface{}{
		"alertId":     alert.ID,
		"patientId":   alert.PatientID,
		"wearableId":  alert.WearableID,
		"alertType":   alert.AlertType,
		"alertLevel":  alert.AlertLevel,
		"alertStatus": alert.AlertStatus,
		"createdAt":   alert.CreatedAt,
		"alertCode":   alertMsg.AlertCode,
		"bpm":         alertMsg.BPM,
	})
}

// handleBattery procesa una actualización de batería
func (s *MQTTSubscriber) handleBattery(mac string, payload []byte) {
	var batteryMsg BatteryMessage
	if err := json.Unmarshal(payload, &batteryMsg); err != nil {
		log.Printf("MQTT: Error deserializando batería de MAC %s: %v", mac, err)
		return
	}

	ctx := context.Background()

	// 1. Buscar dispositivo por MAC
	device, err := s.deviceRepo.GetByMacAddress(ctx, mac)
	if err != nil {
		log.Printf("MQTT: Dispositivo con MAC %s no registrado", mac)
		return
	}

	// 2. Actualizar nivel de batería
	device.BatteryLevel = batteryMsg.Level
	device.BatteryVoltage = batteryMsg.Voltage
	device.IsCharging = batteryMsg.Charging
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		log.Printf("MQTT: Error actualizando batería del dispositivo %s: %v", device.ID, err)
		return
	}

	log.Printf("MQTT: 🔋 Batería actualizada [%s] MAC %s → %d%% (%.2fV, cargando: %t)", device.ID, mac, batteryMsg.Level, batteryMsg.Voltage, batteryMsg.Charging)

	// 3. Notificar por WebSocket
	ws.Broadcast("battery_update", map[string]interface{}{
		"wearableId":     device.ID,
		"macAddress":     mac,
		"batteryLevel":   batteryMsg.Level,
		"batteryVoltage": batteryMsg.Voltage,
		"isCharging":     batteryMsg.Charging,
	})
}

// handleStatus procesa cambios de estado (online/offline via LWT)
func (s *MQTTSubscriber) handleStatus(mac string, payload []byte) {
	var statusMsg StatusMessage
	if err := json.Unmarshal(payload, &statusMsg); err != nil {
		log.Printf("MQTT: Error deserializando status de MAC %s: %v", mac, err)
		return
	}

	ctx := context.Background()

	device, err := s.deviceRepo.GetByMacAddress(ctx, mac)
	if err != nil {
		log.Printf("MQTT: Dispositivo con MAC %s no registrado", mac)
		return
	}

	isActive := statusMsg.Status == "online"
	device.IsActive = isActive
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		log.Printf("MQTT: Error actualizando estado del dispositivo %s: %v", device.ID, err)
		return
	}

	log.Printf("MQTT: 📡 Estado del dispositivo [%s] MAC %s → %s", device.ID, mac, statusMsg.Status)

	ws.Broadcast("device_status", map[string]interface{}{
		"wearableId": device.ID,
		"macAddress": mac,
		"isActive":   isActive,
	})
}

// findPatientByDevice busca el paciente asignado a un dispositivo
func (s *MQTTSubscriber) findPatientByDevice(ctx context.Context, deviceID string) (string, error) {
	// Obtenemos todos los pacientes y buscamos el que tenga el wearable asignado
	patients, err := s.patientRepo.GetAll(ctx)
	if err != nil {
		return "", err
	}

	for _, p := range patients {
		patient, err := s.patientRepo.GetByID(ctx, p.ID)
		if err != nil {
			continue
		}
		for _, pw := range patient.PatientWearables {
			if pw.WearableID == deviceID {
				return patient.ID, nil
			}
		}
	}

	return "", fmt.Errorf("no se encontró paciente con dispositivo %s", deviceID)
}

// findAlertTypeByCode busca un tipo de alerta por su código
func (s *MQTTSubscriber) findAlertTypeByCode(ctx context.Context, code string) (string, error) {
	alertTypes, err := s.alertTypeRepo.GetAll(ctx)
	if err != nil {
		return "", err
	}

	for _, at := range alertTypes {
		if strings.EqualFold(at.Code, code) {
			return at.ID, nil
		}
	}

	return "", fmt.Errorf("tipo de alerta con código '%s' no encontrado", code)
}

// Close cierra la conexión MQTT
func (s *MQTTSubscriber) Close() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(1000) // Esperar 1 segundo para desconectar limpiamente
		log.Println("MQTT Subscriber desconectado")
	}
}
