package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

// AuditEvent es la estructura del mensaje que se va a enviar a Kafka
type AuditEvent struct {
	Action string    `json:"action"`
	User   string    `json:"user"`
	IP     string    `json:"ip"`
	Topic  string    `json:"-"`
	Time   time.Time `json:"time"`
	Status int       `json:"status"`
}

// InitKafkaProducer inicializa el writer de Kafka
func InitKafkaProducer(brokers []string) {
	writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		// Opciones de performance y reintentos:
		BatchTimeout: 10 * time.Millisecond,
		MaxAttempts:  5,
	}
	log.Printf("Kafka producer iniciado conectando a %v (topics dinámicos)", brokers)
}

// ProduceAuditEvent envía de forma asíncrona (fuego y olvido) a Kafka
func ProduceAuditEvent(event AuditEvent) {
	if writer == nil {
		log.Println("Kafka writer no inicializado")
		return
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error emparejando evento de Kafka: %v\n", err)
		return
	}

	msg := kafka.Message{
		Topic: event.Topic,
		Key:   []byte(event.User),
		Value: bytes,
		Time:  time.Now(),
	}

	// Ejecutamos en una goroutine para no bloquear el request http
	go func() {
		var err error
		// Reintento manual: damos tiempo a Kafka de crear el tópico si es la primera vez
		for i := 0; i < 3; i++ {
			err = writer.WriteMessages(context.Background(), msg)
			if err == nil {
				return // Éxito
			}
			time.Sleep(500 * time.Millisecond)
		}
		log.Printf("Fallo definitivo al enviar mensaje a kafka (tópico %s): %v\n", event.Topic, err)
	}()
}

// Close cierra limpiamente la conexión a Kafka
func Close() {
	if writer != nil {
		if err := writer.Close(); err != nil {
			log.Println("Error cerrando Kafka writer:", err)
		}
	}
}
