package authTest

const (
	ep  = "/api/test/:id/user/:name"
	epO = "/api/testo"
)

// Permission Codes
const (
	AllUsers = 999
)

// Set permission codes to end-points
var EndPointPermissions = map[string]int{
	ep:  AllUsers,
	epO: AllUsers,
}
