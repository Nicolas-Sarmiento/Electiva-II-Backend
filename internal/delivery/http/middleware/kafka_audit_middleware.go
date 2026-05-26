package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"ancianato-backend/internal/infrastructure/auth"
	mykafka "ancianato-backend/internal/infrastructure/kafka"
)

// responseRecorder envuelve a http.ResponseWriter para poder capturar el estatus de la respuesta
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func realIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func determineTopic(r *http.Request, statusCode int) string {
	path := r.URL.Path

	if strings.HasPrefix(path, "/login") {
		return "topic-autenticacion"
	}

	// Si es problema de roles o de autenticación lo ponemos en "acceso"
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return "topic-acceso"
	}

	// Servicios (ej. /patients -> topic-servicio-patients)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		return "topic-servicio-" + parts[0]
	}

	return "topic-general"
}

// KafkaAuditMiddleware intercepta cada petición y envía un evento a Kafka
func KafkaAuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Bypassear el responseRecorder para solicitudes WebSocket para permitir hijacking/upgrade
		if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
			next.ServeHTTP(w, r)
			return
		}

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		user := "anonymous"
		ip := realIP(r)

		// Extraemos username anticipadamente
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				claims, err := auth.ParseClaims(parts[1])
				if err == nil {
					user = claims.PreferredUsername
				}
			}
		} else if r.URL.Path == "/login" && r.Method == http.MethodPost {
			// Intento de leer el username del body (JSON)
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restauramos el body
			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				if uname, ok := data["username"].(string); ok {
					user = uname
				}
			}
		}

		next.ServeHTTP(rec, r)

		actionStr := r.Method + " " + r.URL.Path
		topic := determineTopic(r, rec.statusCode)

		event := mykafka.AuditEvent{
			Action: actionStr,
			User:   user,
			IP:     ip,
			Topic:  topic,
			Time:   time.Now(),
			Status: rec.statusCode,
		}

		mykafka.ProduceAuditEvent(event)
	})
}
