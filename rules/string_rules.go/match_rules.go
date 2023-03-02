package string_rules

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/notefan-golang/exceptions"
)

func PasswordMatch(comparePassword string) validation.RuleFunc {
	return func(value interface{}) error {
		password, _ := value.(string)
		if password != comparePassword {
			return exceptions.ValidationPasswordNotMatch
		}
		return nil
	}
}
