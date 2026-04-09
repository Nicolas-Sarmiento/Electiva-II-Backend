package main

import (
	"log"
	"net/http"

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
	customMiddleware "ancianato-backend/internal/delivery/http/middleware"
)

func main() {
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

	// 3. Inicializar Router (Chi)
	r := chi.NewRouter()

	// 4. Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(customMiddleware.KafkaAuditMiddleware)

	// 5. Rutas de prueba
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Servidor funcionando en modo OK"))
	})

	// 6. Inyección de dependencias e Inicialización de Rutas REST
	deliveryHttp.NewAuthHandler(r) // Agrega /login libre de Auth

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

	// Configuramos sub-enrutador protegido
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)

		// Todos los roles (nurse, administrator) pueden ver pacientes, alertas y configurar alertas
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RoleMiddleware("nurse", "administrator"))
			
			// Alertas (Nurse lo ocupa para verlas y resolverlas)
			deliveryHttp.NewAlertHandler(r, alertUseCase)
			// Pacientes (Ver pacientes para las alertas)
			deliveryHttp.NewPatientHandler(r, patientUseCase)
		})

		// Solo administrador puede gestionar catálogo (Room, Devices, AlertTypes, Shifts)
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RoleMiddleware("administrator"))

			deliveryHttp.NewDeviceHandler(r, deviceUseCase)
			deliveryHttp.NewRoomHandler(r, roomUseCase)
			deliveryHttp.NewAlertTypeHandler(r, alertTypeUseCase)
			deliveryHttp.NewShiftHandler(r, shiftUseCase)
		})
	})

	// 7. Levantar Servidor
	port := ":" + cfg.Port
	log.Printf("El servidor del ancianato ha iniciado en http://localhost%s\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
