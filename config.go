package auth

import "time"

const Admin = 1

type Config struct {
	Header              HeaderConfig
	Payload             PayloadConfig
	Redis               RedisConfig
	JwtSecretKey        string
	EndpointPermissions map[string]int
}

func (cfg *Config) init() {
	if cfg.Header.Typ == "" {
		cfg.Header.Typ = "JWT"
	}
	if cfg.Header.Alg == "" {
		cfg.Header.Alg = "HS256"
	}

	if cfg.Payload.ExpiresAt == 0 {
		cfg.Payload.ExpiresAt = time.Now().Add(time.Minute * 30).Unix()
	}
	if cfg.Payload.IssuedAt == 0 {
		cfg.Payload.IssuedAt = time.Now().Unix()
	}

	if cfg.JwtSecretKey == "" {
		panic("Auth Configuration init: JWT secret key is required.")
	}

	/*if cfg.EndpointPermissions == nil {

	}*/
}
