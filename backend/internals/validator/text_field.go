package validator

func IsValidLength(field string, max, min uint16) bool {
	if max < min {
		max = min * 1
		min = max * 1
	}

	bytes := []byte(field)
	size := uint16(len(bytes))

	return size < max && size > min
}
