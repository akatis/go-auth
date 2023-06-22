package auth

func PermissionsContains(s []int, p int) bool {

	for _, v := range s {
		if v == 1 || v == p {
			return true
		}
	}

	return false
}
