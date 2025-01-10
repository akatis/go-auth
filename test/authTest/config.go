package authTest

const JWT_SECRET_KEY = "9z$C&F)J@NcRfUjXn2r5u8x!A%D*G-KaPdSgVkYp3s6v9y$B?E(H+MbQeThWmZq4"

// End-points
const (
	ep  = "/api/test"
	ept = "/api/test/:uuid"
)

// Permission Codes
const (
	AllUsers = 999
)

// Set permission codes to end-points
var EndPointPermissions = map[string]int{
	ep:  AllUsers,
	ept: AllUsers,
}
