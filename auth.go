package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
	"time"
)

type Auth struct {
	Header              string
	Payload             string
	JwtSecretKey        []byte
	AccessToken         string
	RedisClient         *redis.Client
	EndPointPermissions map[string]int
}

func NewConfig(config *Config) *Auth {
	config.init()

	var auth Auth

	auth.JwtSecretKey = []byte(config.JwtSecretKey)

	// Redis: connection
	auth.RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.Redis.RedisAddr,
		Password:     config.Redis.RedisPass,
		DB:           config.Redis.RedisDb,
		PoolSize:     config.Redis.PoolSize,
		MaxIdleConns: config.Redis.MaxIdleConns,
		MinIdleConns: config.Redis.MinIdleConns,
		DialTimeout:  config.Redis.DialTimeout,
	})

	//auth.RedisConn = client.Conn()

	auth.EndPointPermissions = config.EndpointPermissions

	return &auth
}

// Create AccessToken
func (a *Auth) CreateAccessToken(uuid string, roles []int) string {
	var payload PayloadConfig
	payload.Uuid = uuid
	payload.Roles = roles
	payload.ExpiresAt = time.Now().Add(time.Minute * 30).Unix()
	payload.IssuedAt = time.Now().Unix()

	// JSON
	jsonHeader, _ := json.Marshal(a.Header)
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

// REDIS TRANSACTIONS
func (a *Auth) AddToRedis(uuid, userAgent string) error {
	conn := a.RedisClient.Conn()
	defer conn.Close()

	split := strings.Split(a.AccessToken, ".")
	payload := split[1]

	ctx := context.Background()

	_, err := conn.HSet(ctx, uuid, payload, userAgent).Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	return err
}

func (a *Auth) CheckFromRedis(uuid string) error {
	conn := a.RedisClient.Conn()
	defer conn.Close()

	ctx := context.Background()

	_, err := conn.HGet(ctx, uuid, a.Payload).Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	return err
}
func (a *Auth) DeleteFromRedis(authPayload string) error {

	jwtPayload, _ := base64.RawURLEncoding.DecodeString(authPayload)

	var payload PayloadConfig
	err := json.Unmarshal(jwtPayload, &payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn := a.RedisClient.Conn()
	defer conn.Close()

	err = conn.HDel(context.Background(), payload.Uuid, authPayload).Err()
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

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

	// CHECK USER PERMISSION
	pathPermission := a.EndPointPermissions[ctx.Path()]
	hasPermission := PermissionsContains(payload.Roles, pathPermission)

	if !hasPermission && pathPermission != 999 {
		response.Message = "You don't have permission to access this end-point"
		return response.HttpResponse(ctx, 401)
	}

	err = a.CheckFromRedis(payload.Uuid)
	if err != nil {
		response.Message = "User session not found."
		return response.HttpResponse(ctx, 403)
	}

	ctx.Next()

	return nil
}
