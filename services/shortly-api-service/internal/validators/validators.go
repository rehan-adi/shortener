package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type SignupValidator struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=2,max=15"`
	Password string `json:"password" validate:"required,min=6"`
}

type SigninValidator struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreateUrlValidator struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	ShortKey    string `json:"short_key" validate:"omitempty,min=2,max=20,shortkeychars"`
	Title       string `json:"title" validate:"omitempty,max=255"`
}

var (
	validate      = validator.New()
	validShortKey = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func init() {
	_ = validate.RegisterValidation("shortkeychars", func(fl validator.FieldLevel) bool {
		key := fl.Field().String()
		return validShortKey.MatchString(key)
	})
}

func ValidateSignupData(input SignupValidator) map[string]string {
	return validateStruct(input)
}

func ValidateSigninData(input SigninValidator) map[string]string {
	return validateStruct(input)
}

func ValidateCreateUrlData(input CreateUrlValidator) map[string]string {
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
