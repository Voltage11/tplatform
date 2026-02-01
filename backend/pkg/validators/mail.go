package validators

import "regexp"

func IsEmailValid(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Проверяем строку на соответствие шаблону
	isValid, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}

	return isValid
}
