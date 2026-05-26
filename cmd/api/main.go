package main

import (
	"log"
	"net/http"
	"time"

	"ancianato-backend/internal/config"
	deliveryHttp "ancianato-backend/internal/delivery/http"
	"ancianato-backend/internal/infrastructure/cache"
	"ancianato-backend/internal/infrastructure/validation"
	"ancianato-backend/internal/repository"
	"ancianato-backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ancianato-backend/internal/infrastructure/auth"
	mykafka "ancianato-backend/internal/infrastructure/kafka"
	mymqtt "ancianato-backend/internal/infrastructure/mqtt"
	ws "ancianato-backend/internal/infrastructure/websocket"
	customMiddleware "ancianato-backend/internal/delivery/http/middleware"
)

func main() {
	// 0. Configurar zona horaria (Bogotá, Colombia UTC-5)
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Printf("⚠️  No se pudo cargar zona horaria America/Bogota: %v. Usando UTC.", err)
	} else {
		time.Local = loc
		log.Println("🕐 Zona horaria configurada: America/Bogota (UTC-5)")
	}

	// 1. Cargar Configuración
	cfg := config.LoadConfig()

	// 1.5 Inicializar Herramientas (Caché y Validaciones)
	cache.InitCache()
	validation.InitValidator()

	// 2. Conectar a Base de Datos y lanzar AutoMigrate
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	// 2.5 Inicializar Servicios Externos
	mykafka.InitKafkaProducer([]string{cfg.KafkaBrokers})
	defer mykafka.Close()

	auth.InitKeycloak(cfg.KeycloakURL, cfg.KeycloakRealm, cfg.KeycloakClientID, cfg.KeycloakClientSecret)

	// 2.6 Inicializar WebSocket Hub (tiempo real para el frontend)
	ws.InitGlobalHub()

	// 3. Inicializar Router (Chi)
	r := chi.NewRouter()

	// 3.1 WebSocket endpoint (sin NINGÚN middleware para evitar romper el WebSocket Handshake/Hijack)
	r.Get("/ws", ws.HandleWebSocket)

	// 4. Sub-router para la API REST con middlewares
	api := chi.NewRouter()
	api.Use(middleware.Logger)
	api.Use(middleware.Recoverer)
	api.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	api.Use(customMiddleware.KafkaAuditMiddleware)

	// 5. Rutas de prueba en la API
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Servidor funcionando en modo OK"))
	})

	// 6. Inyección de dependencias e Inicialización de Rutas REST
	deliveryHttp.NewAuthHandler(api) // Agrega /login libre de Auth

	// Repositorios
	roomRepo := repository.NewRoomRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	alertTypeRepo := repository.NewAlertTypeRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	shiftRepo := repository.NewShiftRepository(db)

	// UseCases
	patientUseCase := usecase.NewPatientUseCase(patientRepo, roomRepo, deviceRepo)
	deviceUseCase := usecase.NewDeviceUseCase(deviceRepo)
	alertUseCase := usecase.NewAlertUseCase(alertRepo, patientRepo, deviceRepo, alertTypeRepo)
	roomUseCase := usecase.NewRoomUseCase(roomRepo)
	alertTypeUseCase := usecase.NewAlertTypeUseCase(alertTypeRepo)
	shiftUseCase := usecase.NewShiftUseCase(shiftRepo)

	// Configuramos sub-enrutador protegido dentro de la API REST
	api.Group(func(api chi.Router) {
		api.Use(customMiddleware.AuthMiddleware)

		// Todos los roles (nurse, administrator) pueden ver pacientes, alertas y configurar alertas
		api.Group(func(api chi.Router) {
			api.Use(customMiddleware.RoleMiddleware("nurse", "administrator"))
			
			// Alertas (Nurse lo ocupa para verlas y resolverlas)
			deliveryHttp.NewAlertHandler(api, alertUseCase)
			// Pacientes (Ver pacientes para las alertas)
			deliveryHttp.NewPatientHandler(api, patientUseCase)
		})

		// Solo administrador puede gestionar catálogo (Room, Devices, AlertTypes, Shifts)
		api.Group(func(api chi.Router) {
			api.Use(customMiddleware.RoleMiddleware("administrator"))

			deliveryHttp.NewDeviceHandler(api, deviceUseCase)
			deliveryHttp.NewRoomHandler(api, roomUseCase)
			deliveryHttp.NewAlertTypeHandler(api, alertTypeUseCase)
			deliveryHttp.NewShiftHandler(api, shiftUseCase)
		})
	})

	// Montamos las rutas de la API en el root router
	r.Mount("/", api)

	// 6.5 Inicializar MQTT Subscriber (conexión con los dispositivos wearables)
	if cfg.MQTTBroker != "" && cfg.MQTTBroker != "none" {
		mqttSubscriber := mymqtt.NewMQTTSubscriber(
			cfg.MQTTBroker,
			deviceRepo,
			alertRepo,
			alertTypeRepo,
			patientRepo,
		)
		if err := mqttSubscriber.Connect(); err != nil {
			log.Printf("⚠️  MQTT: No se pudo conectar al broker (%s). Los dispositivos IoT no enviarán datos en tiempo real. Error: %v", cfg.MQTTBroker, err)
		} else {
			defer mqttSubscriber.Close()
		}
	} else {
		log.Println("ℹ️  MQTT: Deshabilitado para esta instancia de backend.")
	}


	// 7. Levantar Servidor
	port := ":" + cfg.Port
	log.Printf("El servidor del ancianato ha iniciado en http://localhost%s\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
