package auth

func PermissionsContains(s []int, p int) bool {

	for _, v := range s {
		if v == Admin || v == p {
			return true
		}
	}
	if p == AllUser {
		return true
	}

	return false
}
