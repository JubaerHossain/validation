package main

import (
	"fmt"

	"github.com/JubaerHossain/validator"
)

type User struct {
	Name     string `json:"first_name"`
	Password string `json:"password"`
}

func main() {

	userInput := User{
		Name:     "John",
		Password: "123456",
	}
	validationErrors := ValidateUser(userInput)
	if validationErrors != nil {
		var errorMsgs []string
		for _, validationErr := range validationErrors {
			errorMsgs = append(errorMsgs, validationErr.Field+" : "+validationErr.Message)
		}
		fmt.Println(errorMsgs)
	}
}

func ValidateUser(user User) []validator.ValidationErrorItem {

	rules := []validator.ValidationRule{
		{
			Field:       "name",
			Description: "Name",
			Validations: []func(interface{}) validator.ValidationErrorItem{
				validator.RequiredValidation,
				validator.MinLengthValidation(3),
				validator.MaxLengthValidation(50),
			},
		},
		{
			Field:       "password",
			Description: "Password",
			Validations: []func(interface{}) validator.ValidationErrorItem{
				validator.RequiredValidation,
				validator.MinLengthValidation(6),
				validator.MaxLengthValidation(50),
			},
		},
	}

	return validator.Validate(user, rules)

}
