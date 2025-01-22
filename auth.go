package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	Header              string
	Payload             string
	JwtSecretKey        []byte
	AccessToken         string
	EndPointPermissions map[string]int
}

func New(config *Config) *Auth {
	return &Auth{
		JwtSecretKey:        []byte(config.JwtSecretKey),
		EndPointPermissions: config.EndpointPermissions,
	}
}

// CreateAccessToken generates a new JWT token with the given user information
func (a *Auth) CreateAccessToken(uuid string, roles []int, shopId, companyId *int) string {
	header := HeaderConfig{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := PayloadConfig{
		Uuid:      uuid,
		Roles:     roles,
		ShopID:    shopId,
		CompanyID: companyId,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	// Serialize header and payload to JSON
	jsonHeader, _ := json.Marshal(header)
	jsonPayload, _ := json.Marshal(payload)

	// Encode header and payload using Base64 URL encoding
	a.Header = base64.RawURLEncoding.EncodeToString(jsonHeader)
	a.Payload = base64.RawURLEncoding.EncodeToString(jsonPayload)

	headerPayload := a.Header + "." + a.Payload

	// Generate signature using HMAC with SHA-256
	hasher := hmac.New(sha256.New, a.JwtSecretKey)
	hasher.Write([]byte(headerPayload))
	signature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	token := headerPayload + "." + signature
	a.AccessToken = token

	return token
}

// TokenVerify verifies the token signature
func (a *Auth) TokenVerify(signature string) error {
	sign := a.Header + "." + a.Payload

	hasher := hmac.New(sha256.New, a.JwtSecretKey)
	hasher.Write([]byte(sign))
	expectedSignature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	if signature != expectedSignature {
		return errors.New("failed to verify token signature")
	}

	return nil
}

// GetUUID extracts the UUID from the token
func (a *Auth) GetUUID(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("invalid token")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("malformed token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	authTokenParts := strings.Split(token, ".")
	if len(authTokenParts) != 3 {
		return "", errors.New("malformed token")
	}

	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	var payload PayloadConfig
	if err := json.Unmarshal(payloadBase64, &payload); err != nil {
		return "", errors.New(err.Error())
	}

	return payload.Uuid, nil
}

// GetShopID extracts the ShopID from the token
func (a *Auth) GetShopID(authHeader string) (int, error) {
	if authHeader == "" {
		return 0, errors.New("invalid token")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("malformed token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	authTokenParts := strings.Split(token, ".")
	if len(authTokenParts) != 3 {
		return 0, errors.New("malformed token")
	}

	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	var payload PayloadConfig
	if err := json.Unmarshal(payloadBase64, &payload); err != nil {
		return 0, errors.New(err.Error())
	}

	if payload.ShopID == nil {
		return 0, errors.New("shopID is nil")
	}

	return *payload.ShopID, nil
}

// GetCompanyID extracts the CompanyID from the token
func (a *Auth) GetCompanyID(authHeader string) (int, error) {
	if authHeader == "" {
		return 0, errors.New("invalid token")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("malformed token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	authTokenParts := strings.Split(token, ".")
	if len(authTokenParts) != 3 {
		return 0, errors.New("malformed token")
	}

	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	var payload PayloadConfig
	if err := json.Unmarshal(payloadBase64, &payload); err != nil {
		return 0, errors.New(err.Error())
	}

	if payload.CompanyID == nil {
		return 0, errors.New("companyID is nil")
	}

	return *payload.CompanyID, nil
}

// Middleware performs authentication and authorization
func (a *Auth) Middleware(ctx *fiber.Ctx) error {
	var response Response

	authHeader := ctx.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		response.Message = "invalid token"
		return response.HttpResponse(ctx, fiber.StatusUnauthorized)
	}

	a.AccessToken = strings.TrimPrefix(authHeader, "Bearer ")
	tokenParts := strings.Split(a.AccessToken, ".")
	if len(tokenParts) != 3 {
		response.Message = "malformed token"
		return response.HttpResponse(ctx, fiber.StatusBadRequest)
	}

	a.Header, a.Payload = tokenParts[0], tokenParts[1]

	jwtPayload, _ := base64.RawURLEncoding.DecodeString(a.Payload)
	var payload PayloadConfig
	if err := json.Unmarshal(jwtPayload, &payload); err != nil {
		response.Message = "invalid payload"
		return response.HttpResponse(ctx, fiber.StatusForbidden)
	}

	if err := a.TokenVerify(tokenParts[2]); err != nil {
		response.Message = "invalid signature"
		return response.HttpResponse(ctx, fiber.StatusForbidden)
	}

	if payload.ExpiresAt < time.Now().Unix() {
		response.Message = "token expired"
		return response.HttpResponse(ctx, fiber.StatusForbidden)
	}

	// Check user permission by matching endpoint with required permission
	requestedPath := ctx.Path()
	matchedPermission, matched := matchPathWithPermission(requestedPath, a.EndPointPermissions)
	if !matched {
		response.Message = "access denied"
		return response.HttpResponse(ctx, fiber.StatusForbidden)
	}

	// Check if user's roles include the required permission
	if PermissionsContains(payload.Roles, matchedPermission) {
		return ctx.Next()
	}

	response.Message = "access denied"
	return response.HttpResponse(ctx, fiber.StatusForbidden)
}

// matchPathWithPermission matches the requested path with defined permissions
func matchPathWithPermission(requestedPath string, dynamicPermissions map[string]int) (int, bool) {
	for pattern, permission := range dynamicPermissions {
		if matchRoute(pattern, requestedPath) {
			return permission, true
		}
	}
	return 0, false
}

// matchRoute compares a route pattern with an actual path
//   - routeDef: e.g. /api/test/:id
//   - path: e.g. /api/test/453
func matchRoute(routeDef, path string) bool {
	defParts := strings.Split(strings.Trim(routeDef, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Return false early if segment lengths don't match
	if len(defParts) != len(pathParts) {
		return false
	}

	for i, seg := range defParts {
		// If segment is a parameter (e.g., :id, :xyz)
		if strings.HasPrefix(seg, ":") {
			// // Ensure parameter segment is numeric
			// if !isNumeric(pathParts[i]) {
			// 	return false
			// }
			continue
		}
		// If not a parameter, check for exact match
		if seg != pathParts[i] {
			return false
		}
	}

	return true
}

// isNumeric checks if a string contains only numeric characters
func isNumeric(str string) bool {
	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
