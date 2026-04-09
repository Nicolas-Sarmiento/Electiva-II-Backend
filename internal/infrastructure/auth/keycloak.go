package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// KeycloakConfig almacena la configuración de Keycloak
type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

var keycloakCfg KeycloakConfig

func InitKeycloak(baseURL, realm, clientID, clientSecret string) {
	keycloakCfg = KeycloakConfig{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// Login maneja la petición contra el token endpoint de Keycloak
func Login(username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakCfg.BaseURL, keycloakCfg.Realm)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", keycloakCfg.ClientID)
	data.Set("username", username)
	data.Set("password", password)

	// Opcional: Client Secret si es un cliente confidencial.
	if keycloakCfg.ClientSecret != "" {
		data.Set("client_secret", keycloakCfg.ClientSecret)
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("access_token no encontrado")
	}

	return token, nil
}

// CustomClaims de JWT para sacar roles facilmente
type CustomClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	PreferredUsername string `json:"preferred_username"`
	jwt.RegisteredClaims
}

// ValidateToken parses un token JWT verificando su firma
// Nota para probar: si Keycloak usa RS256, se debe obtener el Public Key del Realm.
// Para simplificar este entorno, Keycloak por defecto lo expone en sus JWKS.
// Por brevedad y para evitar depender del JWKS remoto en desarrollo local, parsearemos
// los claims asumiendo que el proxy del fronted o un API Gateway valida la firma,
// o realizamos un unverified parse.
// RECOMIENDO en prod usar una lib como nerzal/gocloak.
func ParseClaims(tokenString string) (*CustomClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, errors.New("claims inválidos")
}
