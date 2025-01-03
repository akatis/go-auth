package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Auth struct {
	Header              string
	Payload             string
	JwtSecretKey        []byte
	AccessToken         string
	RedisConn           *redis.Conn
	EndPointPermissions map[string]int
}

func New(config *Config) *Auth {
	//config.init()

	var auth Auth

	auth.JwtSecretKey = []byte(config.JwtSecretKey)

	// Redis: connection
	auth.RedisConn = redis.NewClient(&redis.Options{
		Addr:         config.Redis.RedisAddr,
		Password:     config.Redis.RedisPass,
		DB:           config.Redis.RedisDb,
		PoolSize:     config.Redis.PoolSize,
		MaxIdleConns: config.Redis.MaxIdleConns,
		MinIdleConns: config.Redis.MinIdleConns,
		DialTimeout:  config.Redis.DialTimeout,
	}).Conn()

	//auth.RedisConn = client.Conn()

	auth.EndPointPermissions = config.EndpointPermissions

	return &auth
}

// Create AccessToken
func (a *Auth) CreateAccessToken(uuid string, roles []int, shopId, companyId *int) string {
	var header HeaderConfig
	header.Alg = "HS256"
	header.Typ = "JWT"

	var payload = PayloadConfig{
		Uuid:      uuid,
		Roles:     roles,
		ShopID:    shopId,
		CompanyID: companyId,
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	// JSON
	jsonHeader, _ := json.Marshal(header)
	jsonPayload, _ := json.Marshal(payload)

	// BASE64
	a.Header = base64.RawURLEncoding.EncodeToString(jsonHeader)
	a.Payload = base64.RawURLEncoding.EncodeToString(jsonPayload)

	headerPayload := a.Header + "." + a.Payload

	hasher := hmac.New(sha256.New, a.JwtSecretKey)
	hasher.Write([]byte(headerPayload))
	signature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	token := a.Header + "." + a.Payload + "." + signature

	a.AccessToken = token
	return token
}

// Token Verify
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

func (a *Auth) GetUUID(authHeader string) (string, error) {
	var uuid string

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

	var payload PayloadConfig
	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	err := json.Unmarshal(payloadBase64, &payload)
	if err != nil {
		return "", errors.New(err.Error())
	}

	uuid = payload.Uuid

	return uuid, nil
}

func (a *Auth) GetShopID(ctx *fiber.Ctx) (int, error) {
	var shopID int

	authHeader := string(ctx.Request().Header.Peek("Authorization"))
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

	var payload PayloadConfig
	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	err := json.Unmarshal(payloadBase64, &payload)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	if payload.ShopID == nil {
		return 0, errors.New("shopID is nil")
	}

	shopID = *payload.ShopID

	return shopID, nil

}

func (a *Auth) GetCompanyID(ctx *fiber.Ctx) (int, error) {
	var companyID int

	authHeader := string(ctx.Request().Header.Peek("Authorization"))
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

	var payload PayloadConfig
	payloadBase64, _ := base64.RawURLEncoding.DecodeString(authTokenParts[1])
	err := json.Unmarshal(payloadBase64, &payload)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	if payload.CompanyID == nil {
		return 0, errors.New("companyID is nil")
	}

	companyID = *payload.CompanyID

	return companyID, nil
}

// REDIS TRANSACTIONS
// Add user session to Redis
func (a *Auth) AddToRedis(uuid, userAgent string) error {
	//conn := a.RedisClient.Conn()
	//defer conn.Close()
	conn := a.RedisConn

	split := strings.Split(a.AccessToken, ".")
	payload := split[1]

	ctx := context.Background()

	_, err := conn.HSet(ctx, uuid, payload, userAgent).Result()
	if err != nil {
		return errors.New(err.Error())
	}
	return err
}

// Checking user session from Redis
func (a *Auth) CheckFromRedis(uuid string) error {
	//conn := a.RedisClient.Conn()
	//defer conn.Close()
	conn := a.RedisConn

	ctx := context.Background()

	_, err := conn.HGet(ctx, uuid, a.Payload).Result()
	if err != nil {
		return errors.New(err.Error())
	}
	return err
}

// Token expired, delete session from Redis
func (a *Auth) DeleteFromRedis(authPayload string) error {

	jwtPayload, _ := base64.RawURLEncoding.DecodeString(authPayload)

	var payload PayloadConfig
	err := json.Unmarshal(jwtPayload, &payload)
	if err != nil {
		return errors.New(err.Error())
	}

	//conn := a.RedisClient.Conn()
	//defer conn.Close()
	conn := a.RedisConn

	err = conn.HDel(context.Background(), payload.Uuid, authPayload).Err()
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Delete all sessions of deregistered user from Redis
func (a *Auth) DeleteKeyFromRedis(uuid string) error {
	conn := a.RedisConn

	err := conn.Del(context.Background(), uuid).Err()
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// AUTHORIZATION AND AUTHENTICATION MIDDLEWARE
func (a *Auth) Middleware(ctx *fiber.Ctx) error {

	var response Response

	authHeader := string(ctx.Request().Header.Peek("Authorization"))
	if authHeader == "" {
		response.Message = "1-Invalid token."
		return response.HttpResponse(ctx, 401)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Message = "2-Malformed token."
		return response.HttpResponse(ctx, 403)
	}

	a.AccessToken = strings.TrimPrefix(authHeader, "Bearer ")
	tokenParts := strings.Split(a.AccessToken, ".")
	if len(tokenParts) != 3 {
		response.Message = "3-Malformed token."
		return response.HttpResponse(ctx, 403)
	}
	a.Header = tokenParts[0]
	a.Payload = tokenParts[1]

	jwtPayload, _ := base64.RawURLEncoding.DecodeString(a.Payload)

	var payload PayloadConfig
	err := json.Unmarshal(jwtPayload, &payload)
	if err != nil {
		response.Message = "4-Invalid token"
		return response.HttpResponse(ctx, 403)
	}

	err = a.TokenVerify(tokenParts[2])
	if err != nil {
		response.Message = "5-Invalid token"
		return response.HttpResponse(ctx, 403)
	}

	if payload.ExpiresAt < time.Now().Unix() {
		response.Message = "6-Token expired"
		return response.HttpResponse(ctx, 403)
	}
	err = a.CheckFromRedis(payload.Uuid)
	if err != nil {
		response.Message = "User session not found."
		return response.HttpResponse(ctx, 403)
	}

	// CHECK USER PERMISSION

	pathPermission := a.EndPointPermissions[ctx.Path()]
	hasPermission := PermissionsContains(payload.Roles, pathPermission)

	if !hasPermission && pathPermission != 999 {
		response.Message = "You don't have permission to access this end-point"
		return response.HttpResponse(ctx, 401)
	}

	ctx.Next()

	return nil
}
