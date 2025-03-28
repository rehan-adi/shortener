package validators

import "github.com/go-playground/validator/v10"

type SignupValidator struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=2,max=15"`
	Password string `json:"password" validate:"required,min=6"`
}

var validate = validator.New()

func ValidateSignupData(input SignupValidator) map[string]string {
	return validateStruct(input)
}

func validateStruct(input interface{}) map[string]string {
	errs := make(map[string]string)
	if err := validate.Struct(input); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = "Invalid " + e.Field()
		}
	}
	return errs
}
