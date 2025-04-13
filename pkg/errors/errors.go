package errors

func ErrorResponse(code string, message ...string) map[string]interface{} {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	}
	return map[string]interface{}{
		"code":    code,
		"message": msg,
	}
}
