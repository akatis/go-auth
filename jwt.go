package auth

type HeaderConfig struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type PayloadConfig struct {
	Uuid      string `json:"uuid"`
	Roles     []int  `json:"roles"`
	ShopID    *int   `json:"shop_id,omitempty"`
	CompanyID *int   `json:"company_id,omitempty"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

/*type AccessToken struct {
	Header    HeaderConfig
	Payload   PayloadConfig
	SecretKey string
}

type Token struct {
	AccessToken string
}*/
