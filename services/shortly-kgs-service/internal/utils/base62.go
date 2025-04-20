package utils

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encode(n int64) string {
	if n == 0 {
		return padLeft(string(charset[0]), 6)
	}

	result := ""
	for n > 0 {
		remainder := n % 62
		result = string(charset[remainder]) + result
		n /= 62
	}

	return padLeft(result, 6)
}

func padLeft(s string, length int) string {
	for len(s) < length {
		s = string(charset[0]) + s
	}
	return s
}
