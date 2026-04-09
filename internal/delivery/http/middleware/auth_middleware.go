package middleware

import (
	"context"
	"net/http"
	"strings"

	"ancianato-backend/internal/infrastructure/auth"
)

// AuthMiddleware extrae el JWT, extrae los Claims y los pone en el ctx.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Falta cabecera Authorization", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Formato Authorization inválido. Usa 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := auth.ParseClaims(tokenString)
		if err != nil {
			http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Save username y roles en context
		ctx := context.WithValue(r.Context(), "username", claims.PreferredUsername)

		var allRoles []string
		allRoles = append(allRoles, claims.RealmAccess.Roles...)
		for _, clientRoles := range claims.ResourceAccess {
			allRoles = append(allRoles, clientRoles.Roles...)
		}

		ctx = context.WithValue(ctx, "roles", allRoles)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware exige que el usuario tenga uno de los roles permitidos.
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rolesInterf := r.Context().Value("roles")
			if rolesInterf == nil {
				http.Error(w, "Roles no encontrados en el contexto (¿falta AuthMiddleware?)", http.StatusForbidden)
				return
			}

			userRoles := rolesInterf.([]string)
			hasRole := false
			for _, allowed := range allowedRoles {
				for _, userRole := range userRoles {
					if strings.EqualFold(userRole, allowed) {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				http.Error(w, "No tienes permiso para acceder a este recurso", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
